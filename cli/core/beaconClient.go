package core

import (
	"context"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrNoEigenPod               = errors.New("no eigenpod associated with owner")
	ErrBeaconClientNotSupported = errors.New("could not instantiate beacon chain client")
	ErrValidatorNotFound        = errors.New("validator not found")
)

type BeaconClient interface {
	GetBeaconHeader(ctx context.Context, blockId string) (*v1.BeaconBlockHeader, error)
	GetBeaconState(ctx context.Context, stateId string) (*spec.VersionedBeaconState, error)
}

type beaconClient struct {
	eth2client eth2client.Service
}

func NewBeaconClient(endpoint string) (BeaconClient, context.CancelFunc, error) {
	beaconClient := beaconClient{}
	ctx, cancel := context.WithCancel(context.Background())

	client, err := http.New(ctx,
		// WithAddress supplies the address of the beacon node, as a URL.
		http.WithAddress(endpoint),
		http.WithLogLevel(zerolog.WarnLevel),
		http.WithTimeout(300*time.Second),
	)
	if err != nil {
		return nil, cancel, err
	}
	log.Info().Msgf("Connected to %s\n", client.Name())

	beaconClient.eth2client = client
	return &beaconClient, cancel, nil
}

func (b *beaconClient) GetBeaconHeader(ctx context.Context, blockId string) (*v1.BeaconBlockHeader, error) {
	if provider, isProvider := b.eth2client.(eth2client.BeaconBlockHeadersProvider); isProvider {
		opts := &api.BeaconBlockHeaderOpts{Block: blockId}
		response, err := provider.BeaconBlockHeader(ctx, opts)
		if err != nil {
			return nil, err
		}
		return response.Data, nil
	}

	return nil, ErrBeaconClientNotSupported
}

func (b *beaconClient) GetBeaconState(ctx context.Context, stateId string) (*spec.VersionedBeaconState, error) {
	if provider, ok := b.eth2client.(eth2client.BeaconStateProvider); ok {
		log.Info().Msgf("downloading beacon state %s", stateId)
		opts := &api.BeaconStateOpts{State: stateId}
		beaconState, err := provider.BeaconState(ctx, opts)
		if err != nil {
			return nil, err
		}

		if beaconState == nil {
			return nil, errors.New("beacon state is nil")
		}

		log.Info().Msg("finished download")
		return beaconState.Data, nil
	}

	return nil, ErrBeaconClientNotSupported
}
