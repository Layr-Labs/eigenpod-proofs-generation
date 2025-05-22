package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

type consolidationRequestJSON struct {
	SrcPubkey    string `json:"SrcPubkey"`
	TargetPubkey string `json:"TargetPubkey"`
}

func LoadConsolidationRequestFromFile(path string) ([]EigenPod.IEigenPodTypesConsolidationRequest, error) {
	var raw []consolidationRequestJSON

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &raw); err != nil {
		return nil, err
	}

	var result []EigenPod.IEigenPodTypesConsolidationRequest
	for _, r := range raw {
		src, err := hexStringToBytes(r.SrcPubkey)
		if err != nil {
			return nil, fmt.Errorf("invalid SrcPubkey: %w", err)
		}
		tgt, err := hexStringToBytes(r.TargetPubkey)
		if err != nil {
			return nil, fmt.Errorf("invalid TargetPubkey: %w", err)
		}

		result = append(result, EigenPod.IEigenPodTypesConsolidationRequest{
			SrcPubkey:    src,
			TargetPubkey: tgt,
		})
	}

	return result, nil
}

// SubmitConsolidationRequests submits consolidation requests via EigenPod.requestConsolidation
//
// Parameters:
//   - requests are consolidation requests submitted as a single transaction
//   - predeployFee is the msg.value to send with the transaction. This should be the current predeploy
//     fee multiplied by len(requests).
func SubmitConsolidationRequests(
	ctx context.Context,
	owner, eigenpodAddress string,
	chainId *big.Int,
	eth *ethclient.Client,
	requests []EigenPod.IEigenPodTypesConsolidationRequest,
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
		color.Green("calling EigenPod.requestConsolidation()... [%s]", func() string {
			if ownerAccount.TransactionOptions.NoSend {
				return "simulated"
			} else {
				return "live"
			}
		}())
	}

	return eigenPod.RequestConsolidation(
		ownerAccount.TransactionOptions,
		requests,
	)
}

func hexStringToBytes(s string) ([]byte, error) {
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}
