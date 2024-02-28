package eigenpodproofs

import "testing"

func BenchmarkComputeBeaconStateRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = epp.ComputeBeaconStateRoot(&oracleState)
	}

}
