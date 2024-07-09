package core

import (
	"context"
	"math/big"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fatih/color"
)

func GenerateValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient) (*eigenpodproofs.VerifyValidatorFieldsCallParams, []uint64) {
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

	allValidatorsForEigenpod := FindAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	allValidatorInfo := GetOnchainValidatorInfo(eth, eigenpodAddress, allValidatorsForEigenpod)

	var validatorIndices = FilterInactiveValidators(allValidatorsForEigenpod, allValidatorInfo)
	if len(validatorIndices) == 0 {
		color.Red("You have no inactive validators to verify. Everything up-to-date.")
		return nil, nil
	}

	color.Blue("Verifying %d inactive validators", len(validatorIndices))

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	PanicOnError("failled to initialize prover", err)

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, validatorIndices)
	PanicOnError("failed to prove validators.", err)

	return validatorProofs, validatorIndices
}
