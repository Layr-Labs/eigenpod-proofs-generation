package txsubmitter

import eigenpodproofs "github.com/Layr-Labs/eigenpod-proofs-generation"

type EigenPodProofTxSubmitter struct {
	chainClient    *ChainClient
	eigenPodProofs *eigenpodproofs.EigenPodProofs
}

func ConvertProofsToBytes32Array(proof []eigenpodproofs.Bytes32) [][32]byte {
	proofBytes32 := make([][32]byte, len(proof))
	for i, e := range proof {
		proofBytes32[i] = e
	}
	return proofBytes32
}

func NewEigenPodProofTxSubmitter(chainClient ChainClient, epp eigenpodproofs.EigenPodProofs) *EigenPodProofTxSubmitter {

	return &EigenPodProofTxSubmitter{
		chainClient:    &chainClient,
		eigenPodProofs: &epp,
	}
}
