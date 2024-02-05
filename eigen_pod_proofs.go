package eigenpodproofs

import (
	"encoding/json"
	"errors"

	"strconv"
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/hashicorp/golang-lru/v2/expirable"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

const (
	BEACON_STATE_ROOT_PREFIX            = "BEACON_STATE_ROOT_"
	BEACON_STATE_TOP_LEVEL_ROOTS_PREFIX = "BEACON_STATE_TOP_LEVEL_ROOTS_"
	VALIDATOR_TREE_PREFIX               = "VALIDATOR_TREE_"
	MAX_ORACLE_STATE_CACHE_SIZE         = 2000000
)

type EigenPodProofs struct {
	chainID                       uint64
	oracleStateCache              *expirable.LRU[string, []byte]
	oracleStateCacheExpirySeconds int
}

func NewEigenPodProofs(chainID uint64, oracleStateCacheExpirySeconds int) (*EigenPodProofs, error) {
	if chainID != 1 && chainID != 5 {
		return nil, errors.New("chainID not supported")
	}
	// note that TTL applies equally to each entry
	oracleStateCache := expirable.NewLRU[string, []byte](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)

	return &EigenPodProofs{
		chainID:                       chainID,
		oracleStateCache:              oracleStateCache,
		oracleStateCacheExpirySeconds: oracleStateCacheExpirySeconds,
	}, nil
}

func (epp *EigenPodProofs) ComputeBeaconStateRoot(beaconState *deneb.BeaconState) (phase0.Root, error) {
	beaconStateRootSlice, err := epp.loadOrComputeBeaconData(
		BEACON_STATE_ROOT_PREFIX,
		beaconState.Slot,
		func() ([]byte, error) {
			stateRoot, err := beaconState.HashTreeRoot()
			if err != nil {
				return nil, err
			}
			return stateRoot[:], nil
		},
	)
	if err != nil {
		return phase0.Root{}, err
	}
	var beaconStateRoot phase0.Root
	copy(beaconStateRoot[:], beaconStateRootSlice)
	return beaconStateRoot, nil
}

func (epp *EigenPodProofs) ComputeBeaconStateTopLevelRoots(beaconState *spec.VersionedBeaconState) (*beacon.BeaconStateTopLevelRoots, error) {
	//get the versioned beacon state's slot
	slot, err := beaconState.Slot()
	if err != nil {
		return nil, err
	}

	beaconStateTopLevelRootsSlice, err := epp.loadOrComputeBeaconData(
		BEACON_STATE_TOP_LEVEL_ROOTS_PREFIX,
		slot,
		func() ([]byte, error) {
			beaconStateTopLevelRoots, err := epp.ComputeVersionedBeaconStateTopLevelRoots(beaconState)
			if err != nil {
				return nil, err
			}
			return json.Marshal(beaconStateTopLevelRoots)
		},
	)
	if err != nil {
		return nil, err
	}
	beaconStateTopLevelRoots := &beacon.BeaconStateTopLevelRoots{}
	err = json.Unmarshal(beaconStateTopLevelRootsSlice, beaconStateTopLevelRoots)
	return beaconStateTopLevelRoots, err
}

func (epp *EigenPodProofs) ComputeVersionedBeaconStateTopLevelRoots(beaconState *spec.VersionedBeaconState) (*beacon.BeaconStateTopLevelRoots, error) {
	switch beaconState.Version {
	case spec.DataVersionDeneb:
		return beacon.ComputeBeaconStateTopLevelRootsDeneb(beaconState.Deneb)
	case spec.DataVersionCapella:
		return beacon.ComputeBeaconStateTopLevelRootsCapella(beaconState.Capella)
	default:
		return nil, errors.New("unsupported beacon state version")
	}
}

func (epp *EigenPodProofs) ComputeValidatorTree(slot phase0.Slot, validators []*phase0.Validator) ([][]phase0.Root, error) {
	validatorTreeSlice, err := epp.loadOrComputeBeaconData(
		VALIDATOR_TREE_PREFIX,
		slot,
		func() ([]byte, error) {
			// compute the validator tree leaves
			validatorLeaves, err := beacon.ComputeValidatorTreeLeaves(validators)
			if err != nil {
				return nil, err
			}

			// compute the validator tree
			validatorTree, err := common.ComputeMerkleTreeFromLeaves(validatorLeaves, beacon.ValidatorListMerkleSubtreeNumLayers)
			if err != nil {
				return nil, err
			}

			// cache the validator tree
			validatorTreeSlice, err := json.Marshal(validatorTree)
			if err != nil {
				return nil, err
			}
			return validatorTreeSlice, nil
		},
	)
	if err != nil {
		return nil, err
	}

	// unmarshal the validator tree
	validatorTree := [][]phase0.Root{}
	err = json.Unmarshal(validatorTreeSlice, &validatorTree)
	if err != nil {
		return nil, err
	}
	return validatorTree, nil
}

func (epp *EigenPodProofs) loadOrComputeBeaconData(prefix string, slot phase0.Slot, getData func() ([]byte, error)) ([]byte, error) {
	// check if the data is cached
	data, found := epp.oracleStateCache.Get(key(prefix, uint64(slot)))
	// if the data is cached, return it
	if found {
		return data, nil
	}

	// compute the data
	data, err := getData()
	if err != nil {
		return nil, err
	}

	// cache the beacon state root
	_ = epp.oracleStateCache.Add(key(prefix, uint64(slot)), data)
	return data, nil
}

func key(prefix string, slot uint64) string {
	return prefix + strconv.FormatUint(slot, 10)
}
