package core

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func SubmitCheckpointProof(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, proof *eigenpodproofs.VerifyCheckpointProofsCallParams, eth *ethclient.Client, batchSize uint64, noPrompt bool, noSend bool, verbose bool) ([]*types.Transaction, error) {
	tracing := GetContextTracingCallbacks(ctx)

	allProofChunks := chunk(proof.BalanceProofs, batchSize)
	transactions := []*types.Transaction{}
	if verbose {
		color.Green("calling EigenPod.VerifyCheckpointProofs() (using %d txn(s), max(%d) proofs per txn)", len(allProofChunks), batchSize)
	}

	for i := 0; i < len(allProofChunks); i++ {
		balanceProofs := allProofChunks[i]
		tracing.OnStartSection("pepe::proof::checkpoint::batch::submit", map[string]string{
			"chunk": fmt.Sprintf("%d", i),
		})
		txn, err := SubmitCheckpointProofBatch(ctx, owner, eigenpodAddress, chainId, proof.ValidatorBalancesRootProof, balanceProofs, eth, noSend, verbose)
		tracing.OnEndSection()
		if err != nil {
			// failed to submit batch.
			return transactions, err
		}
		transactions = append(transactions, txn)
		if verbose {
			fmt.Printf("Submitted chunk %d/%d -- waiting for transaction...: ", i+1, len(allProofChunks))
		}
		tracing.OnStartSection("pepe::proof::checkpoint::batch::wait", map[string]string{
			"chunk": fmt.Sprintf("%d", i),
		})

		if !noSend {
			bind.WaitMined(ctx, eth, txn)
		}
		tracing.OnEndSection()
		if verbose {
			color.Green("OK")
		}
	}

	if verbose {
		if !noSend {
			color.Green("Complete! re-run with `status` to see the updated Eigenpod state.")
		} else {
			color.Yellow("Submit these proofs to network and re-run with `status` to see the updated Eigenpod state.")
		}
	}
	return transactions, nil
}

func SubmitCheckpointProofBatch(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, proof *eigenpodproofs.ValidatorBalancesRootProof, balanceProofs []*eigenpodproofs.BalanceProof, eth *ethclient.Client, noSend bool, verbose bool) (*types.Transaction, error) {
	tracing := GetContextTracingCallbacks(ctx)

	ownerAccount, err := PrepareAccount(&owner, chainId, noSend)
	if err != nil {
		return nil, err
	}

	if verbose {
		fmt.Printf("Using account(0x%s) to submit onchain\n", common.Bytes2Hex(ownerAccount.FromAddress[:]))
	}

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, err
	}

	tracing.OnStartSection("pepe::proof::checkpoint::onchain::VerifyCheckpointProofs", map[string]string{
		"eigenpod": eigenpodAddress,
	})
	txn, err := eigenPod.VerifyCheckpointProofs(
		ownerAccount.TransactionOptions,
		onchain.BeaconChainProofsBalanceContainerProof{
			BalanceContainerRoot: proof.ValidatorBalancesRoot,
			Proof:                proof.Proof.ToByteSlice(),
		},
		CastBalanceProofs(balanceProofs),
	)
	tracing.OnEndSection()
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

func GenerateCheckpointProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient, verbose bool) (*eigenpodproofs.VerifyCheckpointProofsCallParams, error) {
	tracing := GetContextTracingCallbacks(ctx)

	tracing.OnStartSection("GetCurrentCheckpoint", map[string]string{})
	currentCheckpoint, err := GetCurrentCheckpoint(eigenpodAddress, eth)
	if err != nil {
		return nil, err
	}
	tracing.OnEndSection()

	tracing.OnStartSection("GetCurrentCheckpointBlockRoot", map[string]string{})
	blockRoot, err := GetCurrentCheckpointBlockRoot(eigenpodAddress, eth)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch last checkpoint: %w", err)
	}
	if blockRoot == nil {
		return nil, fmt.Errorf("failed to fetch last checkpoint - nil blockRoot")
	}
	tracing.OnEndSection()

	rootBytes := *blockRoot
	if AllZero(rootBytes[:]) {
		return nil, fmt.Errorf("no checkpoint active. Are you sure you started a checkpoint?")
	}

	headerBlock := "0x" + hex.EncodeToString((*blockRoot)[:])
	tracing.OnStartSection("GetBeaconHeader", map[string]string{})
	header, err := beaconClient.GetBeaconHeader(ctx, headerBlock)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon header (%s): %w", headerBlock, err)
	}
	tracing.OnEndSection()

	tracing.OnStartSection("GetBeaconState", map[string]string{})
	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch beacon state: %w", err)
	}
	tracing.OnEndSection()

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize prover: %w", err)
	}

	return GenerateCheckpointProofForState(ctx, eigenpodAddress, beaconState, header, eth, currentCheckpoint, proofs, verbose)
}

func GenerateCheckpointProofForState(ctx context.Context, eigenpodAddress string, beaconState *spec.VersionedBeaconState, header *v1.BeaconBlockHeader, eth *ethclient.Client, currentCheckpointTimestamp uint64, proofs *eigenpodproofs.EigenPodProofs, verbose bool) (*eigenpodproofs.VerifyCheckpointProofsCallParams, error) {
	tracing := GetContextTracingCallbacks(ctx)

	// filter through the beaconState's validators, and select only ones that have withdrawal address set to `eigenpod`.
	tracing.OnStartSection("FindAllValidatorsForEigenpod", map[string]string{})
	allValidators, err := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	if err != nil {
		return nil, err
	}
	tracing.OnEndSection()

	if verbose {
		color.Yellow("You have a total of %d validators pointed to this pod.", len(allValidators))
	}

	tracing.OnStartSection("SelectCheckpointableValidators", map[string]string{})
	checkpointValidators, err := SelectCheckpointableValidators(eth, eigenpodAddress, allValidators, currentCheckpointTimestamp)
	if err != nil {
		return nil, err
	}
	tracing.OnEndSection()

	validatorIndices := make([]uint64, len(checkpointValidators))
	for i, v := range checkpointValidators {
		validatorIndices[i] = v.Index
	}

	if verbose {
		color.Yellow("Proving validators at indices: %s", asJSON(validatorIndices))
	}

	tracing.OnStartSection("ProveCheckpointProofs", map[string]string{})
	proof, err := proofs.ProveCheckpointProofs(header.Header.Message, beaconState, validatorIndices)
	if err != nil {
		return nil, fmt.Errorf("failed to prove checkpoint: %w", err)
	}
	tracing.OnEndSection()

	return proof, nil
}
