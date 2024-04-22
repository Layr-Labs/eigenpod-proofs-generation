package eigenpodproofs

import (
	"testing"

	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
)

func BenchmarkComputeBeaconStateRoot(b *testing.B) {
	computed, err := epp.ComputeBeaconStateRoot(&oracleState)
	if err != nil {
		b.Fatal(err)
	}

	var cached phase0.Root
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeBeaconStateRoot(&oracleState)
		if err != nil {
			b.Fatal(err)
		}
	}

	assert.Equal(b, computed, cached)
}

func BenchmarkComputeBeaconStateTopLevelRoots(b *testing.B) {
	versionedState := spec.VersionedBeaconState{Deneb: &oracleState}
	versionedState.Version = spec.DataVersionDeneb
	computed, err := epp.ComputeBeaconStateTopLevelRoots(&versionedState)
	if err != nil {
		b.Fatal(err)
	}

	var cached *beacon.BeaconStateTopLevelRoots
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeBeaconStateTopLevelRoots(&versionedState)
		if err != nil {
			b.Fatal(err)
		}

	}
	assert.Equal(b, computed, cached)
}

func BenchmarkComputeValidatorTree(b *testing.B) {
	computed, err := epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)
	if err != nil {
		b.Fatal(err)
	}

	var cached [][]phase0.Root
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeValidatorTree(oracleState.Slot, oracleState.Validators)
		if err != nil {
			b.Fatal(err)
		}
	}

	assert.Equal(b, computed, cached)
}
