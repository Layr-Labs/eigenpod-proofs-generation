package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/DelegationManager"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPod"
	"github.com/Layr-Labs/eigenlayer-contracts/pkg/bindings/EigenPodManager"
	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/utils"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func BeaconStrategy() gethCommon.Address {
	return gethCommon.HexToAddress("0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0")
}

type EigenpodStatus struct {
	Validators map[string]utils.Validator

	ActiveCheckpoint *utils.Checkpoint

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
	utils.PanicOnError("failed to load validator balances", err)

	return validatorBalances
}

func sumValidatorBeaconBalancesGwei(allValidators []utils.ValidatorWithOnchainInfo, allBalances []phase0.Gwei) *big.Int {
	sumGwei := big.NewInt(0)

	for i := 0; i < len(allValidators); i++ {
		validator := allValidators[i]
		sumGwei = sumGwei.Add(sumGwei, new(big.Int).SetUint64(uint64(allBalances[validator.Index])))
	}

	return sumGwei
}

func sumRestakedBalancesGwei(activeValidators []utils.ValidatorWithOnchainInfo) *big.Int {
	sumGwei := big.NewInt(0)

	for i := 0; i < len(activeValidators); i++ {
		validator := activeValidators[i]
		sumGwei = sumGwei.Add(sumGwei, new(big.Int).SetUint64(validator.Info.RestakedBalanceGwei))
	}

	return sumGwei
}

func GetStatus(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, beaconClient utils.BeaconClient) EigenpodStatus {
	validators := map[string]utils.Validator{}
	var activeCheckpoint *utils.Checkpoint = nil

	eigenPod, err := EigenPod.NewEigenPod(gethCommon.HexToAddress(eigenpodAddress), eth)
	utils.PanicOnError("failed to reach eigenpod", err)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	utils.PanicOnError("failed to fetch checkpoint information", err)

	// Fetch the beacon state associated with the checkpoint (or "head" if there is no checkpoint)
	checkpointTimestamp, state, err := utils.GetCheckpointTimestampAndBeaconState(ctx, eigenpodAddress, eth, beaconClient)
	utils.PanicOnError("failed to fetch checkpoint and beacon state", err)

	allValidatorsForEigenpod, err := utils.FindAllValidatorsForEigenpod(eigenpodAddress, state)
	utils.PanicOnError("failed to find validators", err)

	allValidatorsWithInfoForEigenpod, err := utils.FetchMultipleOnchainValidatorInfo(ctx, eth, eigenpodAddress, allValidatorsForEigenpod)
	utils.PanicOnError("failed to fetch validator info", err)

	allBeaconBalancesGwei := getRegularBalancesGwei(state)

	activeValidators, err := utils.SelectActiveValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod)
	utils.PanicOnError("failed to find active validators", err)

	checkpointableValidators, err := utils.SelectCheckpointableValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod, checkpointTimestamp)
	utils.PanicOnError("failed to find checkpointable validators", err)

	sumBeaconBalancesWei := utils.IGweiToWei(sumValidatorBeaconBalancesGwei(activeValidators, allBeaconBalancesGwei))
	sumRestakedBalancesWei := utils.IGweiToWei(sumRestakedBalancesGwei(activeValidators))

	utils.PanicOnError("failed to calculate sum of onchain validator balances", err)

	for _, validator := range allValidatorsWithInfoForEigenpod {

		validators[fmt.Sprintf("%d", validator.Index)] = utils.Validator{
			Index:                               validator.Index,
			Status:                              int(validator.Info.Status),
			Slashed:                             validator.Validator.Slashed,
			PublicKey:                           validator.Validator.PublicKey.String(),
			IsAwaitingActivationQueue:           validator.Validator.ActivationEpoch == utils.FAR_FUTURE_EPOCH,
			IsAwaitingWithdrawalCredentialProof: utils.IsAwaitingWithdrawalCredentialProof(validator.Info, validator.Validator),
			EffectiveBalance:                    uint64(validator.Validator.EffectiveBalance),
			CurrentBalance:                      uint64(allBeaconBalancesGwei[validator.Index]),
			WithdrawalPrefix:                    validator.Validator.WithdrawalCredentials[0],
		}
	}

	eigenpodManagerContractAddress, err := eigenPod.EigenPodManager(nil)
	utils.PanicOnError("failed to get manager address", err)

	eigenPodManager, err := EigenPodManager.NewEigenPodManager(eigenpodManagerContractAddress, eth)
	utils.PanicOnError("failed to get manager instance", err)

	eigenPodOwner, err := eigenPod.PodOwner(nil)
	utils.PanicOnError("failed to get eigenpod owner", err)

	proofSubmitter, err := eigenPod.ProofSubmitter(nil)
	utils.PanicOnError("failed to get eigenpod proof submitter", err)

	delegationManagerAddress, err := eigenPodManager.DelegationManager(nil)
	utils.PanicOnError("failed to read delegationManager", err)

	delegationManager, err := DelegationManager.NewDelegationManager(delegationManagerAddress, eth)
	utils.PanicOnError("failed to reach delegationManager", err)

	shares, err := delegationManager.GetWithdrawableShares(nil, eigenPodOwner, []gethCommon.Address{
		BeaconStrategy(),
	})
	utils.PanicOnError("failed to load owner shares", err)

	currentOwnerSharesETH := utils.IweiToEther(shares.WithdrawableShares[0])
	currentOwnerSharesWei := shares.WithdrawableShares[0]

	withdrawableRestakedExecutionLayerGwei, err := eigenPod.WithdrawableRestakedExecutionLayerGwei(nil)
	utils.PanicOnError("failed to fetch withdrawableRestakedExecutionLayerGwei", err)

	// Estimate the total shares we'll have if we complete an existing checkpoint
	// (or start a new one and complete that).
	//
	// First, we need the change in the pod's native ETH balance since the last checkpoint:
	var nativeETHDeltaWei *big.Int
	mustForceCheckpoint := false

	if checkpointTimestamp != 0 {
		// Change in the pod's native ETH balance (already calculated for us when the checkpoint was started)
		nativeETHDeltaWei = utils.IGweiToWei(new(big.Int).SetUint64(checkpoint.PodBalanceGwei))

		// Remove already-computed delta from an in-progress checkpoint
		sumRestakedBalancesWei = new(big.Int).Sub(
			sumRestakedBalancesWei,
			utils.IGweiToWei(big.NewInt(checkpoint.BalanceDeltasGwei)),
		)

		activeCheckpoint = &utils.Checkpoint{
			ProofsRemaining: checkpoint.ProofsRemaining.Uint64(),
			StartedAt:       checkpointTimestamp,
		}
	} else {
		latestPodBalanceWei, err := eth.BalanceAt(ctx, gethCommon.HexToAddress(eigenpodAddress), nil)
		utils.PanicOnError("failed to fetch pod balance", err)

		// We don't have a checkpoint currently, so we need to calculate what
		// checkpoint.PodBalanceGwei would be if we started one now:
		nativeETHDeltaWei = new(big.Int).Sub(
			latestPodBalanceWei,
			utils.IGweiToWei(new(big.Int).SetUint64(withdrawableRestakedExecutionLayerGwei)),
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

	pendingEth := utils.GweiToEther(utils.WeiToGwei(pendingSharesWei))

	return EigenpodStatus{
		Validators:                     validators,
		ActiveCheckpoint:               activeCheckpoint,
		CurrentTotalSharesETH:          currentOwnerSharesETH,
		TotalSharesAfterCheckpointGwei: utils.WeiToGwei(pendingSharesWei),
		TotalSharesAfterCheckpointETH:  pendingEth,
		NumberValidatorsToCheckpoint:   len(checkpointableValidators),
		PodOwner:                       eigenPodOwner,
		ProofSubmitter:                 proofSubmitter,
		MustForceCheckpoint:            mustForceCheckpoint,
	}
}
