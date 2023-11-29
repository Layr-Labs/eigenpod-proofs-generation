package main

import (
	"math/big"

	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

type VerifyWithdrawalCredentialsCallParams struct {
	OracleTimestamp       uint64          `json:"oracleTimestamp"`
	StateRootProof        *StateRootProof `json:"stateRootProof"`
	ValidatorIndices      []uint64        `json:"validatorIndices"`
	ValidatorFieldsProofs []Proof         `json:"validatorFieldsProofs"`
	ValidatorFields       [][]Bytes32     `json:"validatorFields"`
}

func (epp *EigenPodProofs) ProveValidatorWithdrawalCredentials(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *capella.BeaconState, validatorIndices []uint64) (*VerifyWithdrawalCredentialsCallParams, error) {
	verifyWithdrawalCredentialsCallParams := &VerifyWithdrawalCredentialsCallParams{}
	verifyWithdrawalCredentialsCallParams.StateRootProof = &StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	// Get beacon state root.
	verifyWithdrawalCredentialsCallParams.StateRootProof.BeaconStateRoot, err = epp.ComputeBeaconStateRoot(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	verifyWithdrawalCredentialsCallParams.StateRootProof.StateRootProof, err = ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	verifyWithdrawalCredentialsCallParams.OracleTimestamp = GetSlotTimestamp(oracleBeaconState, oracleBlockHeader)
	verifyWithdrawalCredentialsCallParams.ValidatorIndices = make([]uint64, len(validatorIndices))
	verifyWithdrawalCredentialsCallParams.ValidatorFieldsProofs = make([]Proof, len(validatorIndices))
	verifyWithdrawalCredentialsCallParams.ValidatorFields = make([][]Bytes32, len(validatorIndices))
	for i, validatorIndex := range validatorIndices {
		verifyWithdrawalCredentialsCallParams.ValidatorIndices[i] = validatorIndex
		// prove the validator fields against the beacon state
		verifyWithdrawalCredentialsCallParams.ValidatorFieldsProofs[i], err = epp.ProveValidatorAgainstBeaconState(oracleBeaconState, beaconStateTopLevelRoots, validatorIndex)
		if err != nil {
			return nil, err
		}

		verifyWithdrawalCredentialsCallParams.ValidatorFields[i] = ConvertValidatorToValidatorFields(oracleBeaconState.Validators[validatorIndex])
	}

	return verifyWithdrawalCredentialsCallParams, nil
}

type VerifyBalanceUpdateCallParams struct {
	OracleTimestamp    uint64              `json:"oracleTimestamp"`
	ValidatorIndex     uint64              `json:"validatorIndex"`
	StateRootProof     *StateRootProof     `json:"stateRootProof"`
	BalanceUpdateProof *BalanceUpdateProof `json:"validatorFieldsProofs"`
	ValidatorFields    []Bytes32           `json:"validatorFields"`
}

type BalanceUpdateProof struct {
	ValidatorBalanceProof Proof   `json:"validatorBalanceProof"`
	ValidatorFieldsProof  Proof   `json:"validatorFieldsProof"`
	BalanceRoot           Bytes32 `json:"balanceRoot"`
}

func (epp *EigenPodProofs) ProveValidatorBalance(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *capella.BeaconState, validatorIndex uint64) (*VerifyBalanceUpdateCallParams, error) {
	verifyBalanceUpdateCallParams := &VerifyBalanceUpdateCallParams{}
	verifyBalanceUpdateCallParams.StateRootProof = &StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	verifyBalanceUpdateCallParams.StateRootProof.BeaconStateRoot, err = epp.ComputeBeaconStateRoot(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	verifyBalanceUpdateCallParams.StateRootProof.StateRootProof, err = ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	verifyBalanceUpdateCallParams.OracleTimestamp = GetSlotTimestamp(oracleBeaconState, oracleBlockHeader)
	verifyBalanceUpdateCallParams.ValidatorIndex = validatorIndex
	verifyBalanceUpdateCallParams.BalanceUpdateProof = &BalanceUpdateProof{}
	verifyBalanceUpdateCallParams.BalanceUpdateProof.ValidatorBalanceProof, err = epp.ProveValidatorBalanceAgainstBalanceRoot(oracleBeaconState, validatorIndex)
	if err != nil {
		return nil, err
	}

	verifyBalanceUpdateCallParams.BalanceUpdateProof.ValidatorFieldsProof, err = epp.ProveValidatorAgainstBeaconState(oracleBeaconState, beaconStateTopLevelRoots, validatorIndex)
	if err != nil {
		return nil, err
	}

	verifyBalanceUpdateCallParams.BalanceUpdateProof.BalanceRoot = ConvertUint64ToBytes32(uint64(oracleBeaconState.Balances[validatorIndex]))
	verifyBalanceUpdateCallParams.ValidatorFields = ConvertValidatorToValidatorFields(oracleBeaconState.Validators[validatorIndex])

	return verifyBalanceUpdateCallParams, err
}

func (epp *EigenPodProofs) ProveValidatorFields(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *capella.BeaconState, validatorIndex uint64) (*StateRootProof, Proof, error) {
	stateRootProof := &StateRootProof{}
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, nil, err
	}

	// Get beacon state root. TODO: Combine this cheaply with compute beacon state top level roots
	stateRootProof.BeaconStateRoot, err = epp.ComputeBeaconStateRoot(oracleBeaconState)
	if err != nil {
		return nil, nil, err
	}

	stateRootProof.StateRootProof, err = ProveStateRootAgainstBlockHeader(oracleBlockHeader)

	if err != nil {
		return nil, nil, err
	}

	validatorFieldsProof, err := epp.ProveValidatorAgainstBeaconState(oracleBeaconState, beaconStateTopLevelRoots, validatorIndex)

	if err != nil {
		return nil, nil, err
	}

	return stateRootProof, validatorFieldsProof, nil
}

func (epp *EigenPodProofs) ProveValidatorBalanceAgainstBalanceRoot(oracleBeaconState *capella.BeaconState, validatorIndex uint64) ([][32]byte, error) {
	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	// prove the validator balance list root against the beacon state
	beaconStateProof, err := ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, balanceListIndex)
	if err != nil {
		return nil, err
	}

	// prove the validator balance root against the validator balance list root
	balancesProof, err := ProveValidatorBalanceAgainstValidatorBalanceList(oracleBeaconState.Balances, validatorIndex)
	if err != nil {
		return nil, err
	}

	fullBalanceProof := append(balancesProof, beaconStateProof...)
	return fullBalanceProof, nil
}

func (epp *EigenPodProofs) ProveValidatorAgainstBeaconState(oracleBeaconState *capella.BeaconState, beaconStateTopLevelRoots *BeaconStateTopLevelRoots, validatorIndex uint64) (Proof, error) {
	// prove the validator list against the beacon state
	validatorListProof, err := ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, validatorListIndex)
	if err != nil {
		return nil, err
	}

	// prove the validator root against the validator list root
	validatorProof, err := epp.ProveValidatorAgainstValidatorList(uint64(oracleBeaconState.Slot), oracleBeaconState.Validators, validatorIndex)
	if err != nil {
		return nil, err
	}

	proof := append(validatorProof, validatorListProof...)

	return proof, nil
}

func (epp *EigenPodProofs) ProveValidatorAgainstValidatorList(slot uint64, validators []*phase0.Validator, validatorIndex uint64) (Proof, error) {
	validatorTree, err := epp.ComputeValidatorTree(slot, validators)
	if err != nil {
		return nil, err
	}

	proof, err := ComputeMerkleProofFromTree(validatorTree, validatorIndex, validatorListMerkleSubtreeNumLayers)
	if err != nil {
		return nil, err
	}
	//append the length of the validator array to the proof
	//convert big endian to little endian
	validatorListLenLE := BigToLittleEndian(big.NewInt(int64(len(validators))))

	proof = append(proof, validatorListLenLE)
	return proof, nil
}
