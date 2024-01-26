package eigenpodproofs

import (
	"math/big"

	"github.com/attestantio/go-eth2-client/spec/deneb"
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

func (epp *EigenPodProofs) ProveValidatorContainer(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *deneb.BeaconState, validatorIndices []uint64) (*VerifyValidatorFieldsCallParams, error) {
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

	VerifyValidatorFieldsCallParams.OracleTimestamp = GetSlotTimestamp(oracleBeaconState, oracleBlockHeader)
	VerifyValidatorFieldsCallParams.ValidatorIndices = make([]uint64, len(validatorIndices))
	VerifyValidatorFieldsCallParams.ValidatorFieldsProofs = make([]common.Proof, len(validatorIndices))
	VerifyValidatorFieldsCallParams.ValidatorFields = make([][]Bytes32, len(validatorIndices))
	for i, validatorIndex := range validatorIndices {
		VerifyValidatorFieldsCallParams.ValidatorIndices[i] = validatorIndex
		// prove the validator fields against the beacon state
		VerifyValidatorFieldsCallParams.ValidatorFieldsProofs[i], err = epp.ProveValidatorAgainstBeaconState(oracleBeaconState, beaconStateTopLevelRoots, validatorIndex)
		if err != nil {
			return nil, err
		}

		VerifyValidatorFieldsCallParams.ValidatorFields[i] = ConvertValidatorToValidatorFields(oracleBeaconState.Validators[validatorIndex])
	}

	return VerifyValidatorFieldsCallParams, nil
}

func (epp *EigenPodProofs) ProveValidatorFields(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *deneb.BeaconState, validatorIndex uint64) (*StateRootProof, common.Proof, error) {
	stateRootProof := &StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, nil, err
	}

	// Get beacon state root. TODO: Combine this cheaply with compute beacon state top level roots
	stateRootProof.BeaconStateRoot = oracleBlockHeader.StateRoot
	if err != nil {
		return nil, nil, err
	}

	stateRootProof.StateRootProof, err = beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)

	if err != nil {
		return nil, nil, err
	}

	validatorFieldsProof, err := epp.ProveValidatorAgainstBeaconState(oracleBeaconState, beaconStateTopLevelRoots, validatorIndex)

	if err != nil {
		return nil, nil, err
	}

	return stateRootProof, validatorFieldsProof, nil
}

func (epp *EigenPodProofs) ProveValidatorAgainstBeaconState(oracleBeaconState *deneb.BeaconState, beaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots, validatorIndex uint64) (common.Proof, error) {
	// prove the validator list against the beacon state
	validatorListProof, err := beacon.ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, beacon.ValidatorListIndex)
	if err != nil {
		return nil, err
	}

	// prove the validator root against the validator list root
	validatorProof, err := epp.ProveValidatorAgainstValidatorList(oracleBeaconState.Slot, oracleBeaconState.Validators, validatorIndex)
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
