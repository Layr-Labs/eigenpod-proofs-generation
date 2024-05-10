package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	commonutils "github.com/Layr-Labs/eigenpod-proofs-generation/common_utils"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const tmpPath = "/tmp"

func GenerateValidatorFieldsProofs(oracleBlockHeaderFile string, stateFile string, validatorIndices []uint64, chainID uint64, output string) error {

	var oracleBeaconBlockHeader phase0.BeaconBlockHeader

	oracleBeaconBlockHeader, err := commonutils.ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
		return err
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
		return err
	}

	versionedStatePtr, err := versionedStateUsingCache(stateFile, chainID, uint64(oracleBeaconBlockHeader.Slot))
	if err != nil {
		log.Debug().AnErr("Error with versionedStateUsingCache", err)
		return err
	}

	var versionedState = *versionedStatePtr

	// When a new CL version is released, this will break.
	// original line (uses state directly):
	//	beaconStateRoot, err := state.HashTreeRoot()
	beaconStateRoot, err := versionedState.Deneb.HashTreeRoot()
	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
		return err
	}

	validators, err := versionedState.Validators()
	if err != nil {
		log.Debug().AnErr("Error getting versioned Validators", err)
		return err
	}

	allProofs := make([]commonutils.WithdrawalCredentialProofs, len(validatorIndices))
	results := make(chan commonutils.WithdrawalCredentialProofs, len(validatorIndices))
	var wg sync.WaitGroup

	for _, validatorIndex := range validatorIndices {
		wg.Add(1)
		go func(index uint64) {
			defer wg.Done()
			stateRootProof, validatorFieldsProof, err := eigenpodproofs.ProveValidatorFields(epp, &oracleBeaconBlockHeader, &versionedState, index)
			if err != nil {
				log.Printf("Error with ProveValidatorFields for index %d: %v", index, err)
				return
			}

			proof := commonutils.WithdrawalCredentialProofs{
				StateRootAgainstLatestBlockHeaderProof: commonutils.ConvertBytesToStrings(stateRootProof.StateRootProof),
				BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
				ValidatorIndex:                         index,
				WithdrawalCredentialProof:              commonutils.ConvertBytesToStrings(validatorFieldsProof),
				ValidatorFields:                        commonutils.GetValidatorFields(validators[index]),
			}
			results <- proof
		}(validatorIndex)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	i := 0
	for proof := range results {
		allProofs[i] = proof
		i++
	}

	if len(allProofs) != len(validatorIndices) {
		numMissingProofs := len(validatorIndices) - len(allProofs)
		err = errors.New(fmt.Sprintf("failed to generate %d proofs", numMissingProofs))
		log.Debug().AnErr("missing proofs from results: ", err)
		return err
	}

	proofData, err := json.Marshal(allProofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
		return err
	}

	_ = os.WriteFile(output, proofData, 0644)

	return nil
}

func stateCacheFilePath(chainID uint64, slot uint64) string {
	return fmt.Sprintf("%s/versioned_state_chain_%d_slot_%d", tmpPath, chainID, slot)
}

/**
 * Reading the state file and parsing the JSON can take several seconds. This uses a disk cache to store the parsed versioned state.
 */
func versionedStateUsingCache(stateFile string, chainID uint64, slot uint64) (*spec.VersionedBeaconState, error) {
	var state deneb.BeaconState
	var versionedState spec.VersionedBeaconState

	versionedStateCache := NewCache(stateCacheFilePath(chainID, slot))
	cachedVersionedState, err := versionedStateCache.Get()

	if err != nil {
		log.Debug().AnErr("Error fetching versionedState cache", err)
	} else if cachedVersionedState != nil {
		unmarshaledVersionedState, err := beacon.UnmarshalSSZVersionedBeaconState(cachedVersionedState)
		if err != nil {
			log.Debug().Msg("failed to unmarshal versioned beacon state")
			return nil, err
		}
		versionedState = *unmarshaledVersionedState
	}

	if versionedState.IsEmpty() {
		stateJSON, err := commonutils.ParseDenebStateJSONFile(stateFile)
		if err != nil {
			log.Debug().Msg("error with JSON parsing")
			return nil, err
		}

		commonutils.ParseDenebBeaconStateFromJSON(*stateJSON, &state)

		versionedState, err = beacon.CreateVersionedState(&state)
		if err != nil {
			log.Debug().AnErr("Error with CreateVersionedState", err)
			return nil, err
		}

		ssz, err := beacon.MarshalSSZVersionedBeaconState(versionedState)
		if err != nil {
			log.Debug().AnErr("Error with caching versionedState", err)
			return nil, err
		}
		versionedStateCache.Set(ssz)
	}

	return &versionedState, nil
}

func ClearStateCache(oracleBlockHeaderFile string, chainID uint64) error {
	oracleBeaconBlockHeader, err := commonutils.ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
		return err
	}

	err = DeleteCache(stateCacheFilePath(chainID, uint64(oracleBeaconBlockHeader.Slot)))
	if err != nil {
		log.Debug().AnErr("Failed to delete cache file", err)
		return err
	}
	return nil
}
