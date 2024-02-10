package eigenpodproofs

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/assert"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	contractBeaconChainProofs "github.com/Layr-Labs/eigenpod-proofs-generation/bindings"
	epgcommon "github.com/Layr-Labs/eigenpod-proofs-generation/common"
	common "github.com/ethereum/go-ethereum/common"
)

var (
	ctx                            context.Context
	oracleState                    deneb.BeaconState
	oracleBlockHeader              phase0.BeaconBlockHeader
	blockHeader                    phase0.BeaconBlockHeader
	blockHeaderIndex               uint64
	block                          deneb.BeaconBlock
	validatorIndex                 phase0.ValidatorIndex
	beaconBlockHeaderToVerifyIndex uint64
	executionPayload               deneb.ExecutionPayload
	epp                            *EigenPodProofs
	executionPayloadFieldRoots     []phase0.Root
	chainClient                    *ChainClient
	beaconChainProofs              *contractBeaconChainProofs.BeaconChainProofsTest
)

// var VALIDATOR_INDEX uint64 = 61068 //this is the index of a validator that has a partial withdrawal
var VALIDATOR_INDEX uint64 = 61336           //this is the index of a validator that has a full withdrawal.
var REPOINTED_VALIDATOR_INDEX uint64 = 61511 //this is the index of a validator that we use for the withdrawal credential proofs

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%slotsPerHistoricalRoot in the block_roots.

// var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 656
var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 2262

const DENEB_FORK_TIMESTAMP_GOERLI = uint64(1705473120)

var GOERLI_CHAIN_ID uint64 = 5

func TestMain(m *testing.M) {
	// Setup
	log.Println("Setting up suite")
	setupSuite()

	// Run tests
	code := m.Run()

	// Teardown
	log.Println("Tearing down suite")
	teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

func setupSuite() {
	log.Println("Setting up suite")
	rpc := "https://rpc.ankr.com/eth_goerli"
	privateKey := os.Getenv("PRIVATE_KEY")
	ethClient, err := ethclient.Dial(rpc)
	if err != nil {
		log.Panicf("failed to connect to the Ethereum client: %s", err)
	}

	chainClient, err = NewChainClient(ethClient, privateKey)
	if err != nil {
		log.Panicf("failed to create chain client: %s", err)
	}

	//BeaconChainProofs.sol deployment: https://goerli.etherscan.io/address/0xd132dD701d3980bb5d66A21e2340f263765e4a19#code
	contractAddress := common.HexToAddress("0xd132dD701d3980bb5d66A21e2340f263765e4a19")
	beaconChainProofs, err = contractBeaconChainProofs.NewBeaconChainProofsTest(contractAddress, chainClient)
	if err != nil {
		log.Panicf("failed to create beacon chain proofs contract: %s", err)
	}

	stateFile := "data/deneb_goerli_slot_7413760.json"
	oracleHeaderFile := "data/deneb_goerli_block_header_7413760.json"
	headerFile := "data/deneb_goerli_block_header_7426113.json"
	bodyFile := "data/deneb_goerli_block_7426113.json"

	//ParseCapellaBeaconState(stateFile)

	stateJSON, err := ParseJSONFile(stateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	ParseDenebBeaconStateFromJSON(*stateJSON, &oracleState)

	blockHeader, err = ExtractBlockHeader(headerFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	oracleBlockHeader, err = ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with oracle block header", err)
	}

	block, err = ExtractBlockDeneb(bodyFile)
	if err != nil {
		fmt.Println("error with block body", err)
	}

	executionPayload = *block.Body.ExecutionPayload

	blockHeaderIndex = uint64(blockHeader.Slot) % beacon.SlotsPerHistoricalRoot

	epp, err = NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	if err != nil {
		fmt.Println("error in NewEigenPodProofs", err)
	}

	executionPayloadFieldRoots, _ = beacon.ComputeExecutionPayloadFieldRootsDeneb(block.Body.ExecutionPayload)
}

func teardownSuite() {
	// Any cleanup you want to perform should go here
	fmt.Println("all done!")
}

// verifies that the "ProveValidatorContainers" call, which the backend calls, returns valid proofs
func TestProveValidatorContainers(t *testing.T) {

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(&oracleBlockHeader, &versionedOracleState, []uint64{VALIDATOR_INDEX})
	if err != nil {
		fmt.Println("error", err)
	}

	stateRootProof := verifyValidatorFieldsCallParams.StateRootProof
	validatorFieldsProofs := verifyValidatorFieldsCallParams.ValidatorFieldsProofs
	validatorIndices := verifyValidatorFieldsCallParams.ValidatorIndices

	flag := verifyStateRootAgainstBlockHeaderProof(oracleBlockHeader, oracleState, stateRootProof.StateRootProof)
	assert.True(t, flag, "State Root Proof %v failed")
	flag = verifyValidatorAgainstBeaconState(&oracleState, validatorFieldsProofs[0], validatorIndices[0])
	assert.True(t, flag, "State Root Proof %v failed")
}

func TestProveWithdrawals(t *testing.T) {
	oracleStateFile := "data/deneb_goerli_slot_7431952.json"
	oracleHeaderFile := "data/deneb_goerli_block_header_7431952.json"
	oracleStateJSON, err := ParseJSONFile(oracleStateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	oracleBlockHeader, err = ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	ParseDenebBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error creating versioned state", err)
		return
	}

	historicalSummaryStateJSON, err := ParseJSONFile("data/deneb_goerli_slot_7421952.json")
	if err != nil {
		fmt.Println("error parsing historicalSummaryState JSON")
	}
	var historicalSummaryState deneb.BeaconState
	ParseDenebBeaconStateFromJSON(*historicalSummaryStateJSON, &historicalSummaryState)
	historicalSummaryStateBlockRoots := historicalSummaryState.BlockRoots

	withdrawalBlock, err := ExtractBlockDeneb("data/deneb_goerli_block_7421951.json")
	if err != nil {
		fmt.Println("block.UnmarshalJSON error", err)
	}
	withdrawalBlockHeader, err := ExtractBlockHeader("data/deneb_goerli_block_header_7421951.json")
	if err != nil {
		fmt.Println("blockHeader.UnmarshalJSON error", err)
	}

	versionedWithdrawalBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	withdrawalIndex := uint64(0)
	withdrawalValidatorIndex := uint64(627559) //this is the index of the validator with the first withdrawal in the withdrawalBlock 7421951

	verifyAndProcessWithdrawalCallParams, err := epp.ProveWithdrawals(
		&oracleBlockHeader,
		&versionedOracleState,
		[][]phase0.Root{historicalSummaryStateBlockRoots},
		[]*spec.VersionedSignedBeaconBlock{&versionedWithdrawalBlock},
		[]uint64{withdrawalValidatorIndex},
	)
	if err != nil {
		fmt.Println("error", err)
	}

	executionPayloadRoot, err := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		fmt.Println("error", err)
	}

	oracleStateRoot, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("error", err)
		return
	}

	flag := verifyStateRootAgainstBlockHeaderProof(oracleBlockHeader, oracleState, verifyAndProcessWithdrawalCallParams.StateRootProof.StateRootProof)
	assert.True(t, flag, "State Root Proof %v failed")

	flag = verifyValidatorAgainstBeaconState(&oracleState, verifyAndProcessWithdrawalCallParams.ValidatorFieldsProofs[0], withdrawalValidatorIndex)
	assert.True(t, flag, "Validator Fields Proof %v failed")

	flag = verifyWithdrawalAgainstExecutionPayload(executionPayloadRoot, verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalProof, withdrawalIndex, withdrawalBlock.Body.ExecutionPayload.Withdrawals[0])
	assert.True(t, flag, "Withdrawal Proof %v failed")

	flag = verifyTimestampAgainstExecutionPayload(executionPayloadRoot, verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampProof, withdrawalBlock.Body.ExecutionPayload.Timestamp)
	assert.True(t, flag, "Timestamp Proof %v failed")

	flag = verifyBlockRootAgainstBeaconStateViaHistoricalSummaries(
		oracleStateRoot,
		verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryBlockRootProof,
		withdrawalBlockHeader,
		verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex,
		verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex,
	)
	assert.True(t, flag, "Historical Summary Block Root Proof %v failed")
}

func TestUnmarshalSSZVersionedBeaconStateDeneb(t *testing.T) {
	oracleStateBytes, err := oracleState.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}

	versionedBeaconState, err := beacon.UnmarshalSSZVersionedBeaconState(oracleStateBytes)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconState.Version, spec.DataVersionDeneb, "Version %v failed")

	versionedBeaconStateBytes, err := versionedBeaconState.Deneb.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconStateBytes, oracleStateBytes, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestUnmarshalSSZVersionedBeaconStateCapella(t *testing.T) {
	var capellaState capella.BeaconState
	capellaStateJSON, err := ParseJSONFileCapella("data/goerli_slot_6409723.json")
	if err != nil {
		fmt.Println("error", err)
	}
	ParseCapellaBeaconStateFromJSON(*capellaStateJSON, &capellaState)

	capellaStateBytes, err := capellaState.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}

	versionedBeaconState, err := beacon.UnmarshalSSZVersionedBeaconState(capellaStateBytes)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconState.Version, spec.DataVersionCapella, "Version %v failed")

	versionedBeaconStateBytes, err := versionedBeaconState.Capella.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconStateBytes, capellaStateBytes, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestMarshalSSZVersionedBeaconStateDeneb(t *testing.T) {
	oracleStateBytes, err := oracleState.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}
	versionedBeaconState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
	}

	versionedBeaconStateBytes, err := beacon.MarshalSSZVersionedBeaconState(versionedBeaconState)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconStateBytes, oracleStateBytes, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestMarshalSSZVersionedBeaconStateCapella(t *testing.T) {
	var capellaState capella.BeaconState
	capellaStateJSON, err := ParseJSONFileCapella("data/goerli_slot_6409723.json")
	if err != nil {
		fmt.Println("error", err)
	}
	ParseCapellaBeaconStateFromJSON(*capellaStateJSON, &capellaState)

	capellaStateBytes, err := capellaState.MarshalSSZ()
	if err != nil {
		fmt.Println("error", err)
	}
	versionedBeaconState, err := beacon.CreateVersionedState(&capellaState)
	if err != nil {
		fmt.Println("error", err)
	}

	versionedBeaconStateBytes, err := beacon.MarshalSSZVersionedBeaconState(versionedBeaconState)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBeaconStateBytes, capellaStateBytes, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestGenerateWithdrawalCredentialsProof(t *testing.T) {

	// picking up one random validator index
	validatorIndex := phase0.ValidatorIndex(REPOINTED_VALIDATOR_INDEX)

	beaconStateTopLevelRoots, err := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)
	if err != nil {
		fmt.Println("error reading beaconStateTopLevelRoots")
	}

	proof, err := epp.ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleState.Slot, oracleState.Validators, uint64(validatorIndex))
	if err != nil {
		fmt.Println(err)
	}

	flag := verifyValidatorAgainstBeaconState(&oracleState, proof, uint64(validatorIndex))

	assert.True(t, flag, "Proof %v failed")
}

func TestCreateVersionedSignedBlockDeneb(t *testing.T) {
	block := deneb.BeaconBlock{}
	versionedBlock, err := beacon.CreateVersionedSignedBlock(block)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBlock.Version, spec.DataVersionDeneb, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestCreateVersionedSignedBlockCapella(t *testing.T) {
	block := capella.BeaconBlock{}
	versionedBlock, err := beacon.CreateVersionedSignedBlock(block)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedBlock.Version, spec.DataVersionCapella, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestCreateVersionedSignedBlockAltair(t *testing.T) {
	block := altair.BeaconBlock{}
	_, err := beacon.CreateVersionedSignedBlock(block)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.NotNil(t, err, "error %v was nil")
}

func TestCreateVersionedBeaconStateDeneb(t *testing.T) {
	oracleState := deneb.BeaconState{}
	versionedState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedState.Version, spec.DataVersionDeneb, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestCreateVersionedBeaconStateCapella(t *testing.T) {
	state := capella.BeaconState{}
	versionedState, err := beacon.CreateVersionedState(&state)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Equal(t, versionedState.Version, spec.DataVersionCapella, "Version %v failed")
	assert.Nil(t, err, "Error %v failed")
}

func TestCreateVersionedBeaconStateAltair(t *testing.T) {
	state := altair.BeaconState{}
	_, err := beacon.CreateVersionedState(&state)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.NotNil(t, err, "error %v was nil")
}

func TestProveValidatorBalanceAgainstValidatorBalanceList(t *testing.T) {

	validatorIndex := phase0.ValidatorIndex(REPOINTED_VALIDATOR_INDEX)
	proof, _ := beacon.ProveValidatorBalanceAgainstValidatorBalanceList(oracleState.Balances, uint64(validatorIndex))

	beaconStateTopLevelRoots, _ := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)
	root := beaconStateTopLevelRoots.BalancesRoot

	balanceRootList, err := beacon.GetBalanceRoots(oracleState.Balances)
	if err != nil {
		fmt.Println("error", err)
	}

	balanceIndex := validatorIndex / 4

	leaf := balanceRootList[balanceIndex]

	flag := epgcommon.ValidateProof(*root, proof, leaf, uint64(balanceIndex))
	if flag != true {
		fmt.Println("balance proof failed")
	}
	assert.True(t, flag, "Proof %v failed")
}

func TestProveBeaconTopLevelRootAgainstBeaconState(t *testing.T) {

	// get the oracle state root for a merkle tree with top level roots as the leaves
	beaconStateTopLevelRoots, err := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)
	if err != nil {
		fmt.Println("error")
	}

	// compute the Merkle proof for the inclusion of Validators Root as a leaf
	validatorsRootProof, err := beacon.ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, beacon.ValidatorListIndex)
	if err != nil {
		fmt.Println("error")
	}

	// getting Merkle root of the BeaconStateRoot Merkle tree from attestation's code
	beaconStateRoot, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	// validation of the proof
	// get the leaf denoting the validatorsRoot in the BeaconStateRoot Merkle tree
	leaf := beaconStateTopLevelRoots.ValidatorsRoot
	flag := epgcommon.ValidateProof(beaconStateRoot, validatorsRootProof, *leaf, beacon.ValidatorListIndex)
	if flag != true {
		fmt.Println("error")
	}
	// fmt.Println("flag", flag)

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetHistoricalSummariesBlockRootsProofProof(t *testing.T) {

	//curl -H "Accept: application/json" https://data.spiceai.io/goerli/beacon/eth/v2/debug/beacon/states/7431952 -o deneb_goerli_slot_7431952.json --header 'X-API-Key: 343035|8b6ddd9b31f54c07b3fc18282b30f61c'
	currentBeaconStateJSON, err := ParseJSONFile("data/deneb_goerli_slot_7431952.json")

	if err != nil {
		fmt.Println("error parsing currentBeaconStateJSON")
	}

	//this is not the beacon state of the slot containing the old withdrawal we want to prove but rather
	// its the state that was merkleized to create a historical summary containing the slot that has that withdrawal
	//, ie, 7421952 mod 8192 = 0 and 7421952 - 7421951 < 8192
	oldBeaconStateJSON, err := ParseJSONFile("data/deneb_goerli_slot_7421952.json")
	if err != nil {
		fmt.Println("error parsing oldBeaconStateJSON")
	}

	var blockHeader phase0.BeaconBlockHeader
	blockHeader, err = ExtractBlockHeader("data/deneb_goerli_block_header_7421951.json")

	if err != nil {
		fmt.Println("blockHeader.UnmarshalJSON error", err)
	}

	var currentBeaconState deneb.BeaconState
	var oldBeaconState deneb.BeaconState

	ParseDenebBeaconStateFromJSON(*currentBeaconStateJSON, &currentBeaconState)
	ParseDenebBeaconStateFromJSON(*oldBeaconStateJSON, &oldBeaconState)

	currentBeaconStateTopLevelRoots, _ := beacon.ComputeBeaconStateTopLevelRootsDeneb(&currentBeaconState)

	if err != nil {
		fmt.Println("error")
	}

	historicalSummaryIndex := uint64(271) //7421951 - FIRST_CAPELLA_SLOT_GOERLI // 8192
	beaconBlockHeaderToVerifyIndex = 8191 //(7421951 mod 8192)
	if err != nil {
		fmt.Println("error", err)
	}

	oldBlockRoots := oldBeaconState.BlockRoots

	historicalSummaryBlockHeaderProof, err := beacon.ProveBlockRootAgainstBeaconStateViaHistoricalSummaries(
		currentBeaconStateTopLevelRoots,
		currentBeaconState.HistoricalSummaries,
		oldBlockRoots,
		historicalSummaryIndex,
		beaconBlockHeaderToVerifyIndex,
	)

	if err != nil {
		fmt.Println("error")
	}

	currentBeaconStateRoot, _ := currentBeaconState.HashTreeRoot()

	flag := verifyBlockRootAgainstBeaconStateViaHistoricalSummaries(currentBeaconStateRoot, historicalSummaryBlockHeaderProof, blockHeader, beaconBlockHeaderToVerifyIndex, historicalSummaryIndex)

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetHistoricalSummariesBlockRootsProofProofCapellaAgainstDeneb(t *testing.T) {

	//curl -H "Accept: application/json" https://data.spiceai.io/goerli/beacon/eth/v2/debug/beacon/states/7431952 -o deneb_goerli_slot_7431952.json --header 'X-API-Key: 343035|8b6ddd9b31f54c07b3fc18282b30f61c'
	currentBeaconStateJSON, err := ParseJSONFile("data/deneb_goerli_slot_7431952.json")

	if err != nil {
		fmt.Println("error parsing currentBeaconStateJSON")
	}

	//this is not the beacon state of the slot containing the old withdrawal we want to proof but rather
	// its the state that was merklized to create a historical summary containing the slot that has that withdrawal, ie, 7421952 mod 8192 = 0
	oldBeaconStateJSON, err := ParseJSONFileCapella("data/goerli_slot_6397952.json")
	if err != nil {
		fmt.Println("error parsing oldBeaconStateJSON", err)
	}

	var blockHeader phase0.BeaconBlockHeader
	//blockHeader, err = ExtractBlockHeader("data/goerli_block_header_6397852.json")
	blockHeader, err = ExtractBlockHeader("data/goerli_block_header_6397852.json")

	if err != nil {
		fmt.Println("blockHeader.UnmarshalJSON error", err)
	}

	var currentBeaconState deneb.BeaconState
	var oldBeaconState capella.BeaconState

	ParseDenebBeaconStateFromJSON(*currentBeaconStateJSON, &currentBeaconState)
	ParseCapellaBeaconStateFromJSON(*oldBeaconStateJSON, &oldBeaconState)

	currentBeaconStateTopLevelRoots, _ := beacon.ComputeBeaconStateTopLevelRootsDeneb(&currentBeaconState)
	//oldBeaconStateTopLevelRoots, _ := ComputeBeaconStateTopLevelRoots(&oldBeaconState)

	if err != nil {
		fmt.Println("error")
	}

	historicalSummaryIndex := uint64(146)
	beaconBlockHeaderToVerifyIndex = 8092 //(7421951 mod 8192)
	if err != nil {
		fmt.Println("error", err)
	}

	oldBlockRoots := oldBeaconState.BlockRoots

	historicalSummaryBlockHeaderProof, err := beacon.ProveBlockRootAgainstBeaconStateViaHistoricalSummaries(
		currentBeaconStateTopLevelRoots,
		currentBeaconState.HistoricalSummaries,
		oldBlockRoots,
		historicalSummaryIndex,
		beaconBlockHeaderToVerifyIndex,
	)

	if err != nil {
		fmt.Println("error")
	}

	currentBeaconStateRoot, _ := currentBeaconState.HashTreeRoot()

	flag := verifyBlockRootAgainstBeaconStateViaHistoricalSummaries(currentBeaconStateRoot, historicalSummaryBlockHeaderProof, blockHeader, beaconBlockHeaderToVerifyIndex, historicalSummaryIndex)

	assert.True(t, flag, "Proof %v failed\n")

}

func TestProveValidatorAgainstValidatorList(t *testing.T) {

	// picking up one random validator index
	validatorIndex := phase0.ValidatorIndex(10000)

	// get the validators field
	validators := oracleState.Validators

	// get the Merkle proof for inclusion
	validatorProof, err := epp.ProveValidatorAgainstValidatorList(0, validators, uint64(validatorIndex))
	if err != nil {
		fmt.Println("error")
	}

	// verify the proof
	// get the leaf corresponding to validatorIndex
	leaf, err := validators[validatorIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	// get the oracle state root for a merkle tree with top level roots as the leaves
	beaconStateTopLevelRoots, err := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)
	if err != nil {
		fmt.Println("error")
	}

	// calling the proof verification func
	flag := epgcommon.ValidateProof(*beaconStateTopLevelRoots.ValidatorsRoot, validatorProof, leaf, uint64(validatorIndex))
	if flag != true {
		fmt.Println("error")
	}
	// fmt.Println("flag", flag)

	assert.True(t, flag, "Proof %v failed\n")
}

func TestComputeBlockSlotProof(t *testing.T) {
	// get the proof for slot in the block header
	blockHeaderSlotProof, err := beacon.ProveSlotAgainstBlockHeader(&blockHeader)
	if err != nil {
		fmt.Println("error", err)
	}

	// get the hash of the slot - this will be the leaf of the merkle tree
	var slotHashRoot phase0.Root
	hh := ssz.NewHasher()
	hh.PutUint64(uint64(blockHeader.Slot))
	copy(slotHashRoot[:], hh.Hash())

	// get the block header root which will be used as a root of the Merkle tree
	// Note that the blockHeader was obtained from the actual Block header
	beaconBlockHeaderRoot, err := blockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("error:", err)
	}

	// calling the proof verification function
	flag := epgcommon.ValidateProof(beaconBlockHeaderRoot, blockHeaderSlotProof, slotHashRoot, beacon.SlotIndex)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestProveBlockBodyAgainstBlockHeader(t *testing.T) {

	// get the proof for block body in the block header
	blockHeaderBlockBodyProof, err := beacon.ProveBlockBodyAgainstBlockHeader(&blockHeader)
	if err != nil {
		fmt.Println("error", err)
	}

	// get the hash of the block body root - this will be the leaf of the merkle tree
	var blockBodyHashRoot phase0.Root
	hh := ssz.NewHasher()
	hh.PutBytes(blockHeader.BodyRoot[:])
	copy(blockBodyHashRoot[:], hh.Hash())

	// get the block header root which will be used as a root of the Merkle tree
	// Note that the blockHeader was obtained from the actual Block header
	beaconBlockHeaderRoot, err := blockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("error:", err)
	}

	// calling the proof verification function
	flag := epgcommon.ValidateProof(beaconBlockHeaderRoot, blockHeaderBlockBodyProof, blockBodyHashRoot, beacon.BeaconBlockBodyRootIndex)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestComputeExecutionPayloadHeader(t *testing.T) {

	// get the proof for execution payload in the block body
	beaconBlockBodyProof, _, err := beacon.ProveExecutionPayloadAgainstBlockBodyDeneb(block.Body)
	if err != nil {
		fmt.Println("error", err)
	}

	// get the hash root of the actual execution payload
	var executionPayloadHashRoot phase0.Root
	hh := ssz.NewHasher()
	{
		if err = block.Body.ExecutionPayload.HashTreeRootWith(hh); err != nil {
			fmt.Println("error", err)
		}
		copy(executionPayloadHashRoot[:], hh.Hash())
	}

	// get the body root in the beacon block header -  will be used as the Merkle root
	blockHeaderBodyRoot := blockHeader.BodyRoot

	// calling the proof verification function
	flag := epgcommon.ValidateProof(blockHeaderBodyRoot, beaconBlockBodyProof, executionPayloadHashRoot, beacon.ExecutionPayloadIndex)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestStateRootAgainstLatestBlockHeaderProof(t *testing.T) {

	// this is the state where the latest block header from the oracle was taken.  This is the next slot after
	// the state we want to prove things about (remember latestBlockHeader.state_root = previous slot's state root)
	// oracleStateJSON, err := ParseJSONFile("data/historical_summary_proof/goerli_slot_6399999.json")
	// var oracleState deneb.BeaconState
	// ParseCapellaBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	var blockHeader phase0.BeaconBlockHeader
	blockHeader, err := ExtractBlockHeader("data/deneb_goerli_block_header_7413760.json")
	if err != nil {
		fmt.Println("error with block header", err)
	}

	//the state from the prev slot which contains shit we wanna prove about
	stateToProveJSON, err := ParseJSONFile("data/deneb_goerli_slot_7413760.json")
	if err != nil {
		fmt.Println("error with parsing JSON state file", err)
	}

	var stateToProve deneb.BeaconState
	ParseDenebBeaconStateFromJSON(*stateToProveJSON, &stateToProve)

	proof, err := beacon.ProveStateRootAgainstBlockHeader(&blockHeader)
	if err != nil {
		fmt.Println("Error in generating proof", err)
	}

	flag := verifyStateRootAgainstBlockHeaderProof(blockHeader, stateToProve, proof)
	assert.True(t, flag, "Proof %v failed")
}

func TestGetExecutionPayloadProof(t *testing.T) {

	// get the proof for execution payload in the block body

	executionPayloadProof, _, _ := beacon.ProveExecutionPayloadAgainstBlockHeaderDeneb(&blockHeader, block.Body)

	// get the hash root of the actual execution payload
	var executionPayloadHashRoot, _ = block.Body.ExecutionPayload.HashTreeRoot()

	// get the body root in the beacon block header -  will be used as the Merkle root
	root, _ := blockHeader.HashTreeRoot()

	index := beacon.BeaconBlockBodyRootIndex<<(beacon.BlockBodyMerkleSubtreeNumLayers) | beacon.ExecutionPayloadIndex

	// calling the proof verification function
	flag := epgcommon.ValidateProof(root, executionPayloadProof, executionPayloadHashRoot, index)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed")
}

func TestComputeWithdrawalsListProof(t *testing.T) {

	withdrawalsListProof, err := beacon.ProveWithdrawalListAgainstExecutionPayload(executionPayloadFieldRoots)
	if err != nil {
		fmt.Println("error!", err)
	}

	var withdrawalsHashRoot phase0.Root
	hh := ssz.NewHasher()

	{
		subIndx := hh.Index()
		num := uint64(len(block.Body.ExecutionPayload.Withdrawals))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			fmt.Println("error!", err)
		}
		for _, elem := range block.Body.ExecutionPayload.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				fmt.Println("error 4", err)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(withdrawalsHashRoot[:], hh.Hash())
		hh.Reset()
	}

	var executionPayloadHashRoot phase0.Root
	{
		if err = block.Body.ExecutionPayload.HashTreeRootWith(hh); err != nil {
			fmt.Println("error hel", err)
		}
		copy(executionPayloadHashRoot[:], hh.Hash())
	}
	flag := epgcommon.ValidateProof(executionPayloadHashRoot, withdrawalsListProof, withdrawalsHashRoot, beacon.WithdrawalsIndex)
	if flag != true {
		fmt.Println("Proof Failed")
	}
	assert.True(t, flag, "Proof %v failed\n")

}

func TestComputeIndividualWithdrawalProof(t *testing.T) {

	// picking up one random validator index
	withdrawalIndex := uint8(0)

	// get the validators field
	withdrawals := block.Body.ExecutionPayload.Withdrawals

	// get the Merkle proof for inclusion
	withdrawalProof, err := beacon.ProveWithdrawalAgainstWithdrawalList(withdrawals, withdrawalIndex)
	if err != nil {
		fmt.Println("error")
	}

	// verify the proof
	// get the leaf corresponding to validatorIndex
	leaf, err := withdrawals[withdrawalIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	var withdrawalsHashRoot phase0.Root
	hh := ssz.NewHasher()

	{
		subIndx := hh.Index()
		num := uint64(len(block.Body.ExecutionPayload.Withdrawals))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			fmt.Println("error", err)
		}
		for _, elem := range block.Body.ExecutionPayload.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				fmt.Println("error", err)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(withdrawalsHashRoot[:], hh.Hash())
		hh.Reset()
	}

	// calling the proof verification func
	flag := epgcommon.ValidateProof(withdrawalsHashRoot, withdrawalProof, leaf, uint64(withdrawalIndex))
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetWithdrawalProof(t *testing.T) {

	// picking up one random validator index
	withdrawalIndex := 0

	withdrawalProof, _ := beacon.ProveWithdrawalAgainstExecutionPayload(executionPayloadFieldRoots, block.Body.ExecutionPayload.Withdrawals, uint8(withdrawalIndex))

	executionPayloadRoot, _ := block.Body.ExecutionPayload.HashTreeRoot()

	// calling the proof verification func
	flag := verifyWithdrawalAgainstExecutionPayload(executionPayloadRoot, withdrawalProof, uint64(withdrawalIndex), block.Body.ExecutionPayload.Withdrawals[withdrawalIndex])

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetTimestampProof(t *testing.T) {

	// get the block number
	executionPayloadFields := block.Body.ExecutionPayload

	// get the Merkle proof for inclusion
	timestampProof, _ := beacon.ProveTimestampAgainstExecutionPayload(executionPayloadFieldRoots)

	root, err := block.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	flag := verifyTimestampAgainstExecutionPayload(root, timestampProof, executionPayloadFields.Timestamp)

	assert.True(t, flag, "Proof %v failed")
}

func TestGetValidatorProof(t *testing.T) {
	// picking up one random validator index
	validatorIndex := uint64(VALIDATOR_INDEX)

	beaconStateTopLevelRoots, err := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)
	if err != nil {
		fmt.Println("error reading beaconStateTopLevelRoots")
	}

	validatorProof, _ := epp.ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleState.Slot, oracleState.Validators, uint64(validatorIndex))

	flag := verifyValidatorAgainstBeaconState(&oracleState, validatorProof, validatorIndex)

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetSlotProof(t *testing.T) {
	// picking up one random validator index
	slot := blockHeader.Slot

	buf := make([]byte, 32)
	binary.LittleEndian.PutUint64(buf, uint64(slot))
	var bytes32 [32]byte
	copy(bytes32[:], buf[:32])

	proof, _ := beacon.ProveSlotAgainstBlockHeader(&blockHeader)

	root, _ := blockHeader.HashTreeRoot()

	hh := ssz.NewHasher()
	hh.PutUint64(uint64(slot))

	leaf := ConvertTo32ByteArray(hh.Hash())

	flag := epgcommon.ValidateProof(root, proof, leaf, 0)
	if flag != true {
		fmt.Println("error")
	}
	assert.True(t, flag, "Proof %v failed\n")
}

type Proofs struct {
	Slot                  uint64   `json:"slot"`
	ValidatorIndex        uint64   `json:"validatorIndex"`
	WithdrawalIndex       uint64   `json:"withdrawalIndex"`
	BlockHeaderRootIndex  uint64   `json:"blockHeaderRootIndex"`
	BeaconStateRoot       string   `json:"beaconStateRoot"`
	SlotRoot              string   `json:"slotRoot"`
	BlockNumberRoot       string   `json:"blockNumberRoot"`
	BlockHeaderRoot       string   `json:"blockHeaderRoot"`
	BlockBodyRoot         string   `json:"blockBodyRoot"`
	ExecutionPayloadRoot  string   `json:"executionPayloadRoot"`
	BlockHeaderProof      []string `json:"BlockHeaderProof"`
	SlotProof             []string `json:"SlotProof"`
	WithdrawalProof       []string `json:"WithdrawalProof"`
	ValidatorProof        []string `json:"ValidatorProof"`
	BlockNumberProof      []string `json:"BlockNumberProof"`
	ExecutionPayloadProof []string `json:"ExecutionPayloadProof"`
	ValidatorFields       []string `json:"ValidatorFields"`
	WithdrawalFields      []string `json:"WithdrawalFields"`
}

func ParseJSONFile(filePath string) (*beaconStateJSONDeneb, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("error with reading file")
		return nil, err
	}

	var beaconState beaconStateVersionDeneb
	err = json.Unmarshal(data, &beaconState)
	if err != nil {
		fmt.Println("error with beaconState JSON unmarshalling")
		return nil, err
	}

	actualData := beaconState.Data
	return &actualData, nil
}

func verifyStateRootAgainstBlockHeaderProof(oracleBlockHeader phase0.BeaconBlockHeader, oracleState deneb.BeaconState, proof epgcommon.Proof) bool {
	root, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}
	leaf, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}

	flag := epgcommon.ValidateProof(root, proof, leaf, 3)
	if flag != true {
		fmt.Println("this error")
	}
	return flag
}

func verifyValidatorAgainstBeaconState(oracleState *deneb.BeaconState, proof epgcommon.Proof, validatorIndex uint64) bool {
	leaf, err := oracleState.Validators[validatorIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root")
	}

	root, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root of beacon state")
	}

	index := beacon.ValidatorListIndex<<(beacon.ValidatorListMerkleSubtreeNumLayers+1) | validatorIndex

	flag := epgcommon.ValidateProof(root, proof, leaf, index)
	if flag != true {
		fmt.Println("error")
	}
	return flag
}

func verifyWithdrawalAgainstExecutionPayload(executionPayloadRoot phase0.Root, proof epgcommon.Proof, withdrawalIndex uint64, withdrawal *capella.Withdrawal) bool {
	leaf, err := withdrawal.HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	withdrawalRelativeToELPayloadIndex := beacon.WithdrawalsIndex<<(beacon.WithdrawalListMerkleSubtreeNumLayers+1) | uint64(withdrawalIndex)

	return epgcommon.ValidateProof(executionPayloadRoot, proof, leaf, withdrawalRelativeToELPayloadIndex)

}

func verifyTimestampAgainstExecutionPayload(executionPayloadRoot phase0.Root, proof epgcommon.Proof, timestamp uint64) bool {
	hh := ssz.NewHasher()
	hh.PutUint64(timestamp)
	leaf := ConvertTo32ByteArray(hh.Hash())

	return epgcommon.ValidateProof(executionPayloadRoot, proof, leaf, beacon.TimestampIndex)
}

func verifyBlockRootAgainstBeaconStateViaHistoricalSummaries(oracleBeaconStateRoot phase0.Root, proof epgcommon.Proof, beaconBlockHeaderToVerify phase0.BeaconBlockHeader, beaconBlockHeaderToVerifyIndex uint64, historicalSummaryIndex uint64) bool {
	historicalBlockHeaderIndex := beacon.HistoricalSummaryListIndex<<((beacon.HistoricalSummaryListMerkleSubtreeNumLayers+1)+1+(beacon.BlockRootsMerkleSubtreeNumLayers)) |
		historicalSummaryIndex<<(1+beacon.BlockRootsMerkleSubtreeNumLayers) |
		beacon.BlockSummaryRootIndex<<(beacon.BlockRootsMerkleSubtreeNumLayers) | beaconBlockHeaderToVerifyIndex

	beaconBlockHeaderToVerifyRoot, err := beaconBlockHeaderToVerify.HashTreeRoot()
	if err != nil {
		fmt.Println("beaconBlockHeaderToVerifyRoot error", err)
		return false
	}
	return epgcommon.ValidateProof(oracleBeaconStateRoot, proof, beaconBlockHeaderToVerifyRoot, historicalBlockHeaderIndex)
}
