package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
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
	IsAwaitingActivationQueue           bool
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

func getRegularBalancesGwei(state *spec.VersionedBeaconState) []phase0.Gwei {
	validatorBalances, err := state.ValidatorBalances()
	PanicOnError("failed to load validator balances", err)

	return validatorBalances
}

func sumValidatorBeaconBalancesGwei(allValidators []ValidatorWithOnchainInfo, allBalances []phase0.Gwei) *big.Int {
	sumGwei := big.NewInt(0)

	for i := 0; i < len(allValidators); i++ {
		validator := allValidators[i]
		sumGwei = sumGwei.Add(sumGwei, new(big.Int).SetUint64(uint64(allBalances[validator.Index])))
	}

	return sumGwei
}

func sumRestakedBalancesGwei(activeValidators []ValidatorWithOnchainInfo) *big.Int {
	sumGwei := big.NewInt(0)

	for i := 0; i < len(activeValidators); i++ {
		validator := activeValidators[i]
		sumGwei = sumGwei.Add(sumGwei, new(big.Int).SetUint64(validator.Info.RestakedBalanceGwei))
	}

	return sumGwei
}

func GetStatus(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, beaconClient BeaconClient) EigenpodStatus {
	validators := map[string]Validator{}
	var activeCheckpoint *Checkpoint = nil

	eigenPod, err := onchain.NewEigenPod(gethCommon.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to reach eigenpod", err)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to fetch checkpoint information", err)

	// Fetch the beacon state associated with the checkpoint (or "head" if there is no checkpoint)
	checkpointTimestamp, state, err := GetCheckpointTimestampAndBeaconState(ctx, eigenpodAddress, eth, beaconClient)
	PanicOnError("failed to fetch checkpoint and beacon state", err)

	allValidatorsForEigenpod, err := FindAllValidatorsForEigenpod(eigenpodAddress, state)
	PanicOnError("failed to find validators", err)

	allValidatorsWithInfoForEigenpod, err := FetchMultipleOnchainValidatorInfo(ctx, eth, eigenpodAddress, allValidatorsForEigenpod)
	PanicOnError("failed to fetch validator info", err)

	allBeaconBalancesGwei := getRegularBalancesGwei(state)

	activeValidators, err := SelectActiveValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod)
	PanicOnError("failed to find active validators", err)

	checkpointableValidators, err := SelectCheckpointableValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod, checkpointTimestamp)
	PanicOnError("failed to find checkpointable validators", err)

	sumBeaconBalancesWei := IGweiToWei(sumValidatorBeaconBalancesGwei(activeValidators, allBeaconBalancesGwei))
	sumRestakedBalancesWei := IGweiToWei(sumRestakedBalancesGwei(activeValidators))

	PanicOnError("failed to calculate sum of onchain validator balances", err)

	for _, validator := range allValidatorsWithInfoForEigenpod {

		validators[fmt.Sprintf("%d", validator.Index)] = Validator{
			Index:                               validator.Index,
			Status:                              int(validator.Info.Status),
			Slashed:                             validator.Validator.Slashed,
			PublicKey:                           validator.Validator.PublicKey.String(),
			IsAwaitingActivationQueue:           validator.Validator.ActivationEpoch == FAR_FUTURE_EPOCH,
			IsAwaitingWithdrawalCredentialProof: IsAwaitingWithdrawalCredentialProof(validator.Info, validator.Validator),
			EffectiveBalance:                    uint64(validator.Validator.EffectiveBalance),
			CurrentBalance:                      uint64(allBeaconBalancesGwei[validator.Index]),
		}
	}

	eigenpodManagerContractAddress, err := eigenPod.EigenPodManager(nil)
	PanicOnError("failed to get manager address", err)

	eigenPodManager, err := onchain.NewEigenPodManager(eigenpodManagerContractAddress, eth)
	PanicOnError("failed to get manager instance", err)

	eigenPodOwner, err := eigenPod.PodOwner(nil)
	PanicOnError("failed to get eigenpod owner", err)

	proofSubmitter, err := eigenPod.ProofSubmitter(nil)
	PanicOnError("failed to get eigenpod proof submitter", err)

	currentOwnerShares, err := eigenPodManager.PodOwnerShares(nil, eigenPodOwner)
	// currentOwnerShares = big.NewInt(0)
	PanicOnError("failed to load pod owner shares", err)
	currentOwnerSharesETH := IweiToEther(currentOwnerShares)
	currentOwnerSharesWei := currentOwnerShares

	withdrawableRestakedExecutionLayerGwei, err := eigenPod.WithdrawableRestakedExecutionLayerGwei(nil)
	PanicOnError("failed to fetch withdrawableRestakedExecutionLayerGwei", err)

	// Estimate the total shares we'll have if we complete an existing checkpoint
	// (or start a new one and complete that).
	//
	// First, we need the change in the pod's native ETH balance since the last checkpoint:
	var nativeETHDeltaWei *big.Int
	mustForceCheckpoint := false

	if checkpointTimestamp != 0 {
		// Change in the pod's native ETH balance (already calculated for us when the checkpoint was started)
		fmt.Printf("pod had a checkpoint\n")
		nativeETHDeltaWei = IGweiToWei(new(big.Int).SetUint64(checkpoint.PodBalanceGwei))

		// Remove already-computed delta from an in-progress checkpoint
		sumRestakedBalancesWei = new(big.Int).Sub(
			sumRestakedBalancesWei,
			IGweiToWei(checkpoint.BalanceDeltasGwei),
		)

		activeCheckpoint = &Checkpoint{
			ProofsRemaining: checkpoint.ProofsRemaining.Uint64(),
			StartedAt:       checkpointTimestamp,
		}
	} else {
		fmt.Printf("pod did not have a checkpoint\n")
		latestPodBalanceWei, err := eth.BalanceAt(ctx, gethCommon.HexToAddress(eigenpodAddress), nil)
		PanicOnError("failed to fetch pod balance", err)

		// We don't have a checkpoint currently, so we need to calculate what
		// checkpoint.PodBalanceGwei would be if we started one now:
		nativeETHDeltaWei = new(big.Int).Sub(
			latestPodBalanceWei,
			IGweiToWei(new(big.Int).SetUint64(withdrawableRestakedExecutionLayerGwei)),
		)

		// Determine whether the checkpoint needs to be started with `--force`
		if nativeETHDeltaWei.Sign() == 0 {
			mustForceCheckpoint = true
		}
	}

	// Next, we need the change in the pod's beacon chain balances since the last
	// checkpoint:
	//
	// beaconETHDeltaWei = sumBeaconBalancesWei - sumRestakedBalancesWei
	beaconETHDeltaWei := new(big.Int).Sub(
		sumBeaconBalancesWei,
		sumRestakedBalancesWei,
	)

	// Sum of these two deltas represents the change in shares after this checkpoint
	totalShareDeltaWei := new(big.Int).Add(
		nativeETHDeltaWei,
		beaconETHDeltaWei,
	)

	// Calculate new total shares by applying delta to current shares
	pendingSharesWei := new(big.Int).Add(
		currentOwnerSharesWei,
		totalShareDeltaWei,
	)

	pendingEth := GweiToEther(WeiToGwei(pendingSharesWei))

	return EigenpodStatus{
		Validators:                     validators,
		ActiveCheckpoint:               activeCheckpoint,
		CurrentTotalSharesETH:          currentOwnerSharesETH,
		TotalSharesAfterCheckpointGwei: WeiToGwei(pendingSharesWei),
		TotalSharesAfterCheckpointETH:  pendingEth,
		NumberValidatorsToCheckpoint:   len(checkpointableValidators),
		PodOwner:                       eigenPodOwner,
		ProofSubmitter:                 proofSubmitter,
		MustForceCheckpoint:            mustForceCheckpoint,
	}
}
