package gloasproofs

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pk910/dynamic-ssz/treeproof"
)

var fixtureIndices = []uint64{0, 1, 4, 5, 20, 21, 84, 85}

// Mainnet preset dimensions, synthetic contents. These are Merkle/contract test
// vectors, not a claim that a consensus client accepted this beacon state.
func fixtureState() *gloas.BeaconState {
	s := &gloas.BeaconState{
		Slot: 100, Fork: &phase0.Fork{}, LatestBlockHeader: &phase0.BeaconBlockHeader{}, ETH1Data: &phase0.ETH1Data{BlockHash: make([]byte, 32)},
		BlockRoots: make([]phase0.Root, 8192), StateRoots: make([]phase0.Root, 8192), RANDAOMixes: make([]phase0.Root, 65536), Slashings: make([]phase0.Gwei, 8192),
		PreviousJustifiedCheckpoint: &phase0.Checkpoint{}, CurrentJustifiedCheckpoint: &phase0.Checkpoint{}, FinalizedCheckpoint: &phase0.Checkpoint{}, JustificationBits: []byte{0},
		CurrentSyncCommittee: &altair.SyncCommittee{Pubkeys: make([]phase0.BLSPubKey, 512)}, NextSyncCommittee: &altair.SyncCommittee{Pubkeys: make([]phase0.BLSPubKey, 512)},
		ProposerLookahead: make([]phase0.ValidatorIndex, 64), ExecutionPayloadAvailability: make([]uint8, 1024), BuilderPendingPayments: make([]*gloas.BuilderPendingPayment, 64),
		LatestExecutionPayloadBid: &gloas.ExecutionPayloadBid{BlockHash: phase0.Hash32{0x22}, ParentBlockHash: phase0.Hash32{0x11}},
		PTCWindow:                 make([][]phase0.ValidatorIndex, 96),
	}
	for i := range s.BuilderPendingPayments {
		s.BuilderPendingPayments[i] = &gloas.BuilderPendingPayment{Withdrawal: &gloas.BuilderPendingWithdrawal{}}
	}
	for i := range s.PTCWindow {
		s.PTCWindow[i] = make([]phase0.ValidatorIndex, 512)
	}
	for i := 0; i < 86; i++ {
		var pubkey phase0.BLSPubKey
		binary.LittleEndian.PutUint64(pubkey[:8], uint64(i+1))
		credentials := make([]byte, 32)
		credentials[0] = 1
		credentials[30] = 0xBE
		credentials[31] = 0xEF
		if i%2 == 1 {
			credentials[0] = 2
		}
		s.Validators = append(s.Validators, &phase0.Validator{PublicKey: pubkey, WithdrawalCredentials: credentials, EffectiveBalance: 32e9, ExitEpoch: math.MaxUint64, WithdrawableEpoch: math.MaxUint64})
		s.Balances = append(s.Balances, 32e9)
	}
	return s
}

func generator(t *testing.T, s *gloas.BeaconState) (*Generator, *phase0.BeaconBlockHeader) {
	t.Helper()
	root, err := s.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	h := &phase0.BeaconBlockHeader{Slot: s.Slot, StateRoot: root}
	blockRoot, err := h.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	g, err := NewGenerator(h, s, Hash(blockRoot))
	if err != nil {
		t.Fatal(err)
	}
	return g, h
}

func verify(t *testing.T, root Hash, leaf Hash, index uint64, p []byte) {
	t.Helper()
	hashes := make([][]byte, len(p)/32)
	for i := range hashes {
		hashes[i] = p[i*32 : (i+1)*32]
	}
	ok, err := treeproof.VerifyProof(root[:], &treeproof.Proof{Index: int(index), Leaf: leaf[:], Hashes: hashes})
	if err != nil || !ok {
		t.Fatalf("proof %d invalid: %v", index, err)
	}
}

func TestProofsAndSSZRoundTrip(t *testing.T) {
	s := fixtureState()
	g, h := generator(t, s)
	c, err := g.Credentials(fixtureIndices)
	if err != nil {
		t.Fatal(err)
	}
	verify(t, g.BlockRoot(), c.StateRootProof.BeaconStateRoot, 11, c.StateRootProof.Proof)
	sf, _ := field(11)
	for n, i := range fixtureIndices {
		li, _ := field(i)
		index, _ := concat(sf, li)
		root, _ := s.Validators[i].HashTreeRoot()
		verify(t, Hash(h.StateRoot), Hash(root), index, c.ValidatorFieldsProofs[n])
	}
	cp, err := g.Checkpoint(fixtureIndices, Hash{0x22})
	if err != nil {
		t.Fatal(err)
	}
	bf, _ := field(12)
	stateP, err := g.prove(bf)
	if err != nil {
		t.Fatal(err)
	}
	bi, _ := concat(11, bf)
	verify(t, g.BlockRoot(), cp.BalanceContainerProof.BalanceContainerRoot, bi, append(flatten(stateP.Hashes), g.stateProof.Proof...))
	for n, i := range fixtureIndices {
		index, _ := field(i / 4)
		p := cp.BalanceProofs[n]
		verify(t, cp.BalanceContainerProof.BalanceContainerRoot, p.BalanceRoot, index, p.Proof)
		if binary.LittleEndian.Uint64(p.BalanceRoot[(i%4)*8:]) != 32e9 {
			t.Fatal("balance packing")
		}
	}
	f43, _ := field(43)
	f2, _ := field(2)
	ei, _ := concat(f43, f2)
	ei, _ = concat(11, ei)
	verify(t, g.BlockRoot(), cp.CommittedExecutionHash, ei, cp.ExecutionHashProof)
	if !bytes.Equal(cp.ExecutionHashProof[3*32:4*32], s.LatestExecutionPayloadBid.ParentBlockHash[:]) {
		t.Fatal("wrong parent sibling")
	}
	encoded, err := s.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}
	var round gloas.BeaconState
	if err = round.UnmarshalSSZ(encoded); err != nil {
		t.Fatal(err)
	}
	g2, _ := generator(t, &round)
	if g2.BlockRoot() != g.BlockRoot() {
		t.Fatal("SSZ round trip")
	}
}

func TestRejectUnsettledAndInvalidInputs(t *testing.T) {
	s := fixtureState()
	g, h := generator(t, s)
	for _, indices := range [][]uint64{{86}, {1, 1}, {1 << 40}} {
		if _, err := g.Credentials(indices); err == nil {
			t.Fatalf("accepted %v", indices)
		}
	}
	if _, err := g.Checkpoint([]uint64{0}, Hash{0x33}); err == nil {
		t.Fatal("accepted empty parent")
	}
	if _, err := g.Checkpoint([]uint64{0}, Hash{}); err == nil {
		t.Fatal("accepted missing parent")
	}
	s.LatestExecutionPayloadBid.ParentBlockHash = s.LatestExecutionPayloadBid.BlockHash
	g, _ = generator(t, s)
	if _, err := g.Checkpoint([]uint64{0}, Hash{0x22}); err == nil {
		t.Fatal("accepted alias")
	}
	h.StateRoot[0] ^= 1
	if _, err := NewGenerator(h, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted wrong root")
	}
	s.Slot++
	if _, err := NewGenerator(h, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted wrong slot")
	}
	if _, err := NewGenerator(nil, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted missing header")
	}
	if _, err := NewGenerator(h, nil, g.BlockRoot()); err == nil {
		t.Fatal("accepted missing state")
	}
	s = fixtureState()
	s.Balances = s.Balances[:1]
	h.Slot = s.Slot
	if _, err := NewGenerator(h, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted mismatched validator/balance lengths")
	}
}

func TestGeneratorRequiresTrustedRoot(t *testing.T) {
	s := fixtureState()
	g, h := generator(t, s)
	for _, root := range []Hash{{}, {0xFF}} {
		if _, err := NewGenerator(h, s, root); err == nil {
			t.Fatal("accepted an unauthenticated header")
		}
	}
	s.Balances[0]--
	if _, err := NewGenerator(h, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted an unauthenticated state")
	}
	s.Validators[0] = nil
	if _, err := NewGenerator(h, s, g.BlockRoot()); err == nil {
		t.Fatal("accepted a nil validator")
	}
}

func TestGeneratorOwnsInputsAndOutputs(t *testing.T) {
	s := fixtureState()
	g, h := generator(t, s)
	credentials, err := g.Credentials(fixtureIndices)
	if err != nil {
		t.Fatal(err)
	}
	cp, err := g.Checkpoint(fixtureIndices, Hash{0x22})
	if err != nil {
		t.Fatal(err)
	}
	expectedCredentials, _ := json.Marshal(credentials)
	expectedCheckpoint, _ := json.Marshal(cp)
	expectedRoot := g.BlockRoot()
	// All output fields must come from the authenticated snapshot, including
	// pubkey hashes and the delivery guard. Caller mutations cannot poison it.
	h.StateRoot[0] ^= 1
	s.Validators[0].PublicKey[0] ^= 1
	s.Validators[0].WithdrawalCredentials[0] ^= 1
	s.Validators[0].EffectiveBalance--
	s.Validators[0] = nil
	s.Validators = nil
	s.Balances = nil
	s.LatestExecutionPayloadBid.BlockHash = phase0.Hash32{0x33}
	s.LatestExecutionPayloadBid.ParentBlockHash = phase0.Hash32{0x33}
	credentials.StateRootProof.Proof[0] ^= 1
	credentials.ValidatorIndices[0]++
	credentials.ValidatorFields[0][0][0] ^= 1
	credentials.ValidatorFieldsProofs[0][0] ^= 1
	cp.BalanceContainerProof.Proof[0] ^= 1
	cp.BalanceProofs[0].Proof[0] ^= 1
	cp.ExecutionHashProof[0] ^= 1
	credentials, err = g.Credentials(fixtureIndices)
	if err != nil {
		t.Fatal(err)
	}
	cp, err = g.Checkpoint(fixtureIndices, Hash{0x22})
	if err != nil {
		t.Fatal(err)
	}
	actualCredentials, _ := json.Marshal(credentials)
	actualCheckpoint, _ := json.Marshal(cp)
	if !bytes.Equal(actualCredentials, expectedCredentials) || !bytes.Equal(actualCheckpoint, expectedCheckpoint) || g.BlockRoot() != expectedRoot {
		t.Fatal("caller mutation changed the authenticated generator snapshot")
	}
	if _, err := g.Checkpoint(fixtureIndices, Hash{0x33}); err == nil {
		t.Fatal("caller mutation bypassed the execution-parent guard")
	}
}

func TestProgressiveBoundaries(t *testing.T) {
	for i, want := range map[uint64]uint64{0: 2, 1: 24, 4: 27, 5: 224, 20: 239, 21: 1920, 84: 1983} {
		got, err := ProgressiveGindex(i)
		if err != nil || got != want {
			t.Fatalf("%d: got %d want %d (%v)", i, got, want, err)
		}
	}
	if _, err := ProgressiveGindex(math.MaxUint64); err == nil {
		t.Fatal("overflow")
	}
}

func TestABIEncoding(t *testing.T) {
	for _, n := range []int{0, 1, 31, 32, 33, 352} {
		a := bytes.Repeat([]byte{0x12}, n)
		b := bytes.Repeat([]byte{0x34}, n+7)
		out := EncodeProofPair(a, b)
		offset := int(binary.BigEndian.Uint64(out[56:64]))
		if offset != 96+(n+31)/32*32 || !bytes.Equal(out[96:96+n], a) || !bytes.Equal(out[offset+32:offset+32+len(b)], b) {
			t.Fatal("ABI encoding")
		}
	}
}

func TestCLI(t *testing.T) {
	dir := t.TempDir()
	binaryPath := filepath.Join(dir, "gloasproof")
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binaryPath, "./cmd/gloasproof")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, out)
	}
	s := fixtureState()
	g, h := generator(t, s)
	stateSSZ, err := s.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}
	stateJSON, err := s.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	headerJSON, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	for name, data := range map[string][]byte{"state.ssz": stateSSZ, "state.json": stateJSON, "header.json": headerJSON} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	rootJSON, _ := json.Marshal(g.BlockRoot())
	root := strings.Trim(string(rootJSON), "\"")
	parentJSON, _ := json.Marshal(Hash{0x22})
	parent := strings.Trim(string(parentJSON), "\"")
	base := []string{"-state", filepath.Join(dir, "state.ssz"), "-header", filepath.Join(dir, "header.json"), "-expected-root", root, "-execution-parent", parent, "-indices", "0,1,4,5,20,21,84,85"}
	run := func(args []string, want any, wantError string) {
		t.Helper()
		out, err := exec.Command(binaryPath, args...).CombinedOutput()
		if wantError != "" {
			if err == nil || !strings.Contains(string(out), wantError) {
				t.Fatalf("wanted %q, got %s (%v)", wantError, out, err)
			}
			return
		}
		if err != nil {
			t.Fatalf("CLI: %v\n%s", err, out)
		}
		var actual, expected any
		if err = json.Unmarshal(out, &actual); err != nil {
			t.Fatal(err)
		}
		encoded, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(encoded, &expected); err != nil {
			t.Fatal(err)
		}
		a, _ := json.Marshal(actual)
		b, _ := json.Marshal(expected)
		if !bytes.Equal(a, b) {
			t.Fatal("CLI output differs from library")
		}
	}
	credentials, err := g.Credentials(fixtureIndices)
	if err != nil {
		t.Fatal(err)
	}
	run(append(append([]string(nil), base...), "-mode", "credentials"), credentials, "")
	cp, err := g.Checkpoint(fixtureIndices, Hash{0x22})
	if err != nil {
		t.Fatal(err)
	}
	run(base, cp, "")
	run(append(append([]string(nil), base...), "-format", "ssz"), cp, "")
	jsonArgs := append([]string(nil), base...)
	jsonArgs[1] = filepath.Join(dir, "state.json")
	run(append(jsonArgs, "-format", "json"), nil, "only SSZ beacon-state input is supported")
	run(jsonArgs, nil, "decode Gloas state")
	run(append(append([]string(nil), base...), "-format", "yaml"), nil, "only SSZ beacon-state input is supported")
	wrongRoot := append([]string(nil), base...)
	wrongRoot[5] = "0x" + strings.Repeat("00", 32)
	run(wrongRoot, nil, "expected root")
	// Authentication must reject the source state itself, not only a wrong CLI root flag.
	s.Balances[0] = 0
	changed, err := s.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "state.ssz"), changed, 0600); err != nil {
		t.Fatal(err)
	}
	run(base, nil, "state does not match")
	s.LatestExecutionPayloadBid.BlockHash = phase0.Hash32{0x33}
	s.LatestExecutionPayloadBid.ParentBlockHash = phase0.Hash32{0x22}
	g, h = generator(t, s)
	changed, err = s.MarshalSSZ()
	if err != nil {
		t.Fatal(err)
	}
	headerJSON, err = json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "state.ssz"), changed, 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "header.json"), headerJSON, 0600); err != nil {
		t.Fatal(err)
	}
	rootJSON, _ = json.Marshal(g.BlockRoot())
	base[5] = strings.Trim(string(rootJSON), "\"")
	run(base, nil, "payload was not delivered")
	hash, proof, err := g.CancellationProof()
	if err != nil {
		t.Fatal(err)
	}
	run(append(base, "-mode", "cancel"), struct {
		Hash  Hash     `json:"committedExecutionHash"`
		Proof HexBytes `json:"executionHashProof"`
	}{hash, proof}, "")
}
