package main

import (
	"math"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/commands"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	cli "github.com/urfave/cli/v2"
)

// Destinations for values set by various flags
var eigenpodAddress, beacon, node, sender, eigenpodOwner string
var useJSON = false
var specificValidator uint64 = math.MaxUint64
var estimateGas = false
var slashedValidatorIndex uint64
var amountWei uint64
var verbose = false
var checkFee = false
var feeOverestimateFactor = float64(1.5)
var batchSize uint64

const DefaultHealthcheckTolerance = float64(5.0)

func main() {
	var forceCheckpoint = false
	var disableColor = false
	var noPrompt = false
	var tolerance = DefaultHealthcheckTolerance

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLI",
		HelpName:               "./cli",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator. By default, the unsigned transactions will be printed to stdout as JSON. If you want to sign and broadcast these automatically, pass `--sender <pk>`.",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:      "find-stale-pods",
				Args:      true,
				Usage:     "Locate stale pods, whose balances have deviated by more than 5% due to beacon slashing.",
				UsageText: "./cli find-stale-pods <args>",
				Flags: []cli.Flag{
					VerboseFlag,
					ExecNodeFlag,
					BeaconNodeFlag,
					&cli.Float64Flag{
						Name:        "tolerance",
						Value:       DefaultHealthcheckTolerance, // default: 5
						Usage:       "The percentage balance deviation to tolerate when deciding whether an eigenpod should be corrected. Default is 5% (e.g --tolerance 5).",
						Destination: &tolerance,
					},
				},
				Action: func(_ *cli.Context) error {
					return commands.FindStalePodsCommand(commands.TFindStalePodsCommandArgs{
						EthNode:    node,
						BeaconNode: beacon,
						Verbose:    verbose,
						Tolerance:  tolerance,
					})
				},
			},
			{
				Name:      "correct-stale-pod",
				Args:      true,
				Usage:     "Correct a stale balance on an eigenpod, which has been slashed on the beacon chain.",
				UsageText: "./cli correct-stale-pod [FLAGS] <validatorIndex>",
				Flags: []cli.Flag{
					VerboseFlag,
					PodAddressFlag,
					ExecNodeFlag,
					BeaconNodeFlag,
					BatchBySize(&batchSize, utils.DEFAULT_BATCH_CHECKPOINT),
					Require(SenderPkFlag),
					&cli.Uint64Flag{
						Name:        "validatorIndex",
						Usage:       "The index of a validator slashed that belongs to the pod.",
						Required:    true,
						Destination: &slashedValidatorIndex,
					},
				},
				Action: func(_ *cli.Context) error {
					return commands.FixStaleBalance(commands.TFixStaleBalanceArgs{
						EthNode:               node,
						BeaconNode:            beacon,
						Sender:                sender,
						EigenpodAddress:       eigenpodAddress,
						SlashedValidatorIndex: slashedValidatorIndex,
						Verbose:               verbose,
						CheckpointBatchSize:   batchSize,
						NoPrompt:              noPrompt,
					})
				},
			},
			{
				Name:      "assign-submitter",
				Args:      true,
				Usage:     "Assign a different address to be able to submit your proofs. You'll always be able to submit from your EigenPod owner PK.",
				UsageText: "./cli assign-submitter [FLAGS] <0xsubmitter>",
				Flags: []cli.Flag{
					VerboseFlag,
					PodAddressFlag,
					ExecNodeFlag,
					Require(SenderPkFlag),
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
					VerboseFlag,
					PodAddressFlag,
					BeaconNodeFlag,
					ExecNodeFlag,
					PrintJSONFlag,
				},
				Action: func(_ *cli.Context) error {
					return commands.StatusCommand(commands.TStatusArgs{
						EigenpodAddress: eigenpodAddress,
						DisableColor:    disableColor,
						UseJSON:         useJSON,
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
					VerboseFlag,
					PodAddressFlag,
					BeaconNodeFlag,
					ExecNodeFlag,
					SenderPkFlag,
					EstimateGasFlag,
					BatchBySize(&batchSize, utils.DEFAULT_BATCH_CHECKPOINT),
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "If true, starts a checkpoint even if the pod has no native ETH to award shares",
						Destination: &forceCheckpoint,
					},
				},
				Action: func(_ *cli.Context) error {
					return commands.CheckpointCommand(commands.TCheckpointCommandArgs{
						DisableColor:        disableColor,
						NoPrompt:            noPrompt,
						SimulateTransaction: sender == "" || estimateGas,
						BatchSize:           batchSize,
						ForceCheckpoint:     forceCheckpoint,
						Node:                node,
						BeaconNode:          beacon,
						EigenpodAddress:     eigenpodAddress,
						Verbose:             verbose,
						Sender:              sender,
					})
				},
			},
			{
				Name:    "credentials",
				Aliases: []string{"cr", "creds"},
				Usage:   "Generates a proof for use with EigenPod.verifyWithdrawalCredentials()",
				Flags: []cli.Flag{
					VerboseFlag,
					PodAddressFlag,
					BeaconNodeFlag,
					ExecNodeFlag,
					SenderPkFlag,
					EstimateGasFlag,
					PrintJSONFlag,
					BatchBySize(&batchSize, utils.DEFAULT_BATCH_CREDENTIALS),
					&cli.Uint64Flag{
						Name:        "validatorIndex",
						Usage:       "The `index` of a specific validator to prove (e.g a slashed validator for `verifyStaleBalance()`).",
						Destination: &specificValidator,
					},
				},
				Action: func(_ *cli.Context) error {
					return commands.CredentialsCommand(commands.TCredentialCommandArgs{
						EigenpodAddress:     eigenpodAddress,
						DisableColor:        disableColor,
						UseJSON:             useJSON,
						SimulateTransaction: sender == "" || estimateGas,
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
			{
				Name:  "consolidate",
				Usage: "(EIP-7521) Consolidates eligible validators via EigenPod.requestConsolidation()",
				Subcommands: []*cli.Command{
					{
						Name:    "switch-to-compounding",
						Aliases: []string{"switch"},
						Usage:   "Specify a list of validator indices to switch from 0x01 to 0x02 withdrawal prefix.",
						Flags: append(
							ConsolidationFlags,
							&cli.Uint64SliceFlag{
								Name:     "validators",
								Required: true,
								Usage:    "[required] The list of validator indices to switch from 0x01 to 0x02 withdrawal credentials",
							},
						),
						Action: func(ctx *cli.Context) error {
							return commands.ConsolidateSwitchCommand(commands.TConsolidateSwitchCommandArgs{
								ConsolidateBaseCommandArgs: commands.ConsolidateBaseCommandArgs{
									EigenpodAddress:       eigenpodAddress,
									DisableColor:          disableColor,
									UseJSON:               useJSON,
									SimulateTransaction:   sender == "" || estimateGas,
									Node:                  node,
									BeaconNode:            beacon,
									Sender:                sender,
									BatchSize:             batchSize,
									NoPrompt:              noPrompt,
									Verbose:               verbose,
									CheckFee:              checkFee,
									FeeOverestimateFactor: feeOverestimateFactor,
								},
								Validators: ctx.Uint64Slice("validators"),
							})
						},
					},
					{
						Name:  "source-to-target",
						Usage: "Specify a target validator inbdex and a list of source validator indices to consolidate into the target.",
						Flags: append(
							ConsolidationFlags,
							&cli.Uint64Flag{
								Name:     "target",
								Usage:    "[required] Specify the target validator index for a consolidation",
								Required: true,
							},
							&cli.Uint64SliceFlag{
								Name:     "sources",
								Usage:    "[required] Specify the source validator indices for a consolidation",
								Required: true,
							},
						),
						Action: func(ctx *cli.Context) error {
							return commands.ConsolidateToTargetCommand(commands.TConsolidateToTargetCommandArgs{
								ConsolidateBaseCommandArgs: commands.ConsolidateBaseCommandArgs{
									EigenpodAddress:       eigenpodAddress,
									DisableColor:          disableColor,
									UseJSON:               useJSON,
									SimulateTransaction:   sender == "" || estimateGas,
									Node:                  node,
									BeaconNode:            beacon,
									Sender:                sender,
									BatchSize:             batchSize,
									NoPrompt:              noPrompt,
									Verbose:               verbose,
									CheckFee:              checkFee,
									FeeOverestimateFactor: feeOverestimateFactor,
								},
								TargetValidator:  ctx.Uint64("target"),
								SourceValidators: ctx.Uint64Slice("sources"),
							})
						},
					},
				},
			},
			{
				Name:  "request-withdrawal",
				Usage: "(EIP-7002) Request partial or full exits via EigenPod.requestWithdrawal()",
				Subcommands: []*cli.Command{
					{
						Name:    "full-exit",
						Aliases: []string{"full"},
						Usage:   "Specify a list of validator indices to fully withdraw from the beacon chain.",
						Flags: append(
							RequestWithdrawalFlags,
							&cli.Uint64SliceFlag{
								Name:     "validators",
								Required: true,
								Usage:    "[required] The list of validator indices to exit from the beacon chain",
							},
						),
						Action: func(ctx *cli.Context) error {
							return commands.RequestFullExitCommand(commands.TRequestFullExitCommandArgs{
								WithdrawalBaseCommandArgs: commands.WithdrawalBaseCommandArgs{
									EigenpodAddress:       eigenpodAddress,
									DisableColor:          disableColor,
									UseJSON:               useJSON,
									SimulateTransaction:   sender == "" || estimateGas,
									Node:                  node,
									BeaconNode:            beacon,
									Sender:                sender,
									BatchSize:             batchSize,
									NoPrompt:              noPrompt,
									Verbose:               verbose,
									CheckFee:              checkFee,
									FeeOverestimateFactor: feeOverestimateFactor,
								},
								Validators: ctx.Uint64Slice("validators"),
							})
						},
					},
					{
						Name:  "partial",
						Usage: "Specify a list of validator indices and gwei amounts to request beacon chain withdrawals for.",
						Flags: append(
							RequestWithdrawalFlags,
							&cli.Uint64SliceFlag{
								Name:     "validators",
								Required: true,
								Usage:    "[required] The list of validator indices for which partial withdrawal requests will be submitted",
							},
							&cli.Uint64SliceFlag{
								Name:     "amounts",
								Required: true,
								Usage:    "[required] The amount (in gwei) for each partial withdrawal request",
							},
						),
						Action: func(ctx *cli.Context) error {
							return commands.RequestPartialWithdrawalCommand(commands.TRequestPartialWithdrawalCommandArgs{
								WithdrawalBaseCommandArgs: commands.WithdrawalBaseCommandArgs{
									EigenpodAddress:       eigenpodAddress,
									DisableColor:          disableColor,
									UseJSON:               useJSON,
									SimulateTransaction:   sender == "" || estimateGas,
									Node:                  node,
									BeaconNode:            beacon,
									Sender:                sender,
									BatchSize:             batchSize,
									NoPrompt:              noPrompt,
									Verbose:               verbose,
									CheckFee:              checkFee,
									FeeOverestimateFactor: feeOverestimateFactor,
								},
								Validators: ctx.Uint64Slice("validators"),
								AmtsGwei:   ctx.Uint64Slice("amounts"),
							})
						},
					},
				},
			},
			{
				Name:  "complete-all-withdrawals",
				Args:  true,
				Usage: "Completes all withdrawals queued on the podOwner, for which Native ETH is the sole strategy in the withdrawal. Attempts to execute a group of withdrawals whose sum does not exceed Pod.withdrawableRestakedExecutionLayerGwei() in value.",
				Flags: []cli.Flag{
					VerboseFlag,
					ExecNodeFlag,
					PodAddressFlag,
					SenderPkFlag,
					EstimateGasFlag,
				},
				Action: func(_ *cli.Context) error {
					return commands.CompleteAllWithdrawalsCommand(commands.TCompleteWithdrawalArgs{
						EthNode:     node,
						EigenPod:    eigenpodAddress,
						Sender:      sender,
						EstimateGas: estimateGas,
					})
				},
			},
			{
				Name:  "queue-withdrawal",
				Args:  true,
				Usage: "Queues a withdrawal for shares associated with the native ETH strategy. Queues a withdrawal whose size does not exceed Pod.withdrawableRestakedExecutionLayerGwei() in value.",
				Flags: []cli.Flag{
					VerboseFlag,
					ExecNodeFlag,
					PodAddressFlag,
					SenderPkFlag,
					EstimateGasFlag,
					AmountWeiFlag,
				},
				Action: func(_ *cli.Context) error {
					return commands.QueueWithdrawalCommand(commands.TQueueWithdrawallArgs{
						EthNode:     node,
						EigenPod:    eigenpodAddress,
						Sender:      sender,
						EstimateGas: estimateGas,
						AmountWei:   amountWei,
					})
				},
			},
			{
				Name:  "show-withdrawals",
				Args:  true,
				Usage: "Shows all pending withdrawals for the podOwner.",
				Flags: []cli.Flag{
					VerboseFlag,
					ExecNodeFlag,
					PodAddressFlag,
				},
				Action: func(_ *cli.Context) error {
					return commands.ShowWithdrawalsCommand(commands.TShowWithdrawalArgs{
						EthNode:  node,
						EigenPod: eigenpodAddress,
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
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
