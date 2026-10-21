package gloasproofs

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// The saved state/root came from Lighthouse, independently of this generator.
// Keep the authenticated snapshot compressed rather than storing a synthetic replacement.
func TestRealClientProofs(t *testing.T) {
	var meta struct {
		StateSHA256     string                   `json:"stateSHA256"`
		ExpectedRoot    string                   `json:"expectedRoot"`
		ExecutionParent string                   `json:"executionParent"`
		Header          phase0.BeaconBlockHeader `json:"header"`
		Indices         []uint64                 `json:"indices"`
	}
	b, err := os.ReadFile("testdata/lighthouse.json")
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, &meta); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open("testdata/lighthouse-state.ssz.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	b, err = io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(b)
	if hex.EncodeToString(sum[:]) != meta.StateSHA256 {
		t.Fatal("client state checksum mismatch")
	}
	var state gloas.BeaconState
	if err = state.UnmarshalSSZ(b); err != nil {
		t.Fatal(err)
	}
	root := func(s string) Hash {
		t.Helper()
		b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
		if err != nil || len(b) != 32 {
			t.Fatal("invalid snapshot hash")
		}
		return Hash(b)
	}
	g, err := NewGenerator(&meta.Header, &state, root(meta.ExpectedRoot))
	if err != nil {
		t.Fatal(err)
	}
	credentials, err := g.Credentials(meta.Indices)
	if err != nil {
		t.Fatal(err)
	}
	cp, err := g.Checkpoint(meta.Indices, root(meta.ExecutionParent))
	if err != nil {
		t.Fatal(err)
	}
	for n, index := range meta.Indices {
		leaf, err := state.Validators[index].HashTreeRoot()
		if err != nil {
			t.Fatal(err)
		}
		gi, err := validatorGindex(index)
		if err != nil {
			t.Fatal(err)
		}
		verify(t, Hash(meta.Header.StateRoot), Hash(leaf), gi, credentials.ValidatorFieldsProofs[n])
		gi, err = field(index / 4)
		if err != nil {
			t.Fatal(err)
		}
		p := cp.BalanceProofs[n]
		verify(t, cp.BalanceContainerProof.BalanceContainerRoot, p.BalanceRoot, gi, p.Proof)
		if binary.LittleEndian.Uint64(p.BalanceRoot[(index%4)*8:]) != uint64(state.Balances[index]) {
			t.Fatal("real-client balance mismatch")
		}
	}
	if dir := os.Getenv("GLOAS_FIXTURE_DIR"); dir != "" {
		out := struct {
			BlockRoot   Hash         `json:"blockRoot"`
			Credentials *Credentials `json:"credentials"`
			Checkpoint  *Checkpoint  `json:"checkpoint"`
		}{g.BlockRoot(), credentials, cp}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, "real-lighthouse.json"), append(b, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
}
