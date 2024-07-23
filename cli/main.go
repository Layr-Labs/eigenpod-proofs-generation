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

// maximum number of proofs per txn for each of the following proof types:
const DEFAULT_BATCH_CREDENTIALS = 60
const DEFAULT_BATCH_CHECKPOINT = 80

func main() {
	var eigenpodAddress, beacon, node, owner, output string
	var batchSize uint64
	var checkpointProofPath string
	var forceCheckpoint, disableColor, verbose bool
	var noPrompt bool
	var useJson bool = false
	ctx := context.Background()

	app := &cli.App{
		Name:                   "Eigenlayer Proofs CLi",
		HelpName:               "eigenproofs",
		Usage:                  "Generates proofs to (1) checkpoint your validators, or (2) verify the withdrawal credentials of an inactive validator.",
		EnableBashCompletion:   true,
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "assign-submitter",
				Args:  true,
				Usage: "Assign a different address to be able to submit your proofs. You'll always be able to submit from your EigenPod owner PK.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "json",
						Value:       false,
						Usage:       "print only plain JSON",
						Required:    false,
						Destination: &useJson,
					},
				},
				Action: func(cctx *cli.Context) error {
					targetAddress := cctx.Args().First()
					if len(targetAddress) == 0 {
						return fmt.Errorf("usage: `assign-submitter <0xsubmitter>")
					}

					eth, err := ethclient.Dial(node)
					if err != nil {
						return fmt.Errorf("failed to reach eth --node: %w", err)
					}

					chainId, err := eth.ChainID(ctx)
					if err != nil {
						return fmt.Errorf("failed to reach eth node for chain id: %w", err)
					}

					ownerAccount, err := core.PrepareAccount(&owner, chainId)
					if err != nil {
						return fmt.Errorf("failed to parse --owner: %w", err)
					}

					pod, err := onchain.NewEigenPod(gethCommon.HexToAddress(eigenpodAddress), eth)
					if err != nil {
						return fmt.Errorf("error contacting eigenpod: %w", err)
					}

					if !noPrompt {
						core.PanicIfNoConsent(fmt.Sprintf("This will update your EigenPod to allow %s to submit proofs on its behalf. You can always change this later using the EigenPod's PK.", targetAddress))
					}

					txn, err := pod.SetProofSubmitter(ownerAccount.TransactionOptions, gethCommon.HexToAddress(targetAddress))
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
					&cli.BoolFlag{
						Name:        "json",
						Value:       false,
						Usage:       "print only plain JSON",
						Required:    false,
						Destination: &useJson,
					},
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
						// pretty print everything
						color.New(color.Bold, color.FgBlue).Printf("Eigenpod validators\n")
						for index, validator := range status.Validators {

							var targetColor color.Attribute
							var description string

							if validator.Status == core.ValidatorStatusActive {
								targetColor = color.FgGreen
								description = "active"
							} else if validator.Status == core.ValidatorStatusInactive {
								targetColor = color.FgHiYellow
								description = "inactive"
							} else if validator.Status == core.ValidatorStatusWithdrawn {
								targetColor = color.FgHiRed
								description = "withdrawn"
							}

							if validator.Slashed {
								description = description + " (slashed)"
							}

							publicKey := validator.PublicKey
							if !verbose {
								publicKey = shortenHex(publicKey)
							}

							color.New(targetColor).Printf("\t- #%s (%s) [%s]\n", index, publicKey, description)
						}

						bold := color.New(color.Bold, color.FgBlue)
						ital := color.New(color.Italic, color.FgBlue)
						fmt.Println()

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

							bold.Printf("This will require:\n\t")
							ital.Printf("- 1x startCheckpoint() transaction, and \n\t- %dx EigenPod.verifyCheckpointProofs() transaction(s)\n\n", int(math.Ceil(float64(status.NumberValidatorsToCheckpoint)/float64(80))))
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
						Usage:       "the path to a previous proof generated from this step (via `-o proof.json`). If provided, this proof will submitted to network via the `--owner` flag.",
						Destination: &checkpointProofPath,
					},
					&cli.StringFlag{
						Name:        "out",
						Aliases:     []string{"O", "output"},
						Value:       "",
						Usage:       "Output `path` for the proof. (defaults to stdout). NOTE: If `--out` is supplied along with `--owner`, `--out` takes precedence and the proof will not be broadcast.",
						Destination: &output,
					},
				},
				Action: func(cctx *cli.Context) error {
					if disableColor {
						color.NoColor = true
					}

					var out, owner *string = nil, nil

					if len(cctx.String("out")) > 0 {
						outProp := cctx.String("out")
						out = &outProp
					}

					if len(cctx.String("owner")) > 0 {
						ownerProp := cctx.String("owner")
						owner = &ownerProp
					}

					eth, beaconClient, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to reach ethereum clients", err)

					if len(checkpointProofPath) > 0 {
						// user specified the proof
						if owner == nil || len(*owner) == 0 {
							core.Panic("If using --proof, --owner <privateKey> must also be supplied.")
						}

						// load `proof` from file.
						proof, err := core.LoadCheckpointProofFromFile(checkpointProofPath)
						core.PanicOnError("failed to parse checkpoint proof from file", err)

						txns, err := core.SubmitCheckpointProof(ctx, *owner, eigenpodAddress, chainId, proof, eth, batchSize, noPrompt)
						for _, txn := range txns {
							color.Green("submitted txn: %s", txn.Hash())
						}
						core.PanicOnError("an error occurred while submitting your checkpoint proofs", err)
						return nil
					}

					currentCheckpoint, err := core.GetCurrentCheckpoint(eigenpodAddress, eth)
					core.PanicOnError("failed to load checkpoint", err)

					if currentCheckpoint == 0 {
						if owner != nil {
							if !noPrompt {
								core.PanicIfNoConsent(core.StartCheckpointProofConsent())
							}

							newCheckpoint, err := core.StartCheckpoint(ctx, eigenpodAddress, *owner, chainId, eth, forceCheckpoint)
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
					} else if owner != nil {
						txns, err := core.SubmitCheckpointProof(ctx, *owner, eigenpodAddress, chainId, proof, eth, batchSize, noPrompt)
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
					BatchBySize(&batchSize, DEFAULT_BATCH_CREDENTIALS),
				},
				Action: func(cctx *cli.Context) error {
					if disableColor {
						color.NoColor = true
					}

					var owner *string = nil
					if len(cctx.String("owner")) > 0 {
						ownerProp := cctx.String("owner")
						owner = &ownerProp
					}

					eth, beaconClient, chainId, err := core.GetClients(ctx, node, beacon)
					core.PanicOnError("failed to reach ethereum clients", err)

					validatorProofs, oracleBeaconTimestamp, err := core.GenerateValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient)
					if err != nil || validatorProofs == nil {
						core.PanicOnError("Failed to generate validator proof", err)
						core.Panic("no inactive validators")
					}

					if owner != nil {
						txns, err := core.SubmitValidatorProof(ctx, *owner, eigenpodAddress, chainId, eth, batchSize, validatorProofs, oracleBeaconTimestamp, noPrompt)
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
			&cli.StringFlag{
				Name:        "podAddress",
				Aliases:     []string{"p", "pod"},
				Value:       "",
				Usage:       "[required] The onchain `address` of your eigenpod contract (0x123123123123)",
				Required:    true,
				Destination: &eigenpodAddress,
			},
			&cli.StringFlag{
				Name:        "beaconNode",
				Aliases:     []string{"b"},
				Value:       "",
				Usage:       "[required] `URL` to a functioning beacon node RPC (https://)",
				Required:    true,
				Destination: &beacon,
			},
			&cli.StringFlag{
				Name:        "execNode",
				Aliases:     []string{"e"},
				Value:       "",
				Usage:       "[required] `URL` to a functioning execution-layer RPC (https://)",
				Required:    true,
				Destination: &node,
			},
			&cli.StringFlag{
				Name:        "owner",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "`Private key` of the owner. If set, this will automatically submit the proofs to their corresponding onchain functions after generation. If using checkpoint mode, it will also begin a checkpoint if one hasn't been started already.",
				Destination: &owner,
			},
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
