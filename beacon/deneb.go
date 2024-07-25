package beacon

import (
	"errors"
	"sync"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/prysmaticlabs/gohashtree"
)

// taken from https://github.com/attestantio/go-eth2-client/blob/21f7dd480fed933d8e0b1c88cee67da721c80eb2/spec/deneb/beaconstate_ssz.go#L640
func ComputeBeaconStateTopLevelRootsDeneb(b *deneb.BeaconState) (*BeaconStateTopLevelRoots, error) {
	var err error
	beaconStateTopLevelRoots := &BeaconStateTopLevelRoots{}

	var errs = make(chan error)
	var wg sync.WaitGroup
	wg.Add(28)

	// Field (0) 'GenesisTime'
	go func() {
		defer wg.Done()
		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		hh.PutUint64(b.GenesisTime)
		tmp0 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.GenesisTimeRoot = &tmp0
	}()

	// Field (1) 'GenesisValidatorsRoot'
	go func() {
		defer wg.Done()
		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.GenesisValidatorsRoot); size != 32 {
			err = ssz.ErrBytesLengthFn("BeaconState.GenesisValidatorsRoot", size, 32)
			errs <- err
			return
		}
		hh.PutBytes(b.GenesisValidatorsRoot[:])
		tmp1 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.GenesisValidatorsRoot = &tmp1
	}()

	// Field (2) 'Slot'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		hh.PutUint64(uint64(b.Slot))
		tmp2 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.SlotRoot = &tmp2
	}()

	// Field (3) 'Fork'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.Fork == nil {
			b.Fork = new(phase0.Fork)
		}
		if err = b.Fork.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp3 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ForkRoot = &tmp3
	}()

	// Field (4) 'LatestBlockHeader'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.LatestBlockHeader == nil {
			b.LatestBlockHeader = new(phase0.BeaconBlockHeader)
		}
		if err = b.LatestBlockHeader.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp4 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.LatestBlockHeaderRoot = &tmp4
	}()

	// Field (5) 'BlockRoots'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.BlockRoots); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.BlockRoots", size, 8192)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.BlockRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				errs <- err
				return
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp5 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.BlockRootsRoot = &tmp5
	}()

	// Field (6) 'StateRoots'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.StateRoots); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.StateRoots", size, 8192)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.StateRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				errs <- err
				return
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp6 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.StateRootsRoot = &tmp6
	}()

	// Field (7) 'HistoricalRoots'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.HistoricalRoots); size > 16777216 {
			err = ssz.ErrListTooBigFn("BeaconState.HistoricalRoots", size, 16777216)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.HistoricalRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				errs <- err
				return
			}
			hh.Append(i[:])
		}
		numItems := uint64(len(b.HistoricalRoots))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(16777216, numItems, 32))
		tmp7 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.HistoricalRootsRoot = &tmp7
	}()

	// Field (8) 'ETH1Data'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.ETH1Data == nil {
			b.ETH1Data = new(phase0.ETH1Data)
		}
		if err = b.ETH1Data.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp8 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ETH1DataRoot = &tmp8
	}()

	// Field (9) 'ETH1DataVotes'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		subIndx := hh.Index()
		num := uint64(len(b.ETH1DataVotes))
		if num > 2048 {
			err = ssz.ErrIncorrectListSize
			errs <- err
			return
		}
		for _, elem := range b.ETH1DataVotes {
			if err = elem.HashTreeRootWith(hh); err != nil {
				errs <- err
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 2048)
		tmp9 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ETH1DataVotesRoot = &tmp9
	}()

	// Field (10) 'ETH1DepositIndex'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		hh.PutUint64(b.ETH1DepositIndex)
		tmp10 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ETH1DepositIndexRoot = &tmp10
	}()

	// Field (11) 'Validators'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		subIndx := hh.Index()
		num := uint64(len(b.Validators))
		if num > 1099511627776 {
			err = ssz.ErrIncorrectListSize
			errs <- err
			return
		}
		for _, elem := range b.Validators {
			if err = elem.HashTreeRootWith(hh); err != nil {
				errs <- err
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 1099511627776)
		tmp11 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ValidatorsRoot = &tmp11
	}()

	// Field (12) 'Balances'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.Balances); size > 1099511627776 {
			err = ssz.ErrListTooBigFn("BeaconState.Balances", size, 1099511627776)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.Balances {
			hh.AppendUint64(uint64(i))
		}
		hh.FillUpTo32()
		numItems := uint64(len(b.Balances))

		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
		tmp12 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.BalancesRoot = &tmp12
	}()

	// Field (13) 'RANDAOMixes'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.RANDAOMixes); size != 65536 {
			err = ssz.ErrVectorLengthFn("BeaconState.RANDAOMixes", size, 65536)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.RANDAOMixes {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				errs <- err
				return
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp13 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.RANDAOMixesRoot = &tmp13
	}()

	// Field (14) 'Slashings'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.Slashings); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.Slashings", size, 8192)
			errs <- err
			return
		}
		subIndx := hh.Index()
		for _, i := range b.Slashings {
			hh.AppendUint64(uint64(i))
		}
		hh.Merkleize(subIndx)
		tmp14 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.SlashingsRoot = &tmp14
	}()

	// Field (15) 'PreviousEpochParticipation'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.PreviousEpochParticipation); size > 1099511627776 {
			errs <- ssz.ErrListTooBigFn("BeaconState.PreviousEpochParticipation", size, 1099511627776)
			return
		}
		subIndx := hh.Index()
		for _, i := range b.PreviousEpochParticipation {
			hh.AppendUint8(uint8(i))
		}
		hh.FillUpTo32()
		numItems := uint64(len(b.PreviousEpochParticipation))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 1))
		tmp15 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.PreviousEpochParticipationRoot = &tmp15
	}()

	// Field (16) 'CurrentEpochParticipation'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.CurrentEpochParticipation); size > 1099511627776 {
			errs <- ssz.ErrListTooBigFn("BeaconState.CurrentEpochParticipation", size, 1099511627776)
			return
		}
		subIndx := hh.Index()
		for _, i := range b.CurrentEpochParticipation {
			hh.AppendUint8(uint8(i))
		}
		hh.FillUpTo32()
		numItems := uint64(len(b.CurrentEpochParticipation))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 1))
		tmp16 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.CurrentEpochParticipationRoot = &tmp16
	}()

	// Field (17) 'JustificationBits'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.JustificationBits); size != 1 {
			errs <- ssz.ErrBytesLengthFn("BeaconState.JustificationBits", size, 1)
			return
		}
		hh.PutBytes(b.JustificationBits)
		tmp17 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.JustificationBitsRoot = &tmp17
	}()

	// Field (18) 'PreviousJustifiedCheckpoint'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.PreviousJustifiedCheckpoint == nil {
			b.PreviousJustifiedCheckpoint = new(phase0.Checkpoint)
		}
		if err = b.PreviousJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp18 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.PreviousJustifiedCheckpointRoot = &tmp18
	}()

	// Field (19) 'CurrentJustifiedCheckpoint'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.CurrentJustifiedCheckpoint == nil {
			b.CurrentJustifiedCheckpoint = new(phase0.Checkpoint)
		}
		if err = b.CurrentJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp19 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.CurrentJustifiedCheckpointRoot = &tmp19
	}()

	// Field (20) 'FinalizedCheckpoint'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.FinalizedCheckpoint == nil {
			b.FinalizedCheckpoint = new(phase0.Checkpoint)
		}
		if err = b.FinalizedCheckpoint.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp20 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.FinalizedCheckpointRoot = &tmp20
	}()

	// Field (21) 'InactivityScores'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if size := len(b.InactivityScores); size > 1099511627776 {
			errs <- ssz.ErrListTooBigFn("BeaconState.InactivityScores", size, 1099511627776)
			return
		}
		subIndx := hh.Index()
		for _, i := range b.InactivityScores {
			hh.AppendUint64(i)
		}
		hh.FillUpTo32()
		numItems := uint64(len(b.InactivityScores))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
		tmp21 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.InactivityScoresRoot = &tmp21
	}()

	// Field (22) 'CurrentSyncCommittee'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.CurrentSyncCommittee == nil {
			b.CurrentSyncCommittee = new(altair.SyncCommittee)
		}
		if err = b.CurrentSyncCommittee.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp22 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.CurrentSyncCommitteeRoot = &tmp22
	}()

	// Field (23) 'NextSyncCommittee'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if b.NextSyncCommittee == nil {
			b.NextSyncCommittee = new(altair.SyncCommittee)
		}
		if err = b.NextSyncCommittee.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp23 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.NextSyncCommitteeRoot = &tmp23
	}()

	// Field (24) 'LatestExecutionPayloadHeader'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		if err = b.LatestExecutionPayloadHeader.HashTreeRootWith(hh); err != nil {
			errs <- err
			return
		}
		tmp24 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.LatestExecutionPayloadHeaderRoot = &tmp24
	}()

	// Field (25) 'NextWithdrawalIndex'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		hh.PutUint64(uint64(b.NextWithdrawalIndex))
		tmp25 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.NextWithdrawalIndexRoot = &tmp25
	}()

	// Field (26) 'NextWithdrawalValidatorIndex'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		hh.PutUint64(uint64(b.NextWithdrawalValidatorIndex))
		tmp26 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.NextWithdrawalValidatorIndexRoot = &tmp26
	}()

	// Field (27) 'HistoricalSummaries'
	go func() {
		defer wg.Done()

		hh := ssz.NewHasherWithHashFn(gohashtree.HashByteSlice)
		subIndx := hh.Index()
		num := uint64(len(b.HistoricalSummaries))
		if num > 16777216 {
			errs <- ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range b.HistoricalSummaries {
			if err = elem.HashTreeRootWith(hh); err != nil {
				errs <- err
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16777216)
		tmp27 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.HistoricalSummariesRoot = &tmp27
	}()

	wg.Wait()
	if len(errs) > 0 {
		// something failed.
		return nil, errors.Join(toSlice(errs)...)
	}

	return beaconStateTopLevelRoots, nil
}

func toSlice[T any](c chan T) []T {
	s := make([]T, 0)
	for i := range c {
		s = append(s, i)
	}
	return s
}
