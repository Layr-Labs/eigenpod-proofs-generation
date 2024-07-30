//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"math/big"
	"syscall/js"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getClients(node, beaconNodeUri string) (*ethclient.Client, core.BeaconClient, *big.Int, error) {
	eth, err := ethclient.Dial(node)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach eth --node: %w", err)
	}

	chainId, err := eth.ChainID(context.Background())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to fetch chain id: %w", err)
	}

	beaconClient, err := core.GetBeaconClient(beaconNodeUri)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to reach beacon chain: %w", err)
	}

	return eth, beaconClient, chainId, nil
}

func Checkpoint(this js.Value, args []js.Value) any {
	eigenpodAddress := args[0].String()
	beacon := args[1].String()
	node := args[2].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			var eth, beaconClient, chainId, err = getClients(node, beacon)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			proofs, err := core.GenerateCheckpointProof(context.Background(), eigenpodAddress, eth, chainId, beaconClient)

			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
			} else {
				resolve.Invoke(js.ValueOf(proofs))
			}
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

type TValidatorProofStruct = struct {
	Proof           eigenpodproofs.VerifyValidatorFieldsCallParams
	BeaconTimestamp uint64
}

func Validator(this js.Value, args []js.Value) any {
	eigenpodAddress := args[0].String()
	beacon := args[1].String()
	node := args[2].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			var eth, beaconClient, chainId, err = getClients(node, beacon)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
				return
			}

			proofs, beaconTimestamp, err := core.GenerateValidatorProof(context.Background(), eigenpodAddress, eth, chainId, beaconClient)
			if err != nil {
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)
			} else {
				resolve.Invoke(js.ValueOf(TValidatorProofStruct{
					Proof:           *proofs,
					BeaconTimestamp: beaconTimestamp,
				}))
			}
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func main() {
	eigen := js.Global().Get("Object").New()

	eigen.Set("checkpoint", js.FuncOf(Checkpoint))
	eigen.Set("validator", js.FuncOf(Validator))

	js.Global().Set("_eigen", eigen)

	<-make(chan struct{})
}
