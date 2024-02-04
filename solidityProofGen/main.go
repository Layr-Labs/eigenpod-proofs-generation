package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	ssz "github.com/ferranbt/fastssz"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
)

var (
	b                       deneb.BeaconState
	blockHeader             phase0.BeaconBlockHeader
	blockBody               deneb.BeaconBlockBody
	signedBlock             deneb.SignedBeaconBlock
	executionPayload        deneb.ExecutionPayload
	executionPayloadCapella capella.ExecutionPayload
)

var VALIDATOR_INDEX uint64 = 61068 //this is the index of a validator that has a full withdrawal

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%SLOTS_PER_HISTORICAL_ROOT in the block_roots.

var BEACON_BLOCK_HEADER_TO_VERIFY_INDEX uint64 = 2262

const FIRST_CAPELLA_SLOT_GOERLI = uint64(5193728)
const FIRST_CAPELLA_SLOT_MAINNET = uint64(6209536)

var GOERLI_CHAIN_ID uint64 = 5

func main() {
	// Example optional flag
	var newBalance int
	flag.IntVar(&newBalance, "newBalance", -1, "new balance")

	// Parse the flags
	flag.Parse()

	// Using flag.Args() to get positional arguments
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("No command provided")
		return
	}

	switch args[0] {
	case "WithdrawalFieldsProof":
		// if len(args) < 14 {
		// 	fmt.Println("Incorrect number of arguments for WithdrawalFieldsProof")
		// 	return
		// }

		validatorIndex, err := strconv.ParseUint(args[1], 10, 64)
		historicalSummariesIndex, err := strconv.ParseUint(args[2], 10, 64)
		blockHeaderIndex, err := strconv.ParseUint(args[3], 10, 64)
		modifyStateToIncludeFullWithdrawal, err := strconv.ParseBool(args[4])
		partialWithdrawalProof, err := strconv.ParseBool(args[5])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		oracleBlockHeaderFile := args[6]
		stateFile := args[7]
		historicalSummaryStateFile := args[8]
		headerFile := args[9]
		bodyFile := args[10]
		outputFile := args[11]
		advanceSlotOfWithdrawal, err := strconv.ParseBool(args[12])
		isCapella, err := strconv.ParseBool(args[13])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if isCapella {
			fmt.Print("CAPELLA")
			GenerateWithdrawalFieldsProofCapella(validatorIndex, historicalSummariesIndex, blockHeaderIndex, oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, outputFile, modifyStateToIncludeFullWithdrawal, partialWithdrawalProof, advanceSlotOfWithdrawal)
		} else {
			fmt.Print("DENEB")

			GenerateWithdrawalFieldsProof(validatorIndex, historicalSummariesIndex, blockHeaderIndex, oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, outputFile, modifyStateToIncludeFullWithdrawal, partialWithdrawalProof, advanceSlotOfWithdrawal)
		}

	case "ValidatorFieldsProof":

		index, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		changeBalance, err := strconv.ParseBool(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		oracleStateFile := args[3]
		stateFile := args[4]
		outputFile := args[5]

		GenerateValidatorFieldsProof(oracleStateFile, stateFile, index, changeBalance, uint64(newBalance), outputFile)

	case "BalanceUpdateProof":

		index, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		changeBalance, err := strconv.ParseBool(args[2])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		incrementSlot, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		oracleStateFile := args[4]
		stateFile := args[5]
		outputFile := args[6]

		GenerateBalanceUpdateProof(oracleStateFile, stateFile, index, changeBalance, uint64(newBalance), incrementSlot, outputFile)

	default:
		fmt.Println("Unknown command:", args[0])
	}
}

// TODO: Get this working
func GenerateValidatorFieldsProof(oracleBlockHeaderFile string, stateFile string, index uint64, changeBalance bool, newBalance uint64, output string) {
	// var state deneb.BeaconState
	// var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	// SetupValidatorProof(oracleBlockHeaderFile, stateFile, index, changeBalance, newBalance, 0, &state, &oracleBeaconBlockHeader)

	// validatorIndex := phase0.ValidatorIndex(index)

	// beaconStateRoot, _ := state.HashTreeRoot()

	// balanceRootList, _ := eigenpodproofs.GetBalanceRoots(state.Balances)
	// balanceRoot := balanceRootList[validatorIndex/4]

	// latestBlockHeaderRoot, err := oracleBeaconBlockHeader.HashTreeRoot()
	// if err != nil {
	// 	fmt.Println("Error with HashTreeRoot of latestBlockHeader", err)
	// }

	// epp, err := eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	// if err != nil {
	// 	fmt.Println("Error creating EPP object", err)

	// }
	// balanceProof, _ := epp.ProveValidatorBalance(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))

	// stateRootProof, validatorFieldsProof, _ := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))

	// proofs := WithdrawalCredentialProofs{
	// 	ValidatorIndex:                         uint64(validatorIndex),
	// 	BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
	// 	BalanceRoot:                            "0x" + hex.EncodeToString(balanceRoot[:]),
	// 	LatestBlockHeaderRoot:                  "0x" + hex.EncodeToString(latestBlockHeaderRoot[:]),
	// 	ValidatorBalanceProof:                  ConvertBytesToStrings(balanceProof.BalanceUpdateProof.ValidatorBalanceProof),
	// 	WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
	// 	ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	// 	StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
	// }

	// proofData, err := json.Marshal(proofs)
	// if err != nil {
	// 	fmt.Println("error")
	// }

	// _ = ioutil.WriteFile(output, proofData, 0644)

}

func GenerateWithdrawalFieldsProof(index, historicalSummariesIndex, blockHeaderIndex uint64, oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, outputFile string, modifyStateToIncludeFullWithdrawal bool, partialWithdrawalProof bool, advanceSlotOfWithdrawal bool) {

	//this is the oracle provided state
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	//this is the state with the withdrawal in it
	var oracleState deneb.BeaconState
	var historicalSummaryState deneb.BeaconState
	var withdrawalBlockHeader phase0.BeaconBlockHeader
	var withdrawalBlock deneb.BeaconBlock

	withdrawalToModifyIndex := uint64(0)
	SetUpWithdrawalsProof(oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, &oracleBeaconBlockHeader, &oracleState, &historicalSummaryState, &withdrawalBlockHeader, &withdrawalBlock, modifyStateToIncludeFullWithdrawal, partialWithdrawalProof, index, historicalSummariesIndex, withdrawalToModifyIndex, advanceSlotOfWithdrawal)
	root, _ := withdrawalBlock.Body.HashTreeRoot()
	fmt.Println("blockBody.hashtreeroot()", hex.EncodeToString(root[:]))
	fmt.Println("blockheader.bodyroot)", hex.EncodeToString(withdrawalBlockHeader.BodyRoot[:]))
	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, _ := oracleState.HashTreeRoot()

	fmt.Println("beaconStateRoot", hex.EncodeToString(beaconStateRoot[:]))

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := eigenpodproofs.ConvertTo32ByteArray(hh.Hash())

	latestBlockHeaderRoot, err := oracleBeaconBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("Error with HashTreeRoot of latestBlockHeader", err)
	}

	timestamp := withdrawalBlock.Body.ExecutionPayload.Timestamp
	hh.PutUint64(uint64(timestamp))
	timestampRoot := eigenpodproofs.ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, _ := withdrawalBlockHeader.HashTreeRoot()
	blockBodyRoot, _ := withdrawalBlock.Body.HashTreeRoot()
	executionPayloadRoot, _ := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()

	epp, err := eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	if err != nil {
		fmt.Println("Error creating EPP object", err)
	}
	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&oracleState)
	if err != nil {
		fmt.Println("Error computing beacon state top level roots", err)
	}
	//blockHeaderProof, slotProof, withdrawalProof, validatorProof, timestampProof, executionPayloadProof, stateRootAgainstLatestBlockHeaderProof, historicalSummaryProof, err :=
	// withdrawalProof, stateRootProof, validatorProof, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &oracleState, historicalSummaryState.BlockRoots, &withdrawalBlock, validatorIndex)
	withdrawalProof, err := epp.ProveDenebWithdrawal(&oracleBeaconBlockHeader, &oracleState, oracleBeaconStateTopLevelRoots, historicalSummaryState.BlockRoots, &withdrawalBlock, uint64(validatorIndex))
	if err != nil {
		fmt.Println("ProveWithdrawal error", err)
	}
	stateRootProof, err := beacon.ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		fmt.Println("ProveStateRootAgainstBlockHeader error", err)
	}
	validatorProof, err := epp.ProveValidatorAgainstBeaconState(&oracleState, oracleBeaconStateTopLevelRoots, uint64(validatorIndex))
	if err != nil {
		fmt.Println("ProveValidatorAgainstBeaconState error", err)
	}

	proofs := WithdrawalProofs{
		Slot:                                   uint64(slot),
		HistoricalSummaryIndex:                 uint64(historicalSummariesIndex),
		WithdrawalIndex:                        withdrawalToModifyIndex,
		BlockHeaderRootIndex:                   beaconBlockHeaderToVerifyIndex,
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		SlotRoot:                               "0x" + hex.EncodeToString(slotRoot[:]),
		TimestampRoot:                          "0x" + hex.EncodeToString(timestampRoot[:]),
		BlockHeaderRoot:                        "0x" + hex.EncodeToString(blockHeaderRoot[:]),
		BlockBodyRoot:                          "0x" + hex.EncodeToString(blockBodyRoot[:]),
		ExecutionPayloadRoot:                   "0x" + hex.EncodeToString(executionPayloadRoot[:]),
		LatestBlockHeaderRoot:                  "0x" + hex.EncodeToString(latestBlockHeaderRoot[:]),
		SlotProof:                              ConvertBytesToStrings(withdrawalProof.SlotProof),
		WithdrawalProof:                        ConvertBytesToStrings(withdrawalProof.WithdrawalProof),
		ValidatorProof:                         ConvertBytesToStrings(validatorProof),
		TimestampProof:                         ConvertBytesToStrings(withdrawalProof.TimestampProof),
		ExecutionPayloadProof:                  ConvertBytesToStrings(withdrawalProof.ExecutionPayloadProof),
		ValidatorFields:                        GetValidatorFields(oracleState.Validators[validatorIndex]),
		WithdrawalFields:                       GetWithdrawalFields(withdrawalBlock.Body.ExecutionPayload.Withdrawals[withdrawalToModifyIndex]),
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof),
		HistoricalSummaryProof:                 ConvertBytesToStrings(withdrawalProof.HistoricalSummaryBlockRootProof),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		fmt.Println("error")
	}

	_ = os.WriteFile(outputFile, proofData, 0644)

}

func GenerateWithdrawalFieldsProofCapella(index, historicalSummariesIndex, blockHeaderIndex uint64, oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, outputFile string, modifyStateToIncludeFullWithdrawal bool, partialWithdrawalProof bool, advanceSlotOfWithdrawal bool) {

	//this is the oracle provided state
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	//this is the state with the withdrawal in it
	var oracleState deneb.BeaconState
	var historicalSummaryState capella.BeaconState
	var withdrawalBlockHeader phase0.BeaconBlockHeader
	var withdrawalBlock capella.BeaconBlock

	withdrawalToModifyIndex := uint64(0)
	fmt.Println("hustoricalsummary state file", historicalSummaryStateFile)
	SetUpWithdrawalsProofCapella(oracleBlockHeaderFile, stateFile, historicalSummaryStateFile, headerFile, bodyFile, &oracleBeaconBlockHeader, &oracleState, &historicalSummaryState, &withdrawalBlockHeader, &withdrawalBlock, modifyStateToIncludeFullWithdrawal, partialWithdrawalProof, index, historicalSummariesIndex, withdrawalToModifyIndex, advanceSlotOfWithdrawal)
	root, _ := withdrawalBlock.Body.HashTreeRoot()
	fmt.Println("blockBody.hashtreeroot()", hex.EncodeToString(root[:]))
	fmt.Println("blockheader.bodyroot)", hex.EncodeToString(withdrawalBlockHeader.BodyRoot[:]))
	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, _ := oracleState.HashTreeRoot()

	fmt.Println("beaconStateRoot", hex.EncodeToString(beaconStateRoot[:]))

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := eigenpodproofs.ConvertTo32ByteArray(hh.Hash())

	latestBlockHeaderRoot, err := oracleBeaconBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("Error with HashTreeRoot of latestBlockHeader", err)
	}

	timestamp := withdrawalBlock.Body.ExecutionPayload.Timestamp
	hh.PutUint64(uint64(timestamp))
	timestampRoot := eigenpodproofs.ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, _ := withdrawalBlockHeader.HashTreeRoot()
	blockBodyRoot, _ := withdrawalBlock.Body.HashTreeRoot()
	executionPayloadRoot, _ := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()

	epp, err := eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	if err != nil {
		fmt.Println("Error creating EPP object", err)
	}
	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&oracleState)
	//blockHeaderProof, slotProof, withdrawalProof, validatorProof, timestampProof, executionPayloadProof, stateRootAgainstLatestBlockHeaderProof, historicalSummaryProof, err :=
	// withdrawalProof, stateRootProof, validatorProof, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &oracleState, historicalSummaryState.BlockRoots, &withdrawalBlock, validatorIndex)
	withdrawalProof, err := epp.ProveCapellaWithdrawal(&oracleBeaconBlockHeader, &oracleState, oracleBeaconStateTopLevelRoots, historicalSummaryState.BlockRoots, &withdrawalBlock, uint64(validatorIndex))
	if err != nil {
		fmt.Println("ProveWithdrawal error", err)
	}
	stateRootProof, err := beacon.ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		fmt.Println("ProveStateRootAgainstBlockHeader error", err)
	}
	validatorProof, err := epp.ProveValidatorAgainstBeaconState(&oracleState, oracleBeaconStateTopLevelRoots, uint64(validatorIndex))
	if err != nil {
		fmt.Println("ProveValidatorAgainstBeaconState error", err)
	}
	proofs := WithdrawalProofs{
		Slot:                                   uint64(slot),
		ValidatorIndex:                         uint64(validatorIndex),
		HistoricalSummaryIndex:                 uint64(historicalSummariesIndex),
		WithdrawalIndex:                        withdrawalToModifyIndex,
		BlockHeaderRootIndex:                   beaconBlockHeaderToVerifyIndex,
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		SlotRoot:                               "0x" + hex.EncodeToString(slotRoot[:]),
		TimestampRoot:                          "0x" + hex.EncodeToString(timestampRoot[:]),
		BlockHeaderRoot:                        "0x" + hex.EncodeToString(blockHeaderRoot[:]),
		BlockBodyRoot:                          "0x" + hex.EncodeToString(blockBodyRoot[:]),
		ExecutionPayloadRoot:                   "0x" + hex.EncodeToString(executionPayloadRoot[:]),
		LatestBlockHeaderRoot:                  "0x" + hex.EncodeToString(latestBlockHeaderRoot[:]),
		SlotProof:                              ConvertBytesToStrings(withdrawalProof.SlotProof),
		WithdrawalProof:                        ConvertBytesToStrings(withdrawalProof.WithdrawalProof),
		ValidatorProof:                         ConvertBytesToStrings(validatorProof),
		TimestampProof:                         ConvertBytesToStrings(withdrawalProof.TimestampProof),
		ExecutionPayloadProof:                  ConvertBytesToStrings(withdrawalProof.ExecutionPayloadProof),
		ValidatorFields:                        GetValidatorFields(oracleState.Validators[validatorIndex]),
		WithdrawalFields:                       GetWithdrawalFields(withdrawalBlock.Body.ExecutionPayload.Withdrawals[withdrawalToModifyIndex]),
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof),
		HistoricalSummaryProof:                 ConvertBytesToStrings(withdrawalProof.HistoricalSummaryBlockRootProof),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		fmt.Println("error")
	}

	_ = os.WriteFile(outputFile, proofData, 0644)

}

// TODO: get this working
func GenerateBalanceUpdateProof(oracleBlockHeaderFile string, stateFile string, index uint64, changeBalance bool, newBalance uint64, incrementSlot uint64, output string) {

	// var state deneb.BeaconState
	// var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	// SetupValidatorProof(oracleBlockHeaderFile, stateFile, index, changeBalance, newBalance, incrementSlot, &state, &oracleBeaconBlockHeader)

	// eigenpodproofs.GenerateValidatorFieldsProof()
	// validatorIndex := phase0.ValidatorIndex(index)

	// beaconStateRoot, _ := state.HashTreeRoot()

	// beaconTopLevelRoots, _ := eigenpodproofs.ComputeBeaconStateTopLevelRoots(&state)
	// slotRoot := beaconTopLevelRoots.SlotRoot

	// latestBlockHeaderRoot, err := oracleBeaconBlockHeader.HashTreeRoot()
	// if err != nil {
	// 	fmt.Println("Error with HashTreeRoot of latestBlockHeader", err)
	// }

	// fmt.Println("slotRoot", slotRoot)

	// epp, err := eigenpodproofs.NewEigenPodProofs(GOERLI_CHAIN_ID, 1000)
	// if err != nil {
	// 	fmt.Println("Error creating EPP object", err)
	// }

	// balanceRootList, _ := eigenpodproofs.GetBalanceRoots(state.Balances)
	// balanceRoot := balanceRootList[validatorIndex/4]
	// balanceProof, _ := epp.ProveValidatorBalance(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))

	// stateRootProof, validatorFieldsProof, _ := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))

	// proofs := BalanceUpdateProofs{
	// 	ValidatorIndex:                         uint64(validatorIndex),
	// 	BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
	// 	BalanceRoot:                            "0x" + hex.EncodeToString(balanceRoot[:]),
	// 	SlotRoot:                               "0x" + hex.EncodeToString(slotRoot[:]),
	// 	LatestBlockHeaderRoot:                  "0x" + hex.EncodeToString(latestBlockHeaderRoot[:]),
	// 	ValidatorBalanceProof:                  ConvertBytesToStrings(balanceProof.BalanceUpdateProof.ValidatorBalanceProof),
	// 	ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	// 	StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
	// 	WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
	// }

	// proofData, err := json.Marshal(proofs)
	// if err != nil {
	// 	fmt.Println("this error")
	// }

	// _ = os.WriteFile(output, proofData, 0644)

}
