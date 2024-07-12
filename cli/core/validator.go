package core

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func SubmitValidatorProof(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, eth *ethclient.Client, batchSize uint64, proofs *eigenpodproofs.VerifyValidatorFieldsCallParams, oracleBeaconTimesetamp uint64, noPrompt bool) ([]*types.Transaction, error) {
	ownerAccount, err := PrepareAccount(&owner, chainId)
	if err != nil {
		return nil, err
	}
	PanicOnError("failed to parse private key", err)

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, err
	}

	indices := Uint64ArrayToBigIntArray(proofs.ValidatorIndices)
	validatorIndicesChunks := chunk(indices, batchSize)
	validatorProofsChunks := chunk(proofs.ValidatorFieldsProofs, batchSize)
	validatorFieldsChunks := chunk(proofs.ValidatorFields, batchSize)
	if !noPrompt {
		PanicIfNoConsent(SubmitCredentialsProofConsent(len(validatorFieldsChunks)))
	}

	color.Green("calling EigenPod.VerifyWithdrawalCredentials() (using %d txn(s), max(%d) proofs per txn)", len(indices), batchSize)

	transactions := []*types.Transaction{}
	numChunks := len(validatorIndicesChunks)

	color.Green("Submitting proofs with %d transactions", numChunks)

	for i := 0; i < numChunks; i++ {
		curValidatorIndices := validatorIndicesChunks[i]
		curValidatorProofs := validatorProofsChunks[i]

		var validatorFieldsProofs [][]byte = [][]byte{}
		for i := 0; i < len(curValidatorProofs); i++ {
			pr := curValidatorProofs[i].ToByteSlice()
			validatorFieldsProofs = append(validatorFieldsProofs, pr)
		}
		var curValidatorFields [][][32]byte = CastValidatorFields(validatorFieldsChunks[i])

		fmt.Printf("Submitted chunk %d/%d -- waiting for transaction...: ", i+1, numChunks)
		txn, err := SubmitValidatorProofChunk(ctx, ownerAccount, eigenPod, chainId, eth, curValidatorIndices, curValidatorFields, proofs.StateRootProof, validatorFieldsProofs, oracleBeaconTimesetamp)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, txn)
	}

	return transactions, err
}

func SubmitValidatorProofChunk(ctx context.Context, ownerAccount *Owner, eigenPod *onchain.EigenPod, chainId *big.Int, eth *ethclient.Client, indices []*big.Int, validatorFields [][][32]byte, stateRootProofs *eigenpodproofs.StateRootProof, validatorFieldsProofs [][]byte, oracleBeaconTimesetamp uint64) (*types.Transaction, error) {
	color.Green("submitting onchain...")
	txn, err := eigenPod.VerifyWithdrawalCredentials(
		ownerAccount.TransactionOptions,
		oracleBeaconTimesetamp,
		onchain.BeaconChainProofsStateRootProof{
			Proof:           stateRootProofs.Proof.ToByteSlice(),
			BeaconStateRoot: stateRootProofs.BeaconStateRoot,
		},
		indices,
		validatorFieldsProofs,
		validatorFields,
	)

	return txn, err
}

func GenerateValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient) (*eigenpodproofs.VerifyValidatorFieldsCallParams, uint64, error) {
	latestBlock, err := eth.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load latest block: %w", err)
	}

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to reach eigenpod: %w", err)
	}

	expectedBlockRoot, err := eigenPod.GetParentBlockRoot(nil, latestBlock.Time())
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load parent block root: %w", err)
	}

	header, err := beaconClient.GetBeaconHeader(ctx, "0x"+common.Bytes2Hex(expectedBlockRoot[:]))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch beacon header: %w", err)
	}

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch beacon state: %w", err)
	}

	allValidators := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	awaitingCredentialValidators := SelectAwaitingCredentialValidators(eth, eigenpodAddress, allValidators)

	if len(awaitingCredentialValidators) == 0 {
		color.Red("You have no inactive validators to verify. Everything up-to-date.")
		return nil, 0, nil
	} else {
		color.Blue("Verifying %d inactive validators", len(awaitingCredentialValidators))
	}

	validatorIndices := make([]uint64, len(awaitingCredentialValidators))
	for i, v := range awaitingCredentialValidators {
		validatorIndices[i] = v.Index
	}

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to initialize provider: %w", err)
	}

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, validatorIndices)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to prove validators: %w", err)
	}

	return validatorProofs, latestBlock.Time(), nil
}
