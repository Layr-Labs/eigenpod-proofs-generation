package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func RunCheckpointProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient, out, owner *string) {
	lastCheckpoint := lastCheckpointedForEigenpod(eigenpodAddress, eth)
	if lastCheckpoint == 0 {
		if owner != nil {
			err := startCheckpoint(eigenpodAddress, *owner, chainId, eth)
			PanicOnError("failed to start checkpoint", err)
		} else {
			PanicOnError("no checkpoint active", errors.New("no checkpoint"))
		}
	}

	blockRoot, err := getCurrentCheckpointBlockRoot(eigenpodAddress, eth)
	PanicOnError("failed to fetch last checkpoint.", err)

	if blockRoot != nil {
		rootBytes := *blockRoot
		if AllZero(rootBytes[:]) {
			PanicOnError("No checkpoint active. Are you sure you started a checkpoint?", fmt.Errorf("no checkpoint"))
		}
	}

	headerBlock := "0x" + hex.EncodeToString((*blockRoot)[:])
	header, err := beaconClient.GetBeaconHeader(ctx, headerBlock)
	PanicOnError(fmt.Sprintf("failed to fetch beacon header (%s).", headerBlock), err)

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	PanicOnError("failed to fetch beacon state.", err)

	// filter through the beaconState's validators, and select only ones that have withdrawal address set to `eigenpod`.
	allValidatorsForEigenpod := findAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	allValidatorInfo := getOnchainValidatorInfo(eth, eigenpodAddress, allValidatorsForEigenpod)

	// for each validator, request RPC information from the eigenpod (using the pubKeyHash), and;
	//			- we want all un-checkpointed, non-withdrawn validators that belong to this eigenpoint.
	//			- determine the validator's index.
	var checkpointValidatorIndices = FilterNotCheckpointedOrWithdrawnValidators(allValidatorsForEigenpod, allValidatorInfo, lastCheckpoint)

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	PanicOnError("failled to initialize prover", err)

	proof, err := proofs.ProveCheckpointProofs(header.Header.Message, beaconState, checkpointValidatorIndices)
	PanicOnError("failed to prove checkpoint.", err)

	jsonString, err := json.Marshal(proof)
	PanicOnError("failed to generate JSON proof data.", err)

	WriteOutputToFileOrStdout(jsonString, out)

	if owner != nil {
		// submit the proof onchain
		ownerAccount, err := prepareAccount(owner, chainId)
		PanicOnError("failed to parse private key", err)

		eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
		PanicOnError("failed to reach eigenpod", err)

		color.Green("submitting onchain...")

		txn, err := eigenPod.VerifyCheckpointProofs(
			ownerAccount.TransactionOptions,
			onchain.BeaconChainProofsBalanceContainerProof{
				BalanceContainerRoot: [32]byte(proof.ValidatorBalancesRootProof.Proof.ToByteSlice()),
				Proof:                proof.ValidatorBalancesRootProof.Proof[0][:],
			},
			castBalanceProofs(proof.BalanceProofs),
		)

		PanicOnError("failed to invoke verifyCheckpointProofs", err)
		color.Green("transaction: %s", txn.Hash().Hex())
	}
}
