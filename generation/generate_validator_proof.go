package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	commonutils "github.com/Layr-Labs/eigenpod-proofs-generation/common_utils"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog/log"
)

func GenerateValidatorFieldsProof(oracleBlockHeaderFile string, stateFile string, validatorIndexStr string, chainID uint64, output string) error {
	if validatorIndexStr == "" {
		return fmt.Errorf("validatorIndexStr: can not be empty")
	}

	validatorIndexArray := strings.Split(validatorIndexStr, ",")

	var state deneb.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := commonutils.ParseDenebStateJSONFile(stateFile)
	if err != nil {
		log.Debug().Msg("GenerateValidatorFieldsProof: error with JSON parsing")
		return err
	}
	commonutils.ParseDenebBeaconStateFromJSON(*stateJSON, &state)

	oracleBeaconBlockHeader, err = commonutils.ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
		return err
	}

	beaconStateRoot, err := state.HashTreeRoot()

	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
		return err
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
		return err
	}

	versionedState, err := beacon.CreateVersionedState(&state)
	if err != nil {
		log.Debug().AnErr("Error with CreateVersionedState", err)
		return err
	}

	var proofs = make([]commonutils.WithdrawalCredentialProofs, 0)
	for _, validatorIndex := range validatorIndexArray {
		index, err := strconv.ParseUint(validatorIndex, 10, 64)
		if err != nil {
			log.Debug().AnErr(fmt.Sprintf("Error with ParseUint(%s)", validatorIndex), err)
			return err
		}

		stateRootProof, validatorFieldsProof, err := eigenpodproofs.ProveValidatorFields(epp, &oracleBeaconBlockHeader, &versionedState, uint64(index))
		if err != nil {
			log.Debug().AnErr("Error with ProveValidatorFields", err)
			return err
		}

		proofs = append(proofs, commonutils.WithdrawalCredentialProofs{
			StateRootAgainstLatestBlockHeaderProof: commonutils.ConvertBytesToStrings(stateRootProof.StateRootProof),
			BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
			ValidatorIndex:                         uint64(index),
			WithdrawalCredentialProof:              commonutils.ConvertBytesToStrings(validatorFieldsProof),
			ValidatorFields:                        commonutils.GetValidatorFields(state.Validators[index]),
		})
	}
	var callParams = commonutils.VerifyWithdrawalCredentialsCallParams{
		StateRootProof: commonutils.StateRootProof{
			BeaconStateRoot: proofs[0].BeaconStateRoot,
			Proof:           commonutils.ConvertToHexString(proofs[0].StateRootAgainstLatestBlockHeaderProof),
		},
	}

	for _, proof := range proofs {
		callParams.ValidatorIndices = append(callParams.ValidatorIndices, proof.ValidatorIndex)
		callParams.ValidatorFieldsProofs = append(callParams.ValidatorFieldsProofs, commonutils.ConvertToHexString(proof.WithdrawalCredentialProof))
		callParams.ValidatorFields = append(callParams.ValidatorFields, proof.ValidatorFields)
	}

	proofData, err := json.Marshal(callParams)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
		return err
	}

	_ = os.WriteFile(output, proofData, 0644)

	return nil
}
