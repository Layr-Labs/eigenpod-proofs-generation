package main

import (
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/rs/zerolog/log"
)

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
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing state file", err)
	}
	ParseCapellaBeaconStateFromJSON(*stateJSON, &state)

	historicalSummaryJSON, err := parseStateJSONFile(historicalSummaryStateFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing historical summary state file", err)
	}
	ParseCapellaBeaconStateFromJSON(*historicalSummaryJSON, &historicalSummaryState)

	withdrawalBlockHeader, err = ExtractBlockHeader(blockHeaderFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with parsing header file", err)
	}

	withdrawalBlock, err = ExtractBlock(blockBodyFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with parsing body file", err)
	}

	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	// validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, err := state.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of state", err)
	}

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := ConvertTo32ByteArray(hh.Hash())

	timestamp := withdrawalBlock.Body.ExecutionPayload.Timestamp
	hh.PutUint64(uint64(timestamp))
	timestampRoot := ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, err := withdrawalBlockHeader.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of blockHeader", err)
	}
	executionPayloadRoot, err := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of executionPayload", err)
	}

	epp, err := NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error creating EPP object", err)
	}
	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&state)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ComputeBeaconStateTopLevelRoots", err)
	}
	withdrawalProof, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &state, oracleBeaconStateTopLevelRoots, historicalSummaryState.BlockRoots, &withdrawalBlock, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveWithdrawal", err)
	}
	stateRootProof, err := ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveStateRootAgainstBlockHeader", err)
	}
	validatorProof, err := epp.ProveValidatorAgainstBeaconState(&state, oracleBeaconStateTopLevelRoots, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveValidatorAgainstBeaconState", err)
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
