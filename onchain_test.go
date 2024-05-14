package eigenpodproofs

import (
	"fmt"
	"math/big"
	"testing"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func TestValidatorContainersProofOnChain(t *testing.T) {
	var oracleState deneb.BeaconState
	stateFile := "data/deneb_goerli_slot_7413760.json"
	stateJSON, err := ParseJSONFileDeneb(stateFile)
	if err != nil {
		fmt.Println("error with JSON parsing beacon state")
	}
	ParseDenebBeaconStateFromJSON(*stateJSON, &oracleState)

	versionedOracleState, err := beacon.CreateVersionedState(&oracleState)
	if err != nil {
		fmt.Println("error", err)
	}

	oracleBlockHeader, err = ExtractBlockHeader("data/deneb_goerli_block_header_7413760.json")
	if err != nil {
		fmt.Println("error", err)
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(&oracleBlockHeader, &versionedOracleState, []uint64{VALIDATOR_INDEX})
	if err != nil {
		fmt.Println("error", err)
	}

	validatorFieldsProof := verifyValidatorFieldsCallParams.ValidatorFieldsProofs[0].ToByteSlice()
	validatorIndex := new(big.Int).SetUint64(verifyValidatorFieldsCallParams.ValidatorIndices[0])
	oracleBlockHeaderRoot, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		fmt.Println("error", err)
	}

	err = beaconChainProofs.VerifyStateRootAgainstLatestBlockRoot(
		&bind.CallOpts{},
		oracleBlockHeaderRoot,
		verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
		verifyValidatorFieldsCallParams.StateRootProof.StateRootProof.ToByteSlice(),
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)

	var validatorFields [][32]byte
	for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[0] {
		validatorFields = append(validatorFields, field)
	}

	err = beaconChainProofs.VerifyValidatorFields(
		&bind.CallOpts{},
		verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
		validatorFields,
		validatorFieldsProof,
		validatorIndex,
	)
	if err != nil {
		fmt.Println("error", err)
	}
	assert.Nil(t, err)
}
