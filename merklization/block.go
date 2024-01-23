package merklization

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

func ProveBlockRootAgainstBlockRootsList(blockRoots []phase0.Root, blockRootIndex uint64) (Proof, error) {
	proof, err := GetProof(blockRoots, blockRootIndex, blockRootsMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}
	return proof, nil
}

func ProveBlockBodyAgainstBlockHeader(blockHeader *phase0.BeaconBlockHeader) (Proof, error) {
	blockHeaderContainerRoots, err := GetBlockHeaderFieldRoots(blockHeader)

	if err != nil {
		return nil, err
	}

	return GetProof(blockHeaderContainerRoots, beaconBlockBodyRootIndex, blockHeaderMerkleSubtreeNumLayers)
}

func ProveSlotAgainstBlockHeader(blockHeader *phase0.BeaconBlockHeader) (Proof, error) {
	blockHeaderContainerRoots, err := GetBlockHeaderFieldRoots(blockHeader)
	if err != nil {
		return nil, err
	}

	return GetProof(blockHeaderContainerRoots, slotIndex, blockHeaderMerkleSubtreeNumLayers)
}

func GetBlockHeaderFieldRoots(blockHeader *phase0.BeaconBlockHeader) ([]phase0.Root, error) {
	blockHeaderContainerRoots := make([]phase0.Root, beaconBlockHeaderNumFields)

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
