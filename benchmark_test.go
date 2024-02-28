package eigenpodproofs

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
)

func BenchmarkComputeBeaconStateRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = epp.ComputeBeaconStateRoot(&oracleState)
	}

}

func BenchmarkComputeBeaconStateTopLevelRoots(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = epp.ComputeBeaconStateTopLevelRoots(&spec.VersionedBeaconState{Deneb: &oracleState})
	}
}

func BenchmarkComputeValidatorTree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)
	}
}
