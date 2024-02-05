package beacon

import (
	"errors"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
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
