package main

import (
	"math"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/commands"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	cli "github.com/urfave/cli/v2"
)

// Destinations for values set by various flags
var eigenpodAddress, beacon, node, sender string
var useJson bool = false
var specificValidator uint64 = math.MaxUint64

func main() {
	var batchSize uint64
	var forceCheckpoint, disableColor, verbose bool
	var noPrompt bool

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator. By default, the unsigned transactions will be printed to stdout as JSON. If you want to sign and broadcast these automatically, pass `--sender <pk>`.",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:      "assign-submitter",
				Args:      true,
				Usage:     "Assign a different address to be able to submit your proofs. You'll always be able to submit from your EigenPod owner PK.",
				UsageText: "./cli assign-submitter [FLAGS] <0xsubmitter>",
				Flags: []cli.Flag{
					POD_ADDRESS_FLAG,
					EXEC_NODE_FLAG,
					Require(SENDER_PK_FLAG),
				},
				Action: func(cctx *cli.Context) error {
					return commands.AssignSubmitterCommand(commands.TAssignSubmitterArgs{
						Node:            node,
						TargetAddress:   cctx.Args().First(),
						Sender:          sender,
						EigenpodAddress: eigenpodAddress,
						NoPrompt:        noPrompt,
						Verbose:         verbose,
					})
				},
			},
			{
				Name:  "status",
				Usage: "Checks the status of your eigenpod.",
				Flags: []cli.Flag{
					POD_ADDRESS_FLAG,
					BEACON_NODE_FLAG,
					EXEC_NODE_FLAG,
					PRINT_JSON_FLAG,
				},
				Action: func(cctx *cli.Context) error {
					return commands.StatusCommand(commands.TStatusArgs{
						EigenpodAddress: eigenpodAddress,
						DisableColor:    disableColor,
						UseJSON:         useJson,
						Node:            node,
						BeaconNode:      beacon,
						Verbose:         verbose,
					})
				},
			},
			{
				Name:    "checkpoint",
				Aliases: []string{"cp"},
				Usage:   "Generates a proof for use with EigenPod.verifyCheckpointProofs().",
				Flags: []cli.Flag{
					POD_ADDRESS_FLAG,
					BEACON_NODE_FLAG,
					EXEC_NODE_FLAG,
					SENDER_PK_FLAG,
					BatchBySize(&batchSize, utils.DEFAULT_BATCH_CHECKPOINT),
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "If true, starts a checkpoint even if the pod has no native ETH to award shares",
						Destination: &forceCheckpoint,
					},
				},
				Action: func(cctx *cli.Context) error {
					return commands.CheckpointCommand(commands.TCheckpointCommandArgs{
						DisableColor:        disableColor,
						NoPrompt:            noPrompt,
						SimulateTransaction: len(sender) == 0,
						BatchSize:           batchSize,
						ForceCheckpoint:     forceCheckpoint,
						Node:                node,
						BeaconNode:          beacon,
						EigenpodAddress:     eigenpodAddress,
						Verbose:             verbose,
					})
				},
			},
			{
				Name:    "credentials",
				Aliases: []string{"cr", "creds"},
				Usage:   "Generates a proof for use with EigenPod.verifyWithdrawalCredentials()",
				Flags: []cli.Flag{
					POD_ADDRESS_FLAG,
					BEACON_NODE_FLAG,
					EXEC_NODE_FLAG,
					SENDER_PK_FLAG,
					BatchBySize(&batchSize, utils.DEFAULT_BATCH_CREDENTIALS),
					&cli.Uint64Flag{
						Name:        "validatorIndex",
						Usage:       "The `index` of a specific validator to prove (e.g a slashed validator for `verifyStaleBalance()`).",
						Destination: &specificValidator,
					},
				},
				Action: func(cctx *cli.Context) error {
					return commands.CredentialsCommand(commands.TCredentialCommandArgs{
						EigenpodAddress:     eigenpodAddress,
						DisableColor:        disableColor,
						UseJSON:             useJson,
						SimulateTransaction: len(sender) == 0,
						Node:                node,
						BeaconNode:          beacon,
						Sender:              sender,
						SpecificValidator:   specificValidator,
						BatchSize:           batchSize,
						NoPrompt:            noPrompt,
						Verbose:             verbose,
					})
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-color",
				Value:       false,
				Usage:       "Disables color output for terminals that do not support ANSI color codes.",
				Destination: &disableColor,
			},
			&cli.BoolFlag{
				Name:        "no-prompt",
				Value:       false,
				Usage:       "Disables prompts to approve any transactions occurring (e.g in CI).",
				Destination: &noPrompt,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "Enable verbose output.",
				Destination: &verbose,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
