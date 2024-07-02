package main

import (
	"bytes"
	"context"
	sha256 "crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
	"github.com/pkg/errors"
)

func execute(ctx context.Context, eigenpodAddress, beacon_node_uri, node, command string, out *string, owner *string) {
	eth, err := ethclient.Dial(node)
	PanicOnError("failed to reach eth --node.", err)

	chainId, err := eth.ChainID(ctx)
	PanicOnError("failed to fetch chain id", err)

	beaconClient, err := getBeaconClient(beacon_node_uri)
	PanicOnError("failed to reach beacon chain.", err)

	if command == "checkpoint" {
		proof := RunCheckpointProof(ctx, eigenpodAddress, eth, chainId, beaconClient, owner)

		jsonString, err := json.Marshal(proof)
		PanicOnError("failed to generate JSON proof data.", err)

		WriteOutputToFileOrStdout(jsonString, out)

		if owner != nil {
			// submit the proof onchain
			ownerAccount, err := prepareAccount(owner, chainId)
			PanicOnError("failed to parse private key", err)

			eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
			PanicOnError("failed to reach eigenpod", err)

			color.Green("calling EigenPod.VerifyCheckpointProofs()...")

			txn, err := eigenPod.VerifyCheckpointProofs(
				ownerAccount.TransactionOptions,
				onchain.BeaconChainProofsBalanceContainerProof{
					BalanceContainerRoot: proof.ValidatorBalancesRootProof.ValidatorBalancesRoot,
					Proof:                proof.ValidatorBalancesRootProof.Proof.ToByteSlice(),
				},
				castBalanceProofs(proof.BalanceProofs),
			)

			PanicOnError("failed to invoke verifyCheckpointProofs", err)
			color.Green("transaction: %s", txn.Hash().Hex())
		}

	} else if command == "validator" {
		validatorProofs, validatorIndices, latestBlock := RunValidatorProof(ctx, eigenpodAddress, eth, chainId, beaconClient, owner)

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

			PanicOnError("failed to invoke verifyWithdrawalCredentials", err)

			color.Green("transaction: %s", txn.Hash().Hex())
		}
	} else {
		PanicOnError(fmt.Sprintf("invalid --prove argument. Expected 'checkpoint' or 'validator' (got `%s`)", command), errors.New("invalid command"))
	}
}

func getBeaconClient(beaconUri string) (BeaconClient, error) {
	beaconClient, _, err := NewBeaconClient(beaconUri)
	return beaconClient, err
}

func lastCheckpointedForEigenpod(eigenpodAddress string, client *ethclient.Client) uint64 {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate eigenpod. is your address correct?", err)

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	PanicOnError("failed to locate eigenpod. Is your address correct?", err)

	return timestamp
}

// search through beacon state for validators whose withdrawal address is set to eigenpod.
func findAllValidatorsForEigenpod(eigenpodAddress string, beaconState *spec.VersionedBeaconState) []ValidatorWithIndex {
	allValidators, err := beaconState.Validators()
	PanicOnError("failed to fetch beacon state", err)

	eigenpodAddressBytes := common.FromHex(eigenpodAddress)

	var outputValidators []ValidatorWithIndex = []ValidatorWithIndex{}
	var i uint64 = 0
	maxValidators := uint64(len(allValidators))
	for i = 0; i < maxValidators; i++ {
		validator := allValidators[i]
		if validator == nil || validator.WithdrawalCredentials[0] != 1 { // withdrawalCredentials _need_ their first byte set to 1 to withdraw to execution layer.
			continue
		}
		// we check that the last 20 bytes of expectedCredentials matches validatorCredentials.
		if bytes.Equal(
			eigenpodAddressBytes[:],
			validator.WithdrawalCredentials[12:], // first 12 bytes are not the pubKeyHash, see (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/pods/EigenPod.sol#L663)
		) {
			outputValidators = append(outputValidators, ValidatorWithIndex{
				Validator: validator,
				Index:     i,
			})
		}
	}
	return outputValidators
}

func getOnchainValidatorInfo(client *ethclient.Client, eigenpodAddress string, allValidators []ValidatorWithIndex) []onchain.IEigenPodValidatorInfo {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), client)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	var validatorInfo []onchain.IEigenPodValidatorInfo = []onchain.IEigenPodValidatorInfo{}

	// TODO: batch/multicall
	zeroes := [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for i := 0; i < len(allValidators); i++ {
		// ssz requires values to be 32-byte aligned, which requires 16 bytes of 0's to be added
		// prior to hashing.
		pubKeyHash := sha256.Sum256(
			append(
				(allValidators[i]).Validator.PublicKey[:],
				zeroes[:]...,
			),
		)
		info, err := eigenPod.ValidatorPubkeyHashToInfo(nil, pubKeyHash)
		PanicOnError("failed to fetch validator eigeninfo.", err)
		validatorInfo = append(validatorInfo, info)
	}

	return validatorInfo
}

func getCurrentCheckpointBlockRoot(eigenpodAddress string, eth *ethclient.Client) (*[32]byte, error) {
	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to locate Eigenpod. Is your address correct?", err)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to reach eigenpod.", err)

	return &checkpoint.BeaconBlockRoot, nil
}
