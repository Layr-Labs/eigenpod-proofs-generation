package core

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

type withdrawalRequestJSON struct {
	Pubkey     string `json:"Pubkey"`
	AmountGwei uint64 `json:"AmountGwei"`
}

func LoadWithdrawalRequestFromFile(path string) ([]EigenPod.IEigenPodTypesWithdrawalRequest, error) {
	var raw []withdrawalRequestJSON

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}

	var result []EigenPod.IEigenPodTypesWithdrawalRequest
	for _, r := range raw {
		pubkey, err := hexStringToBytes(r.Pubkey)
		if err != nil {
			return nil, fmt.Errorf("invalid Pubkey: %w", err)
		}

		result = append(result, EigenPod.IEigenPodTypesWithdrawalRequest{
			Pubkey:     pubkey,
			AmountGwei: r.AmountGwei,
		})
	}

	return result, nil
}

// SubmitWithdrawalRequests submits withdrawal requests via EigenPod.requestWithdrawal
//
// Parameters:
//   - requests are withdrawal requests to be submitted as a single transaction
//   - predeployFee is the msg.value to send with the transaction. This should be the current predeploy
//     fee multiplied by len(requests).
func SubmitWithdrawalRequests(
	ctx context.Context,
	owner, eigenpodAddress string,
	chainId *big.Int,
	eth *ethclient.Client,
	requests []EigenPod.IEigenPodTypesWithdrawalRequest,
	predeployFee *big.Int,
	noSend bool,
	verbose bool,
) (*types.Transaction, error) {
	ownerAccount, err := utils.PrepareAccount(&owner, chainId, noSend)
	if err != nil {
		return nil, err
	}
	utils.PanicOnError("failed to parse private key", err)
	ownerAccount.TransactionOptions.Value = predeployFee

	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, err
	}

	if verbose {
		color.Green("calling EigenPod.requestWithdrawal()... [%s]", func() string {
			if ownerAccount.TransactionOptions.NoSend {
				return "simulated"
			} else {
				return "live"
			}
		}())
	}

	return eigenPod.RequestWithdrawal(
		ownerAccount.TransactionOptions,
		requests,
	)
}
