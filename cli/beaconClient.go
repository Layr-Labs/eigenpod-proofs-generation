package main

import (
	"context"
	"time"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const slotsPerEpoch = 32

var (
	ErrNoEigenPod               = errors.New("no eigenpod associated with owner")
	ErrBeaconClientNotSupported = errors.New("could not instantiate beacon chain client")
	ErrValidatorNotFound        = errors.New("validator not found")
)

type BeaconClient interface {
	GetLatestEpoch(ctx context.Context) (phase0.Epoch, error)
	GetLatestSlot(ctx context.Context) (phase0.Slot, error)
	GetBeaconHeader(ctx context.Context, blockId string) (*v1.BeaconBlockHeader, error)
	GetSignedBeaconBlock(ctx context.Context, blockId string) (*spec.VersionedSignedBeaconBlock, error)
	GetBeaconState(ctx context.Context, stateId string) (*spec.VersionedBeaconState, error)
	GetValidator(ctx context.Context, stateID string, validatorIndex phase0.ValidatorIndex) (*v1.Validator, error)
	GetChainGenesisTime(ctx context.Context) (time.Time, error)
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

func (b *beaconClient) GetSignedBeaconBlock(ctx context.Context, blockId string) (*spec.VersionedSignedBeaconBlock, error) {
	if provider, ok := b.eth2client.(eth2client.SignedBeaconBlockProvider); ok {
		opts := &api.SignedBeaconBlockOpts{Block: blockId}
		response, err := provider.SignedBeaconBlock(ctx, opts)
		if err != nil {
			return nil, err
		}
		return response.Data, nil
	}

	return nil, ErrBeaconClientNotSupported
}

func (b *beaconClient) GetValidator(ctx context.Context, stateID string, validatorIndex phase0.ValidatorIndex) (*v1.Validator, error) {
	if provider, ok := b.eth2client.(eth2client.ValidatorsProvider); ok {
		opts := &api.ValidatorsOpts{
			State:   stateID,
			Indices: []phase0.ValidatorIndex{validatorIndex},
		}
		response, err := provider.Validators(ctx, opts)
		if err != nil {
			return nil, err
		}
		validators := response.Data
		var validator *v1.Validator
		if validator, ok = validators[validatorIndex]; !ok {
			return nil, ErrValidatorNotFound
		}
		return validator, nil
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
		return beaconState.Data, nil
	}

	return nil, ErrBeaconClientNotSupported
}

func (b *beaconClient) GetLatestEpoch(ctx context.Context) (phase0.Epoch, error) {
	header, err := b.GetBeaconHeader(ctx, "head")
	if err != nil {
		return 0, err
	}

	return phase0.Epoch(header.Header.Message.Slot / slotsPerEpoch), nil
}

func (b *beaconClient) GetLatestSlot(ctx context.Context) (phase0.Slot, error) {
	header, err := b.GetBeaconHeader(ctx, "head")
	if err != nil {
		return 0, err
	}

	return header.Header.Message.Slot, nil
}

func (b *beaconClient) GetChainGenesisTime(ctx context.Context) (time.Time, error) {
	if provider, ok := b.eth2client.(eth2client.GenesisProvider); ok {
		opts := api.GenesisOpts{}
		response, err := provider.Genesis(ctx, &opts)
		if err != nil {
			return time.Time{}, err
		}
		return response.Data.GenesisTime, nil
	}

	return time.Time{}, ErrBeaconClientNotSupported
}
