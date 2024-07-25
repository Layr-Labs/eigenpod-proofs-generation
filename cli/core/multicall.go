package core

import (
	"log"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/alethio/web3-multicall-go/multicall"
)

// https://www.multicall3.com/deployments
const HoleskyAddress = "0xcA11bde05977b3631167028862bE2a173976CA11"

func NewMulticall(ethNodeUri string) (multicall.Multicall, error) {
	batchLoader, err := httprpc.NewBatchLoader(0, 10*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	provider, err := httprpc.NewWithLoader(ethNodeUri, batchLoader)
	if err != nil {
		log.Fatal(err)
	}
	eth, err := ethrpc.New(provider)
	multi, err := multicall.New(eth, multicall.ContractAddress(HoleskyAddress))
	return multi, err
}

func BatchRequests[A, B *interface{}](mc multicall.Multicall, request A, request2 B) (A, B, error) {
	vcs := multicall.ViewCalls{
		multicall.NewViewCall(
			"1",
			"0x5eb3fa2dfecdde21c950813c665e9364fa609bd2",
			"getLastBlockHash()(bytes32)",
			[]interface{}{},
		),
		multicall.NewViewCall(
			"2",
			"0x6b175474e89094c44da98b954eedeac495271d0f",
			"balanceOf(address)(uint256)",
			[]interface{}{"0x8134d518e0cef5388136c0de43d7e12278701ac5"},
		),
	}
	block := "latest"
	res, err := mc.Call(vcs, block)
	if err != nil {
		return nil, nil, err
	}

	callOneSuccess := res.Calls["1"].Success
	callOneRaw := res.Calls["1"].Raw

	callTwoSuccess := res.Calls["2"].Success
	callTwoRaw := res.Calls["2"].Raw

	return res
}
