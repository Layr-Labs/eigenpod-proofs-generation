package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

const (
	ValidatorStatusInactive  = 0
	ValidatorStatusActive    = 1
	ValidatorStatusWithdrawn = 2
)

func Panic(message string) {
	color.Red(fmt.Sprintf("error: %s\n\n", message))

	os.Exit(1)
}

func PanicOnError(message string, err error) {
	if err != nil {
		color.Red(fmt.Sprintf("error: %s\n\n", message))

		info := color.New(color.FgRed, color.Italic)
		info.Printf(fmt.Sprintf("caused by: %s\n", err))

		os.Exit(1)
	}
}

type ValidatorWithIndex = struct {
	Validator *phase0.Validator
	Index     uint64
}

type Owner = struct {
	FromAddress        gethCommon.Address
	PublicKey          *ecdsa.PublicKey
	TransactionOptions *bind.TransactOpts
}

func startCheckpoint(ctx context.Context, eigenpodAddress string, owner string, chainId *big.Int, eth *ethclient.Client, forceCheckpoint bool) (uint64, error) {
	ownerAccount, err := prepareAccount(&owner, chainId)
	PanicOnError("failed to parse private key", err)

	eigenPod, err := onchain.NewEigenPod(gethCommon.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to reach eigenpod", err)

	// revertIfNoBalance == !forceCheckpoint
	// The CLI exposes this as a `force` parameter for usability
	revertIfNoBalance := true
	if forceCheckpoint {
		revertIfNoBalance = false
	}

	txn, err := eigenPod.StartCheckpoint(ownerAccount.TransactionOptions, revertIfNoBalance)
	PanicOnError("failed to start checkpoint", err)

	color.Green("starting checkpoint: %s.. (waiting for txn to be mined)...", txn.Hash().Hex())

	bind.WaitMined(ctx, eth, txn)

	color.Green("started checkpoint! txn: %s", txn.Hash().Hex())

	currentCheckpoint := getCurrentCheckpoint(eigenpodAddress, eth)
	return currentCheckpoint, nil
}

func castBalanceProofs(proofs []*eigenpodproofs.BalanceProof) []onchain.BeaconChainProofsBalanceProof {
	out := []onchain.BeaconChainProofsBalanceProof{}

	for i := 0; i < len(proofs); i++ {
		proof := proofs[i]
		out = append(out, onchain.BeaconChainProofsBalanceProof{
			PubkeyHash:  proof.PubkeyHash,
			BalanceRoot: proof.BalanceRoot,
			Proof:       proof.Proof.ToByteSlice(),
		})
	}

	return out
}

func prepareAccount(owner *string, chainID *big.Int) (*Owner, error) {
	if owner == nil {
		return nil, fmt.Errorf("no owner")
	}

	privateKey, err := crypto.HexToECDSA(*owner)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}

	return &Owner{
		FromAddress:        fromAddress,
		PublicKey:          publicKeyECDSA,
		TransactionOptions: auth,
	}, nil
}

// golang was a mistake. these types are literally identical :'(
func castValidatorFields(proof [][]eigenpodproofs.Bytes32) [][][32]byte {
	result := make([][][32]byte, len(proof))

	for i, slice := range proof {
		result[i] = make([][32]byte, len(slice))
		for j, bytes := range slice {
			result[i][j] = bytes
		}
	}

	return result
}

func Uint64ArrayToBigIntArray(nums []uint64) []*big.Int {
	out := []*big.Int{}
	for i := 0; i < len(nums); i++ {
		bigInt := new(big.Int).SetUint64(nums[i])
		out = append(out, bigInt)
	}
	return out
}

func AllZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func WriteOutputToFileOrStdout(output []byte, out *string) {
	if out != nil && *out != "" {
		err := os.WriteFile(*out, output, os.ModePerm)
		PanicOnError("failed to write to disk", err)
		color.Green("Wrote output to %s\n", *out)
	} else {
		fmt.Println(string(output))
	}
}

func FilterNotCheckpointedOrWithdrawnValidators(
	allValidatorsForEigenpod []ValidatorWithIndex,
	onchainInfo []onchain.IEigenPodValidatorInfo,
	lastCheckpoint uint64,
) []uint64 {
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := onchainInfo[i]

		notCheckpointed := validatorInfo.LastCheckpointedAt != lastCheckpoint
		notWithdrawn := validatorInfo.Status != ValidatorStatusWithdrawn
		notInactive := validatorInfo.Status != ValidatorStatusInactive

		if notCheckpointed && notWithdrawn && notInactive {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.Index)
		}
	}
	return checkpointValidatorIndices
}

// (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/libraries/BeaconChainProofs.sol#L64)
const FAR_FUTURE_EPOCH = math.MaxUint64

func FilterInactiveValidators(
	allValidatorsForEigenpod []ValidatorWithIndex,
	onchainInfo []onchain.IEigenPodValidatorInfo,
) []uint64 {
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := onchainInfo[i]

		if (validatorInfo.Status == ValidatorStatusInactive) &&
			(validator.Validator.ExitEpoch == FAR_FUTURE_EPOCH) {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.Index)
		}
	}
	return checkpointValidatorIndices
}
