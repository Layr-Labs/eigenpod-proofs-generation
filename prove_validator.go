package eigenpodproofs

import (
	"math/big"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

type VerifyValidatorFieldsCallParams struct {
	OracleTimestamp       uint64          `json:"oracleTimestamp"`
	StateRootProof        *StateRootProof `json:"stateRootProof"`
	ValidatorIndices      []uint64        `json:"validatorIndices"`
	ValidatorFieldsProofs []common.Proof  `json:"validatorFieldsProofs"`
	ValidatorFields       [][]Bytes32     `json:"validatorFields"`
}

// We define these two aliases here, since calls to verifyWithdrawalCredentials and verifyBalanceUpdates use the same inputs
type VerifyWithdrawalCredentialsCallParams = VerifyValidatorFieldsCallParams
type VerifyBalanceUpdatesCallParams = VerifyValidatorFieldsCallParams

func (epp *EigenPodProofs) ProveValidatorContainers(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *spec.VersionedBeaconState, validatorIndices []uint64) (*VerifyValidatorFieldsCallParams, error) {

	oracleBeaconStateSlot, err := oracleBeaconState.Slot()
	if err != nil {
		return nil, err
	}
	oracleBeaconStateValidators, err := oracleBeaconState.Validators()
	if err != nil {
		return nil, err
	}

	VerifyValidatorFieldsCallParams := &VerifyValidatorFieldsCallParams{}
	VerifyValidatorFieldsCallParams.StateRootProof = &StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	// Get beacon state root.
	VerifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot = oracleBlockHeader.StateRoot
	if err != nil {
		return nil, err
	}

	VerifyValidatorFieldsCallParams.StateRootProof.StateRootProof, err = beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	VerifyValidatorFieldsCallParams.OracleTimestamp, err = GetSlotTimestamp(oracleBeaconState, oracleBlockHeader)
	if err != nil {
		return nil, err
	}
	VerifyValidatorFieldsCallParams.ValidatorIndices = make([]uint64, len(validatorIndices))
	VerifyValidatorFieldsCallParams.ValidatorFieldsProofs = make([]common.Proof, len(validatorIndices))
	VerifyValidatorFieldsCallParams.ValidatorFields = make([][]Bytes32, len(validatorIndices))
	for i, validatorIndex := range validatorIndices {
		VerifyValidatorFieldsCallParams.ValidatorIndices[i] = validatorIndex
		// prove the validator fields against the beacon state
		VerifyValidatorFieldsCallParams.ValidatorFieldsProofs[i], err = epp.ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleBeaconStateSlot, oracleBeaconStateValidators, validatorIndex)
		if err != nil {
			return nil, err
		}

		VerifyValidatorFieldsCallParams.ValidatorFields[i] = ConvertValidatorToValidatorFields(oracleBeaconStateValidators[validatorIndex])
	}

	return VerifyValidatorFieldsCallParams, nil
}

func (epp *EigenPodProofs) ProveValidatorAgainstBeaconState(beaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots, oracleBeaconStateSlot phase0.Slot, oracleBeaconStateValidators []*phase0.Validator, validatorIndex uint64) (common.Proof, error) {
	// prove the validator list against the beacon state
	validatorListProof, err := beacon.ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, beacon.ValidatorListIndex)
	if err != nil {
		return nil, err
	}

	// prove the validator root against the validator list root
	validatorProof, err := epp.ProveValidatorAgainstValidatorList(oracleBeaconStateSlot, oracleBeaconStateValidators, validatorIndex)
	if err != nil {
		return nil, err
	}

	proof := append(validatorProof, validatorListProof...)

	return proof, nil
}

func (epp *EigenPodProofs) ProveValidatorAgainstValidatorList(slot phase0.Slot, validators []*phase0.Validator, validatorIndex uint64) (common.Proof, error) {
	validatorTree, err := epp.ComputeValidatorTree(slot, validators)
	if err != nil {
		return nil, err
	}

	proof, err := common.ComputeMerkleProofFromTree(validatorTree, validatorIndex, beacon.ValidatorListMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}
	//append the length of the validator array to the proof
	//convert big endian to little endian
	validatorListLenLE := BigToLittleEndian(big.NewInt(int64(len(validators))))

	proof = append(proof, validatorListLenLE)
	return proof, nil
}
