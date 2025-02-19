package commands

import (
	"context"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	lo "github.com/samber/lo"
)

type TCheckpointCommandArgs struct {
	EigenpodAddress     string
	Node                string
	BeaconNode          string
	Sender              string
	DisableColor        bool
	NoPrompt            bool
	SimulateTransaction bool
	BatchSize           uint64
	ForceCheckpoint     bool
	Verbose             bool
}

func CheckpointCommand(args TCheckpointCommandArgs) error {
	ctx := context.Background()

	if args.DisableColor {
		color.NoColor = true
	}

	isGasEstimate := args.SimulateTransaction && args.Sender != ""
	isVerbose := !args.SimulateTransaction || args.Verbose

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, isVerbose)
	utils.PanicOnError("failed to reach ethereum clients", err)

	currentCheckpoint, err := utils.GetCurrentCheckpoint(args.EigenpodAddress, eth)
	utils.PanicOnError("failed to load checkpoint", err)

	eigenpod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	utils.PanicOnError("failed to connect to eigenpod", err)

	if currentCheckpoint == 0 {
		if len(args.Sender) > 0 || args.SimulateTransaction {
			if !args.NoPrompt && !args.SimulateTransaction {
				utils.PanicIfNoConsent(utils.StartCheckpointProofConsent())
			}

			txn, err := utils.StartCheckpoint(ctx, args.EigenpodAddress, args.Sender, chainId, eth, args.ForceCheckpoint, args.SimulateTransaction)
			utils.PanicOnError("failed to start checkpoint", err)

			if !args.SimulateTransaction {
				color.Green("starting checkpoint: %s.. (waiting for txn to be mined)", txn.Hash().Hex())
				bind.WaitMined(ctx, eth, txn)
				color.Green("started checkpoint! txn: %s", txn.Hash().Hex())
			} else {
				gas := txn.Gas()
				PrintAsJSON([]Transaction{
					{
						Type:     "checkpoint_start",
						To:       txn.To().Hex(),
						CallData: common.Bytes2Hex(txn.Data()),
						GasEstimateGwei: func() *uint64 {
							if isGasEstimate {
								return &gas
							}
							return nil
						}(),
					},
				})

				return nil
			}

			newCheckpoint, err := eigenpod.CurrentCheckpointTimestamp(nil)
			utils.PanicOnError("failed to fetch current checkpoint", err)

			currentCheckpoint = newCheckpoint
		} else {
			utils.PanicOnError("no checkpoint active and no private key provided to start one", errors.New("no checkpoint"))
		}
	}

	if isVerbose {
		color.Green("pod has active checkpoint! checkpoint timestamp: %d", currentCheckpoint)
	}

	proof, err := core.GenerateCheckpointProof(ctx, args.EigenpodAddress, eth, chainId, beaconClient, isVerbose)
	utils.PanicOnError("failed to generate checkpoint proof", err)

	txns, err := core.SubmitCheckpointProof(ctx, args.Sender, args.EigenpodAddress, chainId, proof, eth, args.BatchSize, args.NoPrompt, args.SimulateTransaction, args.Verbose)
	if args.SimulateTransaction {
		printableTxns := lo.Map(txns, func(txn *types.Transaction, _ int) Transaction {
			return Transaction{
				To:       txn.To().Hex(),
				CallData: common.Bytes2Hex(txn.Data()),
				Type:     "checkpoint_proof",
			}
		})
		PrintAsJSON(printableTxns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i, txn.Hash().Hex())
		}
	}
	utils.PanicOnError("an error occurred while submitting your checkpoint proofs", err)

	return nil
}
