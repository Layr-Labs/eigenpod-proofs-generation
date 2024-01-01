package eigenpodproofs

import (
	"encoding/json"
	"errors"

	"strconv"
	"time"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/hashicorp/golang-lru/v2/expirable"
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

func (epp *EigenPodProofs) ComputeBeaconStateRoot(beaconState *capella.BeaconState) (phase0.Root, error) {
	// check if the beacon state root is cached
	beaconStateRootSlice, found := epp.oracleStateCache.Get(key(BEACON_STATE_ROOT_PREFIX, uint64(beaconState.Slot)))
	// if the beacon state root is cached, return it
	if found {
		var beaconStateRoot phase0.Root
		copy(beaconStateRoot[:], beaconStateRootSlice)
		return beaconStateRoot, nil
	}

	// compute the beacon state root
	beaconStateRoot, err := beaconState.HashTreeRoot()
	if err != nil {
		return phase0.Root{}, err
	}

	// cache the beacon state root
	_ = epp.oracleStateCache.Add(key(BEACON_STATE_ROOT_PREFIX, uint64(beaconState.Slot)), beaconStateRoot[:])
	return beaconStateRoot, nil
}

func (epp *EigenPodProofs) ComputeBeaconStateTopLevelRoots(beaconState *capella.BeaconState) (*BeaconStateTopLevelRoots, error) {
	// check if the beacon state top level roots are cached
	beaconStateTopLevelRootsSlice, found := epp.oracleStateCache.Get(key(BEACON_STATE_TOP_LEVEL_ROOTS_PREFIX, uint64(beaconState.Slot)))
	// if the beacon state top level roots are cached, return them
	if found {
		beaconStatbeaconStateTopLevelRoots := &BeaconStateTopLevelRoots{}
		err := json.Unmarshal(beaconStateTopLevelRootsSlice, beaconStatbeaconStateTopLevelRoots)
		return beaconStatbeaconStateTopLevelRoots, err
	}

	// compute the beacon state top level roots
	beaconStateTopLevelRoots, err := ComputeBeaconStateTopLevelRoots(beaconState)
	if err != nil {
		return nil, err
	}

	// cache the beacon state top level roots
	beaconStateTopLevelRootsSlice, err = json.Marshal(beaconStateTopLevelRoots)
	if err != nil {
		return nil, err
	}
	_ = epp.oracleStateCache.Add(key(BEACON_STATE_TOP_LEVEL_ROOTS_PREFIX, uint64(beaconState.Slot)), beaconStateTopLevelRootsSlice)
	return beaconStateTopLevelRoots, nil
}

func (epp *EigenPodProofs) ComputeValidatorTree(slot uint64, validators []*phase0.Validator) ([][]phase0.Root, error) {
	// check if the validator tree leaves are cached
	validatorTreeSlice, found := epp.oracleStateCache.Get(key(VALIDATOR_TREE_PREFIX, slot))

	// if the validator tree leaves are cached, return them
	if found {
		var validatorTree [][]phase0.Root
		err := json.Unmarshal(validatorTreeSlice, &validatorTree)
		return validatorTree, err
	}

	// compute the validator tree leaves
	validatorLeaves, err := ComputeValidatorTreeLeaves(validators)
	if err != nil {
		return nil, err
	}

	// compute the validator tree
	validatorTree, err := ComputeMerkleTreeFromLeaves(validatorLeaves, validatorListMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}

	// cache the validator tree
	validatorTreeSlice, err = json.Marshal(validatorTree)
	if err != nil {
		return nil, err
	}
	_ = epp.oracleStateCache.Add(key(VALIDATOR_TREE_PREFIX, slot), validatorTreeSlice)
	return validatorTree, nil
}

func key(prefix string, slot uint64) string {
	return prefix + strconv.FormatUint(slot, 10)
}
