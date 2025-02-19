package commands

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
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

type TransactionDescription struct {
	Hash string
	Type string
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

	sentTxns := []TransactionDescription{}

	eth, beacon, chainId, err := utils.GetClients(ctx, args.EthNode, args.BeaconNode, args.Verbose)
	utils.PanicOnError("failed to get clients", err)

	validator, err := beacon.GetValidator(ctx, args.SlashedValidatorIndex)
	utils.PanicOnError("failed to fetch validator state", err)

	if !validator.Validator.Slashed {
		utils.Panic("Provided validator was not slashed.")
		return nil
	}

	ownerAccount, err := utils.PrepareAccount(&args.Sender, chainId, false /* noSend */)
	utils.PanicOnError("failed to parse sender PK", err)

	eigenpod, err := EigenPod.NewEigenPod(common.HexToAddress(args.EigenpodAddress), eth)
	utils.PanicOnError("failed to reach eigenpod", err)

	currentCheckpointTimestamp, err := eigenpod.CurrentCheckpointTimestamp(nil)
	utils.PanicOnError("failed to fetch any existing checkpoint info", err)

	if currentCheckpointTimestamp > 0 {
		if !args.NoPrompt {
			utils.PanicIfNoConsent(fmt.Sprintf("This eigenpod has an outstanding checkpoint (since %d). You must complete it before continuing. This will invoke `EigenPod.verifyCheckpointProofs()`, which will end the checkpoint. This may be expensive.", currentCheckpointTimestamp))
		}

		proofs, err := core.GenerateCheckpointProof(ctx, args.EigenpodAddress, eth, chainId, beacon, args.Verbose)
		utils.PanicOnError("failed to generate checkpoint proofs", err)

		txns, err := core.SubmitCheckpointProof(ctx, args.Sender, args.EigenpodAddress, chainId, proofs, eth, args.CheckpointBatchSize, args.NoPrompt, false /* noSend */, args.Verbose)
		utils.PanicOnError("failed to submit checkpoint proofs", err)

		for i, txn := range txns {
			if args.Verbose {
				fmt.Printf("sending txn[%d/%d]: %s (waiting)...", i, len(txns), txn.Hash())
			}
			bind.WaitMined(ctx, eth, txn)
			sentTxns = append(sentTxns, TransactionDescription{Type: "complete_existing_checkpoint", Hash: txn.Hash().Hex()})
		}
	}

	proof, oracleBeaconTimesetamp, err := core.GenerateValidatorProof(ctx, args.EigenpodAddress, eth, chainId, beacon, new(big.Int).SetUint64(args.SlashedValidatorIndex), args.Verbose)
	utils.PanicOnError("failed to generate credential proof for slashed validator", err)

	if !args.NoPrompt {
		utils.PanicIfNoConsent("This will invoke `EigenPod.verifyStaleBalance()` on the given eigenpod, which will start a checkpoint. Once started, this checkpoint must be completed.")
	}

	if args.Verbose {
		color.Black("Calling EigenPod.verifyStaleBalance() to force checkpoint this pod.")
	}

	txn, err := eigenpod.VerifyStaleBalance(
		ownerAccount.TransactionOptions,
		oracleBeaconTimesetamp,
		EigenPod.BeaconChainProofsStateRootProof{
			Proof:           proof.StateRootProof.Proof.ToByteSlice(),
			BeaconStateRoot: proof.StateRootProof.BeaconStateRoot,
		},
		EigenPod.BeaconChainProofsValidatorProof{
			ValidatorFields: proofCast(proof.ValidatorFields[0]),
			Proof:           proof.ValidatorFieldsProofs[0].ToByteSlice(),
		},
	)
	utils.PanicOnError("failed to call verifyStaleBalance()", err)
	sentTxns = append(sentTxns, TransactionDescription{Type: "verify_stale_balance", Hash: txn.Hash().Hex()})

	PrintAsJSON(sentTxns)
	return nil
}
