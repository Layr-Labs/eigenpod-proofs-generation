package beacon

import (
	"errors"
	"reflect"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type VersionedBeaconStateTopLevelRoots struct {
	Version spec.DataVersion
	Deneb   *BeaconStateTopLevelRootsDeneb
	Electra *BeaconStateTopLevelRootsElectra
	Fulu    *BeaconStateTopLevelRootsFulu
}

func (v *VersionedBeaconStateTopLevelRoots) GetBalancesRoot() (*phase0.Root, error) {
	switch v.Version {
	case spec.DataVersionDeneb:
		return v.Deneb.BalancesRoot, nil
	case spec.DataVersionElectra:
		return v.Electra.BalancesRoot, nil
	case spec.DataVersionFulu:
		return v.Fulu.BalancesRoot, nil
	default:
		return nil, errors.New("unsupported beacon state version")
	}
}

type BeaconStateTopLevelRootsDeneb struct {
	GenesisTimeRoot                  *phase0.Root
	GenesisValidatorsRoot            *phase0.Root
	SlotRoot                         *phase0.Root
	ForkRoot                         *phase0.Root
	LatestBlockHeaderRoot            *phase0.Root
	BlockRootsRoot                   *phase0.Root
	StateRootsRoot                   *phase0.Root
	HistoricalRootsRoot              *phase0.Root
	ETH1DataRoot                     *phase0.Root
	ETH1DataVotesRoot                *phase0.Root
	ETH1DepositIndexRoot             *phase0.Root
	ValidatorsRoot                   *phase0.Root
	BalancesRoot                     *phase0.Root
	RANDAOMixesRoot                  *phase0.Root
	SlashingsRoot                    *phase0.Root
	PreviousEpochParticipationRoot   *phase0.Root
	CurrentEpochParticipationRoot    *phase0.Root
	JustificationBitsRoot            *phase0.Root
	PreviousJustifiedCheckpointRoot  *phase0.Root
	CurrentJustifiedCheckpointRoot   *phase0.Root
	FinalizedCheckpointRoot          *phase0.Root
	InactivityScoresRoot             *phase0.Root
	CurrentSyncCommitteeRoot         *phase0.Root
	NextSyncCommitteeRoot            *phase0.Root
	LatestExecutionPayloadHeaderRoot *phase0.Root
	NextWithdrawalIndexRoot          *phase0.Root
	NextWithdrawalValidatorIndexRoot *phase0.Root
	HistoricalSummariesRoot          *phase0.Root
}

type BeaconStateTopLevelRootsElectra struct {
	GenesisTimeRoot                   *phase0.Root
	GenesisValidatorsRoot             *phase0.Root
	SlotRoot                          *phase0.Root
	ForkRoot                          *phase0.Root
	LatestBlockHeaderRoot             *phase0.Root
	BlockRootsRoot                    *phase0.Root
	StateRootsRoot                    *phase0.Root
	HistoricalRootsRoot               *phase0.Root
	ETH1DataRoot                      *phase0.Root
	ETH1DataVotesRoot                 *phase0.Root
	ETH1DepositIndexRoot              *phase0.Root
	ValidatorsRoot                    *phase0.Root
	BalancesRoot                      *phase0.Root
	RANDAOMixesRoot                   *phase0.Root
	SlashingsRoot                     *phase0.Root
	PreviousEpochParticipationRoot    *phase0.Root
	CurrentEpochParticipationRoot     *phase0.Root
	JustificationBitsRoot             *phase0.Root
	PreviousJustifiedCheckpointRoot   *phase0.Root
	CurrentJustifiedCheckpointRoot    *phase0.Root
	FinalizedCheckpointRoot           *phase0.Root
	InactivityScoresRoot              *phase0.Root
	CurrentSyncCommitteeRoot          *phase0.Root
	NextSyncCommitteeRoot             *phase0.Root
	LatestExecutionPayloadHeaderRoot  *phase0.Root
	NextWithdrawalIndexRoot           *phase0.Root
	NextWithdrawalValidatorIndexRoot  *phase0.Root
	HistoricalSummariesRoot           *phase0.Root
	DepositRequestsStartIndexRoot     *phase0.Root
	DepositBalanceToConsumeRoot       *phase0.Root
	ExitBalanceToConsumeRoot          *phase0.Root
	EarliestExitEpochRoot             *phase0.Root
	ConsolidationBalanceToConsumeRoot *phase0.Root
	EarliestConsolidationEpochRoot    *phase0.Root
	PendingDepositsRoot               *phase0.Root
	PendingPartialWithdrawalsRoot     *phase0.Root
	PendingConsolidationsRoot         *phase0.Root
}

type BeaconStateTopLevelRootsFulu struct {
	GenesisTimeRoot                   *phase0.Root
	GenesisValidatorsRoot             *phase0.Root
	SlotRoot                          *phase0.Root
	ForkRoot                          *phase0.Root
	LatestBlockHeaderRoot             *phase0.Root
	BlockRootsRoot                    *phase0.Root
	StateRootsRoot                    *phase0.Root
	HistoricalRootsRoot               *phase0.Root
	ETH1DataRoot                      *phase0.Root
	ETH1DataVotesRoot                 *phase0.Root
	ETH1DepositIndexRoot              *phase0.Root
	ValidatorsRoot                    *phase0.Root
	BalancesRoot                      *phase0.Root
	RANDAOMixesRoot                   *phase0.Root
	SlashingsRoot                     *phase0.Root
	PreviousEpochParticipationRoot    *phase0.Root
	CurrentEpochParticipationRoot     *phase0.Root
	JustificationBitsRoot             *phase0.Root
	PreviousJustifiedCheckpointRoot   *phase0.Root
	CurrentJustifiedCheckpointRoot    *phase0.Root
	FinalizedCheckpointRoot           *phase0.Root
	InactivityScoresRoot              *phase0.Root
	CurrentSyncCommitteeRoot          *phase0.Root
	NextSyncCommitteeRoot             *phase0.Root
	LatestExecutionPayloadHeaderRoot  *phase0.Root
	NextWithdrawalIndexRoot           *phase0.Root
	NextWithdrawalValidatorIndexRoot  *phase0.Root
	HistoricalSummariesRoot           *phase0.Root
	DepositRequestsStartIndexRoot     *phase0.Root
	DepositBalanceToConsumeRoot       *phase0.Root
	ExitBalanceToConsumeRoot          *phase0.Root
	EarliestExitEpochRoot             *phase0.Root
	ConsolidationBalanceToConsumeRoot *phase0.Root
	EarliestConsolidationEpochRoot    *phase0.Root
	PendingDepositsRoot               *phase0.Root
	PendingPartialWithdrawalsRoot     *phase0.Root
	PendingConsolidationsRoot         *phase0.Root
	ProposerLookaheadRoot             *phase0.Root
}

func ProveBeaconTopLevelRootAgainstBeaconState(beaconTopLevelRoots *VersionedBeaconStateTopLevelRoots, index uint64) (common.Proof, error) {
	var v reflect.Value
	var treeHeight uint64
	switch beaconTopLevelRoots.Version {
	case spec.DataVersionDeneb:
		v = reflect.ValueOf(*beaconTopLevelRoots.Deneb)
		treeHeight = BEACON_STATE_TREE_HEIGHT_DENEB
	case spec.DataVersionElectra:
		v = reflect.ValueOf(*beaconTopLevelRoots.Electra)
		treeHeight = BEACON_STATE_TREE_HEIGHT_ELECTRA
	case spec.DataVersionFulu:
		v = reflect.ValueOf(*beaconTopLevelRoots.Fulu)
		treeHeight = BEACON_STATE_TREE_HEIGHT_FULU
	default:
		return nil, errors.New("unsupported beacon state version")
	}

	beaconTopLevelRootsList := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		r := v.Field(i).Interface()
		typedR := r.(*phase0.Root)
		beaconTopLevelRootsList[i] = *typedR
	}

	roots := make([]phase0.Root, len(beaconTopLevelRootsList))
	for i, v := range beaconTopLevelRootsList {
		roots[i] = v.(phase0.Root)
	}

	return common.GetProof(roots, index, treeHeight)
}
