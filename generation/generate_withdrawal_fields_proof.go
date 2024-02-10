package generation

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
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
) error {

	//this is the oracle provided state
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	//this is the state with the withdrawal in it
	var state deneb.BeaconState
	var versionedState spec.VersionedBeaconState
	var historicalSummaryState deneb.BeaconState
	var withdrawalBlockHeader phase0.BeaconBlockHeader
	var withdrawalBlock deneb.BeaconBlock

	oracleBeaconBlockHeader, err := ExtractBlockHeader(oracleBlockHeaderFile)

	root, _ := oracleBeaconBlockHeader.HashTreeRoot()
	fmt.Println("oracleBeaconBlockHeader: ", root)

	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
		return err
	}

	stateJSON, err := ParseDenebStateJSONFile(stateFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing state file", err)
		return err
	}
	ParseDenebBeaconStateFromJSON(*stateJSON, &state)

	historicalSummaryJSON, err := ParseDenebStateJSONFile(historicalSummaryStateFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with JSON parsing historical summary state file", err)
		return err
	}
	ParseDenebBeaconStateFromJSON(*historicalSummaryJSON, &historicalSummaryState)

	withdrawalBlockHeader, err = ExtractBlockHeader(blockHeaderFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with parsing header file", err)
		return err
	}

	withdrawalBlock, err = ExtractBlock(blockBodyFile)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with parsing body file", err)
		return err
	}

	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	// validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, err := state.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of state", err)
		return err
	}

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := common.ConvertTo32ByteArray(hh.Hash())

	timestamp := withdrawalBlock.Body.ExecutionPayload.Timestamp
	hh.PutUint64(uint64(timestamp))
	timestampRoot := common.ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, err := withdrawalBlockHeader.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of blockHeader", err)
		return err
	}
	executionPayloadRoot, err := withdrawalBlock.Body.ExecutionPayload.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with HashTreeRoot of executionPayload", err)
		return err
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error creating EPP object", err)
		return err
	}
	versionedState, err = beacon.CreateVersionedState(&state)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with CreateVersionedState", err)
		return err
	}

	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&versionedState)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ComputeBeaconStateTopLevelRoots", err)
		return err
	}

	versionedSignedBlock, err := beacon.CreateVersionedSignedBlock(withdrawalBlock)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with CreateVersionedSignedBlock", err)
		return err
	}

	withdrawalProof, _, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &versionedState, oracleBeaconStateTopLevelRoots, historicalSummaryState.BlockRoots, &versionedSignedBlock, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveWithdrawal", err)
		return err
	}
	stateRootProofAgainstBlockHeader, err := beacon.ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveStateRootAgainstBlockHeader", err)
		return err
	}
	slotProofAgainstBlockHeader, err := beacon.ProveSlotAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveSlotAgainstBlockHeader", err)
		return err
	}

	validatorProof, err := epp.ProveValidatorAgainstBeaconState(oracleBeaconStateTopLevelRoots, state.Slot, state.Validators, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("GenerateWithdrawalFieldsProof: error with ProveValidatorAgainstBeaconState", err)
		return err
	}
	proofs := WithdrawalProofs{
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProofAgainstBlockHeader),
		SlotAgainstLatestBlockHeaderProof:      ConvertBytesToStrings(slotProofAgainstBlockHeader),
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

	return nil
}
