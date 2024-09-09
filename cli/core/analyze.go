package core

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/core/onchain"
	"github.com/ethereum/go-ethereum/common"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PodInfo struct {
	PodAddress gethCommon.Address
	Owner      gethCommon.Address
	Validators map[string]Validator
}

type PodAnalysis struct {
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

func AnalyzePods(ctx context.Context, pods map[string]PodInfo, eth *ethclient.Client, beaconClient BeaconClient) PodAnalysis {
	fmt.Printf("Analyzing %d pods", len(pods))

	beaconState, err := beaconClient.GetBeaconState(ctx, "head")
	PanicOnError("failed to fetch beacon state: %w", err)

	allValidators, err := beaconState.Validators()
	PanicOnError("failed to fetch state validators: %w", err)

	for _, pod := range pods {
		podValidators, err := FindAllValidatorsForEigenpod()
	}

	validators := map[string]Validator{}
	var activeCheckpoint *Checkpoint = nil

	eigenPod, err := onchain.NewEigenPod(common.HexToAddress(eigenpodAddress), eth)
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

	allBeaconBalances := getRegularBalancesGwei(state)

	activeValidators, err := SelectActiveValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod)
	PanicOnError("failed to find active validators", err)

	checkpointableValidators, err := SelectCheckpointableValidators(eth, eigenpodAddress, allValidatorsWithInfoForEigenpod, checkpointTimestamp)
	PanicOnError("failed to find checkpointable validators", err)

	sumBeaconBalancesGwei := new(big.Float).SetUint64(uint64(sumActiveValidatorBeaconBalancesGwei(activeValidators, allBeaconBalances, state)))

	sumRestakedBalancesU64, err := sumRestakedBalancesGwei(eth, eigenpodAddress, activeValidators)
	PanicOnError("failed to calculate sum of onchain validator balances", err)
	sumRestakedBalancesGwei := new(big.Float).SetUint64(uint64(sumRestakedBalancesU64))

	for _, validator := range allValidatorsWithInfoForEigenpod {

		validators[fmt.Sprintf("%d", validator.Index)] = Validator{
			Index:                               validator.Index,
			Status:                              int(validator.Info.Status),
			Slashed:                             validator.Validator.Slashed,
			PublicKey:                           validator.Validator.PublicKey.String(),
			IsAwaitingActivationQueue:           validator.Validator.ActivationEpoch == FAR_FUTURE_EPOCH,
			IsAwaitingWithdrawalCredentialProof: IsAwaitingWithdrawalCredentialProof(validator.Info, validator.Validator),
			EffectiveBalance:                    uint64(validator.Validator.EffectiveBalance),
			CurrentBalance:                      uint64(allBeaconBalances[validator.Index]),
		}
	}
}
