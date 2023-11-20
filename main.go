package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// this needs to be hand crafted. If you want the root of the header at the slot x,
// then look for entry in (x)%SLOTS_PER_HISTORICAL_ROOT in the block_roots.

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Defining flags for all the parameters
	command := flag.String("command", "", "The command to execute")

	oracleBlockHeaderFile := flag.String("oracleBlockHeaderFile", "", "Oracle block header file")
	stateFile := flag.String("stateFile", "", "State file")
	validatorIndex := flag.Uint64("validatorIndex", 0, "validatorIndex")
	outputFile := flag.String("outputFile", "", "Output file")
	chainID := flag.Uint64("chainID", 0, "Chain ID")

	//WithdrawaProof specific flags
	historicalSummariesIndex := flag.Uint64("historicalSummariesIndex", 0, "Historical summaries index")
	blockHeaderIndex := flag.Uint64("blockHeaderIndex", 0, "Block header index")
	historicalSummaryStateFile := flag.String("historicalSummaryStateFile", "", "Historical summary state file")
	blockHeaderFile := flag.String("blockHeaderFile", "", "Block Header file")
	blockBodyFile := flag.String("blockBodyFile", "", "Block Body file")
	withdrawalIndex := flag.Uint64("withdrawalIndex", 0, "Withdrawal index")

	// Parse the flags
	flag.Parse()

	// Check if the required 'command' flag is provided
	if *command == "" {
		log.Debug().Msg("Error: command flag is required")
		return
	}

	// Handling commands based on the 'command' flag
	switch *command {
	case "ValidatorFieldsProof":
		GenerateValidatorFieldsProof(*oracleBlockHeaderFile, *stateFile, *validatorIndex, *chainID, *outputFile)

	case "WithdrawalFieldsProof":
		GenerateWithdrawalFieldsProof(*oracleBlockHeaderFile, *stateFile, *historicalSummaryStateFile, *blockHeaderFile, *blockBodyFile, *validatorIndex, *withdrawalIndex, *historicalSummariesIndex, *blockHeaderIndex, *chainID, *outputFile)

	case "BalanceUpdateProof":
		GenerateBalanceUpdateProof(*oracleBlockHeaderFile, *stateFile, *validatorIndex, *chainID, *outputFile)

	default:
		log.Debug().Str("Unknown command:", *command)
	}
}

func GenerateValidatorFieldsProof(oracleBlockHeaderFile string, stateFile string, validatorIndex uint64, chainID uint64, output string) {

	var state capella.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := parseStateJSONFile(stateFile)
	if err != nil {
		log.Debug().Msg("error with JSON parsing")
	}
	ParseCapellaBeaconStateFromJSON(*stateJSON, &state)

	oracleBeaconBlockHeader, err = ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)

	}

	beaconStateRoot, err := state.HashTreeRoot()

	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
	}

	epp, err := NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)

	}

	stateRootProof, validatorFieldsProof, err := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
	}

	proofs := WithdrawalCredentialProofs{
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		ValidatorIndex:                         uint64(validatorIndex),
		WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
	}

	_ = os.WriteFile(output, proofData, 0644)

}

func GenerateBalanceUpdateProof(oracleBlockHeaderFile string, stateFile string, validatorIndex uint64, chainID uint64, output string) {

	var state capella.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := parseStateJSONFile(stateFile)
	if err != nil {
		log.Debug().AnErr("error with JSON parsing", err)
	}
	ParseCapellaBeaconStateFromJSON(*stateJSON, &state)

	oracleBeaconBlockHeader, err = ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
	}

	beaconStateRoot, err := state.HashTreeRoot()

	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
	}

	epp, err := NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
	}

	balanceRootList, err := GetBalanceRoots(state.Balances)
	if err != nil {
		log.Debug().AnErr("Error with GetBalanceRoots", err)
	}
	balanceRoot := balanceRootList[validatorIndex/4]
	balanceProof, err := epp.ProveValidatorBalance(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorBalance", err)
	}

	stateRootProof, validatorFieldsProof, err := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
	}
	proofs := BalanceUpdateProofs{
		ValidatorIndex:                         uint64(validatorIndex),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		BalanceRoot:                            "0x" + hex.EncodeToString(balanceRoot[:]),
		ValidatorBalanceProof:                  ConvertBytesToStrings(balanceProof.BalanceUpdateProof.ValidatorBalanceProof),
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
		WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
	}

	_ = os.WriteFile(output, proofData, 0644)

}

func GenerateWithdrawalFieldsProof(
	oracleBlockHeaderFile,
	stateFile,
	historicalSummaryStateFile,
	blockHeaderFile,
	blockBodyFile string,
	validatorIndex,
	withdrawalIndex,
	historicalSummariesIndex,
	blockHeaderIndex,
	chainID uint64,
	outputFile string,
) {

	//this is the oracle provided state
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	//this is the state with the withdrawal in it
	var state capella.BeaconState
	var historicalSummaryState capella.BeaconState
	var withdrawalBlockHeader phase0.BeaconBlockHeader
	var withdrawalBlock capella.BeaconBlock

	oracleBeaconBlockHeader, err := ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
	}

	stateJSON, err := parseStateJSONFile(stateFile)
	if err != nil {
		log.Debug().AnErr("error with JSON parsing state file", err)
	}
	ParseCapellaBeaconStateFromJSON(*stateJSON, &state)

	historicalSummaryJSON, err := parseStateJSONFile(historicalSummaryStateFile)
	if err != nil {
		log.Debug().AnErr("error with JSON parsing historical summary state file", err)
	}
	ParseCapellaBeaconStateFromJSON(*historicalSummaryJSON, &historicalSummaryState)

	withdrawalBlockHeader, err = ExtractBlockHeader(blockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
	}

	withdrawalBlock, err = ExtractBlock(blockBodyFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing body file", err)
	}

	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	// validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, err := state.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
	}

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := ConvertTo32ByteArray(hh.Hash())

	timestamp := withdrawalBlock.Body.ExecutionPayload.Timestamp
	hh.PutUint64(uint64(timestamp))
	timestampRoot := ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, err := withdrawalBlockHeader.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of blockHeader", err)
	}
	executionPayloadRoot, err := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of executionPayload", err)
	}

	epp, err := NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
	}
	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&state)
	if err != nil {
		log.Debug().AnErr("Error with ComputeBeaconStateTopLevelRoots", err)
	}
	withdrawalProof, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &state, oracleBeaconStateTopLevelRoots, historicalSummaryState.BlockRoots, &withdrawalBlock, uint64(validatorIndex), FIRST_CAPELLA_SLOT_GOERLI)
	if err != nil {
		log.Debug().AnErr("Error with ProveWithdrawal", err)
	}
	stateRootProof, err := ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().AnErr("Error with ProveStateRootAgainstBlockHeader", err)
	}
	validatorProof, err := epp.ProveValidatorAgainstBeaconState(&state, oracleBeaconStateTopLevelRoots, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorAgainstBeaconState", err)
	}
	proofs := WithdrawalProofs{
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		WithdrawalProof:                        ConvertBytesToStrings(withdrawalProof.WithdrawalProof),
		SlotProof:                              ConvertBytesToStrings(withdrawalProof.SlotProof),
		ExecutionPayloadProof:                  ConvertBytesToStrings(withdrawalProof.ExecutionPayloadProof),
		TimestampProof:                         ConvertBytesToStrings(withdrawalProof.TimestampProof),
		HistoricalSummaryProof:                 ConvertBytesToStrings(withdrawalProof.HistoricalSummaryBlockRootProof),
		BlockHeaderRootIndex:                   beaconBlockHeaderToVerifyIndex,
		HistoricalSummaryIndex:                 uint64(historicalSummariesIndex),
		WithdrawalIndex:                        withdrawalIndex,
		BlockHeaderRoot:                        "0x" + hex.EncodeToString(blockHeaderRoot[:]),
		SlotRoot:                               "0x" + hex.EncodeToString(slotRoot[:]),
		TimestampRoot:                          "0x" + hex.EncodeToString(timestampRoot[:]),
		ExecutionPayloadRoot:                   "0x" + hex.EncodeToString(executionPayloadRoot[:]),
		ValidatorProof:                         ConvertBytesToStrings(validatorProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
		WithdrawalFields:                       GetWithdrawalFields(withdrawalBlock.Body.ExecutionPayload.Withdrawals[withdrawalIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
	}

	_ = os.WriteFile(outputFile, proofData, 0644)

}
