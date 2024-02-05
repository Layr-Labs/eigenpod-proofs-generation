package main

import (
	"encoding/hex"
	"encoding/json"
	"os"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog/log"
)

func GenerateBalanceUpdateProof(oracleBlockHeaderFile string, stateFile string, validatorIndex uint64, chainID uint64, output string) {

	var state deneb.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := ParseDenebStateJSONFile(stateFile)
	if err != nil {
		log.Debug().AnErr("GenerateBalanceUpdateProof: error with JSON parsing", err)
	}
	ParseDenebBeaconStateFromJSON(*stateJSON, &state)

	oracleBeaconBlockHeader, err = ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)
	}

	beaconStateRoot, err := state.HashTreeRoot()

	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
	}

	epp, err := eigenpodproofs.NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)
	}

	versionedState := CreateVersionedState(state)
	stateRootProof, validatorFieldsProof, err := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &versionedState, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
	}
	proofs := BalanceUpdateProofs{
		ValidatorIndex:                         uint64(validatorIndex),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
		ValidatorFieldsProof:                   ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
	}

	_ = os.WriteFile(output, proofData, 0644)

}
