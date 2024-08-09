package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	multicall "github.com/forta-network/go-multicall"
)

type Cache struct {
	PodOwnerShares map[string]PodOwnerShare
}

type PodOwnerShare struct {
	AsOfBlockNumber uint64
	Shares          uint64
	IsEigenpod      bool
}

var cache Cache

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func isEigenpod(eigenpodAddress []byte) (bool, error) {
	eigenpod := common.Bytes2Hex(eigenpodAddress)

	if val, ok := cache.PodOwnerShares[eigenpod]; ok {
		return val.IsEigenpod, nil
	}

	// Simulate fetching from contracts
	// Implement contract fetching logic here

	cache.PodOwnerShares[eigenpod] = PodOwnerShare{
		AsOfBlockNumber: 123, // replace with actual block number
		Shares:          0,
		IsEigenpod:      false,
	}

	return false, nil
}

func withdrawalCredentialsBelongsToEigenpod(withdrawalCredentials []byte) (bool, error) {
	if withdrawalCredentials[0] != 1 {
		return false, nil
	}

	withdrawalAddress := withdrawalCredentials[12:]
	return isEigenpod(withdrawalAddress)
}

func executionWithdrawlAddress(withdrawalCredentials []byte) *string {
	if withdrawalCredentials[0] != '1' {
		return nil
	}
	addr := common.Bytes2Hex(withdrawalCredentials[12:])
	return &addr
}

func aFilter[T any](coll []T, criteria func(T) bool) []T {
	var result []T
	for _, item := range coll {
		if criteria(item) {
			result = append(result, item)
		}
	}
	return result
}

func aMap[T any, A any](coll []T, mapper func(T, uint64) A) []A {
	var result []A
	for idx, item := range coll {
		result = append(result, mapper(item, uint64(idx)))
	}
	return result
}

// https://www.multicall3.com/deployments
func DeployedAddresses() map[int]string {
	return map[int]string{
		0:     "0xcA11bde05977b3631167028862bE2a173976CA11",
		17000: "0xcA11bde05977b3631167028862bE2a173976CA11",
	}
}

func MulticallEigenpod(eigenpodAddress string) multicall.Contract {
	eigenpodAbi, err := onchain.EigenPodMetaData.GetAbi()
	if err != nil {
		panic(err)
	}

	return multicall.Contract{
		ABI:     eigenpodAbi,
		Address: common.HexToAddress(eigenpodAddress),
	}
}

func scanForUnhealthyEigenpods(ctx context.Context, eth *ethclient.Client, nodeUrl string, beacon BeaconClient, chainId *big.Int) error {
	addr, ok := DeployedAddresses()[int(chainId.Int64())]
	if !ok {
		return fmt.Errorf("no known multicall deployment for chain: %d", chainId.Int64())
	}

	mc, err := multicall.Dial(context.Background(), nodeUrl, addr)
	if err != nil {
		panic(err)
	}

	beaconState, err := beacon.GetBeaconState(ctx, "head")
	if err != nil {
		return fmt.Errorf("error downloading beacon state: %s", err.Error())
	}

	// Simulate fetching validators
	_allValidators, err := beaconState.Validators()
	if err != nil {
		return err
	}

	allValidatorBalances, err := beaconState.ValidatorBalances()
	if err != nil {
		return err
	}

	allValidatorsWithIndices := aMap(_allValidators, func(validator *phase0.Validator, index uint64) ValidatorWithIndex {
		return ValidatorWithIndex{
			Validator: validator,
			Index:     index,
		}
	})

	allWithdrawalAddresses := make(map[string]struct{})
	for _, v := range allValidatorsWithIndices {
		address := executionWithdrawlAddress(v.Validator.WithdrawalCredentials)
		if address != nil {
			allWithdrawalAddresses[*address] = struct{}{}
		}
	}

	allSlashedValidators := aFilter(allValidatorsWithIndices, func(v ValidatorWithIndex) bool {
		if !v.Validator.Slashed {
			return false // we only care about slashed validators.
		}
		if v.Validator.WithdrawalCredentials[0] != 1 {
			return false // not an execution withdrawal address
		}
		return true
	})

	withdrawalAddressesToCheck := make(map[uint64]string)
	for _, validator := range allSlashedValidators {
		withdrawalAddressesToCheck[validator.Index] = *executionWithdrawlAddress(validator.Validator.WithdrawalCredentials)
	}

	if len(withdrawalAddressesToCheck) == 0 {
		log.Println("No EigenValidators were slashed.")
		return nil
	}

	// now, check across $withdrawalAddressesToCheck
	potentialEigenpods := make([]multicall.Contract, len(withdrawalAddressesToCheck))
	numPotentialPods := 0
	for _, withdrawalAddress := range withdrawalAddressesToCheck {
		potentialEigenpods[numPotentialPods] = MulticallEigenpod(withdrawalAddress)
		numPotentialPods++
	}

	loadPodOwners := aMap(potentialEigenpods, func(eigenpod multicall.Contract) *multicall.Call {
		return nil
	})

	// load all of the podOwners
	mc.Call(nil,
		loadPodOwners...,
	)

	log.Printf("%d EigenValidators were slashed\n", len(allSlashedValidatorsBelongingToEigenpods))

	slashedEigenpods := make(map[string][]*phase0.Validator)
	for _, validator := range allSlashedValidatorsBelongingToEigenpods {
		podAddress := executionWithdrawlAddress(validator.WithdrawalCredentials)
		if podAddress != nil {
			slashedEigenpods[*podAddress] = append(slashedEigenpods[*podAddress], validator)
		}
	}

	log.Printf("%d EigenPods were slashed\n", len(slashedEigenpods))

	slashedEigenpodBeaconBalances := make(map[string]phase0.Gwei)
	for idx, validator := range allValidators {
		eigenpod := executionWithdrawlAddress(validator.WithdrawalCredentials)
		if eigenpod != nil {
			isEigenpod := cache.PodOwnerShares[*eigenpod].IsEigenpod
			if isEigenpod {
				slashedEigenpodBeaconBalances[*eigenpod] += allValidatorBalances[idx]
			}
		}
	}

	var unhealthyEigenpods []string
	for pod, balance := range slashedEigenpodBeaconBalances {
		executionBalance := cache.PodOwnerShares[pod].Shares
		if executionBalance == 0 {
			continue
		}
		if balance <= phase0.Gwei(float64(executionBalance)*0.95) {
			unhealthyEigenpods = append(unhealthyEigenpods, pod)
			log.Printf("[%s] %.2f%% deviation (beacon: %d -> execution: %d)\n", pod, 100*(float64(executionBalance)-float64(balance))/float64(executionBalance), balance, executionBalance)
		}
	}

	if len(unhealthyEigenpods) == 0 {
		log.Println("All slashed eigenpods are within 5% of their expected balance.")
		return nil
	}

	log.Printf("%d EigenPods were unhealthy\n", len(unhealthyEigenpods))

	var entries []map[string]interface{}
	for _, val := range unhealthyEigenpods {
		entries = append(entries, map[string]interface{}{
			"eigenpod":          val,
			"slashedValidators": slashedEigenpods[val],
		})
	}

	fmt.Println(entries)
	return nil
}
