package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	commonutils "github.com/Layr-Labs/eigenpod-proofs-generation/common_utils"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog/log"
)

func GenerateBalanceUpdateProof(oracleBlockHeaderFile string, stateFile string, validatorIndexStr string, chainID uint64, output string) error {
	validatorIndex, err := strconv.ParseUint(validatorIndexStr, 10, 64)
	if err != nil {
		log.Debug().AnErr(fmt.Sprintf("Error with ParseUint(%s)", validatorIndexStr), err)
		return err
	}

	var state deneb.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := commonutils.ParseDenebStateJSONFile(stateFile)
	if err != nil {
		log.Debug().AnErr("GenerateBalanceUpdateProof: error with JSON parsing", err)
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

	stateRootProof, validatorFieldsProof, err := eigenpodproofs.ProveValidatorFields(epp, &oracleBeaconBlockHeader, &versionedState, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
		return err
	}
	proofs := commonutils.BalanceUpdateProofs{
		ValidatorIndex:                         uint64(validatorIndex),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		StateRootAgainstLatestBlockHeaderProof: commonutils.ConvertBytesToStrings(stateRootProof.StateRootProof),
		ValidatorFieldsProof:                   commonutils.ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        commonutils.GetValidatorFields(state.Validators[validatorIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
		return err
	}

	_ = os.WriteFile(output, proofData, 0644)

	return nil

}
