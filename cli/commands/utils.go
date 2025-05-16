package commands

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
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

type PredeployRequestTransaction struct {
	Transaction
	Value *big.Int
}

func PrintAsJSON(txns any) {
	out, err := json.MarshalIndent(txns, " ", "   ")
	utils.PanicOnError("failed to serialize", err)
	fmt.Println(string(out))
}
