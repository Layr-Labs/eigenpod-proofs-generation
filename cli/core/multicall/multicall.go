package multicall

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type MultiCallMetaData[T interface{}] struct {
	Address      common.Address
	Data         []byte
	FunctionName string
	Deserialize  func([]byte) (T, error)
}

type Multicall3Result struct {
	Success    bool
	ReturnData []byte
}

type TypedMulticall3Result[A any] struct {
	Success bool
	Value   A
	Error   error
}

type DeserializedMulticall3Result struct {
	Success bool
	Value   any
}

func (md *MultiCallMetaData[T]) Raw() RawMulticall {
	return RawMulticall{
		Address:      md.Address,
		Data:         md.Data,
		FunctionName: md.FunctionName,
		Deserialize: func(data []byte) (any, error) {
			res, err := md.Deserialize(data)
			return any(res), err
		},
	}
}

type RawMulticall struct {
	Address      common.Address
	Data         []byte
	FunctionName string
	Deserialize  func([]byte) (any, error)
}

type MulticallClient struct {
	Contract            *bind.BoundContract
	ABI                 *abi.ABI
	Context             context.Context
	MaxBatchSize        uint64
	OverrideCallOptions *bind.CallOpts
}

type ParamMulticall3Call3 struct {
	Target       common.Address
	AllowFailure bool
	CallData     []byte
}

type TMulticallClientOptions struct {
	OverrideContractAddress *common.Address
	MaxBatchSizeBytes       uint64
	OverrideCallOptions     *bind.CallOpts
}

// maxBatchSizeBytes - 0: no batching.
func NewMulticallClient(ctx context.Context, eth *ethclient.Client, options *TMulticallClientOptions) (*MulticallClient, error) {
	if eth == nil {
		return nil, errors.New("no ethclient passed")
	}

	// taken from: https://www.multicall3.com/
	parsed, err := abi.JSON(strings.NewReader(`[{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"aggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes[]","name":"returnData","type":"bytes[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bool","name":"allowFailure","type":"bool"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call3[]","name":"calls","type":"tuple[]"}],"name":"aggregate3","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bool","name":"allowFailure","type":"bool"},{"internalType":"uint256","name":"value","type":"uint256"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call3Value[]","name":"calls","type":"tuple[]"}],"name":"aggregate3Value","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"blockAndAggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[],"name":"getBasefee","outputs":[{"internalType":"uint256","name":"basefee","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"}],"name":"getBlockHash","outputs":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBlockNumber","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getChainId","outputs":[{"internalType":"uint256","name":"chainid","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockCoinbase","outputs":[{"internalType":"address","name":"coinbase","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockDifficulty","outputs":[{"internalType":"uint256","name":"difficulty","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockGasLimit","outputs":[{"internalType":"uint256","name":"gaslimit","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getCurrentBlockTimestamp","outputs":[{"internalType":"uint256","name":"timestamp","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"addr","type":"address"}],"name":"getEthBalance","outputs":[{"internalType":"uint256","name":"balance","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getLastBlockHash","outputs":[{"internalType":"bytes32","name":"blockHash","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bool","name":"requireSuccess","type":"bool"},{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"tryAggregate","outputs":[{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"bool","name":"requireSuccess","type":"bool"},{"components":[{"internalType":"address","name":"target","type":"address"},{"internalType":"bytes","name":"callData","type":"bytes"}],"internalType":"struct Multicall3.Call[]","name":"calls","type":"tuple[]"}],"name":"tryBlockAndAggregate","outputs":[{"internalType":"uint256","name":"blockNumber","type":"uint256"},{"internalType":"bytes32","name":"blockHash","type":"bytes32"},{"components":[{"internalType":"bool","name":"success","type":"bool"},{"internalType":"bytes","name":"returnData","type":"bytes"}],"internalType":"struct Multicall3.Result[]","name":"returnData","type":"tuple[]"}],"stateMutability":"payable","type":"function"}]`))
	if err != nil {
		return nil, fmt.Errorf("error parsing multicall abi: %s", err.Error())
	}

	contractAddress := func() common.Address {
		if options == nil || options.OverrideContractAddress == nil {
			// also taken from: https://www.multicall3.com/ -- it's deployed at the same addr on most chains
			return common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11")
		}
		return *options.OverrideContractAddress
	}()

	maxBatchSize := func() uint64 {
		if options == nil || options.MaxBatchSizeBytes == 0 {
			return math.MaxUint64
		} else {
			return options.MaxBatchSizeBytes
		}
	}()

	callOptions := func() *bind.CallOpts {
		if options != nil {
			return options.OverrideCallOptions
		}
		return nil
	}()

	return &MulticallClient{OverrideCallOptions: callOptions, MaxBatchSize: maxBatchSize, Context: ctx, ABI: &parsed, Contract: bind.NewBoundContract(contractAddress, parsed, eth, eth, eth)}, nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func MultiCall[T any](contractAddress common.Address, abi abi.ABI, deserialize func([]byte) (T, error), method string, params ...interface{}) (*MultiCallMetaData[T], error) {
	callData, err := abi.Pack(method, params...)
	if err != nil {
		return nil, fmt.Errorf("error packing multicall: %s", err.Error())
	}
	return &MultiCallMetaData[T]{
		Address:      contractAddress,
		Data:         callData,
		FunctionName: method,
		Deserialize:  deserialize,
	}, nil
}

func DoMultiCall[A any, B any](mc MulticallClient, a *MultiCallMetaData[A], b *MultiCallMetaData[B]) (*A, *B, error) {
	res, err := doMultiCallMany(mc, a.Raw(), b.Raw())
	if err != nil {
		return nil, nil, fmt.Errorf("error performing multicall: %s", err.Error())
	}
	return any(res[0].Value).(*A), any(res[1].Value).(*B), nil
}

func DoMultiCallMany[A any](mc MulticallClient, requests ...*MultiCallMetaData[A]) (*[]A, error) {
	res, err := doMultiCallMany(mc, utils.Map(requests, func(mc *MultiCallMetaData[A], index uint64) RawMulticall {
		return mc.Raw()
	})...)
	if err != nil {
		return nil, fmt.Errorf("multicall failed: %s", err.Error())
	}

	anyFailures := utils.Filter(res, func(cur DeserializedMulticall3Result) bool {
		return !cur.Success
	})
	if len(anyFailures) > 0 {
		return nil, errors.New("1 or more calls failed")
	}

	// unwind results
	unwoundResults := utils.Map(res, func(d DeserializedMulticall3Result, i uint64) A {
		// force these back to A
		if !d.Success {
			panic(errors.New("unexpected multicall failure"))
		}
		return any(d.Value).(A)
	})
	return &unwoundResults, nil
}

func DoMultiCallManyReportingFailures[A any](mc MulticallClient, requests ...*MultiCallMetaData[A]) (*[]TypedMulticall3Result[A], error) {
	res, err := doMultiCallMany(mc, utils.Map(requests, func(mc *MultiCallMetaData[A], index uint64) RawMulticall {
		return mc.Raw()
	})...)
	if err != nil {
		return nil, fmt.Errorf("multicall failed: %s", err.Error())
	}

	// unwind results
	unwoundResults := utils.Map(res, func(d DeserializedMulticall3Result, i uint64) TypedMulticall3Result[A] {
		val, ok := any(d.Value).(A)
		if !ok {
			return TypedMulticall3Result[A]{
				Value:   val,
				Success: false,
			}
		}

		return TypedMulticall3Result[A]{
			Value:   val,
			Success: d.Success,
		}
	})
	return &unwoundResults, nil
}

/*
 * Some RPC providers may limit the amount of calldata you can send in one eth_call, which (for those who have 1000's of validators), means
 * you can't just spam one enormous multicall request.
 *
 * This function checks whether the calldata appended exceeds maxBatchSizeBytes
 */
func chunkCalls(allCalls []ParamMulticall3Call3, maxBatchSizeBytes int) [][]ParamMulticall3Call3 {
	// chunk by the maximum size of calldata, which is 1024 per call.
	results := [][]ParamMulticall3Call3{}
	currentBatchSize := 0
	currentBatch := []ParamMulticall3Call3{}

	for _, call := range allCalls {
		if (currentBatchSize + len(call.CallData)) > maxBatchSizeBytes {
			// we can't fit in this batch, so dump the current batch and start a new one
			results = append(results, currentBatch)
			currentBatchSize = 0
			currentBatch = []ParamMulticall3Call3{}
		}

		currentBatch = append(currentBatch, call)
		currentBatchSize += len(call.CallData)
	}

	// check if we forgot to add the last batch
	if len(currentBatch) > 0 {
		results = append(results, currentBatch)
	}

	return results
}

func doMultiCallMany(mc MulticallClient, calls ...RawMulticall) ([]DeserializedMulticall3Result, error) {
	typedCalls := make([]ParamMulticall3Call3, len(calls))
	for i, call := range calls {
		typedCalls[i] = ParamMulticall3Call3{
			Target:       call.Address,
			AllowFailure: true,
			CallData:     call.Data,
		}
	}

	// see if we need to chunk them now
	chunkedCalls := chunkCalls(typedCalls, func() int {
		if mc.MaxBatchSize == 0 {
			return math.MaxInt64
		} else {
			return int(mc.MaxBatchSize)
		}
	}())
	var results = make([]interface{}, len(calls))
	var totalResults = 0

	chunkNumber := 1
	for _, multicalls := range chunkedCalls {
		var res []interface{}
		chunkNumber++
		err := mc.Contract.Call(mc.OverrideCallOptions, &res, "aggregate3", multicalls)
		if err != nil {
			return nil, fmt.Errorf("aggregate3 failed: %s", err)
		}

		multicallResults := *abi.ConvertType(res[0], new([]Multicall3Result)).(*[]Multicall3Result)
		for i := 0; i < len(multicallResults); i++ {
			results[totalResults+i] = multicallResults[i]
		}
		totalResults += len(multicallResults)
	}

	outputs := make([]DeserializedMulticall3Result, len(calls))
	for i, call := range calls {
		res := results[i].(Multicall3Result)
		if res.Success {
			if res.ReturnData != nil {
				val, err := call.Deserialize(res.ReturnData)
				if err != nil {
					outputs[i] = DeserializedMulticall3Result{
						Value:   err,
						Success: false,
					}
				} else {
					outputs[i] = DeserializedMulticall3Result{
						Value:   val,
						Success: res.Success,
					}
				}
			} else {
				outputs[i] = DeserializedMulticall3Result{
					Value:   errors.New("no data returned"),
					Success: false,
				}
			}
		} else {
			outputs[i] = DeserializedMulticall3Result{
				Success: false,
				Value:   errors.New("call failed"),
			}
		}
	}

	return outputs, nil
}
