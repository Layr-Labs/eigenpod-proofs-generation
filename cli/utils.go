package main

import (
	"fmt"
	"math"
	"os"

	"github.com/Layr-Labs/eigenpod-proofs-generation/cli/onchain"
	"github.com/fatih/color"
)

const (
	ValidatorStatusInactive  = 0
	ValidatorStatusActive    = 1
	ValidatorStatusWithdrawn = 2
)

func PanicOnError(message string, err error) {
	if err != nil {
		color.Red(fmt.Sprintf("error: %s\n\n", message))

		info := color.New(color.FgBlack, color.Italic)
		info.Printf(fmt.Sprintf("caused by: %s\n", err))

		os.Exit(1)
	}
}

func AllZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func WriteOutputToFileOrStdout(output []byte, out *string) {
	if out != nil && *out != "" {
		err := os.WriteFile(*out, output, os.ModePerm)
		PanicOnError("failed to write to disk", err)
		color.Green("Wrote output to %s\n", *out)
	} else {
		fmt.Println(string(output))
	}
}

func FilterNotCheckpointedOrWithdrawnValidators(
	allValidatorsForEigenpod []ValidatorWithIndex,
	onchainInfo []onchain.IEigenPodValidatorInfo,
	lastCheckpoint uint64,
) []uint64 {
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := onchainInfo[i]

		notCheckpointed := validatorInfo.LastCheckpointedAt != lastCheckpoint
		notWithdrawn := validatorInfo.Status != ValidatorStatusWithdrawn

		if notCheckpointed && notWithdrawn {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.Index)
		}
	}
	return checkpointValidatorIndices
}

// (https://github.com/Layr-Labs/eigenlayer-contracts/blob/d148952a2942a97a218a2ab70f9b9f1792796081/src/contracts/libraries/BeaconChainProofs.sol#L64)
const FAR_FUTURE_EPOCH = math.MaxUint64

func FilterInactiveValidators(
	allValidatorsForEigenpod []ValidatorWithIndex,
	onchainInfo []onchain.IEigenPodValidatorInfo,
) []uint64 {
	var checkpointValidatorIndices = []uint64{}
	for i := 0; i < len(allValidatorsForEigenpod); i++ {
		validator := allValidatorsForEigenpod[i]
		validatorInfo := onchainInfo[i]

		if (validatorInfo.Status == ValidatorStatusInactive) &&
			(validator.Validator.ExitEpoch == FAR_FUTURE_EPOCH) {
			checkpointValidatorIndices = append(checkpointValidatorIndices, validator.Index)
		}
	}
	return checkpointValidatorIndices
}
