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
)

func PodManagerContracts() map[uint64]string {
	return map[uint64]string{
		0:     "0x91E677b07F7AF907ec9a428aafA9fc14a0d3A338",
		17000: "0x30770d7E3e71112d7A6b7259542D1f680a70e315", //testnet holesky
	}
}

type Cache struct {
	PodOwnerShares map[string]PodOwnerShare
}

type PodOwnerShare struct {
	Shares     uint64
	IsEigenpod bool
}

var cache Cache

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func isEigenpod(eth *ethclient.Client, chainId uint64, eigenpodAddress string) (bool, error) {
	if val, ok := cache.PodOwnerShares[eigenpodAddress]; ok {
		return val.IsEigenpod, nil
	}

	// default to false
	cache.PodOwnerShares[eigenpodAddress] = PodOwnerShare{
		Shares:     0,
		IsEigenpod: false,
	}

	podManAddress, ok := PodManagerContracts()[chainId]
	if !ok {
		return false, fmt.Errorf("chain %d not supported", chainId)
	}
	podMan, err := onchain.NewEigenPodManager(common.HexToAddress(podManAddress), nil)
	if err != nil {
		return false, err
	}

	pod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return false, err
	}

	owner, err := pod.PodOwner(nil)
	if err != nil {
		return false, err
	}

	expectedPod, err := podMan.OwnerToPod(nil, owner)
	if err != nil {
		return false, fmt.Errorf("ownerToPod() failed: %s", err.Error())
	}
	if expectedPod.Cmp(common.HexToAddress(eigenpodAddress)) != 0 {
		return false, nil
	}

	podOwnerShares, err := podMan.PodOwnerShares(nil, owner)
	if err != nil {
		return false, fmt.Errorf("PodOwnerShares() failed: %s", err.Error())
	}

	// Simulate fetching from contracts
	// Implement contract fetching logic here
	cache.PodOwnerShares[eigenpodAddress] = PodOwnerShare{
		Shares:     podOwnerShares.Uint64(),
		IsEigenpod: true,
	}

	return true, nil
}

// func withdrawalCredentialsBelongsToEigenpod(withdrawalCredentials []byte, eth *ethclient.Client, chainId *big.Int) (bool, error) {
// 	if withdrawalCredentials[0] != 1 {
// 		return false, nil
// 	}

// 	withdrawalAddress := withdrawalCredentials[12:]
// 	return isEigenpod(eth, chainId.Uint64(), common.Bytes2Hex(withdrawalAddress))
// }

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

func ScanForUnhealthyEigenpods(ctx context.Context, eth *ethclient.Client, nodeUrl string, beacon BeaconClient, chainId *big.Int) ([]map[string]interface{}, error) {
	beaconState, err := beacon.GetBeaconState(ctx, "head")
	if err != nil {
		return nil, fmt.Errorf("error downloading beacon state: %s", err.Error())
	}

	// Simulate fetching validators
	_allValidators, err := beaconState.Validators()
	if err != nil {
		return nil, err
	}

	allValidatorBalances, err := beaconState.ValidatorBalances()
	if err != nil {
		return nil, err
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
		return []map[string]interface{}{}, nil
	}

	validatorToPod := map[uint64]string{}

	allSlashedValidatorsBelongingToEigenpods := aFilter(allSlashedValidators, func(validator ValidatorWithIndex) bool {
		isPod, err := isEigenpod(eth, chainId.Uint64(), *executionWithdrawlAddress(validator.Validator.WithdrawalCredentials))
		if err != nil {
			return false
		}
		return isPod
	})

	log.Printf("%d EigenValidators were slashed\n", len(allSlashedValidatorsBelongingToEigenpods))

	slashedEigenpods := make(map[string][]ValidatorWithIndex)
	for _, validator := range allSlashedValidatorsBelongingToEigenpods {
		podAddress := executionWithdrawlAddress(validator.Validator.WithdrawalCredentials)
		if podAddress != nil {
			slashedEigenpods[*podAddress] = append(slashedEigenpods[*podAddress], validator)
			validatorToPod[validator.Index] = *podAddress
		}
	}

	log.Printf("%d EigenPods were slashed\n", len(slashedEigenpods))

	var unhealthyEigenpods []string
	for _, validator := range allSlashedValidatorsBelongingToEigenpods {
		balance := allValidatorBalances[validator.Index]
		pod := validatorToPod[validator.Index]
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
		return []map[string]interface{}{}, nil
	}

	log.Printf("%d EigenPods were unhealthy\n", len(unhealthyEigenpods))

	var entries []map[string]interface{}
	for _, val := range unhealthyEigenpods {
		entries = append(entries, map[string]interface{}{
			"eigenpod":          val,
			"slashedValidators": slashedEigenpods[val],
		})
	}

	return entries, nil
}
