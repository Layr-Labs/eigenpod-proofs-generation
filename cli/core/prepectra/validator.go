package prepectra

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"

	lo "github.com/samber/lo"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

type SerializableCredentialProof struct {
	ValidatorProofs       *eigenpodproofs.VerifyValidatorFieldsCallParams
	OracleBeaconTimestamp uint64
}

func LoadValidatorProofFromFile(path string) (*SerializableCredentialProof, error) {
	res := SerializableCredentialProof{}
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

func SubmitValidatorProof(ctx context.Context, owner, eigenpodAddress string, chainId *big.Int, eth *ethclient.Client, batchSize uint64, proofs *eigenpodproofs.VerifyValidatorFieldsCallParams, oracleBeaconTimesetamp uint64, noPrompt bool, noSend bool, verbose bool) ([]*types.Transaction, [][]*big.Int, error) {
	ownerAccount, err := utils.PrepareAccount(&owner, chainId, noSend)
	if err != nil {
		return nil, [][]*big.Int{}, err
	}
	utils.PanicOnError("failed to parse private key", err)

	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return nil, [][]*big.Int{}, err
	}

	indices := utils.Uint64ArrayToBigIntArray(proofs.ValidatorIndices)
	validatorIndicesChunks := utils.Chunk(indices, batchSize)
	validatorProofsChunks := utils.Chunk(proofs.ValidatorFieldsProofs, batchSize)
	validatorFieldsChunks := utils.Chunk(proofs.ValidatorFields, batchSize)
	if !noPrompt && !noSend {
		utils.PanicIfNoConsent(utils.SubmitCredentialsProofConsent(len(validatorFieldsChunks)))
	}

	transactions := []*types.Transaction{}
	numChunks := len(validatorIndicesChunks)

	if verbose {
		color.Green("calling EigenPod.VerifyWithdrawalCredentials() (using %d txn(s), max(%d) proofs per txn [%s])", numChunks, batchSize, func() string {
			if ownerAccount.TransactionOptions.NoSend {
				return "simulated"
			} else {
				return "live"
			}
		}())
		color.Green("Submitting proofs with %d transactions", numChunks)
	}

	for i := 0; i < numChunks; i++ {
		curValidatorIndices := validatorIndicesChunks[i]
		curValidatorProofs := validatorProofsChunks[i]

		var validatorFieldsProofs [][]byte = [][]byte{}
		for i := 0; i < len(curValidatorProofs); i++ {
			pr := curValidatorProofs[i].ToByteSlice()
			validatorFieldsProofs = append(validatorFieldsProofs, pr)
		}
		var curValidatorFields [][][32]byte = utils.CastValidatorFields(validatorFieldsChunks[i])

		if verbose {
			fmt.Printf("Submitted chunk %d/%d -- waiting for transaction...: ", i+1, numChunks)
		}
		txn, err := SubmitValidatorProofChunk(ctx, ownerAccount, eigenPod, chainId, eth, curValidatorIndices, curValidatorFields, proofs.StateRootProof, validatorFieldsProofs, oracleBeaconTimesetamp, verbose)
		if err != nil {
			return transactions, validatorIndicesChunks, err
		}

		transactions = append(transactions, txn)
	}

	return transactions, validatorIndicesChunks, err
}

func SubmitValidatorProofChunk(ctx context.Context, ownerAccount *utils.Owner, eigenPod *EigenPod.EigenPod, chainId *big.Int, eth *ethclient.Client, indices []*big.Int, validatorFields [][][32]byte, stateRootProofs *eigenpodproofs.StateRootProof, validatorFieldsProofs [][]byte, oracleBeaconTimesetamp uint64, verbose bool) (*types.Transaction, error) {
	if verbose {
		color.Green("submitting...")
	}
	txn, err := eigenPod.VerifyWithdrawalCredentials(
		ownerAccount.TransactionOptions,
		oracleBeaconTimesetamp,
		EigenPod.BeaconChainProofsStateRootProof{
			Proof:           stateRootProofs.Proof.ToByteSlice(),
			BeaconStateRoot: stateRootProofs.BeaconStateRoot,
		},
		indices,
		validatorFieldsProofs,
		validatorFields,
	)

	return txn, err
}

/**
 * Generates a .ProveValidatorContainers() proof for all eligible validators on the pod. If `validatorIndex` is set, it will only generate  a proof
 * against that validator, regardless of the validator's state.
 */
func GenerateValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient utils.BeaconClient, validatorIndex *big.Int, verbose bool) (*eigenpodproofs.VerifyValidatorFieldsCallParams, uint64, error) {
	latestBlock, err := eth.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load latest block: %w", err)
	}

	eigenPod, err := EigenPod.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
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

	proofExecutor, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to initialize provider: %w", err)
	}

	proofs, err := GenerateValidatorProofAtState(ctx, proofExecutor, eigenpodAddress, beaconState, eth, chainId, header, latestBlock.Time(), validatorIndex, verbose)
	return proofs, latestBlock.Time(), err
}

func GenerateValidatorProofAtState(ctx context.Context, proofs *eigenpodproofs.EigenPodProofs, eigenpodAddress string, beaconState *spec.VersionedBeaconState, eth *ethclient.Client, chainId *big.Int, header *v1.BeaconBlockHeader, blockTimestamp uint64, forSpecificValidatorIndex *big.Int, verbose bool) (*eigenpodproofs.VerifyValidatorFieldsCallParams, error) {
	allValidators, err := utils.FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	if err != nil {
		return nil, fmt.Errorf("failed to find validators: %w", err)
	}

	var awaitingCredentialValidators []utils.ValidatorWithIndex

	if forSpecificValidatorIndex != nil {
		// prove a specific validator
		for _, v := range allValidators {
			if v.Index == forSpecificValidatorIndex.Uint64() {
				awaitingCredentialValidators = []utils.ValidatorWithIndex{v}
				break
			}
		}
		if len(awaitingCredentialValidators) == 0 {
			return nil, fmt.Errorf("validator at index %d does not exist or does not have withdrawal credentials set to pod %s", forSpecificValidatorIndex.Uint64(), eigenpodAddress)
		}
	} else {
		// default behavior -- load any validators that are inactive / need a credential proof
		allValidatorsWithInfo, err := utils.FetchMultipleOnchainValidatorInfo(ctx, eth, eigenpodAddress, allValidators)
		if err != nil {
			return nil, fmt.Errorf("failed to load validator information: %s", err.Error())
		}

		_awaitingCredentialValidators, err := utils.SelectAwaitingCredentialValidators(eth, eigenpodAddress, allValidatorsWithInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to select awaiting credential validators: %s", err.Error())
		}

		awaitingCredentialValidators = lo.Map(_awaitingCredentialValidators, func(validator utils.ValidatorWithOnchainInfo, index int) utils.ValidatorWithIndex {
			return utils.ValidatorWithIndex{Index: validator.Index, Validator: validator.Validator}
		})
	}

	if len(awaitingCredentialValidators) == 0 {
		if verbose {
			color.Red("You have no inactive validators to verify. Everything up-to-date.")
		}
		return nil, nil
	} else {
		if verbose {
			color.Blue("Verifying %d inactive validators", len(awaitingCredentialValidators))
		}
	}

	validatorIndices := make([]uint64, len(awaitingCredentialValidators))
	for i, v := range awaitingCredentialValidators {
		validatorIndices[i] = v.Index
	}

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, validatorIndices)
	if err != nil {
		return nil, fmt.Errorf("failed to prove validators: %w", err)
	}

	return validatorProofs, nil
}
