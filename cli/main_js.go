//go:build js && wasm

package main

import (
	"context"
	"math/big"
	"syscall/js"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getClients(node, beaconNodeUri string) (*ethclient.Client, BeaconClient, *big.Int) {
	eth, err := ethclient.Dial(node)
	PanicOnError("failed to reach eth --node.", err)

	chainId, err := eth.ChainID(context.Background())
	PanicOnError("failed to fetch chain id", err)

	beaconClient, err := getBeaconClient(beaconNodeUri)
	PanicOnError("failed to reach beacon chain.", err)

	return eth, beaconClient, chainId
}

func Checkpoint(this js.Value, args []js.Value) any {
	eigenpodAddress := args[0].String()
	beacon := args[1].String()
	node := args[2].String()
	owner := args[3].String()

	var eth, beaconClient, chainId = getClients(node, beacon)
	var ownerValue *string = nil
	if len(owner) > 0 {
		ownerValue = &owner
	}

	proofs := RunCheckpointProof(context.Background(), eigenpodAddress, eth, chainId, beaconClient, ownerValue)
	return *proofs
}

type TValidatorProofStruct = struct {
	Proof            eigenpodproofs.VerifyValidatorFieldsCallParams
	ValidatorIndices []uint64
	LatestBlock      types.Block
}

func Validator(this js.Value, args []js.Value) any {
	eigenpodAddress := args[0].String()
	beacon := args[1].String()
	node := args[2].String()
	owner := args[3].String()

	var eth, beaconClient, chainId = getClients(node, beacon)
	var ownerValue *string = nil
	if len(owner) > 0 {
		ownerValue = &owner
	}

	proofs, validatorIndices, block := RunValidatorProof(context.Background(), eigenpodAddress, eth, chainId, beaconClient, ownerValue)
	return TValidatorProofStruct{
		Proof:            *proofs,
		ValidatorIndices: validatorIndices,
		LatestBlock:      *block,
	}
}

func main() {
	eigen := js.Global().Get("Object").New()

	eigen.Set("checkpoint", js.FuncOf(Checkpoint))
	eigen.Set("validator", js.FuncOf(Validator))

	js.Global().Set("_eigen", eigen)

	<-make(chan struct{})
}
