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
