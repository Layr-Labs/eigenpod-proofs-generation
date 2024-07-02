package main

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/ethereum/go-ethereum/ethclient"
)

func RunValidatorProof(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, chainId *big.Int, beaconClient BeaconClient, out *string) {
	header, err := beaconClient.GetBeaconHeader(ctx, "head")
	PanicOnError("failed to fetch latest beacon header.", err)

	beaconState, err := beaconClient.GetBeaconState(ctx, strconv.FormatUint(uint64(header.Header.Message.Slot), 10))
	PanicOnError("failed to fetch beacon state.", err)

	allValidatorsForEigenpod := findAllValidatorsForEigenpod(eigenpodAddress, beaconState)
	allValidatorInfo := getOnchainValidatorInfo(eth, eigenpodAddress, allValidatorsForEigenpod)

	var checkpointValidatorIndices = FilterInactiveValidators(allValidatorsForEigenpod, allValidatorInfo)

	proofs, err := eigenpodproofs.NewEigenPodProofs(chainId.Uint64(), 300 /* oracleStateCacheExpirySeconds - 5min */)
	PanicOnError("failled to initialize prover", err)

	// validator proof
	validatorProofs, err := proofs.ProveValidatorContainers(header.Header.Message, beaconState, checkpointValidatorIndices)
	PanicOnError("failed to prove validators.", err)

	jsonString, err := json.Marshal(validatorProofs)
	PanicOnError("failed to generate JSON proof data.", err)

	WriteOutputToFileOrStdout(jsonString, out)
}
