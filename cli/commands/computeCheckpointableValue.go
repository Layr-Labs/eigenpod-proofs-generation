package commands

import (
	"context"
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/multicall"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type TComputeCheckpointableValueCommandArgs struct {
	Node       string
	BeaconNode string
}

// TODO: this is duplicated
func PodManagerContracts() map[uint64]string {
	return map[uint64]string{
		1:     "0x91E677b07F7AF907ec9a428aafA9fc14a0d3A338",
		17000: "0x30770d7E3e71112d7A6b7259542D1f680a70e315", //testnet holesky
	}
}

//go:embed multicallAbi.json
var multicallAbi string

func ComputeCheckpointableValueCommand(args TComputeCheckpointableValueCommandArgs) error {
	ctx := context.Background()

	eigenpodAbi, err := abi.JSON(strings.NewReader(onchain.EigenPodABI))
	core.PanicOnError("failed to load eigenpod abi", err)

	podManagerAbi, err := abi.JSON(strings.NewReader(onchain.EigenPodManagerABI))
	core.PanicOnError("failed to load eigenpod manager abi", err)

	eth, beaconClient, chainId, err := core.GetClients(ctx, args.Node, args.BeaconNode, true)
	core.PanicOnError("failed to reach ethereum clients", err)

	podManagerAddress, ok := PodManagerContracts()[chainId.Uint64()]
	if !ok {
		core.Panic("unsupported network")
	}

	// fetch latest beacon state.
	beaconState, err := beaconClient.GetBeaconState(ctx, "head")
	core.PanicOnError("failed to load beacon state", err)

	allBalances, err := beaconState.ValidatorBalances()
	core.PanicOnError("failed to parse beacon balances", err)

	_allValidators, err := beaconState.Validators()
	core.PanicOnError("failed to fetch validators", err)
	allValidators := utils.Map(_allValidators, func(validator *phase0.Validator, i uint64) core.ValidatorWithIndex {
		return core.ValidatorWithIndex{
			Validator: validator,
			Index:     i,
		}
	})

	// see which validators are eigenpods
	//
	// 1. can ignore anything that isn't withdrawing to the execution chain.
	executionLayerWithdrawalCredentialValidators := utils.Filter(allValidators, func(validator core.ValidatorWithIndex) bool {
		return validator.Validator.WithdrawalCredentials[0] == 1
	})

	// multicall request a bunch of different validators to check whether they're eigenpods
	mc, err := multicall.NewMulticallClient(ctx, eth, &multicall.TMulticallClientOptions{
		MaxBatchSizeBytes: 8192,
	})
	core.PanicOnError("error initializing mc", err)

	interestingWithdrawalAddresses := getKeys(utils.Reduce(executionLayerWithdrawalCredentialValidators, func(accum map[string]int, next core.ValidatorWithIndex) map[string]int {
		accum[common.Bytes2Hex(next.Validator.WithdrawalCredentials[12:])] = 1
		return accum
	}, map[string]int{}))

	fmt.Printf("Querying %d addresses to see if they may be eigenpods\n", len(interestingWithdrawalAddresses))
	podOwners, err := multicall.DoMultiCallManyReportingFailures[*common.Address](*mc, utils.Map(interestingWithdrawalAddresses, func(address string, index uint64) *multicall.MultiCallMetaData[*common.Address] {
		callMeta, err := multicall.MultiCall(
			common.HexToAddress(address),
			eigenpodAbi,
			func(data []byte) (*common.Address, error) {
				res, err := eigenpodAbi.Unpack("podOwner", data)
				if err != nil {
					return nil, err
				}
				return abi.ConvertType(res[0], new(common.Address)).(*common.Address), nil
			}, "podOwner",
		)
		core.PanicOnError("failed to form mc", err)
		return callMeta
	})...)

	if podOwners == nil || err != nil || len(*podOwners) == 0 {
		core.PanicOnError("failed to fetch podOwners", err)
		core.Panic("loaded no pod owners")
		return err
	}

	// now we can filter by which addresses actually claimed to have podOwner()
	podToPodOwner := map[string]*common.Address{}
	addressesWithPodOwners := utils.FilterI(interestingWithdrawalAddresses, func(address string, i uint64) bool {
		success := (*podOwners)[i].Success
		if success {
			podToPodOwner[address] = (*podOwners)[i].Value
		}
		return success
	})

	printAsJSON(podToPodOwner)

	// array[eigenpods given the owner]
	fmt.Printf("Querying %d addresses (podMan=%s) to see if it knows about these eigenpods\n", len(addressesWithPodOwners), podManagerAddress)

	printAsJSON(podManagerAbi.Methods)
	eigenpodForOwner, err := multicall.DoMultiCallManyReportingFailures(
		*mc,
		utils.Map(addressesWithPodOwners, func(address string, i uint64) *multicall.MultiCallMetaData[*common.Address] {
			claimedOwner := *podToPodOwner[address]
			call, err := multicall.MultiCall(
				common.HexToAddress(podManagerAddress),
				podManagerAbi,
				func(data []byte) (*common.Address, error) {
					res, err := podManagerAbi.Unpack("ownerToPod", data)
					if err != nil {
						return nil, err
					}
					return abi.ConvertType(res[0], new(common.Address)).(*common.Address), nil
				},
				"ownerToPod",
				claimedOwner,
			)
			core.PanicOnError("failed to form multicall", err)
			return call
		})...,
	)
	core.PanicOnError("failed to query", err)

	// now, see which of `addressesWithPodOwners` properly were eigenpods.
	allEigenpods := utils.FilterI(addressesWithPodOwners, func(address string, i uint64) bool {
		return (*eigenpodForOwner)[i].Success && (*eigenpodForOwner)[i].Value.Cmp(common.HexToAddress(addressesWithPodOwners[i])) == 0
	})
	isEigenpodSet := utils.Reduce(allEigenpods, func(allEigenpodSet map[string]int, eigenpod string) map[string]int {
		allEigenpodSet[eigenpod] = 1
		return allEigenpodSet
	}, map[string]int{})

	printAsJSON(allEigenpods)

	fmt.Printf("%d eigenpods discovered on the network", len(allEigenpods))

	// now, determine their pending checkpoint rewards per checkpoint
	pendingRewardsWeiPerEigenpod := map[string]*big.Int{}

	// TODO: compute pending rewards per eigenpod;
	//		see: https://github.com/Layr-Labs/eigenlayer-contracts/blob/dev/src/contracts/pods/EigenPod.sol#L656
	//			return (podBalanceGwei + checkpoint.balanceDeltasGwei)
	//
	//			where
	//
	//				podBalanceGwei = address(pod).balanceGwei - pod.withdrawableRestakedExecutionLayerGwei
	//			and
	//				checkpoint.balanceDeltasGwei = sumBeaconBalancesGwei - sumRestakedBalancesGwei

	multicallAbiRef, err := abi.JSON(strings.NewReader(multicallAbi))
	core.PanicOnError("failed to load multicall abi", err)

	fmt.Printf("Loading address(EigenPod).balance....\n")
	podNativeEthBalances, err := multicall.DoMultiCallMany(
		*mc,
		utils.Map(allEigenpods, func(eigenpod string, index uint64) *multicall.MultiCallMetaData[*big.Int] {
			call, err := multicall.MultiCall(
				common.HexToAddress("0xcA11bde05977b3631167028862bE2a173976CA11"),
				multicallAbiRef,
				func(b []byte) (*big.Int, error) {
					res, err := multicallAbiRef.Unpack("getEthBalance", b)
					if err != nil {
						return nil, err
					}
					out := abi.ConvertType(res[0], new(big.Int)).(*big.Int)
					return out, nil
				},
				"getEthBalance",
				common.HexToAddress(eigenpod),
			)
			core.PanicOnError("failed to form multicall", err)
			return call
		})...,
	)
	if err != nil || podNativeEthBalances == nil {
		core.PanicOnError("failed to multicall eigenpod balances", err)
		core.Panic("failed to load eigenpod balances")
	}

	fmt.Printf("Loading EigenPod.withdrawableRestakedExecutionLayerGwei....\n")
	withdrawableRestakedExecutionLayerGwei, err := multicall.DoMultiCallMany(
		*mc,
		utils.Map(allEigenpods, func(eigenpod string, index uint64) *multicall.MultiCallMetaData[uint64] {
			call, err := multicall.MultiCall(
				common.HexToAddress(eigenpod),
				eigenpodAbi,
				func(b []byte) (uint64, error) {
					res, err := eigenpodAbi.Unpack("withdrawableRestakedExecutionLayerGwei", b)
					if err != nil {
						return 0, err
					}
					out := abi.ConvertType(res[0], new(uint64)).(*uint64)
					return *out, nil
				},
				"withdrawableRestakedExecutionLayerGwei",
			)
			core.PanicOnError("failed to form multicall", err)
			return call
		})...,
	)
	if err != nil || withdrawableRestakedExecutionLayerGwei == nil {
		core.PanicOnError("failed to multicall eigenpod.withdrawableRestakedExecutionLayerGwei", err)
		core.Panic("failed to load eigenpod withdrawableRestakedExecutionLayerGwei")
	}

	pendingExecutionWei := utils.Map(allEigenpods, func(eigenpod string, index uint64) *big.Int {
		podCurrentNativeWei := (*podNativeEthBalances)[index]
		podEigenLayerAwareWei := core.IGweiToWei(new(big.Int).SetUint64((*withdrawableRestakedExecutionLayerGwei)[index]))
		return new(big.Int).Sub(podCurrentNativeWei, podEigenLayerAwareWei)
	})

	allValidatorsForEigenpod := utils.Reduce(allValidators, func(validatorsByPod map[string][]core.ValidatorWithIndex, validator core.ValidatorWithIndex) map[string][]core.ValidatorWithIndex {
		withdrawalAddress := common.BytesToAddress(validator.Validator.WithdrawalCredentials[12:])
		eigenpod := withdrawalAddress.Hex()
		if isEigenpodSet[eigenpod] == 1 {
			if validatorsByPod[eigenpod] == nil {
				validatorsByPod[eigenpod] = []core.ValidatorWithIndex{}
			}

			validatorsByPod[eigenpod] = append(validatorsByPod[eigenpod], validator)
		}
		return validatorsByPod
	}, map[string][]core.ValidatorWithIndex{})

	type ValidatorPodPair struct {
		Validator core.ValidatorWithIndex
		Pod       string
	}

	allEigenlayerValidatorsWithPod := utils.Reduce(getKeys(allValidatorsForEigenpod), func(validators []ValidatorPodPair, eigenpod string) []ValidatorPodPair {
		validators = append(validators, utils.Map(allValidatorsForEigenpod[eigenpod], func(validator core.ValidatorWithIndex, i uint64) ValidatorPodPair {
			return ValidatorPodPair{
				Validator: validator,
				Pod:       eigenpod,
			}
		})...)
		return validators
	}, []ValidatorPodPair{})

	allValidatorInfoRequests := utils.Map(allEigenlayerValidatorsWithPod, func(validatorPodPair ValidatorPodPair, index uint64) *multicall.MultiCallMetaData[*onchain.IEigenPodValidatorInfo] {
		res, err := core.FetchMultipleOnchainValidatorInfoMulticalls(validatorPodPair.Pod, []*phase0.Validator{validatorPodPair.Validator.Validator})
		core.PanicOnError("failed to form mc", err)
		return res[0]
	})
	allValidatorInfo, err := multicall.DoMultiCallMany[*onchain.IEigenPodValidatorInfo](*mc, allValidatorInfoRequests...)
	core.PanicOnError("failed to multicall validator info", err)

	i := 0
	allValidatorInfoLookup := utils.Reduce(*allValidatorInfo, func(validatorInfoLookup map[uint64]*onchain.IEigenPodValidatorInfo, cur *onchain.IEigenPodValidatorInfo) map[uint64]*onchain.IEigenPodValidatorInfo {
		validatorInfoLookup[allEigenlayerValidatorsWithPod[i].Validator.Index] = cur
		i++
		return validatorInfoLookup
	}, map[uint64]*onchain.IEigenPodValidatorInfo{})

	beaconBalancesWei := utils.Map(allEigenlayerValidatorsWithPod, func(validatorPodPair ValidatorPodPair, i uint64) *big.Int {
		validatorInfo := (*allValidatorInfo)[i]
		if validatorInfo.Status != core.ValidatorStatusActive {
			return big.NewInt(0)
		}

		return new(big.Int).SetUint64(uint64(allBalances[validatorPodPair.Validator.Index]))
	})

	sumBeaconBalancesWei := utils.Reduce(beaconBalancesWei, func(sum *big.Int, cur *big.Int) *big.Int {
		return sum.Add(sum, cur)
	}, big.NewInt(0))

	sumRestakedBalancesWei := utils.Reduce(allEigenlayerValidatorsWithPod, func(sum *big.Int, cur ValidatorPodPair) *big.Int {
		info := allValidatorInfoLookup[cur.Validator.Index]
		if info.Status != core.ValidatorStatusActive {
			return big.NewInt(0)
		}

		restakedBalance := core.IGweiToWei(new(big.Int).SetUint64(info.RestakedBalanceGwei))
		return sum.Add(sum, restakedBalance)
	}, big.NewInt(0))

	pendingBeaconWei := big.NewInt(0).Sub(sumBeaconBalancesWei, sumRestakedBalancesWei)

	for i, eigenpod := range allEigenpods {
		pendingRewardsWeiPerEigenpod[eigenpod] = pendingExecutionWei[i]
	}

	// aggregate result
	totalPendingRewards := big.NewInt(0)
	for _, reward := range pendingRewardsWeiPerEigenpod {
		totalPendingRewards.Add(totalPendingRewards, reward)
	}
	totalPendingRewards = totalPendingRewards.Add(totalPendingRewards, pendingBeaconWei)

	totalRewards := map[string]*big.Float{
		"total_pending_shares_wei":  new(big.Float).SetInt(totalPendingRewards),
		"total_pending_shares_gwei": core.WeiToGwei(totalPendingRewards),
		"total_pending_shares_eth":  core.GweiToEther(core.WeiToGwei(totalPendingRewards)),
	}
	printAsJSON(totalRewards)
	return nil
}
