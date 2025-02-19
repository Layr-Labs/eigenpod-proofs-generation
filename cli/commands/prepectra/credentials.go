package prepectra

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/commands"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/prepectra"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	lo "github.com/samber/lo"
)

type TCredentialCommandArgs struct {
	EigenpodAddress string

	DisableColor        bool
	UseJSON             bool
	SimulateTransaction bool
	Node                string
	BeaconNode          string
	Sender              string
	SpecificValidator   uint64
	BatchSize           uint64
	NoPrompt            bool
	Verbose             bool
}

func CredentialsCommand(args TCredentialCommandArgs) error {
	ctx := context.Background()
	if args.DisableColor {
		color.NoColor = true
	}

	isGasEstimate := args.SimulateTransaction && args.Sender != ""
	isVerbose := (!args.UseJSON && !args.SimulateTransaction) || args.Verbose

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, isVerbose)
	utils.PanicOnError("failed to reach ethereum clients", err)

	var specificValidatorIndex *big.Int = nil
	if args.SpecificValidator != math.MaxUint64 && args.SpecificValidator != 0 {
		specificValidatorIndex = new(big.Int).SetUint64(args.SpecificValidator)
		if isVerbose {
			fmt.Printf("Using specific validator: %d", args.SpecificValidator)
		}
	}

	validatorProofs, oracleBeaconTimestamp, err := prepectra.GenerateValidatorProof(ctx, args.EigenpodAddress, eth, chainId, beaconClient, specificValidatorIndex, isVerbose)

	if err != nil || validatorProofs == nil {
		utils.PanicOnError("Failed to generate validator proof", err)
		utils.Panic("no inactive validators")
	}

	if len(args.Sender) != 0 || args.SimulateTransaction {
		txns, indices, err := prepectra.SubmitValidatorProof(ctx, args.Sender, args.EigenpodAddress, chainId, eth, args.BatchSize, validatorProofs, oracleBeaconTimestamp, args.NoPrompt, args.SimulateTransaction, isVerbose)
		utils.PanicOnError(fmt.Sprintf("failed to %s validator proof", func() string {
			if args.SimulateTransaction {
				return "simulate"
			} else {
				return "submit"
			}
		}()), err)

		if args.SimulateTransaction {
			out := lo.Map(txns, func(txn *types.Transaction, _ int) commands.CredentialProofTransaction {
				gas := txn.Gas()
				return commands.CredentialProofTransaction{
					Transaction: commands.Transaction{
						Type:     "credential_proof",
						To:       txn.To().Hex(),
						CallData: common.Bytes2Hex(txn.Data()),
						GasEstimateGwei: func() *uint64 {
							if isGasEstimate {
								return &gas
							}
							return nil
						}(),
					},
					ValidatorIndices: lo.Map(lo.Flatten(indices), func(index *big.Int, _ int) uint64 {
						return index.Uint64()
					}),
				}
			})
			commands.PrintAsJSON(out)
		} else {
			for i, txn := range txns {
				color.Green("transaction(%d): %s", i, txn.Hash().Hex())
			}
		}
	}
	return nil
}
