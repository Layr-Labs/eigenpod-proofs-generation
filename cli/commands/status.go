package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/fatih/color"
)

type TStatusArgs struct {
	EigenpodAddress string
	DisableColor    bool
	UseJSON         bool
	Node            string
	BeaconNode      string
	Verbose         bool
}

func StatusCommand(args TStatusArgs) error {
	ctx := context.Background()
	if args.DisableColor {
		color.NoColor = true
	}

	isVerbose := !args.UseJSON

	eth, beaconClient, _, err := core.GetClients(ctx, args.Node, args.BeaconNode, isVerbose)
	core.PanicOnError("failed to load ethereum clients", err)

	status := core.GetStatus(ctx, args.EigenpodAddress, eth, beaconClient)

	if args.UseJSON {
		bytes, err := json.MarshalIndent(status, "", "      ")
		core.PanicOnError("failed to get status", err)
		statusStr := string(bytes)
		fmt.Println(statusStr)
		return nil
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
		awaitingActivationQueueValidators, inactiveValidators, activeValidators, withdrawnValidators :=
			core.SortByStatus(status.Validators)
		var targetColor *color.Color

		bold.Printf("Eigenpod validators:\n============\n")
		ital.Printf("Format: #ValidatorIndex (pubkey) [effective balance] [current balance]\n")

		// print info on validators who are not yet in the activation queue
		//
		// if these validators have 32 ETH effective balance, they will be
		// activated soon and can then have their credentials verified
		//
		// if these validators do NOT have 32 ETH effective balance yet, the
		// staker needs to deposit more ETH.
		if len(awaitingActivationQueueValidators) != 0 {
			color.New(color.Bold, color.FgHiRed).Printf("- [AWAITING ACTIVATION QUEUE] - These validators have deposited, but either do not meet the minimum balance to be activated, or are awaiting activation:\n")

			for _, validator := range awaitingActivationQueueValidators {
				publicKey := validator.PublicKey
				if !isVerbose {
					publicKey = utils.ShortenHex(publicKey)
				}

				targetColor = color.New(color.FgHiRed)
				if validator.Slashed {
					targetColor.Printf("\t- #%d (%s) [%d] [%d] (slashed on beacon chain)\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
				} else {
					targetColor.Printf("\t- #%d (%s) [%d] [%d]\n", validator.Index, publicKey, validator.EffectiveBalance, validator.CurrentBalance)
				}
			}

			fmt.Println()
		}

		// print info on inactive validators
		// these validators can be added to the pod's active validator set
		// by running the `credentials` command
		if len(inactiveValidators) != 0 {
			targetColor = color.New(color.FgHiYellow)
			color.New(color.Bold, color.FgHiYellow).Printf("- [INACTIVE] - Run `credentials` to verify these %d validators' withdrawal credentials:\n", len(inactiveValidators))
			for _, validator := range inactiveValidators {
				publicKey := validator.PublicKey
				if !isVerbose {
					publicKey = utils.ShortenHex(publicKey)
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
				if !isVerbose {
					publicKey = utils.ShortenHex(publicKey)
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
				if !isVerbose {
					publicKey = utils.ShortenHex(publicKey)
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

			bold.Printf("Batching %d proofs per txn, this will require:\n\t", utils.DEFAULT_BATCH_CHECKPOINT)
			ital.Printf("- 1x startCheckpoint() transaction, and \n\t- %dx EigenPod.verifyCheckpointProofs() transaction(s)\n\n", int(math.Ceil(float64(status.NumberValidatorsToCheckpoint)/float64(utils.DEFAULT_BATCH_CHECKPOINT))))
		}
	}
	return nil
}
