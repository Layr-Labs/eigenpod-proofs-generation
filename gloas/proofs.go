// Package gloasproofs constructs Gloas-only EigenPod proofs from authenticated beacon state.
package gloasproofs

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	dynssz "github.com/pk910/dynamic-ssz"
	"github.com/pk910/dynamic-ssz/treeproof"
)

type Hash [32]byte
type HexBytes []byte

func (h Hash) MarshalJSON() ([]byte, error)     { return json.Marshal("0x" + hex.EncodeToString(h[:])) }
func (b HexBytes) MarshalJSON() ([]byte, error) { return json.Marshal("0x" + hex.EncodeToString(b)) }

type StateRootProof struct {
	BeaconStateRoot Hash     `json:"beaconStateRoot"`
	Proof           HexBytes `json:"proof"`
}
type Credentials struct {
	StateRootProof        StateRootProof `json:"stateRootProof"`
	ValidatorIndices      []uint64       `json:"validatorIndices"`
	ValidatorFieldsProofs []HexBytes     `json:"validatorFieldsProofs"`
	ValidatorFields       [][]Hash       `json:"validatorFields"`
}
type BalanceContainerProof struct {
	BalanceContainerRoot Hash     `json:"balanceContainerRoot"`
	Proof                HexBytes `json:"proof"`
}
type BalanceProof struct {
	PubkeyHash  Hash     `json:"pubkeyHash"`
	BalanceRoot Hash     `json:"balanceRoot"`
	Proof       HexBytes `json:"proof"`
}
type Checkpoint struct {
	BalanceContainerProof  BalanceContainerProof `json:"balanceContainerProof"`
	BalanceProofs          []BalanceProof        `json:"balanceProofs"`
	CommittedExecutionHash Hash                  `json:"committedExecutionHash"`
	ExecutionHashProof     HexBytes              `json:"executionHashProof"`
}

// Generator owns an authenticated, immutable Merkle tree. Inputs are not retained,
// and returned proofs do not share mutable memory with the generator.
// Construct a fresh Generator for each oracle root; there are no slot-only caches.
type Generator struct {
	tree           *treeproof.Node
	validatorCount uint64
	stateProof     StateRootProof
	blockRoot      Hash
}

// NewGenerator authenticates the header against expectedRoot, then the complete
// Gloas state against that header. Read expectedRoot from the EIP-4788 oracle or
// the pod's saved checkpoint; never derive it from the untrusted input header.
// Inputs must not be mutated during construction, but can be released afterwards.
func NewGenerator(header *phase0.BeaconBlockHeader, state *gloas.BeaconState, expectedRoot Hash) (*Generator, error) {
	if header == nil || state == nil || header.Slot != state.Slot {
		return nil, errors.New("header and state must be present and have the same slot")
	}
	if len(state.Validators) != len(state.Balances) || state.LatestExecutionPayloadBid == nil {
		return nil, errors.New("invalid validator/balance lengths or missing execution bid")
	}
	ds := dynssz.NewDynSsz(nil)
	headerTree, err := ds.GetTree(header)
	if err != nil {
		return nil, fmt.Errorf("beacon header tree: %w", err)
	}
	blockRoot := Hash(headerTree.Hash())
	if blockRoot != expectedRoot {
		return nil, errors.New("header does not match the on-chain expected root")
	}
	tree, err := ds.GetTree(state)
	if err != nil {
		return nil, fmt.Errorf("Gloas state tree: %w", err)
	}
	if !bytes.Equal(tree.Hash(), header.StateRoot[:]) {
		return nil, errors.New("beacon state does not match header state root")
	}
	p, err := headerTree.Prove(11)
	if err != nil {
		return nil, err
	}
	return &Generator{
		tree: tree, validatorCount: uint64(len(state.Validators)), blockRoot: blockRoot,
		stateProof: StateRootProof{BeaconStateRoot: Hash(header.StateRoot), Proof: flatten(p.Hashes)},
	}, nil
}

func (g *Generator) BlockRoot() Hash { return g.blockRoot }

// ProgressiveGindex addresses an EIP-7916 chunk before the list/container mix-in.
// Indices requiring more than a signed 64-bit tree address are rejected by this Go tree backend.
func ProgressiveGindex(index uint64) (uint64, error) {
	gindex, size, depth := uint64(1), uint64(1), uint(0)
	for index >= size {
		index -= size
		if depth >= 38 {
			return 0, errors.New("progressive index exceeds proof backend capacity")
		}
		size <<= 2
		depth += 2
		gindex = gindex<<1 | 1
	}
	if bits.Len64(gindex)+int(depth)+1 > 63 {
		return 0, errors.New("progressive index exceeds proof backend capacity")
	}
	return gindex<<(depth+1) | index, nil
}

func concat(parent, child uint64) (uint64, error) {
	depth := bits.Len64(child) - 1
	if child == 0 || bits.Len64(parent)+depth > 63 {
		return 0, errors.New("generalized index exceeds proof backend capacity")
	}
	return parent<<depth | (child ^ (uint64(1) << depth)), nil
}

func field(index uint64) (uint64, error) {
	p, err := ProgressiveGindex(index)
	if err != nil {
		return 0, err
	}
	return concat(2, p)
}

func validatorGindex(index uint64) (uint64, error) {
	listIndex, err := field(index)
	if err != nil {
		return 0, err
	}
	stateIndex, _ := field(11)
	return concat(stateIndex, listIndex)
}

func (g *Generator) prove(index uint64) (*treeproof.Proof, error) {
	if index > math.MaxInt {
		return nil, errors.New("generalized index exceeds platform capacity")
	}
	return g.tree.Prove(int(index))
}

func (g *Generator) validateIndices(indices []uint64) error {
	seen := make(map[uint64]bool, len(indices))
	for _, i := range indices {
		if i >= 1<<40 || i >= g.validatorCount {
			return fmt.Errorf("validator index %d out of range", i)
		}
		if seen[i] {
			return fmt.Errorf("duplicate validator index %d", i)
		}
		seen[i] = true
	}
	return nil
}

func (g *Generator) Credentials(indices []uint64) (*Credentials, error) {
	if err := g.validateIndices(indices); err != nil {
		return nil, err
	}
	out := &Credentials{
		StateRootProof:        StateRootProof{BeaconStateRoot: g.stateProof.BeaconStateRoot, Proof: bytes.Clone(g.stateProof.Proof)},
		ValidatorIndices:      append([]uint64(nil), indices...),
		ValidatorFields:       make([][]Hash, 0, len(indices)),
		ValidatorFieldsProofs: make([]HexBytes, 0, len(indices)),
	}
	for _, i := range indices {
		index, err := validatorGindex(i)
		if err != nil {
			return nil, err
		}
		p, err := g.prove(index)
		if err != nil {
			return nil, err
		}
		validator, err := g.tree.Get(int(index))
		if err != nil {
			return nil, err
		}
		fields := make([]Hash, 8)
		for n := range fields {
			node, err := validator.Get(8 + n)
			if err != nil {
				return nil, err
			}
			fields[n] = Hash(node.Hash())
		}
		out.ValidatorFields = append(out.ValidatorFields, fields)
		out.ValidatorFieldsProofs = append(out.ValidatorFieldsProofs, flatten(p.Hashes))
	}
	return out, nil
}

// Checkpoint constructs both the balance proofs and the authenticated execution-parent guard.
// expectedExecutionParent is the value saved by EigenPodGloas when startCheckpoint executed.
// On an empty-parent checkpoint use CancellationProof, cancel it, and start a fresh checkpoint.
func (g *Generator) Checkpoint(indices []uint64, expectedExecutionParent Hash) (*Checkpoint, error) {
	if err := g.validateIndices(indices); err != nil {
		return nil, err
	}
	executionHash, executionProof, err := g.CancellationProof()
	if err != nil {
		return nil, err
	}
	// The bid's parent hash is the fourth sibling in its block-hash branch.
	// Read both hashes from the authenticated tree, matching the contract guard.
	parentHash := Hash(executionProof[3*32 : 4*32])
	if expectedExecutionParent == (Hash{}) || executionHash != expectedExecutionParent || executionHash == parentHash {
		return nil, errors.New("checkpoint parent payload was not delivered; use cancellation proof")
	}
	return g.checkpointProofs(indices, executionHash, executionProof)
}

func (g *Generator) checkpointProofs(indices []uint64, executionHash Hash, executionProof HexBytes) (*Checkpoint, error) {
	balancesIndex, _ := field(12)
	p, err := g.prove(balancesIndex)
	if err != nil {
		return nil, err
	}
	balances, err := g.tree.Get(int(balancesIndex))
	if err != nil {
		return nil, err
	}
	fullBalanceProof := append(flatten(p.Hashes), g.stateProof.Proof...)
	out := &Checkpoint{BalanceContainerProof: BalanceContainerProof{BalanceContainerRoot: Hash(p.Leaf), Proof: EncodeProofPair(fullBalanceProof, executionProof)}, CommittedExecutionHash: executionHash, ExecutionHashProof: executionProof, BalanceProofs: make([]BalanceProof, 0, len(indices))}
	for _, i := range indices {
		index, err := field(i / 4)
		if err != nil {
			return nil, err
		}
		p, err := balances.Prove(int(index))
		if err != nil {
			return nil, err
		}
		validatorIndex, err := validatorGindex(i)
		if err != nil {
			return nil, err
		}
		validator, err := g.tree.Get(int(validatorIndex))
		if err != nil {
			return nil, err
		}
		// Validator field 0 is the already-hashed 48-byte public key.
		key, err := validator.Get(8)
		if err != nil {
			return nil, err
		}
		out.BalanceProofs = append(out.BalanceProofs, BalanceProof{PubkeyHash: Hash(key.Hash()), BalanceRoot: Hash(p.Leaf), Proof: flatten(p.Hashes)})
	}
	return out, nil
}

// CancellationProof authenticates the committed hash even for withheld/aliasing bids.
// Cancellation is permitted on chain only if this hash disagrees with the saved EL parent,
// or the bid block hash aliases its parent hash. The proof alone never authorizes cancellation.
func (g *Generator) CancellationProof() (Hash, HexBytes, error) {
	stateIndex, _ := field(43)
	bidIndex, _ := field(2)
	index, err := concat(stateIndex, bidIndex)
	if err != nil {
		return Hash{}, nil, err
	}
	p, err := g.prove(index)
	if err != nil {
		return Hash{}, nil, err
	}
	return Hash(p.Leaf), append(flatten(p.Hashes), g.stateProof.Proof...), nil
}

func flatten(hashes [][]byte) HexBytes {
	out := make(HexBytes, 0, len(hashes)*32)
	for _, h := range hashes {
		out = append(out, h...)
	}
	return out
}

// EncodeProofPair is Solidity abi.encode(bytes balanceProof, bytes executionHashProof).
func EncodeProofPair(first, second []byte) HexBytes {
	firstPadded := (len(first) + 31) / 32 * 32
	secondPadded := (len(second) + 31) / 32 * 32
	out := make(HexBytes, 128+firstPadded+secondPadded)
	binary.BigEndian.PutUint64(out[24:32], 64)
	binary.BigEndian.PutUint64(out[56:64], uint64(96+firstPadded))
	binary.BigEndian.PutUint64(out[88:96], uint64(len(first)))
	copy(out[96:], first)
	start := 96 + firstPadded
	binary.BigEndian.PutUint64(out[start+24:start+32], uint64(len(second)))
	copy(out[start+32:], second)
	return out
}
