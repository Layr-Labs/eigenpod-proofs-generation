package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Checkpoint struct {
	ProofsRemaining uint64
	StartedAt       uint64
}

type Validator struct {
	Slashed                             bool
	Index                               uint64
	Status                              int
	PublicKey                           string
	IsAwaitingWithdrawalCredentialProof bool
	EffectiveBalance                    uint64
	CurrentBalance                      uint64
}

type EigenpodStatus struct {
	Validators map[string]Validator

	ActiveCheckpoint *Checkpoint

	NumberValidatorsToCheckpoint int

	CurrentTotalSharesETH *big.Float
	Status                int

	// if you completed a new checkpoint right now, how many shares would you get?
	//
	//  this is computed as:
	// 		- If checkpoint is already started:
	// 			sum(beacon chain balances) + currentCheckpoint.podBalanceGwei + pod.withdrawableRestakedExecutionLayerGwei()
	// 		- If no checkpoint is started:
	// 			total_shares_after_checkpoint = sum(validator[i].regular_balance) + (balanceOf(pod) rounded down to gwei) - withdrawableRestakedExecutionLayerGwei
	TotalSharesAfterCheckpointGwei *big.Float
	TotalSharesAfterCheckpointETH  *big.Float

	PodOwner       gethCommon.Address
	ProofSubmitter gethCommon.Address

	// Whether the checkpoint would need to be started with the `--force` flag.
	// This would be due to the pod not having any uncheckpointed native ETH
	MustForceCheckpoint bool
}

func getRegularBalancesGwei(allValidators []ValidatorWithIndex, state *spec.VersionedBeaconState) []phase0.Gwei {
	validatorBalances, err := state.ValidatorBalances()
	PanicOnError("failed to load validator balances", err)

	return validatorBalances
}

func sumActiveValidatorBalancesGwei(allValidators []ValidatorWithIndex, allBalances []phase0.Gwei, state *spec.VersionedBeaconState) phase0.Gwei {
	var sumGwei phase0.Gwei = 0

	for i := 0; i < len(allValidators); i++ {
		validator := allValidators[i]
		sumGwei = sumGwei + allBalances[validator.Index]
	}

	return sumGwei
}

func GetStatus(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, beaconClient BeaconClient) EigenpodStatus {
	validators := map[string]Validator{}
	var activeCheckpoint *Checkpoint = nil

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to reach eigenpod", err)

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	PanicOnError("failed to fetch current checkpoint timestamp", err)

	state, err := beaconClient.GetBeaconState(ctx, "head")
	PanicOnError("failed to fetch state", err)

	allValidators, err := FindAllValidatorsForEigenpod(eigenpodAddress, state)
	PanicOnError("failed to find validators", err)

	allBalances := getRegularBalancesGwei(allValidators, state)

	activeValidators, err := SelectActiveValidators(eth, eigenpodAddress, allValidators)
	PanicOnError("failed to find active validators", err)

	checkpointableValidators, err := SelectCheckpointableValidators(eth, eigenpodAddress, allValidators, timestamp)
	PanicOnError("failed to find checkpointable validators", err)

	sumRegularBalancesGwei := sumActiveValidatorBalancesGwei(activeValidators, allBalances, state)

	for i := 0; i < len(allValidators); i++ {
		validatorInfo, err := eigenPod.ValidatorPubkeyToInfo(nil, allValidators[i].Validator.PublicKey[:])
		PanicOnError("failed to fetch validator info", err)

		validatorIndex := allValidators[i].Index

		validators[fmt.Sprintf("%d", validatorIndex)] = Validator{
			Index:                               validatorIndex,
			Status:                              int(validatorInfo.Status),
			Slashed:                             allValidators[i].Validator.Slashed,
			PublicKey:                           allValidators[i].Validator.PublicKey.String(),
			IsAwaitingWithdrawalCredentialProof: (validatorInfo.Status == ValidatorStatusInactive) && allValidators[i].Validator.ExitEpoch == FAR_FUTURE_EPOCH,
			EffectiveBalance:                    uint64(allValidators[i].Validator.EffectiveBalance),
			CurrentBalance:                      uint64(allBalances[validatorIndex]),
		}
	}

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to fetch checkpoint information", err)

	eigenpodManagerContractAddress, err := eigenPod.EigenPodManager(nil)
	PanicOnError("failed to get manager address", err)

	eigenPodManager, err := onchain.NewEigenPodManager(eigenpodManagerContractAddress, eth)
	PanicOnError("failed to get manager instance", err)

	eigenPodOwner, err := eigenPod.PodOwner(nil)
	PanicOnError("failed to get eigenpod owner", err)

	proofSubmitter, err := eigenPod.ProofSubmitter(nil)
	PanicOnError("failed to get eigenpod proof submitter", err)

	currentOwnerShares, err := eigenPodManager.PodOwnerShares(nil, eigenPodOwner)
	PanicOnError("failed to load pod owner shares", err)
	currentOwnerSharesETH := IweiToEther(currentOwnerShares)

	withdrawableRestakedExecutionLayerGwei, err := eigenPod.WithdrawableRestakedExecutionLayerGwei(nil)
	PanicOnError("failed to fetch withdrawableRestakedExecutionLayerGwei", err)

	var pendingSharesGwei *big.Float
	mustForceCheckpoint := false
	// If we currently have an active checkpoint, estimate the total shares
	// we'll have when we complete it:
	//
	// pendingSharesGwei = withdrawableRestakedExecutionLayerGwei + checkpoint.PodBalanceGwei + sumRegularBalancesGwei
	if timestamp != 0 {
		pendingSharesGwei = new(big.Float).Add(
			new(big.Float).Add(
				new(big.Float).SetUint64(withdrawableRestakedExecutionLayerGwei),
				new(big.Float).SetUint64(checkpoint.PodBalanceGwei),
			),
			new(big.Float).SetUint64(uint64(sumRegularBalancesGwei)),
		)

		activeCheckpoint = &Checkpoint{
			ProofsRemaining: checkpoint.ProofsRemaining.Uint64(),
			StartedAt:       timestamp,
		}
	} else {
		// If we don't have an active checkpoint, estimate the shares we'd have if
		// we created one and then completed it:
		//
		// pendingSharesGwei = sumRegularBalancesGwei + latestPodBalanceGwei
		latestPodBalanceWei, err := eth.BalanceAt(ctx, common.HexToAddress(eigenpodAddress), nil)
		PanicOnError("failed to fetch pod balance", err)
		latestPodBalanceGwei := WeiToGwei(latestPodBalanceWei)

		pendingSharesGwei = new(big.Float).Add(
			new(big.Float).SetUint64(uint64(sumRegularBalancesGwei)),
			latestPodBalanceGwei,
		)

		// Determine whether the checkpoint needs to be run with `--force`
		checkpointableBalance := new(big.Float).Sub(
			latestPodBalanceGwei,
			new(big.Float).SetUint64(withdrawableRestakedExecutionLayerGwei),
		)

		if checkpointableBalance.Sign() == 0 {
			mustForceCheckpoint = true
		}
	}

	pendingEth := GweiToEther(pendingSharesGwei)

	return EigenpodStatus{
		Validators:                     validators,
		ActiveCheckpoint:               activeCheckpoint,
		CurrentTotalSharesETH:          currentOwnerSharesETH,
		TotalSharesAfterCheckpointGwei: pendingSharesGwei,
		TotalSharesAfterCheckpointETH:  pendingEth,
		NumberValidatorsToCheckpoint:   len(checkpointableValidators),
		PodOwner:                       eigenPodOwner,
		ProofSubmitter:                 proofSubmitter,
		MustForceCheckpoint:            mustForceCheckpoint,
	}
}
