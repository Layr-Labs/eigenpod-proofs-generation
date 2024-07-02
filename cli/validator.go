package main

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func RunValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient, out *string, owner *string) {
	// TODO: where does this come from (latest non-missed-slot block-header)
	oracleBump := "0xdec64b9ad990d457dbf465ee4b9f8ab70a70f6d56e7b6b8b472320a34c7e7b28"

	header, err := beaconClient.GetBeaconHeader(ctx, oracleBump)
	PanicOnError("failed to fetch beacon header.", err)

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	PanicOnError("failed to fetch beacon state.", err)

	allValidatorsForEigenpod := findAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	allValidatorInfo := getOnchainValidatorInfo(eth, eigenpodAddress, allValidatorsForEigenpod)

	var validatorIndices = FilterInactiveValidators(allValidatorsForEigenpod, allValidatorInfo)

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	PanicOnError("failled to initialize prover", err)

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, validatorIndices)
	PanicOnError("failed to prove validators.", err)

	jsonString, err := json.Marshal(validatorProofs)
	PanicOnError("failed to generate JSON proof data.", err)

	WriteOutputToFileOrStdout(jsonString, out)

	if owner != nil {
		ownerAccount, err := prepareAccount(owner, chainId)
		PanicOnError("failed to parse private key", err)

		eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
		PanicOnError("failed to reach eigenpod", err)

		indices := Uint64ArrayToBigIntArray(validatorIndices)

		var validatorFieldsProofs [][]byte = [][]byte{}
		for i := 0; i < len(validatorProofs.ValidatorFieldsProofs); i++ {
			pr := validatorProofs.ValidatorFieldsProofs[i].ToByteSlice()
			validatorFieldsProofs = append(validatorFieldsProofs, pr)
		}

		var validatorFields [][][32]byte = castValidatorFields(validatorProofs.ValidatorFields)

		// TODO: where does this come from?
		oracleTimestamp := 1719941892

		color.Green("submitting onchain...")
		txn, err := eigenPod.VerifyWithdrawalCredentials(
			ownerAccount.TransactionOptions,
			uint64(oracleTimestamp), // TODO: timestamp
			onchain.BeaconChainProofsStateRootProof{
				Proof:           validatorProofs.StateRootProof.Proof.ToByteSlice(),
				BeaconStateRoot: validatorProofs.StateRootProof.BeaconStateRoot,
			},
			indices,
			validatorFieldsProofs,
			validatorFields,
		)

		PanicOnError("failed to invoke verifyWithdrawalCredentials", err)

		color.Green("transaction: %s", txn.Hash().Hex())
	}
}
