package beacon

import (
	"math/big"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

var zeroBytes = [32]byte{}

func ComputeValidatorTreeLeaves(validators []*phase0.Validator) ([]phase0.Root, error) {
	validatorNodeList := make([]phase0.Root, len(validators))
	for i := 0; i < len(validators); i++ {
		validatorRoot, err := validators[i].HashTreeRoot()
		if err != nil {
			return nil, err
		}
		validatorNodeList[i] = phase0.Root(validatorRoot)
	}

	return validatorNodeList, nil
}

func ComputeValidatorBalancesTreeLeaves(balances []phase0.Gwei) []phase0.Root {
	buf := []byte{}

	for i := 0; i < len(balances); i++ {
		buf = ssz.MarshalUint64(buf, uint64(balances[i]))
	}
	// pad the buffer with 0s to make it a multiple of 32 bytes
	if rest := len(buf) % 32; rest != 0 {
		buf = append(buf, zeroBytes[:32-rest]...)
	}

	// //now we divide the buffer into 32 byte chunks
	numLeaves := len(buf) / 32
	balanceRootList := make([]phase0.Root, numLeaves)
	for i := 0; i < numLeaves; i++ {
		copy(balanceRootList[i][:], buf[i*32:(i+1)*32])
	}

	return balanceRootList
}

func GetValidatorBalancesProofDepth(numBalances int) uint64 {
	return uint64(common.GetDepth(ssz.CalculateLimit(1099511627776, uint64(numBalances), 8)))
}

func ProveValidatorBalanceAgainstValidatorBalanceList(balances []phase0.Gwei, validatorIndex uint64) (phase0.Root, common.Proof, error) {
	balanceRootList := ComputeValidatorBalancesTreeLeaves(balances)

	// refer to beaconstate_ssz.go in go-eth2-client
	numLayers := uint64(common.GetDepth(ssz.CalculateLimit(1099511627776, uint64(len(balances)), 8)))
	validatorBalanceIndex := uint64(validatorIndex / 4)
	proof, err := common.GetProof(balanceRootList, validatorBalanceIndex, numLayers)

	if err != nil {
		return phase0.Root{}, nil, err
	}

	// append the length of the balance array to the proof
	// convert big endian to little endian
	balanceListLenLE := common.BigToLittleEndian(big.NewInt(int64(len(balances))))

	proof = append(proof, balanceListLenLE)
	return balanceRootList[validatorBalanceIndex], proof, nil
}

func GetBalanceRoots(balances []phase0.Gwei) ([]phase0.Root, error) {
	buf := []byte{}

	for i := 0; i < len(balances); i++ {
		buf = ssz.MarshalUint64(buf, uint64(balances[i]))
	}

	// //now we divide the buffer into 32 byte chunks
	numLeaves := len(buf) / 32
	balanceRootList := make([]phase0.Root, numLeaves)
	for i := 0; i < numLeaves; i++ {
		copy(balanceRootList[i][:], buf[i*32:(i+1)*32])
	}
	return balanceRootList, nil
}
