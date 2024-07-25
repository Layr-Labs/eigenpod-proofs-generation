package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func SubmitCheckpointProof(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, proof *eigenpodproofs.VerifyCheckpointProofsCallParams, eth *ethclient.Client, batchSize uint64, noPrompt bool) ([]*types.Transaction, error) {
	allProofChunks := chunk(proof.BalanceProofs, batchSize)

	transactions := []*types.Transaction{}

	color.Green("calling EigenPod.VerifyCheckpointProofs() (using %d txn(s), max(%d) proofs per txn)", len(allProofChunks), batchSize)

	for i := 0; i < len(allProofChunks); i++ {
		balanceProofs := allProofChunks[i]
		txn, err := SubmitCheckpointProofBatch(owner, eigenpodAddress, chainId, proof.ValidatorBalancesRootProof, balanceProofs, eth)
		if err != nil {
			// failed to submit batch.
			return transactions, err
		}
		transactions = append(transactions, txn)
		fmt.Printf("Submitted chunk %d/%d -- waiting for transaction...: ", i+1, len(allProofChunks))
		bind.WaitMined(ctx, eth, txn)
		color.Green("OK")
	}

	color.Green("Complete! re-run with `status` to see the updated Eigenpod state.")
	return transactions, nil
}

func SubmitCheckpointProofBatch(owner, eigenpodAddress string, chainId *big.Int, proof *eigenpodproofs.ValidatorBalancesRootProof, balanceProofs []*eigenpodproofs.BalanceProof, eth *ethclient.Client) (*types.Transaction, error) {
	ownerAccount, err := PrepareAccount(&owner, chainId)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Using account(0x%s) to submit onchain\n", common.Bytes2Hex(ownerAccount.FromAddress[:]))

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, err
	}

	txn, err := eigenPod.VerifyCheckpointProofs(
		ownerAccount.TransactionOptions,
		onchain.BeaconChainProofsBalanceContainerProof{
			BalanceContainerRoot: proof.ValidatorBalancesRoot,
			Proof:                proof.Proof.ToByteSlice(),
		},
		CastBalanceProofs(balanceProofs),
	)
	if err != nil {
		return nil, err
	}

	return txn, nil
}

// , out, owner *string, forceCheckpoint bool
func LoadCheckpointProofFromFile(path string) (*eigenpodproofs.VerifyCheckpointProofsCallParams, error) {
	res := eigenpodproofs.VerifyCheckpointProofsCallParams{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func asJSON(obj interface{}) string {
	bytes, _ := json.Marshal(obj)
	return string(bytes)
}

func GenerateCheckpointProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient) (*eigenpodproofs.VerifyCheckpointProofsCallParams, error) {
	currentCheckpoint, err := GetCurrentCheckpointTimestamp(eigenpodAddress, eth)
	if err != nil {
		return nil, err
	}

	blockRoot, err := GetCurrentCheckpointBlockRoot(eigenpodAddress, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch last checkpoint: %w", err)
	}
	if blockRoot == nil {
		return nil, fmt.Errorf("failed to fetch last checkpoint - nil blockRoot")
	}

	rootBytes := *blockRoot
	if AllZero(rootBytes[:]) {
		return nil, fmt.Errorf("no checkpoint active. Are you sure you started a checkpoint?")
	}

	headerBlock := "0x" + hex.EncodeToString((*blockRoot)[:])
	header, err := beaconClient.GetBeaconHeader(ctx, headerBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon header (%s): %w", headerBlock, err)
	}

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}

	// filter through the beaconState's validators, and select only ones that have withdrawal address set to `eigenpod`.
	allValidators, err := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	if err != nil {
		return nil, err
	}

	color.Yellow("You have a total of %d validators pointed to this pod.", len(allValidators))

	checkpointValidators, err := SelectCheckpointableValidators(eth, eigenpodAddress, allValidators, currentCheckpoint)
	if err != nil {
		return nil, err
	}

	validatorIndices := make([]uint64, len(checkpointValidators))
	for i, v := range checkpointValidators {
		validatorIndices[i] = v.Index
	}

	color.Yellow("Proving validators at indices: %s", asJSON(validatorIndices))

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prover: %w", err)
	}

	then := time.Now().UnixMilli()
	proof, err := proofs.ProveCheckpointProofs(header.Header.Message, beaconState, validatorIndices)
	if err != nil {
		return nil, fmt.Errorf("failed to prove checkpoint: %w", err)
	}
	now := time.Now().UnixMilli()
	fmt.Printf("Time spent generating proofs: %dms\n", now-then)

	return proof, nil
}
