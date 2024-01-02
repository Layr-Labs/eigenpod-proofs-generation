package eigenpodproofs

import (
	"encoding/hex"
	"encoding/json"
	"os"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/rs/zerolog/log"
)

func GenerateValidatorFieldsProof(oracleBlockHeaderFile string, stateFile string, validatorIndex uint64, chainID uint64, output string) {

	var state capella.BeaconState
	var oracleBeaconBlockHeader phase0.BeaconBlockHeader
	stateJSON, err := parseStateJSONFile(stateFile)
	if err != nil {
		log.Debug().Msg("GenerateValidatorFieldsProof: error with JSON parsing")
	}
	ParseCapellaBeaconStateFromJSON(*stateJSON, &state)

	oracleBeaconBlockHeader, err = ExtractBlockHeader(oracleBlockHeaderFile)
	if err != nil {
		log.Debug().AnErr("Error with parsing header file", err)

	}

	beaconStateRoot, err := state.HashTreeRoot()

	if err != nil {
		log.Debug().AnErr("Error with HashTreeRoot of state", err)
	}

	epp, err := NewEigenPodProofs(chainID, 1000)
	if err != nil {
		log.Debug().AnErr("Error creating EPP object", err)

	}

	stateRootProof, validatorFieldsProof, err := epp.ProveValidatorFields(&oracleBeaconBlockHeader, &state, uint64(validatorIndex))
	if err != nil {
		log.Debug().AnErr("Error with ProveValidatorFields", err)
	}

	proofs := WithdrawalCredentialProofs{
		StateRootAgainstLatestBlockHeaderProof: ConvertBytesToStrings(stateRootProof.StateRootProof),
		BeaconStateRoot:                        "0x" + hex.EncodeToString(beaconStateRoot[:]),
		ValidatorIndex:                         uint64(validatorIndex),
		WithdrawalCredentialProof:              ConvertBytesToStrings(validatorFieldsProof),
		ValidatorFields:                        GetValidatorFields(state.Validators[validatorIndex]),
	}

	proofData, err := json.Marshal(proofs)
	if err != nil {
		log.Debug().AnErr("JSON marshal error: ", err)
	}

	_ = os.WriteFile(output, proofData, 0644)

}
