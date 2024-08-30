package commands

import (
	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	"github.com/pkg/errors"
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

	eth, beaconClient, chainId, err := core.GetClients(ctx, args.Node, args.BeaconNode, isVerbose)
	core.PanicOnError("failed to reach ethereum clients", err)

	currentCheckpoint, err := core.GetCurrentCheckpoint(args.EigenpodAddress, eth)
	core.PanicOnError("failed to load checkpoint", err)

	eigenpod, err := onchain.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	core.PanicOnError("failed to connect to eigenpod", err)

	if currentCheckpoint == 0 {
		if len(args.Sender) > 0 || args.SimulateTransaction {
			if !args.NoPrompt && !args.SimulateTransaction {
				core.PanicIfNoConsent(core.StartCheckpointProofConsent())
			}

			txn, err := core.StartCheckpoint(ctx, args.EigenpodAddress, args.Sender, chainId, eth, args.ForceCheckpoint, args.SimulateTransaction)
			core.PanicOnError("failed to start checkpoint", err)

			if !args.SimulateTransaction {
				color.Green("starting checkpoint: %s.. (waiting for txn to be mined)", txn.Hash().Hex())
				bind.WaitMined(ctx, eth, txn)
				color.Green("started checkpoint! txn: %s", txn.Hash().Hex())
			} else {
				gas := txn.Gas()
				printAsJSON([]Transaction{
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
			core.PanicOnError("failed to fetch current checkpoint", err)

			currentCheckpoint = newCheckpoint
		} else {
			core.PanicOnError("no checkpoint active and no private key provided to start one", errors.New("no checkpoint"))
		}
	}

	if isVerbose {
		color.Green("pod has active checkpoint! checkpoint timestamp: %d", currentCheckpoint)
	}

	proof, err := core.GenerateCheckpointProof(ctx, args.EigenpodAddress, eth, chainId, beaconClient, isVerbose)
	core.PanicOnError("failed to generate checkpoint proof", err)

	txns, err := core.SubmitCheckpointProof(ctx, args.Sender, args.EigenpodAddress, chainId, proof, eth, args.BatchSize, args.NoPrompt, args.SimulateTransaction, args.Verbose)
	if args.SimulateTransaction {
		printableTxns := utils.Map(txns, func(txn *types.Transaction, _ uint64) Transaction {
			return Transaction{
				To:       txn.To().Hex(),
				CallData: common.Bytes2Hex(txn.Data()),
				Type:     "checkpoint_proof",
			}
		})
		printAsJSON(printableTxns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i, txn.Hash().Hex())
		}
	}
	core.PanicOnError("an error occurred while submitting your checkpoint proofs", err)

	return nil
}
