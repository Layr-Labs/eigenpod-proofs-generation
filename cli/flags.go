package main

import (
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	cli "github.com/urfave/cli/v2"
)

// Required for commands that need an EigenPod's address
var PodAddressFlag = &cli.StringFlag{
	Name:        "podAddress",
	Aliases:     []string{"p", "pod"},
	Value:       "",
	Usage:       "[required] The onchain `address` of your eigenpod contract (0x123123123123)",
	Required:    true,
	Destination: &eigenpodAddress,
}

var PodOwnerFlag = &cli.StringFlag{
	Name:        "podOwner",
	Aliases:     []string{"p", "podOwner"},
	Value:       "",
	Usage:       "[required] The onchain `address` of your eigenpod's owner (0x123123123123)",
	Required:    true,
	Destination: &eigenpodOwner,
}

// Required for commands that need a beacon chain RPC
var BeaconNodeFlag = &cli.StringFlag{
	Name:        "beaconNode",
	Aliases:     []string{"b"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning beacon node RPC (https://)",
	Required:    true,
	Destination: &beacon,
}

// Required for commands that need an execution layer RPC
var ExecNodeFlag = &cli.StringFlag{
	Name:        "execNode",
	Aliases:     []string{"e"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning execution-layer RPC (https://)",
	Required:    true,
	Destination: &node,
}

// Optional commands:

// Optional use for commands that want direct tx submission from a specific private key
var SenderPkFlag = &cli.StringFlag{
	Name:        "sender",
	Aliases:     []string{"s"},
	Value:       "",
	Usage:       "`Private key` of the account that will send any transactions. If set, this will automatically submit the proofs to their corresponding onchain functions after generation. If using checkpoint mode, it will also begin a checkpoint if one hasn't been started already.",
	Destination: &sender,
}

var EstimateGasFlag = &cli.BoolFlag{
	Name:        "gas",
	Aliases:     []string{"g"},
	Value:       false,
	Usage:       "Estimate gas on the transaction (using `--sender` as the sender for simulation). This will NOT send the transaction.",
	Destination: &estimateGas,
}

var AmountWeiFlag = &cli.Uint64Flag{
	Name:        "amountWei",
	Aliases:     []string{},
	Value:       0,
	Usage:       "The amount, in Wei.",
	Destination: &amountWei,
}

var VerboseFlag = &cli.BoolFlag{
	Name:        "verbose",
	Aliases:     []string{"v"},
	Value:       false,
	Usage:       "Enable verbose output.",
	Destination: &verbose,
}

// Optional use for commands that support JSON output
var PrintJSONFlag = &cli.BoolFlag{
	Name:        "json",
	Value:       false,
	Usage:       "print only plain JSON",
	Required:    false,
	Destination: &useJSON,
}

// shared flag --batch
func BatchBySize(destination *uint64, defaultValue uint64) *cli.Uint64Flag {
	return &cli.Uint64Flag{
		Name:        "batch",
		Value:       defaultValue,
		Usage:       "Submit proofs/requests in groups of size `batchSize`, to avoid gas limit.",
		Required:    false,
		Destination: destination,
	}
}

// Flags for each consolidation subcommand
var ConsolidationFlags = []cli.Flag{
	VerboseFlag,
	PodAddressFlag,
	BeaconNodeFlag,
	ExecNodeFlag,
	SenderPkFlag,
	EstimateGasFlag,
	PrintJSONFlag,
	BatchBySize(&batchSize, utils.DEFAULT_BATCH_CONSOLIDATE),
	// &cli.BoolFlag{
	// 	Name: "no-warn",
	// 	Usage: "Turn off warnings for various consolidation failure states. By default, this command will prompt you if a consolidation might fail because of:\n" +
	// 		"* invalid switch request (source already has 0x02 credentials)\n" +
	// 		"* conflicting consolidation already in progress\n" +
	// 		"* conflicting withdrawal requests already in progress\n" +
	// 		"* invalid consolidation source/target\n" +
	// 		"* ... and more.", // TODO
	// 	Destination: &noWarn,
	// },
	&cli.Float64Flag{
		Name:        "fee-overestimate-factor",
		Aliases:     []string{"overestimate", "over"},
		Usage:       "Specify how much to overestimate the predeploy request fee by when sending. Overestimating can help ensure your request succeeds even if the request fee changes before your transaction is included.",
		Value:       feeOverestimateFactor,
		Destination: &feeOverestimateFactor,
	},
}

// Flags for each consolidation subcommand
var RequestWithdrawalFlags = []cli.Flag{
	VerboseFlag,
	PodAddressFlag,
	BeaconNodeFlag,
	ExecNodeFlag,
	SenderPkFlag,
	EstimateGasFlag,
	PrintJSONFlag,
	BatchBySize(&batchSize, utils.DEFAULT_BATCH_WITHDRAWREQUEST),
	// &cli.BoolFlag{
	// 	Name: "no-warn",
	// 	Usage: "Turn off warnings for various consolidation failure states. By default, this command will prompt you if a consolidation might fail because of:\n" +
	// 		"* invalid switch request (source already has 0x02 credentials)\n" +
	// 		"* conflicting consolidation already in progress\n" +
	// 		"* conflicting withdrawal requests already in progress\n" +
	// 		"* invalid consolidation source/target\n" +
	// 		"* ... and more.", // TODO
	// 	Destination: &noWarn,
	// },
	&cli.Float64Flag{
		Name:        "fee-overestimate-factor",
		Aliases:     []string{"overestimate", "over"},
		Usage:       "Specify how much to overestimate the predeploy request fee by when sending. Overestimating can help ensure your request succeeds even if the request fee changes before your transaction is included.",
		Value:       feeOverestimateFactor,
		Destination: &feeOverestimateFactor,
	},
}

// Hack to make a copy of a flag that sets `Required` to true
func Require(flag *cli.StringFlag) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        flag.Name,
		Aliases:     flag.Aliases,
		Value:       flag.Value,
		Usage:       flag.Usage,
		Destination: flag.Destination,
		Required:    true,
	}
}
