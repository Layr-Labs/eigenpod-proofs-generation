package common

//Adapted from https://github.com/ferranbt/fastssz/blob/main/tree.go
import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/minio/sha256-simd"
)

type Proof [][32]byte

var zeroHashes [65][32]byte

func (p Proof) MarshalJSON() ([]byte, error) {
	v := "0x"
	for _, d := range p {
		v += hex.EncodeToString(d[:])
	}
	return json.Marshal(v)
}

func (p *Proof) UnmarshalJSON(data []byte) error {
	// Unmarshal JSON data into a Go string
	var hexString string
	if err := json.Unmarshal(data, &hexString); err != nil {
		return err
	}

	// Remove the "0x" prefix
	if len(hexString) < 2 || hexString[:2] != "0x" {
		return fmt.Errorf("hex string should have '0x' prefix")
	}
	hexString = hexString[2:]

	// Each [32]byte requires 64 hex characters, so calculate how many entries we should have
	numEntries := len(hexString) / 64
	if len(hexString)%64 != 0 {
		return fmt.Errorf("not an even multiple of 32 bytes")
	}

	// Initialize the slice to store our decoded entries
	*p = make([][32]byte, numEntries)

	for i := 0; i < numEntries; i++ {
		start := i * 64
		end := (i + 1) * 64

		// Decode hex string to bytes
		decoded, err := hex.DecodeString(hexString[start:end])
		if err != nil {
			return err
		}

		// Copy into [32]byte array
		copy((*p)[i][:], decoded)
	}

	return nil
}

func (p Proof) ToByteSlice() []byte {
	byteSlice := make([]byte, 0)
	for _, d := range p {
		byteSlice = append(byteSlice, d[:]...)
	}
	return byteSlice
}

func init() {
	tmp := [64]byte{}
	for i := 0; i < 64; i++ {
		copy(tmp[:32], zeroHashes[i][:])
		copy(tmp[32:], zeroHashes[i][:])
		zeroHashes[i+1] = sha256.Sum256(tmp[:])
	}
}

func ComputeMerkleTreeFromLeaves(values []phase0.Root, numLayers uint64) ([][]phase0.Root, error) {
	if len(values) == 0 {
		return nil, errors.New("no values")
	}
	// Initialize the map
	tree := make([][]phase0.Root, numLayers+1)
	tree[0] = values

	for l := 0; l < int(numLayers); l++ {
		if len(tree[l])%2 == 1 {

			zeroHash := phase0.Root(zeroHashes[l])
			tree[l] = append(tree[l], zeroHash)
		}
		nextLevelSize := len(tree[l]) / 2
		values := make([]phase0.Root, nextLevelSize)
		for i := 0; i < len(tree[l]); i += 2 {

			values[i/2] = hashNodes(tree[l][i], tree[l][i+1])

		}
		tree[l+1] = values
	}

	return tree, nil

}

// This proof is from the bottom to the top
func ComputeMerkleProofFromTree(tree [][]phase0.Root, index, numLayers uint64) (Proof, error) {

	var proof [][32]byte
	for l := 0; l < int(numLayers); l++ {
		if len(tree[l]) == 0 {
			return nil, errors.New("no nodes at layer l")
		}
		layerIndex := (index / (uint64(math.Pow(2, float64(l))))) ^ 1
		if layerIndex < uint64(len(tree[l])) {
			proof = append(proof, tree[l][layerIndex])
		} else {
			proof = append(proof, zeroHashes[l])
		}
	}

	return proof, nil
}

func hashNodes(left, right phase0.Root) phase0.Root {
	return phase0.Root(hashFn(append(left[:], right[:]...)))
}

func hashFn(data []byte) [32]byte {
	res := sha256.Sum256(data)
	return res
}

func GetProof(leaves []phase0.Root, index uint64, numLayers uint64) (Proof, error) {

	tree, err := ComputeMerkleTreeFromLeaves(leaves, numLayers)
	if err != nil {
		fmt.Println("error in get proof tree", err)

		return nil, err
	}

	//LogTreeByLevel(tree)

	proof, err := ComputeMerkleProofFromTree(tree, index, numLayers)

	// fmt.Println("PROOF: ", proof)

	if err != nil {
		fmt.Println("error in get proof from tree", err)

		return nil, err
	}

	return proof, nil
}

func LogTreeByLevel(tree [][]phase0.Root) {
	for i := 0; i < len(tree); i++ {
		fmt.Printf("Layer:%d %v \n", i, tree[i])
	}
	fmt.Printf("\n")
}

func ValidateProof(root phase0.Root, proof [][32]byte, element phase0.Root, index uint64) bool {
	target_hash := element
	layer_index := index

	for i := 0; i < len(proof); i++ {
		if layer_index%2 == 0 {
			target_hash = hashNodes(target_hash, phase0.Root(proof[i]))
		} else {
			target_hash = hashNodes(phase0.Root(proof[i]), target_hash)
		}
		layer_index = layer_index / 2
	}

	return target_hash == root
}
