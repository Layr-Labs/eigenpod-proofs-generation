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
	MAX_ORACLE_STATE_CACHE_SIZE = 2000000
)

type EigenPodProofs struct {
	chainID                               uint64
	oracleStateRootCache                  *expirable.LRU[uint64, phase0.Root]
	oracleStateTopLevelRootsCache         *expirable.LRU[uint64, *beacon.BeaconStateTopLevelRoots]
	oracleStateValidatorTreeCache         *expirable.LRU[uint64, [][]phase0.Root]
	oracleStateValidatorBalancesTreeCache *expirable.LRU[uint64, [][]phase0.Root]
	oracleStateCacheExpirySeconds         int
}

// NewEigenPodProofs creates a new EigenPodProofs instance.
// chainID is the chain ID of the chain that the EigenPodProofs instance will be used for.
// oracleStateCacheExpirySeconds is the expiry time for the oracle state cache in seconds. After this time caches of beacon state roots, validator trees and validator balances trees will be evicted.
func NewEigenPodProofs(chainID uint64, oracleStateCacheExpirySeconds int) (*EigenPodProofs, error) {
	if chainID != 1 && chainID != 17000 {
		return nil, errors.New("chainID not supported")
	}

	oracleStateRootCache := expirable.NewLRU[uint64, phase0.Root](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)
	oracleStateTopLevelRootsCache := expirable.NewLRU[uint64, *beacon.BeaconStateTopLevelRoots](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)
	oracleStateValidatorTreeCache := expirable.NewLRU[uint64, [][]phase0.Root](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)
	oracleStateValidatorBalancesTreeCache := expirable.NewLRU[uint64, [][]phase0.Root](MAX_ORACLE_STATE_CACHE_SIZE, nil, time.Duration(oracleStateCacheExpirySeconds)*time.Second)

	return &EigenPodProofs{
		chainID:                               chainID,
		oracleStateRootCache:                  oracleStateRootCache,
		oracleStateTopLevelRootsCache:         oracleStateTopLevelRootsCache,
		oracleStateValidatorTreeCache:         oracleStateValidatorTreeCache,
		oracleStateCacheExpirySeconds:         oracleStateCacheExpirySeconds,
		oracleStateValidatorBalancesTreeCache: oracleStateValidatorBalancesTreeCache,
	}, nil
}

func (epp *EigenPodProofs) PrecomputeCache(state *spec.VersionedBeaconState) error {
	slot, err := state.Slot()
	if err != nil {
		return err
	}
	validators, err := state.Validators()
	if err != nil {
		return err
	}

	balances, err := state.ValidatorBalances()
	if err != nil {
		return err
	}

	epp.ComputeBeaconStateRoot(state.Deneb)
	epp.ComputeBeaconStateTopLevelRoots(state)
	epp.ComputeVersionedBeaconStateTopLevelRoots(state)
	epp.ComputeValidatorTree(slot, validators)
	epp.ComputeValidatorBalancesTree(slot, balances)
	return nil
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
	// get the versioned beacon state's slot
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
	default:
		return nil, errors.New("unsupported beacon state version")
	}
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
			validatorTree, err := common.ComputeMerkleTreeFromLeaves(validatorLeaves, beacon.VALIDATOR_TREE_HEIGHT)
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

func (epp *EigenPodProofs) ComputeValidatorBalancesTree(slot phase0.Slot, balances []phase0.Gwei) ([][]phase0.Root, error) {
	validatorBalancesTree, err := epp.loadOrComputeValidatorBalancesTree(
		slot,
		func() ([][]phase0.Root, error) {
			// compute the validator balances tree leaves
			balanceRoots := beacon.ComputeValidatorBalancesTreeLeaves(balances)

			// compute the validator balances tree
			validatorBalancesTree, err := common.ComputeMerkleTreeFromLeaves(balanceRoots, beacon.GetValidatorBalancesProofDepth(len(balances)))
			if err != nil {
				return nil, err
			}

			// cache the validator balances tree
			return validatorBalancesTree, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return validatorBalancesTree, nil
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

func (epp *EigenPodProofs) loadOrComputeValidatorBalancesTree(slot phase0.Slot, getData func() ([][]phase0.Root, error)) ([][]phase0.Root, error) {
	balancesTree, found := epp.oracleStateValidatorBalancesTreeCache.Get(uint64(slot))
	if found {
		return balancesTree, nil
	}

	// compute the data
	balancesTree, err := getData()
	if err != nil {
		return nil, err
	}

	// cache the beacon state root
	epp.oracleStateValidatorBalancesTreeCache.Add(uint64(slot), balancesTree)
	return balancesTree, nil
}
