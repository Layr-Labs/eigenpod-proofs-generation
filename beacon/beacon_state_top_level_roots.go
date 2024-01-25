package beacon

import (
	"reflect"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type BeaconStateTopLevelRoots struct {
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

func ProveBeaconTopLevelRootAgainstBeaconState(beaconTopLevelRoots *BeaconStateTopLevelRoots, index uint64) (common.Proof, error) {
	v := reflect.ValueOf(*beaconTopLevelRoots)
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

	return common.GetProof(roots, index, beaconStateMerkleSubtreeNumLayers)
}
