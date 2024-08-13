package main

import (
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	cli "github.com/urfave/cli/v2"
)

type Transaction struct {
	Type     string `json:"type"`
	To       string `json:"to"`
	CallData string `json:"calldata"`
}
type TransactionList = []Transaction

type CredentialProofTransaction struct {
	Transaction
	ValidatorIndices []uint64 `json:"validator_indices"`
}

func printProofs(txns any) {
	out, err := json.Marshal(txns)
	core.PanicOnError("failed to serialize proofs", err)
	fmt.Println(string(out))
}

// imagine if golang had a standard library
func aMap[A any, B any](coll []A, mapper func(i A) B) []B {
	out := make([]B, len(coll))
	for i, item := range coll {
		out[i] = mapper(item)
	}
	return out
}

func aFlatten[A any](coll [][]A) []A {
	out := []A{}
	for _, arr := range coll {
		for _, item := range arr {
			out = append(out, item)
		}
	}
	return out
}

func shortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}

// shared flag --batch
func BatchBySize(destination *uint64, defaultValue uint64) *cli.Uint64Flag {
	return &cli.Uint64Flag{
		Name:        "batch",
		Value:       defaultValue,
		Usage:       "Submit proofs in groups of size `batchSize`, to avoid gas limit.",
		Required:    false,
		Destination: destination,
	}
}

// Hack to make a copy of a flag that sets `Required` to true
func Require(flag *cli.StringFlag) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        flag.Name,
		Aliases:     flag.Aliases,
		Value:       flag.Value,
		Usage:       flag.Usage,
		Destination: flag.Destination,
		Required:    true,
	}
}
