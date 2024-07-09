package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Checkpoint struct {
	PendingSharesGwei *big.Float
	ProofsRemaining   uint64
	StartedAt         uint64
}

type Validator struct {
	Slashed                             bool
	Index                               uint64
	Status                              int
	PublicKey                           string
	IsAwaitingWithdrawalCredentialProof bool
}

type EigenpodStatus struct {
	Validators map[string]Validator

	ActiveCheckpoint *Checkpoint

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
}

func sumBeaconChainRegularBalancesGwei(allValidators []ValidatorWithIndex, state *spec.VersionedBeaconState) phase0.Gwei {
	var sumGwei phase0.Gwei = 0

	validatorBalances, err := state.ValidatorBalances()
	PanicOnError("failed to load validator balances", err)

	for i := 0; i < len(allValidators); i++ {
		validator := allValidators[i]
		sumGwei = sumGwei + validatorBalances[validator.Index]
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

	allValidators := FindAllValidatorsForEigenpod(eigenpodAddress, state)
	sumRegularBalancesGwei := sumBeaconChainRegularBalancesGwei(allValidators, state)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to fetch checkpoint information", err)

	eigenpodManagerContractAddress, err := eigenPod.EigenPodManager(nil)
	PanicOnError("failed to get manager address", err)

	eigenPodManager, err := onchain.NewEigenPodManager(eigenpodManagerContractAddress, eth)
	PanicOnError("failed to get manager instance", err)

	eigenPodOwner, err := eigenPod.PodOwner(nil)
	PanicOnError("failed to get eigenpod owner", err)

	currentOwnerShares, err := eigenPodManager.PodOwnerShares(nil, eigenPodOwner)
	PanicOnError("failed to load pod owner shares", err)
	currentOwnerSharesETH := IweiToEther(currentOwnerShares)

	for i := 0; i < len(allValidators); i++ {
		validatorInfo, err := eigenPod.ValidatorPubkeyToInfo(nil, allValidators[i].Validator.PublicKey[:])
		PanicOnError("failed to fetch validator info", err)
		validators[fmt.Sprintf("%d", allValidators[i].Index)] = Validator{
			Index:                               allValidators[i].Index,
			Status:                              int(validatorInfo.Status),
			Slashed:                             allValidators[i].Validator.Slashed,
			PublicKey:                           allValidators[i].Validator.PublicKey.String(),
			IsAwaitingWithdrawalCredentialProof: (validatorInfo.Status == ValidatorStatusInactive) && allValidators[i].Validator.ExitEpoch == FAR_FUTURE_EPOCH,
		}
	}

	withdrawableRestakedExecutionLayerGwei, err := eigenPod.WithdrawableRestakedExecutionLayerGwei(nil)
	PanicOnError("failed to fetch withdrawableRestakedExecutionLayerGwei", err)

	pendingSharesGwei := new(big.Float).Add(
		new(big.Float).Add(
			new(big.Float).SetUint64(withdrawableRestakedExecutionLayerGwei),
			new(big.Float).SetUint64(checkpoint.PodBalanceGwei),
		),
		new(big.Float).SetUint64(uint64(sumRegularBalancesGwei)),
	) // pendingSharesGwei = withdrawableRestakedExecutionLayerGwei + checkpoint.PodBalanceGwei + sumRegularBalancesGwei

	if err == nil && timestamp != 0 {
		activeCheckpoint = &Checkpoint{
			PendingSharesGwei: pendingSharesGwei,
			ProofsRemaining:   checkpoint.ProofsRemaining.Uint64(),
			StartedAt:         timestamp,
		}
	}

	latestPodBalanceWei, err := eth.BalanceAt(ctx, common.HexToAddress(eigenpodAddress), nil)
	PanicOnError("failed to fetch pod balance", err)
	latestPodBalanceGwei := WeiToGwei(latestPodBalanceWei)
	pendingGwei :=
		new(big.Float).Sub(
			new(big.Float).Add(
				new(big.Float).SetUint64(uint64(sumRegularBalancesGwei)),
				latestPodBalanceGwei),
			new(big.Float).SetUint64(withdrawableRestakedExecutionLayerGwei),
		)
	pendingEth := GweiToEther(pendingGwei)

	return EigenpodStatus{
		Validators:                     validators,
		ActiveCheckpoint:               activeCheckpoint,
		CurrentTotalSharesETH:          currentOwnerSharesETH,
		TotalSharesAfterCheckpointGwei: pendingGwei,
		TotalSharesAfterCheckpointETH:  pendingEth,
	}
}
