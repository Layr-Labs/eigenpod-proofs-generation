package eigenpodproofs

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/assert"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

var (
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
)

// var VALIDATOR_INDEX uint64 = 61068 //this is the index of a validator that has a partial withdrawal
var VALIDATOR_INDEX uint64 = 61336           //this is the index of a validator that has a full withdrawal.
var REPOINTED_VALIDATOR_INDEX uint64 = 61511 //this is the index of a validator that we use for the withdrawal credential proofs

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%slotsPerHistoricalRoot in the block_roots.

// var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 656
var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 2262

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
	stateFile := "data/deneb_goerli_slot_7413760.json"
	oracleHeaderFile := "data/deneb_goerli_block_header_7413760.json"
	headerFile := "data/deneb_goerli_block_header_7426113.json"
	bodyFile := "data/deneb_goerli_block_7426113.json"

	//ParseCapellaBeaconState(stateFile)

	stateJSON, err := parseJSONFile(stateFile)
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

	block, err = ExtractBlock(bodyFile)
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

	root, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}
	leaf, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}

	flag := common.ValidateProof(root, stateRootProof.StateRootProof, leaf, 3)
	if flag != true {
		fmt.Println("this error")
	}
	assert.True(t, flag, "Proof %v failed")

	leaf, err = oracleState.Validators[validatorIndices[0]].HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root")
	}

	root, err = oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root of beacon state")
	}

	index := beacon.ValidatorListIndex<<(beacon.ValidatorListMerkleSubtreeNumLayers+1) | validatorIndices[0]

	flag = common.ValidateProof(root, validatorFieldsProofs[0], leaf, index)
	if flag != true {
		fmt.Println("error")
	}
}

func TestProveWithdrawals(t *testing.T) {

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error creating versioned state", err)
		return
	}

	historicalSummaryStateJSON, err := parseJSONFile("data/deneb_goerli_slot_7421952.json")
	if err != nil {
		fmt.Println("error parsing historicalSummaryState JSON")
	}
	fmt.Println("historicalSummaryStateJSON", historicalSummaryStateJSON.Slot)
	fmt.Println("versionedOracleState", versionedOracleState.Deneb.Slot)
	var historicalSummaryState deneb.BeaconState
	historicalSummaryStateBlockRoots := historicalSummaryState.BlockRoots
	ParseDenebBeaconStateFromJSON(*historicalSummaryStateJSON, &historicalSummaryState)

	fmt.Println("historicalSummaryStateJSON", len(historicalSummaryStateBlockRoots))

	var withdrawalBlock deneb.BeaconBlock
	withdrawalBlock, err = ExtractBlock("data/deneb_goerli_block_7421951.json")
	fmt.Println("withdrawalBlock", withdrawalBlock.Slot)

	if err != nil {
		fmt.Println("block.UnmarshalJSON error", err)
	}

	versionedWithdrawalBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	// fmt.Print("versionedWithdrawalBlock", versionedWithdrawalBlock.Deneb.Message.Slot)

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

	fmt.Println("verifyAndProcessWithdrawalCallParams", verifyAndProcessWithdrawalCallParams.OracleTimestamp)

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
	leaf, err := oracleState.Validators[validatorIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root")
	}

	root, err := oracleState.HashTreeRoot()
	if err != nil {
		fmt.Println("error with hash tree root of beacon state")
	}

	index := beacon.ValidatorListIndex<<(beacon.ValidatorListMerkleSubtreeNumLayers+1) | uint64(validatorIndex)

	flag := common.ValidateProof(root, proof, leaf, index)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed")
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

	flag := common.ValidateProof(*root, proof, leaf, uint64(balanceIndex))
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
	flag := common.ValidateProof(beaconStateRoot, validatorsRootProof, *leaf, beacon.ValidatorListIndex)
	if flag != true {
		fmt.Println("error")
	}
	// fmt.Println("flag", flag)

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetHistoricalSummariesBlockRootsProofProof(t *testing.T) {

	//curl -H "Accept: application/json" https://data.spiceai.io/goerli/beacon/eth/v2/debug/beacon/states/7431952 -o deneb_goerli_slot_7431952.json --header 'X-API-Key: 343035|8b6ddd9b31f54c07b3fc18282b30f61c'
	currentBeaconStateJSON, err := parseJSONFile("data/deneb_goerli_slot_7431952.json")

	if err != nil {
		fmt.Println("error parsing currentBeaconStateJSON")
	}

	//this is not the beacon state of the slot containing the old withdrawal we want to prove but rather
	// its the state that was merkleized to create a historical summary containing the slot that has that withdrawal
	//, ie, 7421952 mod 8192 = 0 and 7421952 - 7421951 < 8192
	oldBeaconStateJSON, err := parseJSONFile("data/deneb_goerli_slot_7421952.json")
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
	beaconBlockHeaderToVerify, err := blockHeader.HashTreeRoot()
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

	historicalBlockHeaderIndex := beacon.HistoricalSummaryListIndex<<((beacon.HistoricalSummaryListMerkleSubtreeNumLayers+1)+1+(beacon.BlockRootsMerkleSubtreeNumLayers)) |
		historicalSummaryIndex<<(1+beacon.BlockRootsMerkleSubtreeNumLayers) |
		beacon.BlockSummaryRootIndex<<(beacon.BlockRootsMerkleSubtreeNumLayers) | beaconBlockHeaderToVerifyIndex

	flag := common.ValidateProof(currentBeaconStateRoot, historicalSummaryBlockHeaderProof, beaconBlockHeaderToVerify, historicalBlockHeaderIndex)
	if flag != true {
		fmt.Println("error 2")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetHistoricalSummariesBlockRootsProofProofCapellaAgainstDeneb(t *testing.T) {

	//curl -H "Accept: application/json" https://data.spiceai.io/goerli/beacon/eth/v2/debug/beacon/states/7431952 -o deneb_goerli_slot_7431952.json --header 'X-API-Key: 343035|8b6ddd9b31f54c07b3fc18282b30f61c'
	currentBeaconStateJSON, err := parseJSONFile("data/deneb_goerli_slot_7431952.json")

	if err != nil {
		fmt.Println("error parsing currentBeaconStateJSON")
	}

	//this is not the beacon state of the slot containing the old withdrawal we want to proof but rather
	// its the state that was merklized to create a historical summary containing the slot that has that withdrawal, ie, 7421952 mod 8192 = 0
	oldBeaconStateJSON, err := parseJSONFileCapella("data/goerli_slot_6397952.json")
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
	beaconBlockHeaderToVerify, err := blockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("error", err)
	}

	// fmt.Println("THESE SHOULD BE", hex.EncodeToString(beaconBlockHeaderToVerify[:]))
	// fmt.Println("THE SAME", hex.EncodeToString(beaconBlockHeaderToVerify[:]))
	// fmt.Println("THESE SHOULD BE", hex.EncodeToString(oldBeaconStateTopLevelRoots.BlockRootsRoot[:]))
	// fmt.Println("THE SAME", hex.EncodeToString(currentBeaconState.HistoricalSummaries[146].BlockSummaryRoot[:]))

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

	historicalBlockHeaderIndex := beacon.HistoricalSummaryListIndex<<((beacon.HistoricalSummaryListMerkleSubtreeNumLayers+1)+1+(beacon.BlockRootsMerkleSubtreeNumLayers)) |
		historicalSummaryIndex<<(1+beacon.BlockRootsMerkleSubtreeNumLayers) |
		beacon.BlockSummaryRootIndex<<(beacon.BlockRootsMerkleSubtreeNumLayers) | beaconBlockHeaderToVerifyIndex

	flag := common.ValidateProof(currentBeaconStateRoot, historicalSummaryBlockHeaderProof, beaconBlockHeaderToVerify, historicalBlockHeaderIndex)
	if flag != true {
		fmt.Println("error 2")
	}

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
	flag := common.ValidateProof(*beaconStateTopLevelRoots.ValidatorsRoot, validatorProof, leaf, uint64(validatorIndex))
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
	flag := common.ValidateProof(beaconBlockHeaderRoot, blockHeaderSlotProof, slotHashRoot, beacon.SlotIndex)
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
	flag := common.ValidateProof(beaconBlockHeaderRoot, blockHeaderBlockBodyProof, blockBodyHashRoot, beacon.BeaconBlockBodyRootIndex)
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
	flag := common.ValidateProof(blockHeaderBodyRoot, beaconBlockBodyProof, executionPayloadHashRoot, beacon.ExecutionPayloadIndex)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestStateRootAgainstLatestBlockHeaderProof(t *testing.T) {

	// this is the state where the latest block header from the oracle was taken.  This is the next slot after
	// the state we want to prove things about (remember latestBlockHeader.state_root = previous slot's state root)
	// oracleStateJSON, err := parseJSONFile("data/historical_summary_proof/goerli_slot_6399999.json")
	// var oracleState deneb.BeaconState
	// ParseCapellaBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	var blockHeader phase0.BeaconBlockHeader
	blockHeader, err := ExtractBlockHeader("data/deneb_goerli_block_header_7413760.json")
	if err != nil {
		fmt.Println("error with block header", err)
	}

	//the state from the prev slot which contains shit we wanna prove about
	stateToProveJSON, err := parseJSONFile("data/deneb_goerli_slot_7413760.json")
	if err != nil {
		fmt.Println("error with parsing JSON state file", err)
	}

	var stateToProve deneb.BeaconState
	ParseDenebBeaconStateFromJSON(*stateToProveJSON, &stateToProve)

	proof, err := beacon.ProveStateRootAgainstBlockHeader(&blockHeader)
	if err != nil {
		fmt.Println("Error in generating proof", err)
	}
	root, err := blockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}
	leaf, err := stateToProve.HashTreeRoot()
	if err != nil {
		fmt.Println("this error", err)
	}

	flag := common.ValidateProof(root, proof, leaf, 3)
	if flag != true {
		fmt.Println("this error")
	}
	assert.True(t, flag, "Proof %v failed")
}

func TestGetExecutionPayloadProof(t *testing.T) {

	// get the proof for execution payload in the block body

	exectionPayloadProof, _, _ := beacon.ProveExecutionPayloadAgainstBlockHeaderDeneb(&blockHeader, block.Body)

	// get the hash root of the actual execution payload
	var executionPayloadHashRoot, _ = block.Body.ExecutionPayload.HashTreeRoot()

	// get the body root in the beacon block header -  will be used as the Merkle root
	root, _ := blockHeader.HashTreeRoot()

	index := beacon.BeaconBlockBodyRootIndex<<(beacon.BlockBodyMerkleSubtreeNumLayers) | beacon.ExecutionPayloadIndex

	// calling the proof verification function
	flag := common.ValidateProof(root, exectionPayloadProof, executionPayloadHashRoot, index)
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
	flag := common.ValidateProof(executionPayloadHashRoot, withdrawalsListProof, withdrawalsHashRoot, beacon.WithdrawalsIndex)
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
	flag := common.ValidateProof(withdrawalsHashRoot, withdrawalProof, leaf, uint64(withdrawalIndex))
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetWithdrawalProof(t *testing.T) {

	// picking up one random validator index
	withdrawalIndex := uint8(0)

	withdrawalProof, _ := beacon.ProveWithdrawalAgainstExecutionPayload(executionPayloadFieldRoots, block.Body.ExecutionPayload.Withdrawals, withdrawalIndex)

	executionPayloadRoot, _ := block.Body.ExecutionPayload.HashTreeRoot()

	leaf, err := block.Body.ExecutionPayload.Withdrawals[withdrawalIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}
	// withdrawalIndex = beacon.BeaconBlockBodyRootIndex<<( beacon.BlockBodyMerkleSubtreeNumLayers+ executionPayloadMerkleSubtreeNumLayers+( beacon.WithdrawalListMerkleSubtreeNumLayers+1)) | beacon.ExecutionPayloadIndex<<( executionPayloadMerkleSubtreeNumLayers+( beacon.WithdrawalListMerkleSubtreeNumLayers+1)) | beacon.WithdrawalsIndex<<( beacon.WithdrawalListMerkleSubtreeNumLayers+1) | withdrawalIndex

	withdrawalRelativeToELPayloadIndex := beacon.WithdrawalsIndex<<(beacon.WithdrawalListMerkleSubtreeNumLayers+1) | uint64(withdrawalIndex)

	// calling the proof verification func
	flag := common.ValidateProof(executionPayloadRoot, withdrawalProof, leaf, withdrawalRelativeToELPayloadIndex)
	if flag != true {
		fmt.Println("error")
	}

	assert.True(t, flag, "Proof %v failed\n")
}

func TestGetTimestampProof(t *testing.T) {

	// get the block number
	executionPayloadFields := block.Body.ExecutionPayload

	// get the Merkle proof for inclusion
	timestampProof, _ := beacon.ProveTimestampAgainstExecutionPayload(executionPayloadFieldRoots)

	hh := ssz.NewHasher()
	hh.PutUint64(uint64(executionPayloadFields.Timestamp))

	leaf := ConvertTo32ByteArray(hh.Hash())

	root, err := block.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	// calling the proof verification func
	flag := common.ValidateProof(root, timestampProof, leaf, beacon.TimestampIndex)
	if flag != true {
		fmt.Println("proof failed")
	}

	assert.True(t, flag, "Proof %v failed")
}

func TestGetValidatorProof(t *testing.T) {
	// picking up one random validator index
	validatorIndex := uint64(VALIDATOR_INDEX)

	// get the validators field
	validators := oracleState.Validators

	beaconStateTopLevelRoots, err := beacon.ComputeBeaconStateTopLevelRootsDeneb(&oracleState)

	validatorProof, _ := epp.ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleState.Slot, oracleState.Validators, uint64(validatorIndex))

	// verify the proof
	// get the leaf corresponding to validatorIndex
	leaf, err := validators[validatorIndex].HashTreeRoot()
	if err != nil {
		fmt.Println("error")
	}

	// calling the proof verification func
	beaconRoot, _ := oracleState.HashTreeRoot()

	validatorIndex = beacon.ValidatorListIndex<<(beacon.ValidatorListMerkleSubtreeNumLayers+1) | uint64(validatorIndex)

	flag := common.ValidateProof(beaconRoot, validatorProof, leaf, validatorIndex)
	if flag != true {
		fmt.Println("error")
	}

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

	flag := common.ValidateProof(root, proof, leaf, 0)
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

func parseJSONFile(filePath string) (*beaconStateJSONDeneb, error) {
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

func parseJSONFileCapella(filePath string) (*beaconStateJSONCapella, error) {
	data, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("error with reading file")
		return nil, err
	}

	var beaconState beaconStateVersionCapella
	err = json.Unmarshal(data, &beaconState)
	if err != nil {
		fmt.Println("error with beaconState JSON unmarshalling")
		return nil, err
	}

	actualData := beaconState.Data
	return &actualData, nil
}
