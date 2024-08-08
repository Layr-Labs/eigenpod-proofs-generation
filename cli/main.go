package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"time"

	"context"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
		Usage:       "Submit proofs in groups of size `batchSize`, to avoid gas limit.",
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
var specificValidator uint64 = math.MaxUint64
var slashedValidatorIndex uint64 = math.MaxUint64
var proofPath string
var noPrompt bool

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

var NO_PROMPT_FLAG = &cli.BoolFlag{
	Name:        "no-prompt",
	Value:       false,
	Usage:       "Disables prompts to approve any transactions occurring (e.g in CI).",
	Destination: &noPrompt,
}

// Required for commands that need an execution layer RPC
var SLASHED_VALIDATOR_INDEX_FLAG = &cli.Uint64Flag{
	Name:        "slashedValidatorIndex",
	Value:       0,
	Usage:       "[required] The index of a validator belonging to this pod that was slashed.",
	Required:    true,
	Destination: &slashedValidatorIndex,
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

var PROOF_PATH_FLAG = &cli.StringFlag{
	Name:        "proof",
	Value:       "",
	Usage:       "the `path` to a previous proof generated from this step (via -o proof.json). If provided, this proof will submitted to network via the --sender flag.",
	Destination: &proofPath,
}

// maximum number of proofs per txn for each of the following proof types:
const DEFAULT_BATCH_CREDENTIALS = 60
const DEFAULT_BATCH_CHECKPOINT = 80

func proofCast(proof []eigenpodproofs.Bytes32) [][32]byte {
	result := make([][32]byte, len(proof))
	for i, p := range proof {
		result[i] = p
	}
	return result
}

func main() {
	var batchSize uint64
	var forceCheckpoint, disableColor, verbose bool
	ctx := context.Background()

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator.",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "stale-balance",
				Args:  true,
				Usage: "If needed, calls `verifyStaleBalances` to correct the balance on a pod you don't own. This will attempt to (1) conclude any existing checkpoint on the pod, and then (2) invoke verifyStaleBalances() to start another, more up-to-date, checkpoint.",
				Flags: []cli.Flag{
					POD_ADDRESS_FLAG,
					EXEC_NODE_FLAG,
					BEACON_NODE_FLAG,
					Require(SENDER_PK_FLAG),
					SLASHED_VALIDATOR_INDEX_FLAG,
					NO_PROMPT_FLAG,
					BatchBySize(&batchSize, DEFAULT_BATCH_CHECKPOINT),
				},
				Action: func(cctx *cli.Context) error {
					ctx := context.Background()

					eth, beacon, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to get clients", err)

					ownerAccount, err := core.PrepareAccount(&sender, chainId)
					core.PanicOnError("failed to parse sender PK", err)

					eigenpod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
					core.PanicOnError("failed to reach eigenpod", err)

					currentCheckpointTimestamp, err := eigenpod.CurrentCheckpointTimestamp(nil)
					core.PanicOnError("failed to fetch any existing checkpoint info", err)

					if currentCheckpointTimestamp > 0 {
						// TODO: complete current checkpoint
						fmt.Printf("This eigenpod has an outstanding checkpoint (since %d). You must complete it before continuing.", currentCheckpointTimestamp)

						proofs, err := core.GenerateCheckpointProof(ctx, eigenpodAddress, eth, chainId, beacon)
						core.PanicOnError("failed to generate checkpoint proofs", err)

						txns, err := core.SubmitCheckpointProof(ctx, sender, eigenpodAddress, chainId, proofs, eth, batchSize, noPrompt)
						core.PanicOnError("failed to submit checkpoint proofs", err)

						for i, txn := range txns {
							fmt.Printf("sending txn[%d/%d]: %s (waiting)...", i, len(txns), txn.Hash())
							bind.WaitMined(ctx, eth, txn)
						}
					}

					proof, oracleBeaconTimesetamp, err := core.GenerateValidatorProof(ctx, eigenpodAddress, eth, chainId, beacon, new(big.Int).SetUint64(slashedValidatorIndex))
					core.PanicOnError("failed to generate credential proof for slashed validator", err)

					if !noPrompt {
						core.PanicIfNoConsent("This will invoke `EigenPod.verifyStaleBalance()` on the given eigenpod, which will start a checkpoint. Once started, this checkpoint must be completed.")
					}

					txn, err := eigenpod.VerifyStaleBalance(
						ownerAccount.TransactionOptions,
						oracleBeaconTimesetamp,
						onchain.BeaconChainProofsStateRootProof{
							Proof:           proof.StateRootProof.Proof.ToByteSlice(),
							BeaconStateRoot: proof.StateRootProof.BeaconStateRoot,
						},
						onchain.BeaconChainProofsValidatorProof{
							ValidatorFields: proofCast(proof.ValidatorFields[0]),
							Proof:           proof.ValidatorFieldsProofs[0].ToByteSlice(),
						},
					)
					core.PanicOnError("failed to call verifyStaleBalance()", err)

					fmt.Printf("txn: %s\n", txn.Hash())

					return nil
				},
			},
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
					} else if !common.IsHexAddress(targetAddress) {
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

					pod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
					if err != nil {
						return fmt.Errorf("error contacting eigenpod: %w", err)
					}

					// Check that the existing submitter is not the current submitter
					newSubmitter := common.HexToAddress(targetAddress)
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
					PROOF_PATH_FLAG,
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "If true, starts a checkpoint even if the pod has no native ETH to award shares",
						Destination: &forceCheckpoint,
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

					if len(proofPath) > 0 {
						// user specified the proof
						if len(sender) == 0 {
							core.Panic("If using --proof, --sender <privateKey> must also be supplied.")
						}

						// load `proof` from file.
						proof, err := core.LoadCheckpointProofFromFile(proofPath)
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
					&cli.Uint64Flag{
						Name:        "validatorIndex",
						Usage:       "The `index` of a specific validator to prove (e.g a slashed validator for `verifyStaleBalance()`).",
						Destination: &specificValidator,
					},
					PROOF_PATH_FLAG,
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

					eth, beaconClient, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to reach ethereum clients", err)

					var specificValidatorIndex *big.Int = nil
					if specificValidator != math.MaxUint64 {
						specificValidatorIndex = new(big.Int).SetUint64(specificValidator)
					}

					if len(proofPath) > 0 {
						if len(sender) == 0 {
							core.Panic("If using --proof, --sender <privateKey> must also be supplied.")
						}

						proof, err := core.LoadValidatorProofFromFile(proofPath)
						core.PanicOnError("failed to parse checkpoint proof from file", err)

						txns, err := core.SubmitValidatorProof(ctx, sender, eigenpodAddress, chainId, eth, batchSize, proof.ValidatorProofs, proof.OracleBeaconTimestamp, noPrompt)
						for _, txn := range txns {
							color.Green("submitted txn: %s", txn.Hash())
						}
						core.PanicOnError("an error occurred while submitting your credential proofs", err)
						return nil
					}

					validatorProofs, oracleBeaconTimestamp, err := core.GenerateValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient, specificValidatorIndex)

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
						proof := core.SerializableCredentialProof{
							ValidatorProofs:       validatorProofs,
							OracleBeaconTimestamp: oracleBeaconTimestamp,
						}
						out, err := json.MarshalIndent(proof, "", "   ")
						core.PanicOnError("failed to process proof", err)

						core.WriteOutputToFileOrStdout(out, &output)
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
			NO_PROMPT_FLAG,
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
