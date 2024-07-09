package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"context"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func shortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}

func main() {
	var eigenpodAddress, beacon, node, owner, output string
	var forceCheckpoint, disableColor, verbose bool
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

					eth, err := ethclient.Dial(node)
					core.PanicOnError("failed to reach eth --node.", err)

					beaconClient, err := core.GetBeaconClient(beacon)
					core.PanicOnError("failed to reach beacon chain.", err)

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
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "If true, starts a checkpoint even if the pod has no native ETH to award shares",
						Destination: &forceCheckpoint,
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

					eth, beaconClient, chainId := core.GetClients(ctx, node, beacon)

					currentCheckpoint := core.GetCurrentCheckpoint(eigenpodAddress, eth)
					if currentCheckpoint == 0 {
						if owner != nil {
							newCheckpoint, err := core.StartCheckpoint(ctx, eigenpodAddress, *owner, chainId, eth, forceCheckpoint)
							core.PanicOnError("failed to start checkpoint", err)
							currentCheckpoint = newCheckpoint
						} else {
							core.PanicOnError("no checkpoint active and no private key provided to start one", errors.New("no checkpoint"))
						}
					}
					color.Green("pod has active checkpoint! checkpoint timestamp: %d", currentCheckpoint)

					proof := core.GenerateCheckpointProof(ctx, eigenpodAddress, eth, chainId, beaconClient)

					jsonString, err := json.Marshal(proof)
					core.PanicOnError("failed to generate JSON proof data.", err)

					core.WriteOutputToFileOrStdout(jsonString, out)

					if owner != nil {
						// submit the proof onchain
						ownerAccount, err := core.PrepareAccount(owner, chainId)
						core.PanicOnError("failed to parse private key", err)

						eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
						core.PanicOnError("failed to reach eigenpod", err)

						color.Green("calling EigenPod.VerifyCheckpointProofs()...")
						txn, err := eigenPod.VerifyCheckpointProofs(
							ownerAccount.TransactionOptions,
							onchain.BeaconChainProofsBalanceContainerProof{
								BalanceContainerRoot: proof.ValidatorBalancesRootProof.ValidatorBalancesRoot,
								Proof:                proof.ValidatorBalancesRootProof.Proof.ToByteSlice(),
							},
							core.CastBalanceProofs(proof.BalanceProofs),
						)

						core.PanicOnError("failed to invoke verifyCheckpointProofs", err)
						color.Green("transaction: %s", txn.Hash().Hex())
					}

					return nil
				},
			},
			{
				Name:    "credentials",
				Aliases: []string{"cr", "creds"},
				Usage:   "Generates a proof for use with EigenPod.verifyWithdrawalCredentials()",
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

					eth, beaconClient, chainId := core.GetClients(ctx, node, beacon)
					validatorProofs, validatorIndices := core.GenerateValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient)
					if validatorProofs == nil || validatorIndices == nil {
						return nil
					}

					jsonString, err := json.Marshal(validatorProofs)
					core.PanicOnError("failed to generate JSON proof data.", err)

					core.WriteOutputToFileOrStdout(jsonString, out)

					if owner != nil {
						ownerAccount, err := core.PrepareAccount(owner, chainId)
						core.PanicOnError("failed to parse private key", err)

						eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
						core.PanicOnError("failed to reach eigenpod", err)

						indices := core.Uint64ArrayToBigIntArray(validatorIndices)

						var validatorFieldsProofs [][]byte = [][]byte{}
						for i := 0; i < len(validatorProofs.ValidatorFieldsProofs); i++ {
							pr := validatorProofs.ValidatorFieldsProofs[i].ToByteSlice()
							validatorFieldsProofs = append(validatorFieldsProofs, pr)
						}

						var validatorFields [][][32]byte = core.CastValidatorFields(validatorProofs.ValidatorFields)

						latestBlock, err := eth.BlockByNumber(ctx, nil)
						core.PanicOnError("failed to load latest block", err)

						color.Green("submitting onchain...")
						txn, err := eigenPod.VerifyWithdrawalCredentials(
							ownerAccount.TransactionOptions,
							latestBlock.Time(),
							onchain.BeaconChainProofsStateRootProof{
								Proof:           validatorProofs.StateRootProof.Proof.ToByteSlice(),
								BeaconStateRoot: validatorProofs.StateRootProof.BeaconStateRoot,
							},
							indices,
							validatorFieldsProofs,
							validatorFields,
						)

						core.PanicOnError("failed to invoke verifyWithdrawalCredentials", err)

						color.Green("transaction: %s", txn.Hash().Hex())
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
				Name:        "out",
				Aliases:     []string{"O", "output"},
				Value:       "",
				Usage:       "Output `path` for the proof. (defaults to stdout)",
				Destination: &output,
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
