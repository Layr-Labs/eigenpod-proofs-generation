package commands

import (
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
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
		out = append(out, arr...)
	}
	return out
}

func shortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}
