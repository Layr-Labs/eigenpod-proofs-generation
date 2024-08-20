package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
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

func keys[A comparable, B any](coll map[A]B) []A {
	if len(coll) == 0 {
		return []A{}
	}
	out := make([]A, len(coll))
	for key, _ := range coll {
		out[len(out)] = key
	}
	return out
}

type PodOwnerShare struct {
	Shares     uint64
	NativeETH  phase0.Gwei
	IsEigenpod bool
}

const ACCEPTABLE_BALANCE_DEVIATION = float64(0.95)

var cache Cache

func isEigenpod(eth *ethclient.Client, chainId uint64, eigenpodAddress string) (bool, error) {
	if cache.PodOwnerShares == nil {
		cache.PodOwnerShares = make(map[string]PodOwnerShare)
	}

	if val, ok := cache.PodOwnerShares[eigenpodAddress]; ok {
		return val.IsEigenpod, nil
	}

	// default to false
	cache.PodOwnerShares[eigenpodAddress] = PodOwnerShare{
		Shares:     0,
		NativeETH:  phase0.Gwei(0),
		IsEigenpod: false,
	}

	podManAddress, ok := PodManagerContracts()[chainId]
	if !ok {
		return false, fmt.Errorf("chain %d not supported", chainId)
	}
	podMan, err := onchain.NewEigenPodManager(common.HexToAddress(podManAddress), eth)
	if err != nil {
		return false, err
	}

	if podMan == nil {
		return false, errors.New("failed to find eigenpod manager")
	}

	pod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	if err != nil {
		return false, err
	}

	owner, err := pod.PodOwner(nil)
	if err != nil {
		return false, err
	}

	expectedPod, err := podMan.OwnerToPod(&bind.CallOpts{}, owner)
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

	balance, err := eth.BalanceAt(context.Background(), common.HexToAddress(eigenpodAddress), nil)
	if err != nil {
		return false, fmt.Errorf("balance check failed: %s", err.Error())
	}

	// Simulate fetching from contracts
	// Implement contract fetching logic here
	cache.PodOwnerShares[eigenpodAddress] = PodOwnerShare{
		Shares:     podOwnerShares.Uint64(),
		NativeETH:  phase0.Gwei(balance.Uint64()),
		IsEigenpod: true,
	}

	return true, nil
}

func executionWithdrawalAddress(withdrawalCredentials []byte) *string {
	if withdrawalCredentials[0] != 1 {
		return nil
	}
	addr := common.Bytes2Hex(withdrawalCredentials[12:])
	return &addr
}

func FindStaleEigenpods(ctx context.Context, eth *ethclient.Client, nodeUrl string, beacon BeaconClient, chainId *big.Int, verbose bool) (map[string][]ValidatorWithIndex, error) {
	beaconState, err := beacon.GetBeaconState(ctx, "head")
	if err != nil {
		return nil, fmt.Errorf("error downloading beacon state: %s", err.Error())
	}

	// Simulate fetching validators
	_allValidators, err := beaconState.Validators()
	if err != nil {
		return nil, err
	}

	allValidatorsWithIndices := utils.Map(_allValidators, func(validator *phase0.Validator, index uint64) ValidatorWithIndex {
		return ValidatorWithIndex{
			Validator: validator,
			Index:     index,
		}
	})

	// TODO(pectra): this logic changes after the pectra upgrade.
	allSlashedValidators := utils.Filter(allValidatorsWithIndices, func(v ValidatorWithIndex) bool {
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
		withdrawalAddressesToCheck[validator.Index] = *executionWithdrawalAddress(validator.Validator.WithdrawalCredentials)
	}

	if len(withdrawalAddressesToCheck) == 0 {
		log.Println("No EigenValidators were slashed.")
		return map[string][]ValidatorWithIndex{}, nil
	}

	validatorToPod := map[uint64]string{}

	allSlashedValidatorsBelongingToEigenpods := utils.Filter(allSlashedValidators, func(validator ValidatorWithIndex) bool {
		isPod, err := isEigenpod(eth, chainId.Uint64(), *executionWithdrawalAddress(validator.Validator.WithdrawalCredentials))
		if err != nil {
			return false
		}
		return isPod
	})

	allValidatorInfo := make(map[uint64]onchain.IEigenPodValidatorInfo)
	for _, validator := range allSlashedValidatorsBelongingToEigenpods {
		eigenpodAddress := *executionWithdrawalAddress(validator.Validator.WithdrawalCredentials)
		pod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
		if err != nil {
			// failed to load validator info.
			return map[string][]ValidatorWithIndex{}, fmt.Errorf("failed to dial eigenpod: %s", err.Error())
		}

		info, err := pod.ValidatorPubkeyToInfo(nil, validator.Validator.PublicKey[:])
		if err != nil {
			// failed to load validator info.
			return map[string][]ValidatorWithIndex{}, err
		}
		allValidatorInfo[validator.Index] = info
	}

	allActiveSlashedValidatorsBelongingToEigenpods := utils.Filter(allSlashedValidatorsBelongingToEigenpods, func(validator ValidatorWithIndex) bool {
		validatorInfo := allValidatorInfo[validator.Index]
		return validatorInfo.Status == 1 // "ACTIVE"
	})

	if verbose {
		log.Printf("%d EigenValidators were slashed\n", len(allActiveSlashedValidatorsBelongingToEigenpods))
	}

	slashedEigenpods := utils.Reduce(allActiveSlashedValidatorsBelongingToEigenpods, func(pods map[string][]ValidatorWithIndex, validator ValidatorWithIndex) map[string][]ValidatorWithIndex {
		podAddress := executionWithdrawalAddress(validator.Validator.WithdrawalCredentials)
		if podAddress != nil {
			if pods[*podAddress] == nil {
				pods[*podAddress] = []ValidatorWithIndex{}
			}
			pods[*podAddress] = append(pods[*podAddress], validator)
			validatorToPod[validator.Index] = *podAddress
		}
		return pods
	}, map[string][]ValidatorWithIndex{})

	allValidatorBalances, err := beaconState.ValidatorBalances()
	if err != nil {
		return nil, err
	}

	totalAssetsGweiByEigenpod := utils.Reduce(keys(slashedEigenpods), func(allBalances map[string]phase0.Gwei, eigenpod string) map[string]phase0.Gwei {
		// total assets of an eigenpod are determined as;
		//	SUM(
		//		- native ETH in the pod
		//		- any active validators and their associated balances
		// 	)
		allEigenpodsForValidator := utils.Filter(allValidatorsWithIndices, func(v ValidatorWithIndex) bool {
			withdrawal := executionWithdrawalAddress(v.Validator.WithdrawalCredentials)
			return withdrawal != nil && strings.EqualFold(*withdrawal, eigenpod)
		})

		allValidatorBalancesSummed := utils.Reduce(allEigenpodsForValidator, func(accum phase0.Gwei, validator ValidatorWithIndex) phase0.Gwei {
			return accum + allValidatorBalances[validator.Index]
		}, phase0.Gwei(0))

		balance := phase0.Gwei(cache.PodOwnerShares[eigenpod].NativeETH) + allValidatorBalancesSummed
		allBalances[eigenpod] = phase0.Gwei(balance)
		return allBalances
	}, map[string]phase0.Gwei{})

	if verbose {
		log.Printf("%d EigenPods were slashed\n", len(slashedEigenpods))
	}

	unhealthyEigenpods := utils.Filter(keys(slashedEigenpods), func(eigenpod string) bool {
		balance, ok := totalAssetsGweiByEigenpod[eigenpod]
		if !ok {
			return false
		}
		executionBalance := cache.PodOwnerShares[eigenpod].Shares
		if balance <= phase0.Gwei(float64(executionBalance)*ACCEPTABLE_BALANCE_DEVIATION) {
			if verbose {
				log.Printf("[%s] %.2f%% deviation (beacon: %d -> execution: %d)\n", eigenpod, 100*(float64(executionBalance)-float64(balance))/float64(executionBalance), balance, executionBalance)
			}
			return true
		}

		return false
	})

	if len(unhealthyEigenpods) == 0 {
		if verbose {
			log.Println("All slashed eigenpods are within 5% of their expected balance.")
		}
		return map[string][]ValidatorWithIndex{}, nil
	}

	if verbose {
		log.Printf("%d EigenPods were unhealthy\n", len(unhealthyEigenpods))
	}

	var entries map[string][]ValidatorWithIndex = make(map[string][]ValidatorWithIndex)
	for _, val := range unhealthyEigenpods {
		entries[val] = slashedEigenpods[val]
	}

	return entries, nil
}
