package eigenpodproofs_test

import (
	"math/big"
	"testing"

	BeaconChainProofsWrapper "github.com/Layr-Labs/eigenpod-proofs-generation/bindings/BeaconChainProofsWrapper"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/assert"
)

func TestValidatorContainersProofOnChain(t *testing.T) {
	validators, err := beaconState.Validators()
	if err != nil {
		t.Fatal(err)
	}

	validatorIndices := []uint64{}
	for i := int(0); i < len(validators); i += 100000 {
		validatorIndices = append(validatorIndices, uint64(i))
	}

	verifyValidatorFieldsCallParams, err := epp.ProveValidatorContainers(beaconHeader, beaconState, validatorIndices)
	if err != nil {
		t.Fatal(err)
	}

	blockRoot, err := beaconHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}

	err = beaconChainProofsWrapper.VerifyStateRoot(
		&bind.CallOpts{},
		blockRoot,
		BeaconChainProofsWrapper.BeaconChainProofsStateRootProof{
			BeaconStateRoot: verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
			Proof:           verifyValidatorFieldsCallParams.StateRootProof.Proof.ToByteSlice(),
		},
	)
	assert.Nil(t, err)

	// Update the proof timestamp depending on the beacon state version
	var proofTimestamp uint64
	if beaconState.Version == spec.DataVersionElectra {
		proofTimestamp = uint64(1730822401) // 1 second after mekong genesis
	} else {
		proofTimestamp = uint64(0)
	}

	for i := 0; i < len(verifyValidatorFieldsCallParams.ValidatorFields); i++ {
		validatorFields := [][32]byte{}
		for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[i] {
			validatorFields = append(validatorFields, field)
		}

		err = beaconChainProofsWrapper.VerifyValidatorFields(
			&bind.CallOpts{},
			proofTimestamp,
			verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
			validatorFields,
			verifyValidatorFieldsCallParams.ValidatorFieldsProofs[i].ToByteSlice(),
			new(big.Int).SetUint64(verifyValidatorFieldsCallParams.ValidatorIndices[i]),
		)
		assert.Nil(t, err)
	}
}

func TestValidatorBalancesProofOnChain(t *testing.T) {
	validators, err := beaconState.Validators()
	if err != nil {
		t.Fatal(err)
	}

	validatorIndices := []uint64{}
	for i := int(0); i < len(validators); i += 100000 {
		validatorIndices = append(validatorIndices, uint64(i))
	}

	verifyCheckpointProofsCallParams, err := epp.ProveCheckpointProofs(beaconHeader, beaconState, validatorIndices)
	if err != nil {
		t.Fatal(err)
	}

	blockRoot, err := beaconHeader.HashTreeRoot()
	if err != nil {
		t.Fatal(err)
	}

	// Update the proof timestamp depending on the beacon state version
	var proofTimestamp uint64
	if beaconState.Version == spec.DataVersionElectra {
		proofTimestamp = uint64(1730822401) // 1 second after mekong genesis
	} else {
		proofTimestamp = uint64(0)
	}

	err = beaconChainProofsWrapper.VerifyBalanceContainer(
		&bind.CallOpts{},
		proofTimestamp,
		blockRoot,
		BeaconChainProofsWrapper.BeaconChainProofsBalanceContainerProof{
			BalanceContainerRoot: verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot,
			Proof:                verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.Proof.ToByteSlice(),
		},
	)
	assert.Nil(t, err)

	for i := 0; i < len(verifyCheckpointProofsCallParams.BalanceProofs); i++ {
		_, err = beaconChainProofsWrapper.VerifyValidatorBalance(
			&bind.CallOpts{},
			verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot,
			new(big.Int).SetUint64(validatorIndices[i]),
			BeaconChainProofsWrapper.BeaconChainProofsBalanceProof{
				PubkeyHash:  verifyCheckpointProofsCallParams.BalanceProofs[i].PubkeyHash,
				BalanceRoot: verifyCheckpointProofsCallParams.BalanceProofs[i].BalanceRoot,
				Proof:       verifyCheckpointProofsCallParams.BalanceProofs[i].Proof.ToByteSlice(),
			},
		)
		assert.Nil(t, err)
	}
}
