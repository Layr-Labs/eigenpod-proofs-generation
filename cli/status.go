package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
	- how many validators you have on the beacon chain pointed at your pod
	- how many of those validators are awaiting a withdrawal credential proof
	- how many extra shares you'd get if you completed a checkpoint right now
	- whether you have an active checkpoint
	- how many proofs remaining in your active checkpoint
	- how many shares you'll get when you finish the checkpoint
*/

type Checkpoint struct {
	PendingSharesGWei uint64
	ProofsRemaining   uint64
	StartedAt         uint64
}

type Validator struct {
	Index                               uint64
	Status                              int
	PublicKey                           string
	IsAwaitingWithdrawalCredentialProof bool
}

type EigenpodStatus struct {
	Validators map[string]Validator

	ActiveCheckpoint *Checkpoint

	// if you completed a new checkpoint right now, how many shares would you get?
	//
	//  this is computed as:
	// 		- If checkpoint is already started:
	// 			sum(beacon chain balances) + currentCheckpoint.podBalanceGwei + pod.withdrawableRestakedExecutionLayerGwei()
	// 		- If no checkpoint is started:
	// 			sum(beacon chain balances) + native ETH balance of pod
	SharesPendingCheckpointGwei float64
	SharesPendingCheckpointETH  float64
}

func sumBeaconChainBalancesGwei(allValidators []struct {
	Validator *phase0.Validator
	Index     uint64
}) phase0.Gwei {
	var sumGwei phase0.Gwei = 0

	for i := 0; i < len(allValidators); i++ {
		validator := allValidators[i]
		sumGwei = sumGwei + validator.Validator.EffectiveBalance
	}

	return sumGwei
}

func getStatus(ctx context.Context, eigenpodAddress string, eth *ethclient.Client, beaconClient BeaconClient) EigenpodStatus {
	validators := map[string]Validator{}
	var activeCheckpoint *Checkpoint = nil

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
	PanicOnError("failed to reach eigenpod", err)

	timestamp, err := eigenPod.CurrentCheckpointTimestamp(nil)
	PanicOnError("failed to fetch current checkpoint timestamp", err)

	state, err := beaconClient.GetBeaconState(ctx, "head")
	PanicOnError("failed to fetch state", err)

	allValidators := findAllValidatorsForEigenpod(eigenpodAddress, state)
	sumBalancesGwei := sumBeaconChainBalancesGwei(allValidators)

	checkpoint, err := eigenPod.CurrentCheckpoint(nil)
	PanicOnError("failed to fetch checkpoint information", err)
	if err == nil && timestamp != 0 {
		withdrawableRestakedGwei, err := eigenPod.WithdrawableRestakedExecutionLayerGwei(nil)
		PanicOnError("failed to fetch gwei info", err)

		activeCheckpoint = &Checkpoint{
			PendingSharesGWei: uint64(sumBalancesGwei) + checkpoint.PodBalanceGwei + withdrawableRestakedGwei,
			ProofsRemaining:   checkpoint.ProofsRemaining.Uint64(),
			StartedAt:         timestamp,
		}
	}

	for i := 0; i < len(allValidators); i++ {
		validatorInfo, err := eigenPod.ValidatorPubkeyToInfo(nil, allValidators[i].Validator.PublicKey[:])
		PanicOnError("failed to fetch validator info", err)

		validators[fmt.Sprintf("%d", allValidators[i].Index)] = Validator{
			Index:                               allValidators[i].Index,
			Status:                              int(validatorInfo.Status),
			PublicKey:                           allValidators[i].Validator.PublicKey.String(),
			IsAwaitingWithdrawalCredentialProof: validatorInfo.Status == ValidatorStatusInactive,
		}
	}

	latestPodBalanceWei, err := eth.BalanceAt(ctx, common.HexToAddress(eigenpodAddress), nil)
	PanicOnError("failed to fetch pod balance", err)

	pendingGwei := float64(sumBalancesGwei) + weiToGwei(latestPodBalanceWei)
	pendingEth := gweiToEther(new(big.Float).SetFloat64(pendingGwei))

	return EigenpodStatus{
		Validators:                  validators,
		ActiveCheckpoint:            activeCheckpoint,
		SharesPendingCheckpointGwei: pendingGwei,
		SharesPendingCheckpointETH:  pendingEth,
	}
}
