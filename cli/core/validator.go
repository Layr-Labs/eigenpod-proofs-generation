package core

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func getAllFromChannel[T any](ch chan T) []T {
	elements := []T{}
	for s := range ch {
		elements = append(elements, s)
	}
	return elements
}

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

	numChunks := len(validatorIndicesChunks)

	color.Green("Submitting proofs with %d transactions", numChunks)

	var wg sync.WaitGroup
	txns := make(chan *types.Transaction)
	errs := make(chan error)

	color.Green("submitting onchain...")

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

		wg.Add(1)

		go SubmitValidatorProofChunk(ctx, txns, errs, &wg, ownerAccount, eigenPod, chainId, eth, curValidatorIndices, curValidatorFields, proofs, validatorFieldsProofs, latestBlock.Time())
	}

	wg.Wait()

	var resultErr error = nil
	allTxns := getAllFromChannel(txns)
	allErrs := getAllFromChannel(errs)

	if len(allErrs) > 0 {
		resultErr = fmt.Errorf("%d error(s) occurred while submitting transactions: %w", len(allErrs), errors.Join(allErrs...))
	}

	return allTxns, resultErr
}

func SubmitValidatorProofChunk(ctx context.Context, txOut chan *types.Transaction, errOut chan error, wg *sync.WaitGroup, ownerAccount *Owner, eigenPod *onchain.EigenPod, chainId *big.Int, eth *ethclient.Client, indices []*big.Int, validatorFields [][][32]byte, proofs *eigenpodproofs.VerifyValidatorFieldsCallParams, validatorFieldsProofs [][]byte, oracleBeaconTimesetamp uint64) {
	defer wg.Done()
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

	if txn != nil {
		txOut <- txn
	}
	if err != nil {
		errOut <- err
	}
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
