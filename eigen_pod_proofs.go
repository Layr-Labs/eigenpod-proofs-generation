package eigenpodproofs

import (
	"errors"
	"time"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	expirable "github.com/hashicorp/golang-lru/v2/expirable"

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
	oracleStateRootCache          *expirable.LRU[uint64, phase0.Root]
	oracleStateTopLevelRootsCache *expirable.LRU[uint64, *beacon.BeaconStateTopLevelRoots]
	oracleStateValidatorTreeCache *expirable.LRU[uint64, [][]phase0.Root]
	oracleStateCacheExpirySeconds int
}

func NewEigenPodProofs(chainID uint64, oracleStateCacheExpirySeconds int) (*EigenPodProofs, error) {
	if chainID != 1 && chainID != 5 && chainID != 17000 {
		return nil, errors.New("chainID not supported")
	}

	oracleStateRootCache := expirable.NewLRU[uint64, phase0.Root](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)
	oracleStateTopLevelRootsCache := expirable.NewLRU[uint64, *beacon.BeaconStateTopLevelRoots](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)
	oracleStateValidatorTreeCache := expirable.NewLRU[uint64, [][]phase0.Root](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)

	return &EigenPodProofs{
		chainID:                       chainID,
		oracleStateRootCache:          oracleStateRootCache,
		oracleStateTopLevelRootsCache: oracleStateTopLevelRootsCache,
		oracleStateValidatorTreeCache: oracleStateValidatorTreeCache,
		oracleStateCacheExpirySeconds: oracleStateCacheExpirySeconds,
	}, nil
}

func (epp *EigenPodProofs) ComputeBeaconStateRoot(beaconState *deneb.BeaconState) (phase0.Root, error) {
	beaconStateRoot, err := epp.loadOrComputeBeaconStateRoot(
		beaconState.Slot,
		func() (phase0.Root, error) {
			stateRoot, err := beaconState.HashTreeRoot()
			if err != nil {
				return phase0.Root{}, err
			}
			return stateRoot, nil
		},
	)
	if err != nil {
		return phase0.Root{}, err
	}

	return beaconStateRoot, nil
}

func (epp *EigenPodProofs) ComputeBeaconStateTopLevelRoots(beaconState *spec.VersionedBeaconState) (*beacon.BeaconStateTopLevelRoots, error) {
	//get the versioned beacon state's slot
	slot, err := beaconState.Slot()
	if err != nil {
		return nil, err
	}

	beaconStateTopLevelRoots, err := epp.loadOrComputeBeaconStateTopLevelRoots(
		slot,
		func() (*beacon.BeaconStateTopLevelRoots, error) {
			beaconStateTopLevelRoots, err := epp.ComputeVersionedBeaconStateTopLevelRoots(beaconState)
			if err != nil {
				return nil, err
			}
			return beaconStateTopLevelRoots, nil
		},
	)
	if err != nil {
		return nil, err
	}
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

func (epp *EigenPodProofs) ComputeValidatorTreeCustomComputer(slot phase0.Slot, computer func() ([][]phase0.Root, error)) ([][]phase0.Root, error) {
	validatorTree, err := epp.loadOrComputeValidatorTree(
		slot,
		computer,
	)
	if err != nil {
		return nil, err
	}

	return validatorTree, nil
}

func (epp *EigenPodProofs) ComputeValidatorTree(slot phase0.Slot, validators []*phase0.Validator) ([][]phase0.Root, error) {
	validatorTree, err := epp.loadOrComputeValidatorTree(
		slot,
		func() ([][]phase0.Root, error) {
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
			return validatorTree, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return validatorTree, nil
}

func (epp *EigenPodProofs) loadOrComputeBeaconStateRoot(slot phase0.Slot, getData func() (phase0.Root, error)) (phase0.Root, error) {
	root, found := epp.oracleStateRootCache.Get(uint64(slot))
	if found {
		return root, nil
	}

	// compute the data
	root, err := getData()
	if err != nil {
		return phase0.Root{}, err
	}

	// cache the beacon state root
	epp.oracleStateRootCache.Add(uint64(slot), root)
	return root, nil
}

func (epp *EigenPodProofs) loadOrComputeBeaconStateTopLevelRoots(slot phase0.Slot, getData func() (*beacon.BeaconStateTopLevelRoots, error)) (*beacon.BeaconStateTopLevelRoots, error) {
	topLevelRoots, found := epp.oracleStateTopLevelRootsCache.Get(uint64(slot))
	if found {
		return topLevelRoots, nil
	}

	// compute the data
	topLevelRoots, err := getData()
	if err != nil {
		return nil, err
	}

	// cache the beacon state root
	epp.oracleStateTopLevelRootsCache.Add(uint64(slot), topLevelRoots)
	return topLevelRoots, nil
}

func (epp *EigenPodProofs) loadOrComputeValidatorTree(slot phase0.Slot, getData func() ([][]phase0.Root, error)) ([][]phase0.Root, error) {
	validatorTree, found := epp.oracleStateValidatorTreeCache.Get(uint64(slot))
	if found {
		return validatorTree, nil
	}

	// compute the data
	validatorTree, err := getData()
	if err != nil {
		return nil, err
	}

	// cache the beacon state root
	epp.oracleStateValidatorTreeCache.Add(uint64(slot), validatorTree)
	return validatorTree, nil
}
