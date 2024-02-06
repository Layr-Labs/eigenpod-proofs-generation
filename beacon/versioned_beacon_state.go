package beacon

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
)

func HistoricalSummaries(state *spec.VersionedBeaconState) ([]*capella.HistoricalSummary, error) {
	switch state.Version {
	case spec.DataVersionCapella:
		return state.Capella.HistoricalSummaries, nil
	case spec.DataVersionDeneb:
		return state.Deneb.HistoricalSummaries, nil
	default:
		return nil, errors.New("unsupported beacon state version")
	}
}

func GenesisTime(state *spec.VersionedBeaconState) (uint64, error) {
	switch state.Version {
	case spec.DataVersionCapella:
		return state.Capella.GenesisTime, nil
	case spec.DataVersionDeneb:
		return state.Deneb.GenesisTime, nil
	default:
		return 0, errors.New("unsupported beacon state version")
	}
}
func CreateVersionedSignedBlock(block interface{}) spec.VersionedSignedBeaconBlock {
	var versionedBlock spec.VersionedSignedBeaconBlock
	switch s := block.(type) {
	case *deneb.BeaconBlock:
		versionedBlock.Deneb.Message = s
		versionedBlock.Version = spec.DataVersionDeneb
	case *capella.BeaconBlock:
		versionedBlock.Capella.Message = s
		versionedBlock.Version = spec.DataVersionCapella
	}
	return versionedBlock
}

func CreateVersionedState(state interface{}) spec.VersionedBeaconState {
	var versionedState spec.VersionedBeaconState

	switch s := state.(type) {
	case *deneb.BeaconState:
		versionedState.Deneb = s
		versionedState.Version = spec.DataVersionDeneb
	case *capella.BeaconState:
		versionedState.Capella = s
		versionedState.Version = spec.DataVersionCapella
	}
	return versionedState
}
