package beacon

import (
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

// refer to this: https://github.com/attestantio/go-eth2-client/blob/654ac05b4c534d96562329f988655e49e5743ff5/spec/phase0/beaconblockheader_encoding.go
func ProveStateRootAgainstBlockHeader(b *phase0.BeaconBlockHeader) (common.Proof, error) {
	beaconBlockHeaderContainerRoots, err := GetBlockHeaderFieldRoots(b)
	if err != nil {
		return nil, err
	}

	return common.GetProof(beaconBlockHeaderContainerRoots, STATE_ROOT_INDEX, BEACON_BLOCK_HEADER_TREE_HEIGHT)
}

func GetBlockHeaderFieldRoots(blockHeader *phase0.BeaconBlockHeader) ([]phase0.Root, error) {
	blockHeaderContainerRoots := make([]phase0.Root, BEACON_BLOCK_HEADER_NUM_FIELDS)

	hh := ssz.NewHasher()

	hh.PutUint64(uint64(blockHeader.Slot))
	copy(blockHeaderContainerRoots[0][:], hh.Hash())
	hh.Reset()

	hh.PutUint64(uint64(blockHeader.ProposerIndex))
	copy(blockHeaderContainerRoots[1][:], hh.Hash())
	hh.Reset()

	hh.PutBytes(blockHeader.ParentRoot[:])
	copy(blockHeaderContainerRoots[2][:], hh.Hash())
	hh.Reset()

	hh.PutBytes(blockHeader.StateRoot[:])
	copy(blockHeaderContainerRoots[3][:], hh.Hash())
	hh.Reset()

	hh.PutBytes(blockHeader.BodyRoot[:])
	copy(blockHeaderContainerRoots[4][:], hh.Hash())
	hh.Reset()

	return blockHeaderContainerRoots, nil
}
