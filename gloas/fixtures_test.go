package gloasproofs

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Export reproducible Go-generated vectors consumed by Foundry. No hand-authored Merkle siblings.
func TestExportFixtures(t *testing.T) {
	dir := os.Getenv("GLOAS_FIXTURE_DIR")
	if dir == "" {
		t.Skip("set GLOAS_FIXTURE_DIR to regenerate contract vectors")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	decodeEnv := func(name string, size int) []byte {
		t.Helper()
		value := os.Getenv(name)
		if value == "" {
			return nil
		}
		decoded, err := hex.DecodeString(strings.TrimPrefix(value, "0x"))
		if err != nil || len(decoded) != size {
			t.Fatalf("%s must be %d-byte hex", name, size)
		}
		return decoded
	}
	podAddress := decodeEnv("GLOAS_POD_ADDRESS", 20)
	executionHash := decodeEnv("GLOAS_EXECUTION_HASH", 32)
	parentHash := decodeEnv("GLOAS_EXECUTION_PARENT_HASH", 32)
	for _, name := range []string{"initial", "rewarded", "loss", "full", "partial", "empty", "empty_partial", "alias", "recovery", "slashed"} {
		s := fixtureState()
		if podAddress != nil {
			for _, v := range s.Validators {
				copy(v.WithdrawalCredentials[12:], podAddress)
			}
		}
		if executionHash != nil {
			copy(s.LatestExecutionPayloadBid.BlockHash[:], executionHash)
		}
		if parentHash != nil {
			copy(s.LatestExecutionPayloadBid.ParentBlockHash[:], parentHash)
		}
		committed := s.LatestExecutionPayloadBid.BlockHash
		switch name {
		case "full", "empty", "alias", "recovery":
			s.Balances[0] = 0
		case "rewarded":
			s.Balances[0] = 33e9
		case "loss":
			s.Balances[0] = 31e9
		case "slashed":
			s.Validators[0].Slashed = true
		}
		if name == "empty" || name == "empty_partial" {
			s.LatestExecutionPayloadBid.BlockHash = phase0.Hash32{0x33}
			s.LatestExecutionPayloadBid.ParentBlockHash = committed
		}
		if name == "alias" {
			s.LatestExecutionPayloadBid.ParentBlockHash = committed
		}
		if name == "recovery" && executionHash == nil {
			s.LatestExecutionPayloadBid.BlockHash = phase0.Hash32{0x44}
			s.LatestExecutionPayloadBid.ParentBlockHash = phase0.Hash32{0x22}
		}
		g, _ := generator(t, s)
		creds, err := g.Credentials(fixtureIndices)
		if err != nil {
			t.Fatal(err)
		}
		executionHash, executionProof, err := g.CancellationProof()
		if err != nil {
			t.Fatal(err)
		}
		cp, err := g.checkpointProofs(fixtureIndices, executionHash, executionProof)
		if err != nil {
			t.Fatal(err)
		}
		out := struct {
			BlockRoot   Hash         `json:"blockRoot"`
			Credentials *Credentials `json:"credentials"`
			Checkpoint  *Checkpoint  `json:"checkpoint"`
		}{g.BlockRoot(), creds, cp}
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join(dir, name+".json"), append(b, '\n'), 0644); err != nil {
			t.Fatal(err)
		}
	}
}
