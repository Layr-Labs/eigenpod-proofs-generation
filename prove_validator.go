package eigenpodproofs

import (
	"crypto/sha256"
	"math/big"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"

	beacon "github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
)

type StateRootProof struct {
	BeaconStateRoot phase0.Root  `json:"beaconStateRoot"`
	Proof           common.Proof `json:"stateRootProof"`
}

type VerifyValidatorFieldsCallParams struct {
	StateRootProof        *StateRootProof `json:"stateRootProof"`
	ValidatorIndices      []uint64        `json:"validatorIndices"`
	ValidatorFieldsProofs []common.Proof  `json:"validatorFieldsProofs"`
	ValidatorFields       [][]Bytes32     `json:"validatorFields"`
}

type ValidatorBalancesRootProof struct {
	ValidatorBalancesRoot phase0.Root  `json:"validatorBalanceRoot"`
	Proof                 common.Proof `json:"proof"`
}

type BalanceProof struct {
	PubkeyHash  [32]byte     `json:"pubkeyHash"`
	BalanceRoot phase0.Root  `json:"balanceRoot"`
	Proof       common.Proof `json:"proof"`
}

type VerifyCheckpointProofsCallParams struct {
	ValidatorBalancesRootProof *ValidatorBalancesRootProof `json:"validatorBalancesRootProof"`
	BalanceProofs              []*BalanceProof             `json:"balanceProofs"`
}

// ProveValidatorContainers generates proofs for the validator containers.
// oracleBlockHeader is the block header of block whose state root will be looked up from the EIP-4788 precompile
// oracleBeaconState is the beacon state corresponding to the oracleBlockHeader
// validatorIndices is the list of validator indices for which the proofs are to be generated
func (epp *EigenPodProofs) ProveValidatorContainers(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *spec.VersionedBeaconState, validatorIndices []uint64) (*VerifyValidatorFieldsCallParams, error) {
	oracleBeaconStateSlot, err := oracleBeaconState.Slot()
	if err != nil {
		return nil, err
	}
	oracleBeaconStateValidators, err := oracleBeaconState.Validators()
	if err != nil {
		return nil, err
	}

	verifyValidatorFieldsCallParams := &VerifyValidatorFieldsCallParams{}

	// Get the state root proof
	verifyValidatorFieldsCallParams.StateRootProof = &StateRootProof{}
	verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot = oracleBlockHeader.StateRoot
	verifyValidatorFieldsCallParams.StateRootProof.Proof, err = beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	verifyValidatorFieldsCallParams.ValidatorIndices = make([]uint64, len(validatorIndices))
	verifyValidatorFieldsCallParams.ValidatorFieldsProofs = make([]common.Proof, len(validatorIndices))
	verifyValidatorFieldsCallParams.ValidatorFields = make([][]Bytes32, len(validatorIndices))
	for i, validatorIndex := range validatorIndices {
		verifyValidatorFieldsCallParams.ValidatorIndices[i] = validatorIndex
		// prove the validator fields against the beacon state
		verifyValidatorFieldsCallParams.ValidatorFieldsProofs[i], err = epp.proveValidatorAgainstBeaconState(beaconStateTopLevelRoots, oracleBeaconStateSlot, oracleBeaconStateValidators, validatorIndex)
		if err != nil {
			return nil, err
		}

		verifyValidatorFieldsCallParams.ValidatorFields[i] = ConvertValidatorToValidatorFields(oracleBeaconStateValidators[validatorIndex])
	}

	return verifyValidatorFieldsCallParams, nil
}

func (epp *EigenPodProofs) ProveCheckpointProofs(oracleBlockHeader *phase0.BeaconBlockHeader, oracleBeaconState *spec.VersionedBeaconState, validatorIndices []uint64) (*VerifyCheckpointProofsCallParams, error) {
	oracleBeaconStateSlot, err := oracleBeaconState.Slot()
	if err != nil {
		return nil, err
	}

	oracleBeaconStateValidators, err := oracleBeaconState.Validators()
	if err != nil {
		return nil, err
	}

	oracleBeaconStateValidatorBalances, err := oracleBeaconState.ValidatorBalances()
	if err != nil {
		return nil, err
	}

	verifyCheckpointProofsCallParams := &VerifyCheckpointProofsCallParams{}

	// Get beacon state top level roots
	beaconStateTopLevelRoots, err := epp.ComputeBeaconStateTopLevelRoots(oracleBeaconState)
	if err != nil {
		return nil, err
	}

	// Get state root proof
	verifyCheckpointProofsCallParams.ValidatorBalancesRootProof = &ValidatorBalancesRootProof{}
	stateRootProof, err := beacon.ProveStateRootAgainstBlockHeader(oracleBlockHeader)
	if err != nil {
		return nil, err
	}

	// prove the validator balances root against the beacon state root
	balancesRootProof, err := beacon.ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, beacon.BALANCES_INDEX)
	if err != nil {
		return nil, err
	}

	verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot = *beaconStateTopLevelRoots.BalancesRoot
	verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.Proof = append(balancesRootProof, stateRootProof...)

	verifyCheckpointProofsCallParams.BalanceProofs = make([]*BalanceProof, len(validatorIndices))
	for i, validatorIndex := range validatorIndices {
		balanceRoot, balanceProof, err := epp.proveValidatorBalanceAgainstBeaconState(beaconStateTopLevelRoots, oracleBeaconStateSlot, oracleBeaconStateValidatorBalances, validatorIndex)
		if err != nil {
			return nil, err
		}

		var pubkeyHash [32]byte
		pubkeyHashVariable := computePubkeyHash(oracleBeaconStateValidators[validatorIndex].PublicKey[:])
		copy(pubkeyHash[:], pubkeyHashVariable[:])

		verifyCheckpointProofsCallParams.BalanceProofs[i] = &BalanceProof{
			PubkeyHash:  pubkeyHash,
			BalanceRoot: balanceRoot,
			Proof:       balanceProof,
		}
	}

	return verifyCheckpointProofsCallParams, nil
}

func computePubkeyHash(publicKey []byte) [32]byte {
	// ssz requires 32 byte alignment, and publicKey is 48 bytes.
	// tack on 16 bytes of 0's and compute the sha256 sum.
	zeroes := [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	return sha256.Sum256(
		append(
			publicKey,
			zeroes[:]...,
		),
	)
}

func (epp *EigenPodProofs) proveValidatorAgainstBeaconState(beaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots, oracleBeaconStateSlot phase0.Slot, oracleBeaconStateValidators []*phase0.Validator, validatorIndex uint64) (common.Proof, error) {
	// prove the validator list against the beacon state
	validatorListProof, err := beacon.ProveBeaconTopLevelRootAgainstBeaconState(beaconStateTopLevelRoots, beacon.VALIDATORS_INDEX)
	if err != nil {
		return nil, err
	}

	// prove the validator root against the validator list root
	validatorProof, err := epp.proveValidatorAgainstValidatorList(oracleBeaconStateSlot, oracleBeaconStateValidators, validatorIndex)
	if err != nil {
		return nil, err
	}

	proof := append(validatorProof, validatorListProof...)

	return proof, nil
}

func (epp *EigenPodProofs) proveValidatorAgainstValidatorList(slot phase0.Slot, validators []*phase0.Validator, validatorIndex uint64) (common.Proof, error) {
	validatorTree, err := epp.ComputeValidatorTree(slot, validators)
	if err != nil {
		return nil, err
	}

	proof, err := common.ComputeMerkleProofFromTree(validatorTree, validatorIndex, beacon.VALIDATOR_TREE_HEIGHT)
	if err != nil {
		return nil, err
	}
	// append the length of the validator array to the proof
	// convert big endian to little endian
	validatorListLenLE := BigToLittleEndian(big.NewInt(int64(len(validators))))

	proof = append(proof, validatorListLenLE)
	return proof, nil
}

func (epp *EigenPodProofs) proveValidatorBalanceAgainstBeaconState(beaconStateTopLevelRoots *beacon.BeaconStateTopLevelRoots, oracleBeaconStateSlot phase0.Slot, oracleBeaconStateValidatorBalances []phase0.Gwei, validatorIndex uint64) (phase0.Root, common.Proof, error) {
	// prove the validator root against the validator list root
	balanceRoot, balanceProof, err := epp.proveValidatorBalanceAgainstValidatorBalancesList(oracleBeaconStateSlot, oracleBeaconStateValidatorBalances, validatorIndex)
	if err != nil {
		return phase0.Root{}, nil, err
	}

	return balanceRoot, balanceProof, nil
}

func (epp *EigenPodProofs) proveValidatorBalanceAgainstValidatorBalancesList(slot phase0.Slot, balances []phase0.Gwei, validatorIndex uint64) (phase0.Root, common.Proof, error) {
	validatorBalancesTree, err := epp.ComputeValidatorBalancesTree(slot, balances)
	if err != nil {
		return phase0.Root{}, nil, err
	}

	// 4 balances per leaf
	validatorBalancesIndex := validatorIndex / 4

	proof, err := common.ComputeMerkleProofFromTree(validatorBalancesTree, validatorBalancesIndex, beacon.GetValidatorBalancesProofDepth(len(balances)))
	if err != nil {
		return phase0.Root{}, nil, err
	}
	// append the little endian length of the balances array to the proof
	validatorBalancesListLenLE := BigToLittleEndian(big.NewInt(int64(len(balances))))

	proof = append(proof, validatorBalancesListLenLE)
	return validatorBalancesTree[0][validatorBalancesIndex], proof, nil
}
