package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"time"

	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func shortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}

// shared flag --batch
func BatchBySize(destination *uint64, defaultValue uint64) *cli.Uint64Flag {
	return &cli.Uint64Flag{
		Name:        "batch",
		Value:       defaultValue,
		Usage:       "Submit proofs in groups of size `--batch <batchSize>`, to avoid gas limit.",
		Required:    false,
		Destination: destination,
	}
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

// Destinations for values set by various flags
var eigenpodAddress, beacon, node, sender, output string
var useJson bool = false

// Required flags:

// Required for commands that need an EigenPod's address
var POD_ADDRESS_FLAG = &cli.StringFlag{
	Name:        "podAddress",
	Aliases:     []string{"p", "pod"},
	Value:       "",
	Usage:       "[required] The onchain `address` of your eigenpod contract (0x123123123123)",
	Required:    true,
	Destination: &eigenpodAddress,
}

// Required for commands that need a beacon chain RPC
var BEACON_NODE_FLAG = &cli.StringFlag{
	Name:        "beaconNode",
	Aliases:     []string{"b"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning beacon node RPC (https://)",
	Required:    true,
	Destination: &beacon,
}

// Required for commands that need an execution layer RPC
var EXEC_NODE_FLAG = &cli.StringFlag{
	Name:        "execNode",
	Aliases:     []string{"e"},
	Value:       "",
	Usage:       "[required] `URL` to a functioning execution-layer RPC (https://)",
	Required:    true,
	Destination: &node,
}

// Optional commands:

// Optional use for commands that want direct tx submission from a specific private key
var SENDER_PK_FLAG = &cli.StringFlag{
	Name:        "sender",
	Aliases:     []string{"s"},
	Value:       "",
	Usage:       "`Private key` of the account that will send any transactions. If set, this will automatically submit the proofs to their corresponding onchain functions after generation. If using checkpoint mode, it will also begin a checkpoint if one hasn't been started already.",
	Destination: &sender,
}

// Optional use for commands that support JSON output
var PRINT_JSON_FLAG = &cli.BoolFlag{
	Name:        "json",
	Value:       false,
	Usage:       "print only plain JSON",
	Required:    false,
	Destination: &useJson,
}

// maximum number of proofs per txn for each of the following proof types:
const DEFAULT_BATCH_CREDENTIALS = 60
const DEFAULT_BATCH_CHECKPOINT = 80

func main() {
	var batchSize uint64
	var checkpointProofPath string
	var forceCheckpoint, disableColor, verbose bool
	var noPrompt bool
	ctx := context.Background()

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator.",
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
					targetAddress := cctx.Args().First()
					if len(targetAddress) == 0 {
						return fmt.Errorf("usage: `assign-submitter <0xsubmitter>`")
					} else if !gethCommon.IsHexAddress(targetAddress) {
						return fmt.Errorf("invalid address for 0xsubmitter: %s", targetAddress)
					}

					eth, err := ethclient.Dial(node)
					if err != nil {
						return fmt.Errorf("failed to reach eth --node: %w", err)
					}

					chainId, err := eth.ChainID(ctx)
					if err != nil {
						return fmt.Errorf("failed to reach eth node for chain id: %w", err)
					}

					ownerAccount, err := core.PrepareAccount(&sender, chainId)
					if err != nil {
						return fmt.Errorf("failed to parse --sender: %w", err)
					}

					pod, err := onchain.NewEigenPod(gethCommon.HexToAddress(eigenpodAddress), eth)
					if err != nil {
						return fmt.Errorf("error contacting eigenpod: %w", err)
					}

					// Check that the existing submitter is not the current submitter
					newSubmitter := gethCommon.HexToAddress(targetAddress)
					currentSubmitter, err := pod.ProofSubmitter(nil)
					if err != nil {
						return fmt.Errorf("error fetching current proof submitter: %w", err)
					} else if currentSubmitter.Cmp(newSubmitter) == 0 {
						return fmt.Errorf("error: new proof submitter is existing proof submitter (%s)", currentSubmitter)
					}

					if !noPrompt {
						fmt.Printf("Your pod's current proof submitter is %s.\n", currentSubmitter)
						core.PanicIfNoConsent(fmt.Sprintf("This will update your EigenPod to allow %s to submit proofs on its behalf. As the EigenPod's owner, you can always change this later.", newSubmitter))
					}

					txn, err := pod.SetProofSubmitter(ownerAccount.TransactionOptions, newSubmitter)
					if err != nil {
						return fmt.Errorf("error updating submitter role: %w", err)
					}

					color.Green("submitted txn: %s", txn.Hash())
					color.Green("updated!")

					return nil
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
					if disableColor {
						color.NoColor = true
					}

					eth, beaconClient, _, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to load ethereum clients", err)

					status := core.GetStatus(ctx, eigenpodAddress, eth, beaconClient)

					if useJson {
						bytes, err := json.MarshalIndent(status, "", "      ")
						core.PanicOnError("failed to get status", err)
						statusStr := string(bytes)
						fmt.Println(statusStr)
					} else {
						bold := color.New(color.Bold, color.FgBlue)
						ital := color.New(color.Italic, color.FgBlue)
						ylw := color.New(color.Italic, color.FgHiYellow)

						bold.Printf("Eigenpod Status\n")
						ital.Printf("- Pod owner address: ")
						ylw.Printf("%s\n", status.PodOwner)
						ital.Printf("- Proof submitter address: ")
						ylw.Printf("%s\n", status.ProofSubmitter)
						fmt.Println()

						// sort validators by status
						inactiveValidators, activeValidators, withdrawnValidators := core.SortByStatus(status.Validators)
						var targetColor *color.Color

						bold.Printf("Eigenpod validators:\n============\n")
						ital.Printf("Format: #ValidatorIndex (pubkey) [effective balance] [current balance]\n")

						// print info on inactive validators
						// these validators can be added to the pod's active validator set
						// by running the `credentials` command
						if len(inactiveValidators) != 0 {
							targetColor = color.New(color.FgHiYellow)

							color.New(color.Bold, color.FgHiYellow).Printf("- [INACTIVE] - Run `credentials` to verify these %d validators' withdrawal credentials:\n", len(inactiveValidators))

							for _, validator := range inactiveValidators {
								publicKey := validator.PublicKey
								if !verbose {
									publicKey = shortenHex(publicKey)
								}

								if validator.Slashed {
									targetColor.Printf("\t- #%d (%s) [%d] [%d] (slashed on beacon chain)\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								} else {
									targetColor.Printf("\t- #%d (%s) [%d] [%d]\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								}

							}

							fmt.Println()
						}

						// print info on active validators
						// these validators can be checkpointed using the `checkpoint` command
						if len(activeValidators) != 0 {
							targetColor = color.New(color.FgGreen)

							color.New(color.Bold, color.FgGreen).Printf("- [ACTIVE] - Run `checkpoint` to update these %d validators' balances:\n", len(activeValidators))

							for _, validator := range activeValidators {
								publicKey := validator.PublicKey
								if !verbose {
									publicKey = shortenHex(publicKey)
								}

								if validator.Slashed {
									targetColor.Printf("\t- #%d (%s) [%d] [%d] (slashed on beacon chain)\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								} else {
									targetColor.Printf("\t- #%d (%s) [%d] [%d]\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								}
							}

							fmt.Println()
						}

						// print info on withdrawn validators
						// no further action is required to manage these validators in the pod
						if len(withdrawnValidators) != 0 {
							targetColor = color.New(color.FgHiRed)

							color.New(color.Bold, color.FgHiRed).Printf("- [WITHDRAWN] - %d validators:\n", len(withdrawnValidators))

							for _, validator := range withdrawnValidators {
								publicKey := validator.PublicKey
								if !verbose {
									publicKey = shortenHex(publicKey)
								}

								if validator.Slashed {
									targetColor.Printf("\t- #%d (%s) [%d] [%d] (slashed on beacon chain)\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								} else {
									targetColor.Printf("\t- #%d (%s) [%d] [%d]\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
								}
							}

							fmt.Println()
						}

						// Calculate the change in shares for completing a checkpoint
						deltaETH := new(big.Float).Sub(
							status.TotalSharesAfterCheckpointETH,
							status.CurrentTotalSharesETH,
						)

						if status.ActiveCheckpoint != nil {
							startTime := time.Unix(int64(status.ActiveCheckpoint.StartedAt), 0)

							bold.Printf("!NOTE: There is a checkpoint active! (started at: %s)\n", startTime.String())

							ital.Printf("\t- If you finish it, you may receive up to %f shares. (%f -> %f)\n", deltaETH, status.CurrentTotalSharesETH, status.TotalSharesAfterCheckpointETH)

							ital.Printf("\t- %d proof(s) remaining until completion.\n", status.ActiveCheckpoint.ProofsRemaining)
						} else {
							bold.Printf("Running a `checkpoint` right now will result in: \n")

							ital.Printf("\t%f new shares issued (%f ==> %f)\n", deltaETH, status.CurrentTotalSharesETH, status.TotalSharesAfterCheckpointETH)

							if status.MustForceCheckpoint {
								ylw.Printf("\tNote: pod does not have checkpointable native ETH. To checkpoint anyway, run `checkpoint` with the `--force` flag.\n")
							}

							bold.Printf("Batching %d proofs per txn, this will require:\n\t", DEFAULT_BATCH_CHECKPOINT)
							ital.Printf("- 1x startCheckpoint() transaction, and \n\t- %dx EigenPod.verifyCheckpointProofs() transaction(s)\n\n", int(math.Ceil(float64(status.NumberValidatorsToCheckpoint)/float64(DEFAULT_BATCH_CHECKPOINT))))
						}
					}
					return nil
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
					BatchBySize(&batchSize, DEFAULT_BATCH_CHECKPOINT),
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "If true, starts a checkpoint even if the pod has no native ETH to award shares",
						Destination: &forceCheckpoint,
					},
					&cli.StringFlag{
						Name:        "proof",
						Value:       "",
						Usage:       "the path to a previous proof generated from this step (via `-o proof.json`). If provided, this proof will submitted to network via the `--sender` flag.",
						Destination: &checkpointProofPath,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"O", "output"},
						Value:       "",
						Usage:       "Output `path` for the proof. (defaults to stdout). NOTE: If `--out` is supplied along with `--sender`, `--out` takes precedence and the proof will not be broadcast.",
						Destination: &output,
					},
				},
				Action: func(cctx *cli.Context) error {
					if disableColor {
						color.NoColor = true
					}

					var out *string = nil
					if len(cctx.String("out")) > 0 {
						outProp := cctx.String("out")
						out = &outProp
					}

					eth, beaconClient, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to reach ethereum clients", err)

					if len(checkpointProofPath) > 0 {
						// user specified the proof
						if len(sender) == 0 {
							core.Panic("If using --proof, --sender <privateKey> must also be supplied.")
						}

						// load `proof` from file.
						proof, err := core.LoadCheckpointProofFromFile(checkpointProofPath)
						core.PanicOnError("failed to parse checkpoint proof from file", err)

						txns, err := core.SubmitCheckpointProof(ctx, sender, eigenpodAddress, chainId, proof, eth, batchSize, noPrompt)
						for _, txn := range txns {
							color.Green("submitted txn: %s", txn.Hash())
						}
						core.PanicOnError("an error occurred while submitting your checkpoint proofs", err)
						return nil
					}

					currentCheckpoint, err := core.GetCurrentCheckpoint(eigenpodAddress, eth)
					core.PanicOnError("failed to load checkpoint", err)

					if currentCheckpoint == 0 {
						if len(sender) != 0 {
							if !noPrompt {
								core.PanicIfNoConsent(core.StartCheckpointProofConsent())
							}

							newCheckpoint, err := core.StartCheckpoint(ctx, eigenpodAddress, sender, chainId, eth, forceCheckpoint)
							core.PanicOnError("failed to start checkpoint", err)
							currentCheckpoint = newCheckpoint
						} else {
							core.PanicOnError("no checkpoint active and no private key provided to start one", errors.New("no checkpoint"))
						}
					}
					color.Green("pod has active checkpoint! checkpoint timestamp: %d", currentCheckpoint)

					proof, err := core.GenerateCheckpointProof(ctx, eigenpodAddress, eth, chainId, beaconClient)
					core.PanicOnError("failed to generate checkpoint proof", err)

					jsonString, err := json.Marshal(proof)
					core.PanicOnError("failed to generate JSON proof data.", err)

					if out != nil {
						core.WriteOutputToFileOrStdout(jsonString, out)
					} else if len(sender) != 0 {
						txns, err := core.SubmitCheckpointProof(ctx, sender, eigenpodAddress, chainId, proof, eth, batchSize, noPrompt)
						for _, txn := range txns {
							color.Green("submitted txn: %s", txn.Hash())
						}
						core.PanicOnError("an error occurred while submitting your checkpoint proofs", err)
					}

					return nil
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
					BatchBySize(&batchSize, DEFAULT_BATCH_CREDENTIALS),
				},
				Action: func(cctx *cli.Context) error {
					if disableColor {
						color.NoColor = true
					}

					eth, beaconClient, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to reach ethereum clients", err)

					validatorProofs, oracleBeaconTimestamp, err := core.GenerateValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient)
					if err != nil || validatorProofs == nil {
						core.PanicOnError("Failed to generate validator proof", err)
						core.Panic("no inactive validators")
					}

					if len(sender) != 0 {
						txns, err := core.SubmitValidatorProof(ctx, sender, eigenpodAddress, chainId, eth, batchSize, validatorProofs, oracleBeaconTimestamp, noPrompt)
						for i, txn := range txns {
							color.Green("transaction(%d): %s", i, txn.Hash().Hex())
						}
						core.PanicOnError("failed to invoke verifyWithdrawalCredentials", err)
					} else {
						data := map[string]any{
							"validatorProofs": validatorProofs,
						}
						out, err := json.MarshalIndent(data, "", "   ")
						core.PanicOnError("failed to process proof", err)

						fmt.Printf("%s\n", out)
					}
					return nil
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
