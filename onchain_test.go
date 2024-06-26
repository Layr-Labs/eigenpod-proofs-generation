package eigenpodproofs_test

import (
	"math/big"
	"testing"

	contractBeaconChainProofsWrapper "github.com/Layr-Labs/eigenpod-proofs-generation/bindings/BeaconChainProofsWrapper"
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
		contractBeaconChainProofsWrapper.BeaconChainProofsStateRootProof{
			BeaconStateRoot: verifyValidatorFieldsCallParams.StateRootProof.BeaconStateRoot,
			Proof:           verifyValidatorFieldsCallParams.StateRootProof.Proof.ToByteSlice(),
		},
	)
	assert.Nil(t, err)

	for i := 0; i < len(verifyValidatorFieldsCallParams.ValidatorFields); i++ {
		validatorFields := [][32]byte{}
		for _, field := range verifyValidatorFieldsCallParams.ValidatorFields[i] {
			validatorFields = append(validatorFields, field)
		}

		err = beaconChainProofsWrapper.VerifyValidatorFields(
			&bind.CallOpts{},
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

	err = beaconChainProofsWrapper.VerifyBalanceContainer(
		&bind.CallOpts{},
		blockRoot,
		contractBeaconChainProofsWrapper.BeaconChainProofsBalanceContainerProof{
			BalanceContainerRoot: verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot,
			Proof:                verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.Proof.ToByteSlice(),
		},
	)
	assert.Nil(t, err)

	for i := 0; i < len(verifyCheckpointProofsCallParams.BalanceProofs); i++ {
		err = beaconChainProofsWrapper.VerifyValidatorBalance(
			&bind.CallOpts{},
			verifyCheckpointProofsCallParams.ValidatorBalancesRootProof.ValidatorBalancesRoot,
			new(big.Int).SetUint64(validatorIndices[i]),
			contractBeaconChainProofsWrapper.BeaconChainProofsBalanceProof{
				PubkeyHash:  verifyCheckpointProofsCallParams.BalanceProofs[i].PubkeyHash,
				BalanceRoot: verifyCheckpointProofsCallParams.BalanceProofs[i].BalanceRoot,
				Proof:       verifyCheckpointProofsCallParams.BalanceProofs[i].Proof.ToByteSlice(),
			},
		)
		assert.Nil(t, err)
	}
}
