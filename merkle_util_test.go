package eigenpodproofs

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/assert"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	contractBeaconChainProofs "github.com/Layr-Labs/eigenpod-proofs-generation/bindings"
	epgcommon "github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

var (
	oracleState       deneb.BeaconState
	oracleBlockHeader phase0.BeaconBlockHeader
	blockHeader       phase0.BeaconBlockHeader
	blockHeaderIndex  uint64
	block             deneb.BeaconBlock
	executionPayload  deneb.ExecutionPayload
	epp               *EigenPodProofs
	chainClient       *ChainClient
	beaconChainProofs *contractBeaconChainProofs.BeaconChainProofsTest
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

	epp.ComputeBeaconStateTopLevelRoots(&spec.VersionedBeaconState{Deneb: &oracleState})
	epp.ComputeBeaconStateRoot(&oracleState)
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
	_, proof, _ := beacon.ProveValidatorBalanceAgainstValidatorBalanceList(oracleState.Balances, uint64(validatorIndex))

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

func TestProveValidatorAgainstValidatorList(t *testing.T) {
	epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)

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

func TestStateRootAgainstLatestBlockHeaderProof(t *testing.T) {
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
