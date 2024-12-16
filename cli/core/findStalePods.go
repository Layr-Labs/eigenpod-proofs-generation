package core

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"

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

func validEigenpodsOnly(candidateAddresses []*common.Address, mc *multicall.MulticallClient, chainId uint64, eth *ethclient.Client) ([]*common.Address, error) {
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

	calls := utils.Map(candidateAddresses, func(addr *common.Address, i uint64) *multicall.MultiCallMetaData[common.Address] {
		mc, err := multicall.Describe[common.Address](
			*addr,
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

	////// step 2: using the pod manager, check `ownerToPod` and validate which ones point back at the same address.
	authoritativeOwnerToPodCalls := utils.Map(candidateAddresses, func(addr *common.Address, i uint64) *multicall.MultiCallMetaData[common.Address] {
		mc, err := multicall.Describe[common.Address](
			common.HexToAddress(podManagerAddress),
			EigenPodManagerAbi,
			"ownerToPod",
			(*reportedPodOwners)[i],
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

	authoritativeOwnerToPod, err := multicall.DoManyAllowFailures(mc, authoritativeOwnerToPodCalls...)

	////// step 3: the valid eigenrestpods are the ones where authoritativeOwnerToPod[i] == candidateAddresses[i].
	return utils.FilterI(candidateAddresses, func(candidate *common.Address, i uint64) bool {
		return (candidate.String() == (*authoritativeOwnerToPod)[i].Value.Hex())
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
		return addr != nil && *addr == eigenpod.Hex()
	})

	validatorBalances, err := state.ValidatorBalances()
	PanicOnError("failed to read beacon state validator balances", err)

	validatorInfo, err := FetchMultipleOnchainValidatorInfo(ctx, eth, eigenpod.Hex(), podValidators)
	if err != nil {
		return nil, err
	}

	podBalanceWei, err := eth.BalanceAt(ctx, eigenpod, nil)
	if err != nil {
		return nil, err
	}

	sumCurrentBeaconBalancesGwei := utils.BigSum(utils.Map(podValidators, func(v ValidatorWithIndex, i uint64) *big.Int {
		if validatorInfo[i].Info.Status == 1 /* ACTIVE */ {
			return new(big.Int).SetUint64(uint64(validatorBalances[v.Index]))
		}
		return big.NewInt(0)
	}))

	sumPreviousBeaconBalancesGwei := utils.BigSum(utils.Map(podValidators, func(v ValidatorWithIndex, i uint64) *big.Int {
		if validatorInfo[i].Info.Status == 1 /* ACTIVE */ {
			return new(big.Int).SetUint64(validatorInfo[i].Info.RestakedBalanceGwei)
		}
		return big.NewInt(0)
	}))

	// TODO: when bindings are updated, call `restakedExecutionLayerGwei`.
	regGwei, err := pod.WithdrawableRestakedExecutionLayerGwei(nil)
	PanicOnError("failed to load restakedExecutionLayerGwei", err)

	return new(big.Float).Sub(
		big.NewFloat(1),
		new(big.Float).Quo(
			new(big.Float).Add(WeiToGwei(podBalanceWei), new(big.Float).SetInt(sumCurrentBeaconBalancesGwei)),
			new(big.Float).Add(new(big.Float).SetUint64(regGwei), new(big.Float).SetInt(sumPreviousBeaconBalancesGwei)),
		),
	), nil
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

	slashedEigenpods, err := validEigenpodsOnly(utils.Map(allSlashedValidators, func(v ValidatorWithIndex, i uint64) *common.Address {
		addr := common.HexToAddress(*executionWithdrawalAddress(v.Validator.WithdrawalCredentials))
		return &addr
	}), mc, chainId.Uint64(), eth)

	if len(slashedEigenpods) == 0 {
		log.Println("No eigenpods were slashed.")
		return map[string][]ValidatorWithIndex{}, nil
	}

	// 2. given the set of slashed eigenpods, determine which are unhealthy.

	if verbose {
		log.Printf("%d EigenPods were slashed\n", len(slashedEigenpods))
	}

	unhealthyEigenpods := utils.Filter(slashedEigenpods, func(eigenpod *common.Address) bool {
		deviation, err := ComputeBalanceDeviationSync(ctx, eth, beaconState, *eigenpod)
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
			return *executionWithdrawalAddress(v.Validator.WithdrawalCredentials) == val.Hex()
		})
	}

	return entries, nil
}
