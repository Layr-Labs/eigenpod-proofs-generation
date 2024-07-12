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

func SubmitValidatorProof(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, eth *ethclient.Client, batchSize uint64, proofs *eigenpodproofs.VerifyValidatorFieldsCallParams, noPrompt bool) ([]*types.Transaction, error) {
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
		PanicIfNoConsent(fmt.Sprintf("This will call EigenPod.VerifyWithdrawalCredentials() %d times, to link your validator to your eigenpod.", len(validatorIndicesChunks)))
	}

	color.Green("calling EigenPod.VerifyWithdrawalCredentials() (using %d txn(s), max(%d) proofs per txn)", len(indices), batchSize)

	latestBlock, err := eth.BlockByNumber(ctx, nil)
	if err != nil {
		return []*types.Transaction{}, err
	}

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
		txn, err := SubmitValidatorProofChunk(ctx, ownerAccount, eigenPod, chainId, eth, curValidatorIndices, curValidatorFields, proofs, validatorFieldsProofs, latestBlock.Time())
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, txn)
	}

	return transactions, err
}

func SubmitValidatorProofChunk(ctx context.Context, ownerAccount *Owner, eigenPod *onchain.EigenPod, chainId *big.Int, eth *ethclient.Client, indices []*big.Int, validatorFields [][][32]byte, proofs *eigenpodproofs.VerifyValidatorFieldsCallParams, validatorFieldsProofs [][]byte, oracleBeaconTimesetamp uint64) (*types.Transaction, error) {
	color.Green("submitting onchain...")
	txn, err := eigenPod.VerifyWithdrawalCredentials(
		ownerAccount.TransactionOptions,
		oracleBeaconTimesetamp,
		onchain.BeaconChainProofsStateRootProof{
			Proof:           proofs.StateRootProof.Proof.ToByteSlice(),
			BeaconStateRoot: proofs.StateRootProof.BeaconStateRoot,
		},
		indices,
		validatorFieldsProofs,
		validatorFields,
	)

	return txn, err
}

func GenerateValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient) *eigenpodproofs.VerifyValidatorFieldsCallParams {
	latestBlock, err := eth.BlockByNumber(ctx, nil)
	PanicOnError("failed to load latest block", err)

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to reach eigenpod", err)

	expectedBlockRoot, err := eigenPod.GetParentBlockRoot(nil, latestBlock.Time())
	PanicOnError("failed to load parent block root", err)

	header, err := beaconClient.GetBeaconHeader(ctx, "0x"+common.Bytes2Hex(expectedBlockRoot[:]))
	PanicOnError("failed to fetch beacon header.", err)

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	PanicOnError("failed to fetch beacon state.", err)

	allValidators := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	awaitingCredentialValidators := SelectAwaitingCredentialValidators(eth, eigenpodAddress, allValidators)

	if len(awaitingCredentialValidators) == 0 {
		color.Red("You have no inactive validators to verify. Everything up-to-date.")
		return nil
	} else {
		color.Blue("Verifying %d inactive validators", len(awaitingCredentialValidators))
	}

	validatorIndices := make([]uint64, len(awaitingCredentialValidators))
	for i, v := range awaitingCredentialValidators {
		validatorIndices[i] = v.Index
	}

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	PanicOnError("failled to initialize prover", err)

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, validatorIndices)
	PanicOnError("failed to prove validators.", err)

	return validatorProofs
}
