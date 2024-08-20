package commands

import (
	"context"
	"fmt"
	"math/big"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fatih/color"
)

type TFixStaleBalanceArgs struct {
	EthNode               string
	BeaconNode            string
	Sender                string
	EigenpodAddress       string
	SlashedValidatorIndex uint64
	Verbose               bool
	CheckpointBatchSize   uint64
	NoPrompt              bool
}

// another fun cast brought to you by golang!
func proofCast(proof []eigenpodproofs.Bytes32) [][32]byte {
	res := make([][32]byte, len(proof))
	for i, elt := range proof {
		res[i] = elt
	}
	return res
}

func FixStaleBalance(args TFixStaleBalanceArgs) error {
	ctx := context.Background()

	eth, beacon, chainId, err := core.GetClients(ctx, args.EthNode, args.BeaconNode, args.Verbose)
	core.PanicOnError("failed to get clients", err)

	validator, err := beacon.GetValidator(ctx, args.SlashedValidatorIndex)
	core.PanicOnError("failed to fetch validator state", err)

	if !validator.Validator.Slashed {
		core.Panic("Provided validator was not slashed.")
		return nil
	}

	ownerAccount, err := core.PrepareAccount(&args.Sender, chainId, false /* noSend */)
	core.PanicOnError("failed to parse sender PK", err)

	eigenpod, err := onchain.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	core.PanicOnError("failed to reach eigenpod", err)

	currentCheckpointTimestamp, err := eigenpod.CurrentCheckpointTimestamp(nil)
	core.PanicOnError("failed to fetch any existing checkpoint info", err)

	if currentCheckpointTimestamp > 0 {
		if args.Verbose {
			color.Red("This eigenpod has an outstanding checkpoint (since %d). You must complete it before continuing.", currentCheckpointTimestamp)
		}

		proofs, err := core.GenerateCheckpointProof(ctx, args.EigenpodAddress, eth, chainId, beacon)
		core.PanicOnError("failed to generate checkpoint proofs", err)

		txns, err := core.SubmitCheckpointProof(ctx, args.Sender, args.EigenpodAddress, chainId, proofs, eth, args.CheckpointBatchSize, args.NoPrompt /* noSend */, false)
		core.PanicOnError("failed to submit checkpoint proofs", err)

		for i, txn := range txns {
			if args.Verbose {
				fmt.Printf("sending txn[%d/%d]: %s (waiting)...", i, len(txns), txn.Hash())
			}
			bind.WaitMined(ctx, eth, txn)
		}
	}

	proof, oracleBeaconTimesetamp, err := core.GenerateValidatorProof(ctx, args.EigenpodAddress, eth, chainId, beacon, new(big.Int).SetUint64(args.SlashedValidatorIndex), args.Verbose)
	core.PanicOnError("failed to generate credential proof for slashed validator", err)

	if !args.NoPrompt {
		core.PanicIfNoConsent("This will invoke `EigenPod.verifyStaleBalance()` on the given eigenpod, which will start a checkpoint. Once started, this checkpoint must be completed.")
	}

	if args.Verbose {
		color.Black("Calling verifyStaleBalance() to update pod.")
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
}
