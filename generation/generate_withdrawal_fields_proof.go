package main

import (
	"encoding/hex"
	"encoding/json"
	"os"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	commonutils "github.com/Layr-Labs/eigenpod-proofs-generation/common_utils"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/rs/zerolog/log"
)

const (
	MAINNET_DENEB_FORK_SLOT = 8626176
	MAINNET_CHAIN_ID        = 1
	HOLESKY_DENEB_FORK_SLOT = 950272
	HOLESKY_CHAIN_ID        = 17000
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
	historicalSummaryStateSlot uint64,
	withdrawalSlot uint64,
) error {

	//this is the oracle provided state
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	//this is the state with the withdrawal in it
	var state deneb.BeaconState
	var versionedState spec.VersionedBeaconState
	var historicalSummaryStateBlockRoots []phase0.Root
	var withdrawalBlockHeader phase0.BeaconBlockHeader

	oracleBeaconBlockHeader, err := commonutils.ExtractBlockHeader(oracleBlockHeaderFile)

	if err != nil {
		log.Debug().Msgf("Error with parsing header file: %s", err)
		return err
	}

	stateJSON, err := commonutils.ParseDenebStateJSONFile(stateFile)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with JSON parsing state file: %s", err)
		return err
	}
	commonutils.ParseDenebBeaconStateFromJSON(*stateJSON, &state)

	if chainID == MAINNET_CHAIN_ID && historicalSummaryStateSlot >= MAINNET_DENEB_FORK_SLOT || chainID == HOLESKY_CHAIN_ID && historicalSummaryStateSlot >= HOLESKY_DENEB_FORK_SLOT {
		historicalSummaryJSON, err := commonutils.ParseDenebStateJSONFile(historicalSummaryStateFile)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with JSON parsing historical summary state file: %s", err)
			return err
		}
		var historicalSummaryStateDeneb deneb.BeaconState
		commonutils.ParseDenebBeaconStateFromJSON(*historicalSummaryJSON, &historicalSummaryStateDeneb)
		historicalSummaryStateBlockRoots = historicalSummaryStateDeneb.BlockRoots
	} else {
		historicalSummaryJSON, err := commonutils.ParseCapellaStateJSONFile(historicalSummaryStateFile)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with JSON parsing historical summary state file: %s", err)
			return err
		}
		var historicalSummaryStateCapella capella.BeaconState
		commonutils.ParseCapellaBeaconStateFromJSON(*historicalSummaryJSON, &historicalSummaryStateCapella)
		historicalSummaryStateBlockRoots = historicalSummaryStateCapella.BlockRoots
	}

	withdrawalBlockHeader, err = commonutils.ExtractBlockHeader(blockHeaderFile)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with parsing header file: %s", err)
		return err
	}

	var timestamp uint64
	var executionPayloadRoot phase0.Root
	var withdrawal *capella.Withdrawal
	var versionedSignedBlock spec.VersionedSignedBeaconBlock
	if chainID == MAINNET_CHAIN_ID && historicalSummaryStateSlot >= MAINNET_DENEB_FORK_SLOT || chainID == HOLESKY_CHAIN_ID && historicalSummaryStateSlot >= HOLESKY_DENEB_FORK_SLOT {
		var withdrawalBlockDeneb deneb.BeaconBlock
		withdrawalBlockDeneb, err = commonutils.ExtractBlockDeneb(blockBodyFile)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with parsing body file: %s", err)
			return err
		}
		timestamp = withdrawalBlockDeneb.Body.ExecutionPayload.Timestamp
		executionPayloadRoot, err = withdrawalBlockDeneb.Body.ExecutionPayload.HashTreeRoot()
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with HashTreeRoot of executionPayload: %s", err)
			return err
		}
		withdrawal = withdrawalBlockDeneb.Body.ExecutionPayload.Withdrawals[withdrawalIndex]

		versionedSignedBlock, err = beacon.CreateVersionedSignedBlock(withdrawalBlockDeneb)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with CreateVersionedSignedBlock: %s", err)
			return err
		}
	} else {
		var withdrawalBlockCapella capella.BeaconBlock
		withdrawalBlockCapella, err = commonutils.ExtractBlockCapella(blockBodyFile)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with parsing body file: %s", err)
			return err
		}
		timestamp = withdrawalBlockCapella.Body.ExecutionPayload.Timestamp
		executionPayloadRoot, err = withdrawalBlockCapella.Body.ExecutionPayload.HashTreeRoot()
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with HashTreeRoot of executionPayload: %s", err)
			return err
		}
		withdrawal = withdrawalBlockCapella.Body.ExecutionPayload.Withdrawals[withdrawalIndex]
		versionedSignedBlock, err = beacon.CreateVersionedSignedBlock(withdrawalBlockCapella)
		if err != nil {
			log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with CreateVersionedSignedBlock: %s", err)
			return err
		}
	}

	hh := ssz.NewHasher()

	beaconBlockHeaderToVerifyIndex := blockHeaderIndex

	// validatorIndex := phase0.ValidatorIndex(index)
	beaconStateRoot, err := state.HashTreeRoot()
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with HashTreeRoot of state: %s", err)
		return err
	}

	slot := withdrawalBlockHeader.Slot
	hh.PutUint64(uint64(slot))
	slotRoot := common.ConvertTo32ByteArray(hh.Hash())

	hh.PutUint64(uint64(timestamp))
	timestampRoot := common.ConvertTo32ByteArray(hh.Hash())

	blockHeaderRoot, err := withdrawalBlockHeader.HashTreeRoot()
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with HashTreeRoot of blockHeader: %s", err)
		return err
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error creating EPP object: %s", err)
		return err
	}
	versionedState, err = beacon.CreateVersionedState(&state)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with CreateVersionedState: %s", err)
		return err
	}

	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(&versionedState)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with ComputeBeaconStateTopLevelRoots: %s", err)
		return err
	}

	withdrawalProof, _, err := epp.ProveWithdrawal(&oracleBeaconBlockHeader, &versionedState, oracleBeaconStateTopLevelRoots, historicalSummaryStateBlockRoots, &versionedSignedBlock, uint64(validatorIndex))
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with ProveWithdrawal: %s", err)
		return err
	}
	stateRootProofAgainstBlockHeader, err := beacon.ProveStateRootAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with ProveStateRootAgainstBlockHeader: %s", err)
		return err
	}
	slotProofAgainstBlockHeader, err := beacon.ProveSlotAgainstBlockHeader(&oracleBeaconBlockHeader)
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with ProveSlotAgainstBlockHeader: %s", err)
		return err
	}

	validatorProof, err := epp.ProveValidatorAgainstBeaconState(oracleBeaconStateTopLevelRoots, state.Slot, state.Validators, uint64(validatorIndex))
	if err != nil {
		log.Debug().Msgf("GenerateWithdrawalFieldsProof: error with ProveValidatorAgainstBeaconState: %s", err)
		return err
	}
	proofs := commonutils.WithdrawalProofs{
		StateRootAgainstLatestBlockHeaderProof: commonutils.ConvertBytesToStrings(stateRootProofAgainstBlockHeader),
		SlotAgainstLatestBlockHeaderProof:      commonutils.ConvertBytesToStrings(slotProofAgainstBlockHeader),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		WithdrawalProof:                        commonutils.ConvertBytesToStrings(withdrawalProof.WithdrawalProof),
		SlotProof:                              commonutils.ConvertBytesToStrings(withdrawalProof.SlotProof),
		ExecutionPayloadProof:                  commonutils.ConvertBytesToStrings(withdrawalProof.ExecutionPayloadProof),
		TimestampProof:                         commonutils.ConvertBytesToStrings(withdrawalProof.TimestampProof),
		HistoricalSummaryProof:                 commonutils.ConvertBytesToStrings(withdrawalProof.HistoricalSummaryBlockRootProof),
		BlockHeaderRootIndex:                   beaconBlockHeaderToVerifyIndex,
		HistoricalSummaryIndex:                 uint64(historicalSummariesIndex),
		WithdrawalIndex:                        withdrawalIndex,
		BlockHeaderRoot:                        "0x" + hex.EncodeToString(blockHeaderRoot[:]),
		SlotRoot:                               "0x" + hex.EncodeToString(slotRoot[:]),
		TimestampRoot:                          "0x" + hex.EncodeToString(timestampRoot[:]),
		ExecutionPayloadRoot:                   "0x" + hex.EncodeToString(executionPayloadRoot[:]),
		ValidatorProof:                         commonutils.ConvertBytesToStrings(validatorProof),
		ValidatorFields:                        commonutils.GetValidatorFields(state.Validators[validatorIndex]),
		WithdrawalFields:                       commonutils.GetWithdrawalFields(withdrawal),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().Msgf("JSON marshal error: : %s", err)
	}

	_ = os.WriteFile(outputFile, proofData, 0644)

	return nil
}
