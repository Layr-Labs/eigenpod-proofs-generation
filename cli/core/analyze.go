package core

import (
	"context"
	"fmt"
	"math"
	"math/big"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type PodInfo struct {
	PodAddress gethCommon.Address
	Owner      gethCommon.Address
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
	fmt.Printf("Analyzing %d pods\n", len(pods))

	beaconState, err := beaconClient.GetBeaconState(ctx, "head")
	PanicOnError("failed to fetch beacon state: %w", err)

	var inactiveFound bool
	var problemsFound bool
	podsAnalyzed := 0
	validatorsAnalyzed := 0
	numWithExitEpochs := 0

	for addr, _ := range pods {
		podValidators, err := FindAllValidatorsForEigenpod(addr, beaconState)
		PanicOnError("failed to fetch validators for pod: %w", err)

		podsAnalyzed++
		validatorsAnalyzed += len(podValidators)

		if podsAnalyzed%100 == 0 {
			fmt.Printf("Analyzed %d/%d pods (%d total validators | %d total exited)...\n", podsAnalyzed, len(pods), validatorsAnalyzed, numWithExitEpochs)
		}

		var inactiveValidators []ValidatorWithIndex
		for _, validator := range podValidators {
			if validator.Validator.ActivationEpoch == math.MaxUint64 {
				inactiveFound = true
				inactiveValidators = append(inactiveValidators, validator)
			}

			if validator.Validator.ExitEpoch != math.MaxUint64 {
				numWithExitEpochs++
			}
		}

		if len(inactiveValidators) == 0 {
			continue
		}

		fmt.Printf("Found %d inactive validators in pod %s\n", len(inactiveValidators), addr)

		inactiveValidatorsWithInfo, err := FetchMultipleOnchainValidatorInfo(context.Background(), eth, addr, inactiveValidators)
		PanicOnError("failed to fetch onchain info for pod: %w", err)

		var problemValidators []ValidatorWithOnchainInfo
		for _, validator := range inactiveValidatorsWithInfo {
			if validator.Info.Status == ValidatorStatusActive {
				problemValidators = append(problemValidators, validator)
			}
		}

		if len(problemValidators) == 0 {
			continue
		}

		fmt.Printf("Found %d problematic validators in pod %s\n", len(problemValidators), addr)
	}

	if !inactiveFound {
		fmt.Printf("Didn't find any inactive validators!\n")
	}

	if !problemsFound {
		fmt.Printf("Didn't find any problematic validators!\n")
	}

	return PodAnalysis{}
}
