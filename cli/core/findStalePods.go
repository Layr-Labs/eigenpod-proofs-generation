package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/DelegationManager"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPodManager"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/IEigenPod"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/utils"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jbrower95/multicall-go"
)

func PodManagerContracts() map[uint64]string {
	return map[uint64]string{
		1:     "0x91E677b07F7AF907ec9a428aafA9fc14a0d3A338",
		17000: "0x30770d7E3e71112d7A6b7259542D1f680a70e315", //testnet holesky
	}
}

// multiply by a fraction
func FracMul(a *big.Int, x *big.Int, y *big.Int) *big.Int {
	_a := new(big.Int).Mul(a, x)
	return _a.Div(_a, y)
}

func executionWithdrawalAddress(withdrawalCredentials []byte) *string {
	if withdrawalCredentials[0] != 1 {
		return nil
	}
	addr := common.Bytes2Hex(withdrawalCredentials[12:])
	return &addr
}

func validEigenpodsOnly(candidateAddresses []common.Address, mc *multicall.MulticallClient, chainId uint64, eth *ethclient.Client) ([]common.Address, error) {
	EigenPodAbi, err := abi.JSON(strings.NewReader(EigenPod.EigenPodABI))
	if err != nil {
		return nil, fmt.Errorf("failed to load eigenpod abi: %s", err)
	}
	EigenPodManagerAbi, err := abi.JSON(strings.NewReader(EigenPodManager.EigenPodManagerABI))
	if err != nil {
		return nil, fmt.Errorf("failed to load eigenpod manager abi: %s", err)
	}

	podManagerAddress, ok := PodManagerContracts()[chainId]
	if !ok {
		return nil, fmt.Errorf("Unsupported chainId: %d", chainId)
	}

	////// step 1: cast all addresses to EigenPod, and attempt to read the pod owner.
	var lastError error

	calls := utils.Map(candidateAddresses, func(addr common.Address, i uint64) *multicall.MultiCallMetaData[common.Address] {
		mc, err := multicall.Describe[common.Address](
			addr,
			EigenPodAbi,
			"podOwner",
		)
		if err != nil {
			lastError = err
			return nil
		}
		return mc
	})
	if lastError != nil {
		return nil, lastError
	}

	reportedPodOwners, err := multicall.DoManyAllowFailures(
		mc,
		calls...,
	)
	if err != nil || reportedPodOwners == nil {
		return nil, fmt.Errorf("failed to load podOwners: %w", err)
	}

	type PodOwnerResult struct {
		Query    common.Address
		Response multicall.TypedMulticall3Result[*common.Address]
	}

	podOwnerPairs := utils.Filter(utils.Map(*reportedPodOwners, func(res multicall.TypedMulticall3Result[*common.Address], i uint64) PodOwnerResult {
		return PodOwnerResult{
			Query:    candidateAddresses[i],
			Response: res,
		}
	}), func(m PodOwnerResult) bool {
		return m.Response.Success
	})

	////// step 2: using the pod manager, check `ownerToPod` and validate which ones point back at the same address.
	authoritativeOwnerToPodCalls := utils.Map(podOwnerPairs, func(res PodOwnerResult, i uint64) *multicall.MultiCallMetaData[common.Address] {
		mc, err := multicall.Describe[common.Address](
			common.HexToAddress(podManagerAddress),
			EigenPodManagerAbi,
			"ownerToPod",
			res.Response.Value,
		)
		if err != nil {
			lastError = err
			return nil
		}
		return mc
	})
	if lastError != nil {
		return nil, lastError
	}

	authoritativeOwnerToPod, err := multicall.DoMany(mc, authoritativeOwnerToPodCalls...)
	nullAddress := common.BigToAddress(big.NewInt(0))

	////// step 3: the valid eigenrestpods are the ones where authoritativeOwnerToPod[i] == candidateAddresses[i].
	return utils.Map(utils.FilterI(podOwnerPairs, func(res PodOwnerResult, i uint64) bool {
		return (res.Query.Cmp(*(*authoritativeOwnerToPod)[i]) == 0) && (*authoritativeOwnerToPod)[i].Cmp(nullAddress) != 0
	}), func(v PodOwnerResult, i uint64) common.Address {
		return v.Query
	}), nil
}

func ComputeBalanceDeviationSync(ctx context.Context, eth *ethclient.Client, state *spec.VersionedBeaconState, eigenpod common.Address) (*big.Float, error) {
	pod, err := IEigenPod.NewIEigenPod(eigenpod, eth)
	if err != nil {
		return nil, err
	}

	allValidators, err := state.Validators()
	PanicOnError("failed to read validators", err)

	allValidatorsWithIndexes := utils.Map(allValidators, func(v *phase0.Validator, i uint64) ValidatorWithIndex {
		return ValidatorWithIndex{
			Validator: v,
			Index:     i,
		}
	})

	podValidators := utils.FilterI[ValidatorWithIndex](allValidatorsWithIndexes, func(v ValidatorWithIndex, u uint64) bool {
		addr := executionWithdrawalAddress(v.Validator.WithdrawalCredentials)
		return addr != nil && eigenpod.Cmp(common.HexToAddress(*addr)) == 0
	})

	validatorBalances, err := state.ValidatorBalances()
	PanicOnError("failed to read beacon state validator balances", err)

	validatorInfo, err := FetchMultipleOnchainValidatorInfoWithFailures(ctx, eth, eigenpod.Hex(), podValidators)
	if err != nil {
		return nil, err
	}

	podBalanceWei, err := eth.BalanceAt(ctx, eigenpod, nil)
	if err != nil {
		return nil, err
	}

	sumCurrentBeaconBalancesGwei := utils.BigSum(
		utils.Map(podValidators, func(v ValidatorWithIndex, i uint64) *big.Int {
			if validatorInfo[i].Info != nil && validatorInfo[i].Info.Status == 1 /* ACTIVE */ {
				return new(big.Int).SetUint64(uint64(validatorBalances[v.Index]))
			}
			return big.NewInt(0)
		}),
	)

	eigenPodManagerAddr, err := pod.EigenPodManager(nil)
	PanicOnError("failed to load eigenpod manager", err)

	eigenPodManager, err := EigenPodManager.NewEigenPodManager(eigenPodManagerAddr, eth)
	PanicOnError("failed to load eigenpod manager", err)

	delegationManagerAddress, err := eigenPodManager.DelegationManager(nil)
	PanicOnError("failed to read delegationManager", err)

	delegationManager, err := DelegationManager.NewDelegationManager(delegationManagerAddress, eth)
	PanicOnError("failed to reach delegationManager", err)

	podOwner, err := pod.PodOwner(nil)
	PanicOnError("failed to load pod owner", err)

	activeShares, err := delegationManager.GetWithdrawableShares(nil, podOwner, []common.Address{
		common.HexToAddress(NATIVE_ETH_STRATEGY),
	})
	PanicOnError("failed to load owner shares", err)

	var sharesPendingWithdrawal *big.Int
	withdrawalInfo, err := delegationManager.GetQueuedWithdrawals(nil, podOwner)
	PanicOnError("failed to load queued withdrawals", err)

	for i, withdrawal := range withdrawalInfo.Withdrawals {
		for j, strategy := range withdrawal.Strategies {
			if strategy.Cmp(common.HexToAddress(NATIVE_ETH_STRATEGY)) == 0 {
				sharesPendingWithdrawal = new(big.Int).Add(sharesPendingWithdrawal, withdrawalInfo.Shares[i][j])
			}
		}
	}

	totalSharesInEigenLayer := new(big.Int).Add(activeShares.WithdrawableShares[0], sharesPendingWithdrawal)

	// fmt.Printf("# validators: %d\n", len(podValidators))
	// fmt.Printf("# active validators: %d\n", len(activeValidators))
	// fmt.Printf("delta := 1 - ((podBalanceGwei + sumCurrentBeaconBalancesGwei) / (regGwei + sumPreviousBeaconBalancesGwei)\n")
	// fmt.Printf("delta := 1 - ((%s + %s) / (%d + %s)\n", WeiToGwei(podBalanceWei).String(), sumCurrentBeaconBalancesGwei.String(), regGwei, sumPreviousBeaconBalancesGwei.String())

	currentState := new(big.Float).Add(WeiToGwei(podBalanceWei), new(big.Float).SetInt(sumCurrentBeaconBalancesGwei))
	prevState := WeiToGwei(totalSharesInEigenLayer)

	var delta *big.Float

	if prevState.Cmp(big.NewFloat(0)) == 0 {
		delta = big.NewFloat(0)
	} else {
		delta = new(big.Float).Sub(
			big.NewFloat(1),
			new(big.Float).Quo(
				currentState,
				prevState,
			),
		)
	}

	// fmt.Printf("(delta=%s%%)\n", new(big.Float).Mul(delta, big.NewFloat(100)).String())
	// fmt.Printf("-----------------------------------\n\n")

	return delta, nil
}

func FindStaleEigenpods(ctx context.Context, eth *ethclient.Client, nodeUrl string, beacon BeaconClient, chainId *big.Int, verbose bool, tolerance float64) (map[string][]ValidatorWithIndex, error) {
	beaconState, err := beacon.GetBeaconState(ctx, "head")
	if err != nil {
		return nil, fmt.Errorf("error downloading beacon state: %s", err.Error())
	}

	mc, err := multicall.NewMulticallClient(ctx, eth, nil)
	if err != nil {
		return nil, err
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

	allSlashedWithdrawalAddresses := utils.Unique(
		utils.Map(allSlashedValidators, func(v ValidatorWithIndex, i uint64) common.Address {
			return common.HexToAddress(*executionWithdrawalAddress(v.Validator.WithdrawalCredentials))
		}),
	)

	// fmt.Printf("Checking %d slashed withdrawal addresses for eigenpod status\n", len(allSlashedWithdrawalAddresses))

	slashedEigenpods, err := validEigenpodsOnly(allSlashedWithdrawalAddresses, mc, chainId.Uint64(), eth)

	if len(slashedEigenpods) == 0 {
		log.Println("No eigenpods were slashed.")
		return map[string][]ValidatorWithIndex{}, nil
	}

	// 2. given the set of slashed eigenpods, determine which are unhealthy.

	if verbose {
		log.Printf("%d EigenPods were slashed\n", len(slashedEigenpods))
	}

	unhealthyEigenpods := utils.Filter(slashedEigenpods, func(eigenpod common.Address) bool {
		deviation, err := ComputeBalanceDeviationSync(ctx, eth, beaconState, eigenpod)
		PanicOnError("failed to compute balance deviation for eigenpod", err)

		return deviation.Cmp(big.NewFloat(tolerance)) > 0
	})

	if len(unhealthyEigenpods) == 0 {
		if verbose {
			log.Printf("All slashed eigenpods are within %f%% of their expected balance.\n", tolerance)
		}
		return map[string][]ValidatorWithIndex{}, nil
	}

	if verbose {
		log.Printf("%d EigenPods were unhealthy\n", len(unhealthyEigenpods))
	}

	var entries map[string][]ValidatorWithIndex = make(map[string][]ValidatorWithIndex)
	for _, val := range unhealthyEigenpods {
		entries[val.Hex()] = utils.Filter(allValidatorsWithIndices, func(v ValidatorWithIndex) bool {
			execAddr := executionWithdrawalAddress(v.Validator.WithdrawalCredentials)
			return execAddr != nil && common.HexToAddress(*execAddr).Cmp(val) == 0
		})
	}

	return entries, nil
}
