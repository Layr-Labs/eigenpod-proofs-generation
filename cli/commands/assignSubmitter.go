package commands

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

type TAssignSubmitterArgs struct {
	Node            string
	TargetAddress   string
	Sender          string
	EigenpodAddress string
	NoPrompt        bool
	Verbose         bool
}

func AssignSubmitterCommand(args TAssignSubmitterArgs) error {
	ctx := context.Background()

	if len(args.TargetAddress) == 0 {
		return fmt.Errorf("usage: `assign-submitter <0xsubmitter>`")
	} else if !common.IsHexAddress(args.TargetAddress) {
		return fmt.Errorf("invalid address for 0xsubmitter: %s", args.TargetAddress)
	}

	eth, err := ethclient.Dial(args.Node)
	if err != nil {
		return fmt.Errorf("failed to reach eth --node: %w", err)
	}

	chainId, err := eth.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to reach eth node for chain id: %w", err)
	}

	ownerAccount, err := utils.PrepareAccount(&args.Sender, chainId, false /* noSend */)
	if err != nil {
		return fmt.Errorf("failed to parse --sender: %w", err)
	}

	pod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	if err != nil {
		return fmt.Errorf("error contacting eigenpod: %w", err)
	}

	// Check that the existing submitter is not the current submitter
	newSubmitter := common.HexToAddress(args.TargetAddress)
	currentSubmitter, err := pod.ProofSubmitter(nil)
	if err != nil {
		return fmt.Errorf("error fetching current proof submitter: %w", err)
	} else if currentSubmitter.Cmp(newSubmitter) == 0 {
		return fmt.Errorf("error: new proof submitter is existing proof submitter (%s)", currentSubmitter)
	}

	if !args.NoPrompt {
		fmt.Printf("Your pod's current proof submitter is %s.\n", currentSubmitter)
		utils.PanicIfNoConsent(fmt.Sprintf("This will update your EigenPod to allow %s to submit proofs on its behalf. As the EigenPod's owner, you can always change this later.", newSubmitter))
	}

	txn, err := pod.SetProofSubmitter(ownerAccount.TransactionOptions, newSubmitter)
	if err != nil {
		return fmt.Errorf("error updating submitter role: %w", err)
	}

	color.Green("submitted txn: %s", txn.Hash())
	color.Green("updated!")

	return nil
}
