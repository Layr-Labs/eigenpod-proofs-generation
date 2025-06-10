package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/fatih/color"
	lo "github.com/samber/lo"
)

type ConsolidateBaseCommandArgs struct {
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

type TConsolidateSwitchCommandArgs struct {
	ConsolidateBaseCommandArgs

	Validators []uint64
}

type TConsolidateToTargetCommandArgs struct {
	ConsolidateBaseCommandArgs

	TargetValidator  uint64
	SourceValidators []uint64
}

func ConsolidateSwitchCommand(args TConsolidateSwitchCommandArgs) error {
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
		return fmt.Errorf("usage: consolidate switch --validators <validatorIndexA>, <validatorIndexB>, ...")
	}

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, enableLogs)
	utils.PanicOnError("failed to reach ethereum clients", err)

	headState, err := utils.GetBeaconHeadState(ctx, beaconClient)
	utils.PanicOnError("failed to fetch beacon chain head state", err)

	eigenpodValidators, err := utils.GetEigenPodValidatorsByIndex(args.EigenpodAddress, headState)
	utils.PanicOnError("failed to fetch validators for eigenpod", err)

	// Form requests, filtering duplicates and validators that aren't pointed at the pod
	requests := make([]EigenPod.IEigenPodTypesConsolidationRequest, 0)
	seen := make(map[uint64]bool)
	for _, vIndex := range args.Validators {
		if v, exists := eigenpodValidators[vIndex]; !exists {
			return fmt.Errorf("validator index %d is not pointed at this eigenpod", vIndex)
		} else if seen[vIndex] {
			return fmt.Errorf("validator index %d is included twice in input args", vIndex)
		} else {
			seen[vIndex] = true

			requests = append(requests, EigenPod.IEigenPodTypesConsolidationRequest{
				SrcPubkey:    v.PublicKey[:],
				TargetPubkey: v.PublicKey[:],
			})
		}
	}

	requestChunks := utils.Chunk(requests, args.BatchSize)
	txns := make([]*types.Transaction, 0)

	for i, chunk := range requestChunks {
		feeInfo, err := utils.GetConsolidationFeeInfoForRequest(eth, chunk, args.FeeOverestimateFactor)
		utils.PanicOnError("error getting consolidation fee info", err)

		isSimulatedStr := ""
		if args.SimulateTransaction {
			isSimulatedStr = "[SIMULATED]"
		} else {
			isSimulatedStr = "[LIVE]"
		}

		// Prompt the user for consent.
		// We prompt for individual chunks because the predeploy includes exponential fee growth depending
		// on the size of the queue, and we want to make sure the user is aware of the rising fee.
		utils.PanicIfNoConsent(utils.SubmitSwitchRequestConsent(
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

		txn, err := core.SubmitConsolidationRequests(
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
				fmt.Println("Error submitting consolidation request. Printing successful requests:")
				printConsolidateTxnsAsJSON(txns)
			}

			utils.PanicOnError("consolidation request submission failed", err)
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
		printConsolidateTxnsAsJSON(txns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i, txn.Hash().Hex())
		}
	}

	return nil
}

func ConsolidateToTargetCommand(args TConsolidateToTargetCommandArgs) error {
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

	if len(args.SourceValidators) == 0 {
		return fmt.Errorf("usage: consolidate source-to-target --target <validatorIndexA> --sources <validatorIndexB>, <validatorIndexC>, ...")
	}

	eth, beaconClient, chainId, err := utils.GetClients(ctx, args.Node, args.BeaconNode, enableLogs)
	utils.PanicOnError("failed to reach ethereum clients", err)

	headState, err := utils.GetBeaconHeadState(ctx, beaconClient)
	utils.PanicOnError("failed to fetch beacon chain head state", err)

	eigenpodValidators, err := utils.GetEigenPodValidatorsByIndex(args.EigenpodAddress, headState)
	utils.PanicOnError("failed to fetch source validators for eigenpod", err)

	targetValidator, exists := eigenpodValidators[args.TargetValidator]
	if !exists {
		return fmt.Errorf("target validator (index %d) is not pointed at this eigenpod", args.TargetValidator)
	}

	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	utils.PanicOnError("failed to locate eigenpod. is your address correct?", err)

	status, err := eigenPod.ValidatorStatus(nil, targetValidator.PublicKey[:])
	utils.PanicOnError("failed to fetch target validator status", err)

	// Target validator must be in ACTIVE state in pod (verified withdrawal credentials; not withdrawn)
	if status == utils.ValidatorStatusInactive {
		return fmt.Errorf("target validator must have verified withdrawal credentials and be in the ACTIVE status. got status: INACTIVE")
	} else if status == utils.ValidatorStatusWithdrawn {
		return fmt.Errorf("target validator must have verified withdrawal credentials and be in the ACTIVE status. got status: WITHDRAWN")
	}

	// Form requests, filtering duplicate source validators and validators that aren't pointed at the pod
	requests := make([]EigenPod.IEigenPodTypesConsolidationRequest, 0)
	seen := make(map[uint64]bool)
	for _, vIndex := range args.SourceValidators {
		if v, exists := eigenpodValidators[vIndex]; !exists {
			return fmt.Errorf("source validator (index %d) is not pointed at this eigenpod", vIndex)
		} else if seen[vIndex] {
			return fmt.Errorf("source validator (index %d) is included twice in input args", vIndex)
		} else {
			seen[vIndex] = true

			requests = append(requests, EigenPod.IEigenPodTypesConsolidationRequest{
				SrcPubkey:    v.PublicKey[:],
				TargetPubkey: targetValidator.PublicKey[:],
			})
		}
	}

	requestChunks := utils.Chunk(requests, args.BatchSize)
	txns := make([]*types.Transaction, 0)

	for i, chunk := range requestChunks {
		feeInfo, err := utils.GetConsolidationFeeInfoForRequest(eth, chunk, args.FeeOverestimateFactor)
		utils.PanicOnError("error getting consolidation fee info", err)

		isSimulatedStr := ""
		if args.SimulateTransaction {
			isSimulatedStr = "[SIMULATED]"
		} else {
			isSimulatedStr = "[LIVE]"
		}

		// Prompt the user for consent.
		// We prompt for individual chunks because the predeploy includes exponential fee growth depending
		// on the size of the queue, and we want to make sure the user is aware of the rising fee.
		utils.PanicIfNoConsent(utils.SubmitSourceToTargetRequestConsent(
			len(chunk),
			feeInfo.CurrentQueueSize,
			toPrintableUnits(feeInfo.FeePerRequest),
			toPrintableUnits(feeInfo.TotalFee),
			toPrintableUnits(feeInfo.OverestimateFee),
			args.TargetValidator,
			len(args.SourceValidators),
			isSimulatedStr,
		))

		if isVerbose {
			color.Green("Submitting chunk %d/%d (msg.value: %s)", i+1, len(requestChunks), toPrintableUnits(feeInfo.OverestimateFee))
		}

		txn, err := core.SubmitConsolidationRequests(
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
				fmt.Println("Error submitting consolidation request. Printing successful requests:")
				printConsolidateTxnsAsJSON(txns)
			}

			utils.PanicOnError("consolidation request submission failed", err)
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
		printConsolidateTxnsAsJSON(txns)
	} else {
		for i, txn := range txns {
			color.Green("transaction(%d): %s", i, txn.Hash().Hex())
		}
	}

	return nil
}

// If the amount is greater than 0.0001 ETH, print as ETH
// If amount is less than 100_000 Wei, print as Wei
// Otherwise, print as Gwei
func toPrintableUnits(weiAmt *big.Int) string {
	printAsETHThreshhold := new(big.Int).Mul(
		big.NewInt(100_000),
		big.NewInt(params.GWei),
	)

	printAsWeiThreshhold := new(big.Int).Mul(
		big.NewInt(100_000),
		big.NewInt(params.Wei),
	)

	if weiAmt.Cmp(printAsETHThreshhold) > 0 {
		return fmt.Sprintf("%f ETH", utils.IweiToEther(weiAmt))
	} else if weiAmt.Cmp(printAsWeiThreshhold) < 0 {
		return fmt.Sprintf("%d Wei", weiAmt)
	} else {
		return fmt.Sprintf("%f Gwei", utils.WeiToGwei(weiAmt))
	}
}

func printConsolidateTxnsAsJSON(txns []*types.Transaction) {
	printableTxns := lo.Map(txns, func(txn *types.Transaction, _ int) PredeployRequestTransaction {
		gas := txn.Gas()
		return PredeployRequestTransaction{
			Transaction: Transaction{
				To:              txn.To().Hex(),
				CallData:        common.Bytes2Hex(txn.Data()),
				Type:            "consolidation_request",
				GasEstimateGwei: &gas,
			},
			Value: txn.Value(),
		}
	})
	PrintAsJSON(printableTxns)
}
