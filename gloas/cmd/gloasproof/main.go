// gloasproof generates transaction arguments only; it never submits transactions.
package main

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	proofs "github.com/Layr-Labs/eigenpod-proofs-generation/gloas"
	"github.com/attestantio/go-eth2-client/spec/gloas"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseHash(s string) (proofs.Hash, error) {
	var h proofs.Hash
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil || len(b) != 32 {
		return h, errors.New("expected a 32-byte hex hash")
	}
	copy(h[:], b)
	return h, nil
}

func run() error {
	statePath := flag.String("state", "", "Gloas post-block beacon state in SSZ encoding")
	headerPath := flag.String("header", "", "matching BeaconBlockHeader JSON file (unwrapped)")
	format := flag.String("format", "ssz", "state encoding: ssz only; mainnet preset")
	mode := flag.String("mode", "checkpoint", "credentials, checkpoint, or cancel")
	indicesArg := flag.String("indices", "", "comma-separated validator indices")
	rootArg := flag.String("expected-root", "", "required EIP-4788 root or saved checkpoint beaconBlockRoot")
	parentArg := flag.String("execution-parent", "", "checkpointExecutionParentHash saved by the pod")
	flag.Parse()
	if *format != "ssz" {
		return errors.New("only SSZ beacon-state input is supported; use -format ssz")
	}
	root, err := parseHash(*rootArg)
	if err != nil {
		return fmt.Errorf("expected-root: %w", err)
	}
	stateBytes, err := os.ReadFile(*statePath)
	if err != nil {
		return err
	}
	headerBytes, err := os.ReadFile(*headerPath)
	if err != nil {
		return err
	}
	var state gloas.BeaconState
	if err = state.UnmarshalSSZ(stateBytes); err != nil {
		return fmt.Errorf("decode Gloas state: %w", err)
	}
	var header phase0.BeaconBlockHeader
	if err = json.Unmarshal(headerBytes, &header); err != nil {
		return err
	}
	g, err := proofs.NewGenerator(&header, &state, root)
	if err != nil {
		return err
	}
	var indices []uint64
	if *indicesArg != "" {
		for _, s := range strings.Split(*indicesArg, ",") {
			i, e := strconv.ParseUint(strings.TrimSpace(s), 10, 40)
			if e != nil {
				return e
			}
			indices = append(indices, i)
		}
	}
	var out any
	switch *mode {
	case "credentials":
		if len(indices) == 0 {
			return errors.New("credentials requires indices")
		}
		out, err = g.Credentials(indices)
	case "checkpoint":
		if len(indices) == 0 {
			return errors.New("checkpoint requires indices")
		}
		parent, e := parseHash(*parentArg)
		if e != nil {
			return fmt.Errorf("execution-parent: %w", e)
		}
		out, err = g.Checkpoint(indices, parent)
	case "cancel":
		h, p, e := g.CancellationProof()
		err = e
		out = struct {
			CommittedExecutionHash proofs.Hash     `json:"committedExecutionHash"`
			ExecutionHashProof     proofs.HexBytes `json:"executionHashProof"`
		}{h, p}
	default:
		return errors.New("mode must be credentials, checkpoint, or cancel")
	}
	if err != nil {
		return err
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
