package eigenpodproofs

import (
	"testing"

	"github.com/attestantio/go-eth2-client/spec"
	"github.com/stretchr/testify/assert"
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
	computed, err := epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		cached, err := epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)
		if err != nil {
			b.Fatal(err)
		}
		assert.Equal(b, computed, cached)
	}
}
