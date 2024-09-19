package commands

import (
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
)

type Transaction struct {
	Type            string  `json:"type"`
	To              string  `json:"to"`
	CallData        string  `json:"calldata"`
	GasEstimateGwei *uint64 `json:"gas_estimate_gwei,omitempty"`
}
type TransactionList = []Transaction

type CredentialProofTransaction struct {
	Transaction
	ValidatorIndices []uint64 `json:"validator_indices"`
}

func printAsJSON(txns any) {
	out, err := json.MarshalIndent(txns, " ", "   ")
	core.PanicOnError("failed to serialize", err)
	fmt.Println(string(out))
}

func getKeys[A comparable, B any](data map[A]B) []A {
	values := make([]A, len(data))
	i := 0
	for key, _ := range data {
		values[i] = key
		i++
	}
	return values
}

func getValues[A comparable, B any](data map[A]B) []B {
	values := make([]B, len(data))
	i := 0
	for _, value := range data {
		values[i] = value
		i++
	}
	return values
}
