package gloasproofs

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// Synthetic mainnet-preset states with a populated validator registry and its
// per-validator lists. This measures resource use, not consensus validity.
func BenchmarkGenerator(b *testing.B) {
	for _, count := range []int{1_000, 2_000_000} {
		b.Run(fmt.Sprint(count), func(b *testing.B) {
			s := fixtureState()
			base := *s.Validators[0]
			validators := make([]phase0.Validator, count)
			s.Validators = make([]*phase0.Validator, count)
			s.Balances = make([]phase0.Gwei, count)
			s.PreviousEpochParticipation = make([]altair.ParticipationFlags, count)
			s.CurrentEpochParticipation = make([]altair.ParticipationFlags, count)
			s.InactivityScores = make([]uint64, count)
			for i := range validators {
				validators[i] = base
				binary.LittleEndian.PutUint64(validators[i].PublicKey[:8], uint64(i+1))
				s.Validators[i] = &validators[i]
				s.Balances[i] = phase0.Gwei(32e9 + i + 1)
			}
			stateRoot, err := s.HashTreeRoot()
			if err != nil {
				b.Fatal(err)
			}
			header := &phase0.BeaconBlockHeader{Slot: s.Slot, StateRoot: stateRoot}
			root, err := header.HashTreeRoot()
			if err != nil {
				b.Fatal(err)
			}
			indices := []uint64{0, uint64(count - 1)}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				g, err := NewGenerator(header, s, Hash(root))
				if err != nil {
					b.Fatal(err)
				}
				if _, err := g.Credentials(indices); err != nil {
					b.Fatal(err)
				}
				if _, err := g.Checkpoint(indices, Hash{0x22}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
