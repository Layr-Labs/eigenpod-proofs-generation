package main

import (
	"context"
	"fmt"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ProofServer struct {
	blockHeaderProvider eth2client.BeaconBlockHeadersProvider
	stateProvider       eth2client.BeaconStateProvider
}

func NewProofServer(provider string) (server *ProofServer) {
	// Provide a cancellable context to the creation function.
	ctx, _ := context.WithCancel(context.Background())
	client, err := http.New(ctx,
		// WithAddress supplies the address of the beacon node, as a URL.
		http.WithAddress(provider),
		// LogLevel supplies the level of logging to carry out.
		http.WithLogLevel(zerolog.WarnLevel),
	)
	if err != nil {
		panic(err)
	}

	if provider, isProvider := client.(eth2client.BeaconBlockHeadersProvider); isProvider {
		server.blockHeaderProvider = provider
	} else {
		panic("not a beacon block header provider")
	}

	if provider, isProvider := client.(eth2client.BeaconStateProvider); isProvider {
		server.stateProvider = provider
	} else {
		panic("not a beacon state provider")
	}

	return server
}

func (s *ProofServer) GetValidatorProof(ctx context.Context, req *ValidatorProofRequest) (*ValidatorProofResponse, error) {
	// TODO: check slot is after deneb fork

	var beaconBlockHeader *v1.BeaconBlockHeader
	var versionedState *spec.VersionedBeaconState
	var ok bool

	blockHeaderResponse, err := s.blockHeaderProvider.BeaconBlockHeader(ctx, &api.BeaconBlockHeaderOpts{Block: strconv.FormatUint(req.Slot, 10)})
	if err != nil {
		return nil, err
	}
	if beaconBlockHeader, ok = any(blockHeaderResponse.Data).(*v1.BeaconBlockHeader); !ok {
		return nil, fmt.Errorf("invalid block header type: %T", blockHeaderResponse.Data)
	}

	beaconStateResponse, err := s.stateProvider.BeaconState(ctx, &api.BeaconStateOpts{State: strconv.FormatUint(req.Slot, 10)})
	if err != nil {
		return nil, err
	}
	if versionedState, ok = any(beaconStateResponse.Data).(*spec.VersionedBeaconState); !ok {
		return nil, fmt.Errorf("invalid beacon state type: %T", beaconStateResponse.Data)
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(req.ChainId, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
		return nil, err
	}

	stateRootProof, validatorContainerProof, err := eigenpodproofs.ProveValidatorFields(epp, beaconBlockHeader.Header.Message, versionedState, req.ValidatorIndex)
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
		return nil, err
	}

	return &ValidatorProofResponse{
		StateRootProof:          stateRootProof.SlotRootProof.ToBytesSlice(),
		ValidatorContainerProof: validatorContainerProof.ToBytesSlice(),
	}, nil
}

func (s *ProofServer) GetWithdrawalProof(ctx context.Context, req *WithdrawalProofRequest) (*WithdrawalProofResponse, error) {
	// TODO: Implement the logic to generate and return withdrawal proof
	return nil, nil
}
