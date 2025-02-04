package eigenpodproofs_test

import (
	"testing"

	"github.com/Layr-Labs/eigenpod-proofs-generation/beacon"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/assert"
)

func BenchmarkComputeBeaconStateRoot(b *testing.B) {
	computed, err := epp.ComputeBeaconStateRoot(beaconState)
	if err != nil {
		b.Fatal(err)
	}

	var cached phase0.Root
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeBeaconStateRoot(beaconState)
		if err != nil {
			b.Fatal(err)
		}
	}

	assert.Equal(b, computed, cached)
}

func BenchmarkComputeBeaconStateTopLevelRoots(b *testing.B) {
	computed, err := epp.ComputeBeaconStateTopLevelRoots(beaconState)
	if err != nil {
		b.Fatal(err)
	}

	var cached *beacon.VersionedBeaconStateTopLevelRoots
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeBeaconStateTopLevelRoots(beaconState)
		if err != nil {
			b.Fatal(err)
		}
	}
	assert.Equal(b, computed, cached)
}

func BenchmarkComputeValidatorTree(b *testing.B) {
	// If beacon state is not Deneb, then skip this test
	if beaconState.Version != spec.DataVersionDeneb {
		b.Skip("skipping test for non-Deneb beacon state")
	}

	computed, err := epp.ComputeValidatorTree(beaconState.Deneb.Slot, beaconState.Deneb.Validators)
	if err != nil {
		b.Fatal(err)
	}

	var cached [][]phase0.Root
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeValidatorTree(beaconState.Deneb.Slot, beaconState.Deneb.Validators)
		if err != nil {
			b.Fatal(err)
		}
	}

	assert.Equal(b, computed, cached)
}

func BenchmarkComputeValidatorBalancesTree(b *testing.B) {
	// If beacon state is not Deneb, then skip this test
	if beaconState.Version != spec.DataVersionDeneb {
		b.Skip("skipping test for non-Deneb beacon state")
	}

	computed, err := epp.ComputeValidatorBalancesTree(beaconState.Deneb.Slot, beaconState.Deneb.Balances)
	if err != nil {
		b.Fatal(err)
	}

	var cached [][]phase0.Root
	for i := 0; i < b.N; i++ {
		cached, err = epp.ComputeValidatorBalancesTree(beaconState.Deneb.Slot, beaconState.Deneb.Balances)
		if err != nil {
			b.Fatal(err)
		}
	}

	assert.Equal(b, computed, cached)
}
