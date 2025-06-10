package commands

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/fatih/color"
	lo "github.com/samber/lo"
)

type WithdrawalBaseCommandArgs struct {
	EigenpodAddress string

	DisableColor        bool
	UseJSON             bool
	SimulateTransaction bool
	Node                string
	BeaconNode          string
	Sender              string
	BatchSize           uint64
	NoPrompt            bool
	Verbose             bool

	CheckFee              bool
	NoWarn                bool
	FeeOverestimateFactor float64
}

type TRequestFullExitCommandArgs struct {
	WithdrawalBaseCommandArgs

	Validators []uint64
}

type TRequestPartialWithdrawalCommandArgs struct {
	WithdrawalBaseCommandArgs

	Validators []uint64
	AmtsGwei   []uint64
}

func RequestFullExitCommand(args TRequestFullExitCommandArgs) error {
	ctx := context.Background()
	if args.DisableColor {
		color.NoColor = true
	}

	// "verbosity" in this case refers to validator info printouts.
	// As long as we don't have UseJSON enabled, we keep logs enabled.
	// TODO - we should move to a -v vs -vv vs -vvv system
	isVerbose := args.Verbose
	enableLogs := true
	if args.UseJSON {
		isVerbose = false
		enableLogs = false
	}

	if len(args.Validators) == 0 {
		return fmt.Errorf("usage: request-withdrawal full --validators <validatorIndexA>, <validatorIndexB>, ...")
	}

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, enableLogs)
	utils.PanicOnError("failed to reach ethereum clients", err)

	headState, err := utils.GetBeaconHeadState(ctx, beaconClient)
	utils.PanicOnError("failed to fetch beacon chain head state", err)

	eigenpodValidators, err := utils.GetEigenPodValidatorsByIndex(args.EigenpodAddress, headState)
	utils.PanicOnError("failed to fetch validators for eigenpod", err)

	requests := make([]EigenPod.IEigenPodTypesWithdrawalRequest, 0)
	for _, vIndex := range args.Validators {
		if v, exists := eigenpodValidators[vIndex]; !exists {
			return fmt.Errorf("validator index %d is not pointed at this eigenpod", vIndex)
		} else {
			requests = append(requests, EigenPod.IEigenPodTypesWithdrawalRequest{
				Pubkey:     v.PublicKey[:],
				AmountGwei: 0,
			})
		}
	}

	requestChunks := utils.Chunk(requests, args.BatchSize)
	txns := make([]*types.Transaction, 0)

	for i, chunk := range requestChunks {
		feeInfo, err := utils.GetWithdrawalFeeInfoForRequest(eth, chunk, args.FeeOverestimateFactor)
		utils.PanicOnError("error getting withdrawal fee info", err)

		isSimulatedStr := ""
		if args.SimulateTransaction {
			isSimulatedStr = "[SIMULATED]"
		} else {
			isSimulatedStr = "[LIVE]"
		}

		// Prompt the user for consent.
		// We prompt for individual chunks because the predeploy includes exponential fee growth depending
		// on the size of the queue, and we want to make sure the user is aware of the rising fee.
		utils.PanicIfNoConsent(utils.SubmitFullExitRequestConsent(
			len(chunk),
			feeInfo.CurrentQueueSize,
			toPrintableUnits(feeInfo.FeePerRequest),
			toPrintableUnits(feeInfo.TotalFee),
			toPrintableUnits(feeInfo.OverestimateFee),
			isSimulatedStr,
		))

		if isVerbose {
			color.Green("Submitting chunk %d/%d (msg.value: %s)", i+1, len(requestChunks), toPrintableUnits(feeInfo.OverestimateFee))
		}

		txn, err := core.SubmitWithdrawalRequests(
			ctx,
			args.Sender,
			args.EigenpodAddress,
			chainId,
			eth,
			chunk,
			feeInfo.OverestimateFee,
			args.SimulateTransaction,
			isVerbose,
		)

		// If submission fails, print any successful requests before exiting with the error message
		if err != nil {
			if len(txns) != 0 {
				fmt.Println("Error submitting withdrawal request. Printing successful requests:")
				printWithdrawalTxnsAsJSON(txns)
			}

			utils.PanicOnError("withdrawal request submission failed", err)
		} else {
			if isVerbose {
				color.Green("transaction %d/%d succeeded: %s", i+1, len(requestChunks), txn.Hash().Hex())
			}

			txns = append(txns, txn)
		}
	}

	if isVerbose {
		color.Green("All requests succeeded.")
	}

	// If all submissions succeeded, print transactions
	if args.SimulateTransaction {
		printWithdrawalTxnsAsJSON(txns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i+1, txn.Hash().Hex())
		}
	}

	return nil
}

func RequestPartialWithdrawalCommand(args TRequestPartialWithdrawalCommandArgs) error {
	ctx := context.Background()
	if args.DisableColor {
		color.NoColor = true
	}

	// "verbosity" in this case refers to validator info printouts.
	// As long as we don't have UseJSON enabled, we keep logs enabled.
	// TODO - we should move to a -v vs -vv vs -vvv system
	isVerbose := args.Verbose
	enableLogs := true
	if args.UseJSON {
		isVerbose = false
		enableLogs = false
	}

	if len(args.Validators) == 0 || len(args.Validators) != len(args.AmtsGwei) {
		return fmt.Errorf("usage: request-withdrawal partial --validators <validatorIndexA>, <validatorIndexB> --amounts <amtGweiA>, <amtGweiB>")
	}

	for _, amtGwei := range args.AmtsGwei {
		if amtGwei == 0 {
			return fmt.Errorf("input contains full exit request (amtGwei == 0). Aborting; use full-exit for full exits.")
		}
	}

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, enableLogs)
	utils.PanicOnError("failed to reach ethereum clients", err)

	headState, err := utils.GetBeaconHeadState(ctx, beaconClient)
	utils.PanicOnError("failed to fetch beacon chain head state", err)

	eigenpodValidators, err := utils.GetEigenPodValidatorsByIndex(args.EigenpodAddress, headState)
	utils.PanicOnError("failed to fetch validators for eigenpod", err)

	requests := make([]EigenPod.IEigenPodTypesWithdrawalRequest, 0)
	for i, vIndex := range args.Validators {
		if v, exists := eigenpodValidators[vIndex]; !exists {
			return fmt.Errorf("validator index %d is not pointed at this eigenpod", vIndex)
		} else {
			requests = append(requests, EigenPod.IEigenPodTypesWithdrawalRequest{
				Pubkey:     v.PublicKey[:],
				AmountGwei: args.AmtsGwei[i],
			})
		}
	}

	requestChunks := utils.Chunk(requests, args.BatchSize)
	txns := make([]*types.Transaction, 0)

	for i, chunk := range requestChunks {
		feeInfo, err := utils.GetWithdrawalFeeInfoForRequest(eth, chunk, args.FeeOverestimateFactor)
		utils.PanicOnError("error getting withdrawal fee info", err)

		isSimulatedStr := ""
		if args.SimulateTransaction {
			isSimulatedStr = "[SIMULATED]"
		} else {
			isSimulatedStr = "[LIVE]"
		}

		// Prompt the user for consent.
		// We prompt for individual chunks because the predeploy includes exponential fee growth depending
		// on the size of the queue, and we want to make sure the user is aware of the rising fee.
		utils.PanicIfNoConsent(utils.SubmitPartialExitRequestConsent(
			len(chunk),
			feeInfo.CurrentQueueSize,
			toPrintableUnits(feeInfo.FeePerRequest),
			toPrintableUnits(feeInfo.TotalFee),
			toPrintableUnits(feeInfo.OverestimateFee),
			isSimulatedStr,
		))

		if isVerbose {
			color.Green("Submitting chunk %d/%d (msg.value: %s)", i+1, len(requestChunks), toPrintableUnits(feeInfo.OverestimateFee))
		}

		txn, err := core.SubmitWithdrawalRequests(
			ctx,
			args.Sender,
			args.EigenpodAddress,
			chainId,
			eth,
			chunk,
			feeInfo.OverestimateFee,
			args.SimulateTransaction,
			isVerbose,
		)

		// If submission fails, print any successful requests before exiting with the error message
		if err != nil {
			if len(txns) != 0 {
				fmt.Println("Error submitting withdrawal request. Printing successful requests:")
				printWithdrawalTxnsAsJSON(txns)
			}

			utils.PanicOnError("withdrawal request submission failed", err)
		} else {
			if isVerbose {
				color.Green("transaction %d/%d succeeded: %s", i+1, len(requestChunks), txn.Hash().Hex())
			}

			txns = append(txns, txn)
		}
	}

	if isVerbose {
		color.Green("All requests succeeded.")
	}

	// If all submissions succeeded, print transactions
	if args.SimulateTransaction {
		printWithdrawalTxnsAsJSON(txns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i+1, txn.Hash().Hex())
		}
	}

	return nil
}

func printWithdrawalTxnsAsJSON(txns []*types.Transaction) {
	printableTxns := lo.Map(txns, func(txn *types.Transaction, _ int) PredeployRequestTransaction {
		gas := txn.Gas()
		return PredeployRequestTransaction{
			Transaction: Transaction{
				To:              txn.To().Hex(),
				CallData:        common.Bytes2Hex(txn.Data()),
				Type:            "withdrawal_request",
				GasEstimateGwei: &gas,
			},
			Value: txn.Value(),
		}
	})
	PrintAsJSON(printableTxns)
}
