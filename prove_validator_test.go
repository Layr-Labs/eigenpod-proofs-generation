package eigenpodproofs_test

import (
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
)

func TestProveValidatorContainers(t *testing.T) {
	validators, err := beaconState.Validators()
	if err != nil {
		t.Fatal(err)
	}
	// loop through the beacon state and get every thousandth validator index
	validatorIndices := []uint64{}
	for i := int(0); i < len(validators); i += 1000 {
		validatorIndices = append(validatorIndices, uint64(i))
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(beaconHeader, beaconState, validatorIndices)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, verifyStateRootAgainstBlockHeader(t, epp, beaconHeader, beaconState.Deneb, verifyValidatorFieldsCallParams.StateRootProof.Proof))

	for i := 0; i < len(verifyValidatorFieldsCallParams.ValidatorFields); i++ {
		assert.True(t, verifyValidatorAgainstBeaconState(t, epp, beaconState.Deneb, verifyValidatorFieldsCallParams.ValidatorFieldsProofs[i], validatorIndices[i]))
	}
}

func TestProveValidatorBalances(t *testing.T) {
	validators, err := beaconState.Validators()
	if err != nil {
		t.Fatal(err)
	}
	// loop through the beacon state and get every thousandth validator index
	validatorIndices := []uint64{}
	for i := int(0); i < len(validators); i += 1000 {
		validatorIndices = append(validatorIndices, uint64(i))
	}

	verifyCheckpointProofsCallParams, err := epp.ProveCheckpointProofs(beaconHeader, beaconState, validatorIndices)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, verifyValidatorBalancesRootAgainstBlockHeader(t, epp, beaconHeader, verifyCheckpointProofsCallParams.ValidatorBalancesRootProof))

	for i := 0; i < len(verifyCheckpointProofsCallParams.BalanceProofs); i++ {
		assert.True(t, verifyValidatorBalanceAgainstValidatorBalancesRoot(t, epp, beaconState.Deneb, verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot, verifyCheckpointProofsCallParams.BalanceProofs[i], validatorIndices[i]))
	}
}

func verifyStateRootAgainstBlockHeader(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleBlockHeader *phase0.BeaconBlockHeader, oracleState *deneb.BeaconState, proof common.Proof) bool {
	root, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}

	leaf, err := epp.ComputeBeaconStateRoot(oracleState)
	if err != nil {
		t.Fatal(err)
	}

	return common.ValidateProof(root, proof, leaf, beacon.STATE_ROOT_INDEX)
}

func verifyValidatorAgainstBeaconState(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleState *deneb.BeaconState, proof common.Proof, validatorIndex uint64) bool {
	leaf, err := oracleState.Validators[validatorIndex].HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}

	root, err := epp.ComputeBeaconStateRoot(oracleState)
	if err != nil {
		t.Fatal(err)
	}

	index := beacon.VALIDATORS_INDEX<<(beacon.VALIDATOR_TREE_HEIGHT+1) | validatorIndex
	return common.ValidateProof(root, proof, leaf, index)
}

func verifyValidatorBalancesRootAgainstBlockHeader(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleBlockHeader *phase0.BeaconBlockHeader, proof *eigenpodproofs.ValidatorBalancesRootProof) bool {
	root, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}

	return common.ValidateProof(root, proof.Proof, proof.ValidatorBalancesRoot, beacon.STATE_ROOT_INDEX<<beacon.BEACON_STATE_TREE_HEIGHT|beacon.BALANCES_INDEX)
}

func verifyValidatorBalanceAgainstValidatorBalancesRoot(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleState *deneb.BeaconState, validatorBalancesRoot phase0.Root, proof *eigenpodproofs.BalanceProof, validatorIndex uint64) bool {
	index := beacon.BALANCES_INDEX<<(beacon.GetValidatorBalancesProofDepth(len(oracleState.Balances))+1) | (validatorIndex / 4)

	return common.ValidateProof(validatorBalancesRoot, proof.Proof, proof.BalanceRoot, index)
}
