package eigenpodproofs

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	contractBeaconChainProofs "github.com/Layr-Labs/eigenpod-proofs-generation/bindings"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestValidatorContainersProofOnChain(t *testing.T) {
	var oracleState deneb.BeaconState
	stateFile := "data/deneb_goerli_slot_7413760.json"
	stateJSON, err := ParseJSONFileDeneb(stateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	ParseDenebBeaconStateFromJSON(*stateJSON, &oracleState)

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
	}

	oracleBlockHeader, err = ExtractBlockHeader("data/deneb_goerli_block_header_7413760.json")
	if err != nil {
		fmt.Println("error", err)
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(&oracleBlockHeader, &versionedOracleState, []uint64{VALIDATOR_INDEX})
	if err != nil {
		fmt.Println("error", err)
	}

	validatorFieldsProof := verifyValidatorFieldsCallParams.ValidatorFieldsProofs[0].ToByteSlice()
	validatorIndex := new(big.Int).SetUint64(verifyValidatorFieldsCallParams.ValidatorIndices[0])
	oracleBlockHeaderRoot, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("error", err)
	}

	err = beaconChainProofs.VerifyStateRootAgainstLatestBlockRoot(
		&bind.CallOpts{},
		oracleBlockHeaderRoot,
		verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
		verifyValidatorFieldsCallParams.StateRootProof.StateRootProof.ToByteSlice(),
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)

	var validatorFields [][32]byte
	for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[0] {
		validatorFields = append(validatorFields, field)
	}

	err = beaconChainProofs.VerifyValidatorFields(
		&bind.CallOpts{},
		verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
		validatorFields,
		validatorFieldsProof,
		validatorIndex,
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)
}

// TODO: get these tests working
func TestProvingDenebWithdrawalAgainstDenebStateOnChain(t *testing.T) {

	oracleStateFile := "data/deneb_goerli_slot_7431952.json"
	oracleStateJSON, err := ParseJSONFileDeneb(oracleStateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	oracleState := deneb.BeaconState{}
	ParseDenebBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	oracleHeaderFile := "data/deneb_goerli_block_header_7431952.json"
	oracleBlockHeader, err = ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error creating versioned state", err)
	}

	historicalSummaryStateJSON, err := ParseJSONFileDeneb("data/deneb_goerli_slot_7421952.json")
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

	versionedWithdrawalBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		fmt.Println("error", err)
	}

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

	var withdrawalFields [][32]byte
	for _, field := range verifyAndProcessWithdrawalCallParams.WithdrawalFields[0] {
		withdrawalFields = append(withdrawalFields, field)
	}

	withdrawalProof := contractBeaconChainProofs.BeaconChainProofsContractWithdrawalProof{
		WithdrawalProof:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalProof.ToByteSlice(),
		SlotProof:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotProof.ToByteSlice(),
		ExecutionPayloadProof:           verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadProof.ToByteSlice(),
		TimestampProof:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampProof.ToByteSlice(),
		HistoricalSummaryBlockRootProof: verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryBlockRootProof.ToByteSlice(),
		BlockRootIndex:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex,
		HistoricalSummaryIndex:          verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex,
		WithdrawalIndex:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalIndex,
		BlockRoot:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRoot,
		SlotRoot:                        verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotRoot,
		TimestampRoot:                   verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampRoot,
		ExecutionPayloadRoot:            verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadRoot,
	}

	fmt.Println("historicalSummaryndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex)
	fmt.Println("blockRootIndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex)
	fmt.Println("beacon state root: ", hex.EncodeToString(verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot[:]))

	err = beaconChainProofs.VerifyWithdrawal(
		&bind.CallOpts{},
		verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot,
		withdrawalFields,
		withdrawalProof,
		DENEB_FORK_TIMESTAMP_GOERLI,
	)

	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)
}

func TestProvingCapellaWithdrawalAgainstDenebStateOnChain(t *testing.T) {

	oracleStateFile := "data/deneb_goerli_slot_7431952.json"
	oracleStateJSON, err := ParseJSONFileDeneb(oracleStateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	oracleState := deneb.BeaconState{}
	ParseDenebBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	oracleHeaderFile := "data/deneb_goerli_block_header_7431952.json"
	oracleBlockHeader, err = ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error creating versioned state", err)
	}

	historicalSummaryStateJSON, err := ParseJSONFileCapella("data/goerli_slot_6397952.json")
	if err != nil {
		fmt.Println("error parsing historicalSummaryState JSON")
	}
	var historicalSummaryState capella.BeaconState
	ParseCapellaBeaconStateFromJSON(*historicalSummaryStateJSON, &historicalSummaryState)
	historicalSummaryStateBlockRoots := historicalSummaryState.BlockRoots

	withdrawalBlock, err := ExtractBlockCapella("data/goerli_block_6397852.json")
	if err != nil {
		fmt.Println("block.UnmarshalJSON error", err)
	}

	versionedWithdrawalBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		fmt.Println("error", err)
	}

	withdrawalValidatorIndex := uint64(200240) //this is the index of the validator with the first withdrawal in the withdrawalBlock 7421951

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

	var withdrawalFields [][32]byte
	for _, field := range verifyAndProcessWithdrawalCallParams.WithdrawalFields[0] {
		withdrawalFields = append(withdrawalFields, field)
	}

	withdrawalProof := contractBeaconChainProofs.BeaconChainProofsContractWithdrawalProof{
		WithdrawalProof:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalProof.ToByteSlice(),
		SlotProof:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotProof.ToByteSlice(),
		ExecutionPayloadProof:           verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadProof.ToByteSlice(),
		TimestampProof:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampProof.ToByteSlice(),
		HistoricalSummaryBlockRootProof: verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryBlockRootProof.ToByteSlice(),
		BlockRootIndex:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex,
		HistoricalSummaryIndex:          verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex,
		WithdrawalIndex:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalIndex,
		BlockRoot:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRoot,
		SlotRoot:                        verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotRoot,
		TimestampRoot:                   verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampRoot,
		ExecutionPayloadRoot:            verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadRoot,
	}

	fmt.Println("historicalSummaryndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex)
	fmt.Println("blockRootIndex ", verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex)

	err = beaconChainProofs.VerifyWithdrawal(
		&bind.CallOpts{},
		verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot,
		withdrawalFields,
		withdrawalProof,
		DENEB_FORK_TIMESTAMP_GOERLI,
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)
}

func TestProvingCapellaWithdrawalAgainstCapellaStateOnChain(t *testing.T) {
	oracleStateFile := "data/goerli_slot_6409723.json"
	oracleStateJSON, err := ParseJSONFileCapella(oracleStateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	oracleState := capella.BeaconState{}
	ParseCapellaBeaconStateFromJSON(*oracleStateJSON, &oracleState)

	oracleHeaderFile := "data/goerli_block_header_6409723.json"
	oracleBlockHeader, err = ExtractBlockHeader(oracleHeaderFile)
	if err != nil {
		fmt.Println("error with block header", err)
	}

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error creating versioned state", err)
	}

	historicalSummaryStateJSON, err := ParseJSONFileCapella("data/goerli_slot_6397952.json")
	if err != nil {
		fmt.Println("error parsing historicalSummaryState JSON")
	}
	var historicalSummaryState capella.BeaconState
	ParseCapellaBeaconStateFromJSON(*historicalSummaryStateJSON, &historicalSummaryState)
	historicalSummaryStateBlockRoots := historicalSummaryState.BlockRoots

	withdrawalBlock, err := ExtractBlockCapella("data/goerli_block_6397852.json")
	if err != nil {
		fmt.Println("block.UnmarshalJSON error", err)
	}

	versionedWithdrawalBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		fmt.Println("error", err)
	}

	withdrawalValidatorIndex := uint64(200240) //this is the index of the validator with the first withdrawal in the withdrawalBlock 7421951

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

	var withdrawalFields [][32]byte
	for _, field := range verifyAndProcessWithdrawalCallParams.WithdrawalFields[0] {
		withdrawalFields = append(withdrawalFields, field)
	}

	withdrawalProof := contractBeaconChainProofs.BeaconChainProofsContractWithdrawalProof{
		WithdrawalProof:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalProof.ToByteSlice(),
		SlotProof:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotProof.ToByteSlice(),
		ExecutionPayloadProof:           verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadProof.ToByteSlice(),
		TimestampProof:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampProof.ToByteSlice(),
		HistoricalSummaryBlockRootProof: verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryBlockRootProof.ToByteSlice(),
		BlockRootIndex:                  verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRootIndex,
		HistoricalSummaryIndex:          verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].HistoricalSummaryIndex,
		WithdrawalIndex:                 verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].WithdrawalIndex,
		BlockRoot:                       verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].BlockRoot,
		SlotRoot:                        verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].SlotRoot,
		TimestampRoot:                   verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].TimestampRoot,
		ExecutionPayloadRoot:            verifyAndProcessWithdrawalCallParams.WithdrawalProofs[0].ExecutionPayloadRoot,
	}

	err = beaconChainProofs.VerifyWithdrawal(
		&bind.CallOpts{},
		verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot,
		withdrawalFields,
		withdrawalProof,
		DENEB_FORK_TIMESTAMP_GOERLI,
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)

}
