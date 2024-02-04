package eigenpodproofs

import (
	"errors"
	"math"
	"time"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog/log"
)

// uint64 oracleTimestamp,
// BeaconChainProofs.StateRootProof calldata stateRootProof,
// BeaconChainProofs.WithdrawalProof[] calldata withdrawalProofs,
// bytes[] calldata validatorFieldsProofs,
// bytes32[][] calldata validatorFields,
// bytes32[][] calldata withdrawalFields

type VerifyAndProcessWithdrawalCallParams struct {
	OracleTimestamp       uint64             `json:"oracleTimestamp"`
	StateRootProof        *StateRootProof    `json:"stateRootProof"`
	WithdrawalProofs      []*WithdrawalProof `json:"withdrawalProofs"`
	ValidatorFieldsProofs []common.Proof     `json:"validatorFieldsProofs"`
	ValidatorFields       [][]Bytes32        `json:"validatorFields"`
	WithdrawalFields      [][]Bytes32        `json:"withdrawalFields"`
}

type WithdrawalProof struct {
	WithdrawalProof                 common.Proof `json:"withdrawalProof"`
	SlotProof                       common.Proof `json:"slotProof"`
	ExecutionPayloadProof           common.Proof `json:"executionPayloadProof"`
	TimestampProof                  common.Proof `json:"timestampProof"`
	HistoricalSummaryBlockRootProof common.Proof `json:"historicalSummaryBlockRootProof"`
	BlockRootIndex                  uint64       `json:"blockRootIndex"`
	HistoricalSummaryIndex          uint64       `json:"historicalSummaryIndex"`
	WithdrawalIndex                 uint64       `json:"withdrawalIndex"`
	BlockRoot                       phase0.Root  `json:"blockRoot"`
	SlotRoot                        phase0.Root  `json:"slotRoot"`
	TimestampRoot                   phase0.Root  `json:"timestampRoot"`
	ExecutionPayloadRoot            phase0.Root  `json:"executionPayloadRoot"`
}

type StateRootProof struct {
	BeaconStateRoot phase0.Root  `json:"beaconStateRoot"`
	StateRootProof  common.Proof `json:"stateRootProof"`
	Slot            phase0.Slot  `json:"slot"`
	SlotRootProof   common.Proof `json:"slotRootProof"` //Note:  this slot root is oracle block root being used to prove partial withdrawals is after the specified range of blocks requested by the user
}

const FIRST_CAPELLA_SLOT_GOERLI = uint64(5193728)
const FIRST_CAPELLA_SLOT_MAINNET = uint64(6209536)

func IsProvableWithdrawal(latestOracleBeaconSlot, withdrawalSlot uint64) bool {
	return latestOracleBeaconSlot > beacon.SlotsPerHistoricalRoot+withdrawalSlot
}

func (epp *EigenPodProofs) GetWithdrawalProofParams(latestOracleBeaconSlot, withdrawalSlot uint64) (uint64, error) {
	if withdrawalSlot > latestOracleBeaconSlot {
		return 0, errors.New("withdrawal slot is after than the latest oracle beacon slot")
	} else if latestOracleBeaconSlot-withdrawalSlot < beacon.SlotsPerHistoricalRoot {
		return 0, errors.New("oracle beacon slot does not have enough historical summaries to prove withdrawal")
	}

	var FIRST_CAPELLA_SLOT uint64
	if epp.chainID == 5 {
		FIRST_CAPELLA_SLOT = FIRST_CAPELLA_SLOT_GOERLI
	} else if epp.chainID == 1 {
		FIRST_CAPELLA_SLOT = FIRST_CAPELLA_SLOT_MAINNET
	}
	// index of the historical summary in the array of historical_summaries
	historicalSummaryIndex := (withdrawalSlot - FIRST_CAPELLA_SLOT) / beacon.SlotsPerHistoricalRoot

	// slot of which the beacon state is retrieved for getting the block roots array containing the old block with the old withdrawal
	historicalSummarySlot := FIRST_CAPELLA_SLOT + (historicalSummaryIndex+1)*beacon.SlotsPerHistoricalRoot

	return historicalSummarySlot, nil
}

func (epp *EigenPodProofs) ProveWithdrawals(
	oracleBlockHeader *phase0.BeaconBlockHeader,
	oracleBeaconState *spec.VersionedBeaconState,
	historicalSummaryStateBlockRoots [][]phase0.Root,
	withdrawalBlocks []*spec.VersionedBeaconBlock,
	validatorIndices []uint64,
) (*VerifyAndProcessWithdrawalCallParams, error) {
	oracleBeaconStateValidators, err := oracleBeaconState.Validators()
	if err != nil {
		return nil, err
	}
	oracleBeaconStateSlot, err := oracleBeaconState.Slot()
	if err != nil {
		return nil, err
	}
	verifyAndProcessWithdrawalCallParams := &VerifyAndProcessWithdrawalCallParams{}
	verifyAndProcessWithdrawalCallParams.StateRootProof = &StateRootProof{}
	// Get beacon state top level roots
	oracleBeaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	// Get beacon state root.
	verifyAndProcessWithdrawalCallParams.StateRootProof.BeaconStateRoot = oracleBlockHeader.StateRoot
	if err != nil {
		return nil, err
	}

	verifyAndProcessWithdrawalCallParams.StateRootProof.StateRootProof, err = beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	// Note: this slot and slot root proof are used for partial withdrawal proofs to ensure that the oracle root slot is after the specified range of blocks requested by the user
	verifyAndProcessWithdrawalCallParams.StateRootProof.Slot = oracleBlockHeader.Slot

	verifyAndProcessWithdrawalCallParams.StateRootProof.SlotRootProof, err = beacon.ProveSlotAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	verifyAndProcessWithdrawalCallParams.OracleTimestamp = GetSlotTimestamp(oracleBeaconState, oracleBlockHeader)

	numWithdrawals := len(withdrawalBlocks)

	verifyAndProcessWithdrawalCallParams.WithdrawalProofs = make([]*WithdrawalProof, numWithdrawals)
	verifyAndProcessWithdrawalCallParams.WithdrawalFields = make([][]Bytes32, numWithdrawals)

	verifyAndProcessWithdrawalCallParams.ValidatorFieldsProofs = make([]common.Proof, numWithdrawals)
	verifyAndProcessWithdrawalCallParams.ValidatorFields = make([][]Bytes32, numWithdrawals)

	for i := 0; i < numWithdrawals; i++ {
		start := time.Now()

		if withdrawalBlocks[i].Version == spec.DataVersionDeneb {
			verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i], err = epp.ProveDenebWithdrawal(oracleBlockHeader, oracleBeaconState, oracleBeaconStateTopLevelRoots, historicalSummaryStateBlockRoots[i], withdrawalBlocks[i].Deneb, validatorIndices[i])
			if err != nil {
				return nil, err
			}
			verifyAndProcessWithdrawalCallParams.WithdrawalFields[i] = ConvertWithdrawalToWithdrawalFields(withdrawalBlocks[i].Deneb.Body.ExecutionPayload.Withdrawals[verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].WithdrawalIndex])
			log.Info().Msgf("time to prove withdrawal: %s", time.Since(start))
		} else {
			verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i], err = epp.ProveCapellaWithdrawal(oracleBlockHeader, oracleBeaconState, oracleBeaconStateTopLevelRoots, historicalSummaryStateBlockRoots[i], withdrawalBlocks[i].Capella, validatorIndices[i])
			if err != nil {
				return nil, err
			}
			verifyAndProcessWithdrawalCallParams.WithdrawalFields[i] = ConvertWithdrawalToWithdrawalFields(withdrawalBlocks[i].Capella.Body.ExecutionPayload.Withdrawals[verifyAndProcessWithdrawalCallParams.WithdrawalProofs[i].WithdrawalIndex])
			log.Info().Msgf("time to prove withdrawal: %s", time.Since(start))
		}

		start = time.Now()
		// prove validator
		verifyAndProcessWithdrawalCallParams.ValidatorFieldsProofs[i], err = epp.ProveValidatorAgainstBeaconState(oracleBeaconStateSlot, oracleBeaconStateValidators, oracleBeaconStateTopLevelRoots, validatorIndices[i])
		if err != nil {
			return nil, err
		}
		verifyAndProcessWithdrawalCallParams.ValidatorFields[i] = ConvertValidatorToValidatorFields(oracleBeaconStateValidators[validatorIndices[i]])
		log.Info().Msgf("time to prove validator: %s", time.Since(start))
	}

	return verifyAndProcessWithdrawalCallParams, nil
}

// ProveWithdrawal generates the proofs required to prove a withdrawal
// oracleBlockHeader: the root of this is provided by the oracle, we prove the state root against this
// oracleBeaconState: the state of the block header provided by the oracle
// historicalSummaryState: the state whose slot at which historicalSummaryState.block_roots was hashed and added to historical_summaries
// withdrawalBlock: the block containing the withdrawal
// validatorIndex: the index of the validator that the withdrawal happened for
func (epp *EigenPodProofs) ProveDenebWithdrawal(
	oracleBlockHeader *phase0.BeaconBlockHeader,
	oracleBeaconState *spec.VersionedBeaconState,
	oracleBeaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots,
	historicalSummaryStateBlockRoots []phase0.Root,
	withdrawalBlock *deneb.BeaconBlock,
	validatorIndex uint64,
) (*WithdrawalProof, error) {
	start := time.Now()
	// compute the withdrawal body root
	blockBodyRoot, err := withdrawalBlock.Body.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("time to compute block body root: %s", time.Since(start))

	// setup the withdrawal block header
	withdrawalBlockHeader := &phase0.BeaconBlockHeader{
		Slot:          withdrawalBlock.Slot,
		ProposerIndex: withdrawalBlock.ProposerIndex,
		ParentRoot:    withdrawalBlock.ParentRoot,
		StateRoot:     withdrawalBlock.StateRoot,
		BodyRoot:      blockBodyRoot,
	}

	// initialize the withdrawal proof
	withdrawalProof := &WithdrawalProof{}

	start = time.Now()
	// prove the execution payload against the withdrawal block header
	withdrawalProof.ExecutionPayloadProof, withdrawalProof.ExecutionPayloadRoot, err = beacon.ProveExecutionPayloadAgainstBlockHeaderDeneb(withdrawalBlockHeader, withdrawalBlock.Body)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("time to prove execution payload against block header: %s", time.Since(start))

	start = time.Now()
	// calculate execution payload field roots
	withdrawalExecutionPayloadFieldRoots, err := beacon.ComputeExecutionPayloadFieldRootsDeneb(withdrawalBlock.Body.ExecutionPayload)
	if err != nil {
		return nil, err
	}

	oracleBeaconStateHistoricalSummaries, err := beacon.HistoricalSummaries(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	err = epp.proveWithdrawal(
		withdrawalProof,
		oracleBlockHeader,
		oracleBeaconStateHistoricalSummaries,
		oracleBeaconStateTopLevelRoots,
		historicalSummaryStateBlockRoots,
		withdrawalBlockHeader,
		withdrawalExecutionPayloadFieldRoots,
		withdrawalBlock.Body.ExecutionPayload.Withdrawals,
		withdrawalBlock.Body.ExecutionPayload.Timestamp,
		validatorIndex,
	)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("time to prove withdrawal: %s", time.Since(start))

	return withdrawalProof, nil
}

func (epp *EigenPodProofs) ProveCapellaWithdrawal(
	oracleBlockHeader *phase0.BeaconBlockHeader,
	oracleBeaconState *spec.VersionedBeaconState,
	oracleBeaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots,
	historicalSummaryStateBlockRoots []phase0.Root,
	withdrawalBlock *capella.BeaconBlock,
	validatorIndex uint64,
) (*WithdrawalProof, error) {
	start := time.Now()
	// compute the withdrawal body root
	blockBodyRoot, err := withdrawalBlock.Body.HashTreeRoot()
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("time to compute block body root: %s", time.Since(start))

	// setup the withdrawal block header
	withdrawalBlockHeader := &phase0.BeaconBlockHeader{
		Slot:          withdrawalBlock.Slot,
		ProposerIndex: withdrawalBlock.ProposerIndex,
		ParentRoot:    withdrawalBlock.ParentRoot,
		StateRoot:     withdrawalBlock.StateRoot,
		BodyRoot:      blockBodyRoot,
	}

	// initialize the withdrawal proof
	withdrawalProof := &WithdrawalProof{}

	start = time.Now()
	// prove the execution payload against the withdrawal block header
	withdrawalProof.ExecutionPayloadProof, withdrawalProof.ExecutionPayloadRoot, err = beacon.ProveExecutionPayloadAgainstBlockHeaderCapella(withdrawalBlockHeader, withdrawalBlock.Body)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("time to prove execution payload against block header: %s", time.Since(start))

	start = time.Now()
	// calculate execution payload field roots
	withdrawalExecutionPayloadFieldRoots, err := beacon.ComputeExecutionPayloadFieldRootsCapella(withdrawalBlock.Body.ExecutionPayload)
	if err != nil {
		return nil, err
	}

	oracleBeaconStateHistoricalSummaries, err := beacon.HistoricalSummaries(oracleBeaconState)
	if err != nil {
		return nil, err
	}
	err = epp.proveWithdrawal(
		withdrawalProof,
		oracleBlockHeader,
		oracleBeaconStateHistoricalSummaries,
		oracleBeaconStateTopLevelRoots,
		historicalSummaryStateBlockRoots,
		withdrawalBlockHeader,
		withdrawalExecutionPayloadFieldRoots,
		withdrawalBlock.Body.ExecutionPayload.Withdrawals,
		withdrawalBlock.Body.ExecutionPayload.Timestamp,
		validatorIndex,
	)
	if err != nil {
		return nil, err
	}

	return withdrawalProof, nil
}

func (epp *EigenPodProofs) proveWithdrawal(
	withdrawalProof *WithdrawalProof,
	oracleBlockHeader *phase0.BeaconBlockHeader,
	oracleBeaconStateHistoricalSummaries []*capella.HistoricalSummary,
	oracleBeaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots,
	historicalSummaryStateBlockRoots []phase0.Root,
	withdrawalBlockHeader *phase0.BeaconBlockHeader,
	withdrawalExecutionPayloadFieldRoots []phase0.Root,
	withdrawals []*capella.Withdrawal,
	withdrawalTimestamp uint64,
	validatorIndex uint64,
) error {
	withdrawalProof.WithdrawalIndex = math.MaxUint64 // max uint 64 value
	for i := 0; i < len(withdrawals); i++ {
		if uint64(withdrawals[i].ValidatorIndex) == validatorIndex {
			withdrawalProof.WithdrawalIndex = uint64(i)
			break
		}
	}
	if withdrawalProof.WithdrawalIndex == math.MaxUint64 {
		return errors.New("validator index not found in withdrawal block")
	}

	var FIRST_CAPELLA_SLOT uint64
	if epp.chainID == 5 {
		FIRST_CAPELLA_SLOT = FIRST_CAPELLA_SLOT_GOERLI
	} else if epp.chainID == 1 {
		FIRST_CAPELLA_SLOT = FIRST_CAPELLA_SLOT_MAINNET
	}

	withdrawalSlotUint64 := uint64(withdrawalBlockHeader.Slot)

	// index of the historical summary in the array of historical_summaries
	withdrawalProof.HistoricalSummaryIndex = (withdrawalSlotUint64 - FIRST_CAPELLA_SLOT) / beacon.SlotsPerHistoricalRoot

	// index of the block containing the target withdrawal in the block roots array
	withdrawalProof.BlockRootIndex = withdrawalSlotUint64 % beacon.SlotsPerHistoricalRoot
	withdrawalProof.BlockRoot = historicalSummaryStateBlockRoots[withdrawalProof.BlockRootIndex]

	// make sure the withdrawal index is in range
	if len(withdrawals) <= int(withdrawalProof.WithdrawalIndex) {
		return errors.New("withdrawal index is out of range")
	}

	// log the time it takes to compute each proof
	log.Info().Msg("computing withdrawal proof")

	var err error
	start := time.Now()
	// prove the withdrawal against the execution payload
	withdrawalProof.WithdrawalProof, err = beacon.ProveWithdrawalAgainstExecutionPayload(
		withdrawalExecutionPayloadFieldRoots,
		withdrawals,
		uint8(withdrawalProof.WithdrawalIndex),
	)
	if err != nil {
		return err
	}
	log.Info().Msgf("time to prove withdrawal against execution payload: %s", time.Since(start))

	start = time.Now()
	// prove the slot against the withdrawal block header
	withdrawalProof.SlotProof, err = beacon.ProveSlotAgainstBlockHeader(withdrawalBlockHeader)
	if err != nil {
		return err
	}
	log.Info().Msgf("time to prove slot against block header: %s", time.Since(start))
	withdrawalProof.SlotRoot = ConvertUint64ToRoot(uint64(withdrawalBlockHeader.Slot))

	start = time.Now()
	// prove the timestamp against the execution payload
	withdrawalProof.TimestampProof, err = beacon.ProveTimestampAgainstExecutionPayload(withdrawalExecutionPayloadFieldRoots)
	if err != nil {
		return err
	}
	withdrawalProof.TimestampRoot = ConvertUint64ToRoot(withdrawalTimestamp)
	log.Info().Msgf("time to prove timestamp against execution payload: %s", time.Since(start))

	start = time.Now()
	// prove the withdrawal block root against the oracle state root
	withdrawalProof.HistoricalSummaryBlockRootProof, err = beacon.ProveBlockRootAgainstBeaconStateViaHistoricalSummaries(oracleBeaconStateTopLevelRoots, oracleBeaconStateHistoricalSummaries, historicalSummaryStateBlockRoots, withdrawalProof.HistoricalSummaryIndex, withdrawalProof.BlockRootIndex)
	if err != nil {
		return err
	}
	log.Info().Msgf("time to prove block root against beacon state via historical summaries: %s", time.Since(start))

	return nil
}
