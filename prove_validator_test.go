package eigenpodproofs_test

import (
	"fmt"
	"os"
	"testing"

	eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"
	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
)

var beaconHeader *phase0.BeaconBlockHeader
var beaconState *spec.VersionedBeaconState
var epp *eigenpodproofs.EigenPodProofs

// before all
func TestMain(m *testing.M) {
	var err error

	beaconHeaderFileName := "data/deneb_holesky_beacon_header_1650726.json"
	beaconHeaderBytes, err := common.ReadFile(beaconHeaderFileName)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	beaconStateFileName := "data/deneb_holesky_beacon_state_1650726.ssz"
	beaconStateBytes, err := common.ReadFile(beaconStateFileName)
	if err != nil {
		fmt.Println("error", err)
		return
	}

	beaconHeader = &phase0.BeaconBlockHeader{}
	err = beaconHeader.UnmarshalJSON(beaconHeaderBytes)
	if err != nil {
		fmt.Println("error", err)
	}

	beaconState, err = beacon.UnmarshalSSZVersionedBeaconState(beaconStateBytes)
	if err != nil {
		fmt.Println("error", err)
	}

	epp, err = eigenpodproofs.NewEigenPodProofs(17000, 600)
	if err != nil {
		fmt.Println("error", err)
	}

	code := m.Run()

	// // Teardown
	// log.Println("Tearing down suite")
	// teardownSuite()

	// Exit with test result code
	os.Exit(code)
}

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

	assert.True(t, verifyStateRootAgainstBlockHeaderProof(t, epp, beaconHeader, beaconState.Deneb, verifyValidatorFieldsCallParams.StateRootProof.Proof))

	for i := 0; i < len(verifyValidatorFieldsCallParams.ValidatorFields); i++ {
		assert.True(t, verifyValidatorAgainstBeaconState(t, epp, beaconState.Deneb, verifyValidatorFieldsCallParams.ValidatorFieldsProofs[i], validatorIndices[i]))
	}
}

func TestProveValidatorBalances(t *testing.T) {
	epp, err := eigenpodproofs.NewEigenPodProofs(17000, 600)
	if err != nil {
		t.Fatal(err)
	}

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

	assert.True(t, verifyStateRootAgainstBlockHeaderProof(t, epp, beaconHeader, beaconState.Deneb, verifyCheckpointProofsCallParams.StateRootProof.Proof))

	for i := 0; i < len(verifyCheckpointProofsCallParams.BalanceProofs); i++ {
		assert.True(t, verifyValidatorBalanceAgainstBeaconState(t, epp, beaconState.Deneb, verifyCheckpointProofsCallParams.BalanceProofs[i], validatorIndices[i]))
	}
}

func verifyStateRootAgainstBlockHeaderProof(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleBlockHeader *phase0.BeaconBlockHeader, oracleState *deneb.BeaconState, proof common.Proof) bool {
	root, err := oracleBlockHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := epp.ComputeBeaconStateRoot(oracleState)
	if err != nil {
		t.Fatal(err)
	}

	flag := common.ValidateProof(root, proof, leaf, beacon.StateRootIndex)
	return flag
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

	index := beacon.ValidatorListIndex<<(beacon.ValidatorListMerkleSubtreeNumLayers+1) | validatorIndex
	return common.ValidateProof(root, proof, leaf, index)
}

func verifyValidatorBalanceAgainstBeaconState(t *testing.T, epp *eigenpodproofs.EigenPodProofs, oracleState *deneb.BeaconState, proof *eigenpodproofs.BalanceProof, validatorIndex uint64) bool {
	root, err := epp.ComputeBeaconStateRoot(oracleState)
	if err != nil {
		t.Fatal(err)
	}

	index := beacon.ValidatorBalancesListIndex<<(beacon.GetValidatorBalancesProofDepth(len(oracleState.Balances))+1) | (validatorIndex / 4)

	return common.ValidateProof(root, proof.Proof, proof.BalanceRoot, index)
}
