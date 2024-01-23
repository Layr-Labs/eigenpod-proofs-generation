package beacon

import (
	"fmt"

	"github.com/Layr-Labs/eigenpod-proofs-generation/common"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
)

var zeroBytes = make([]byte, 32)

func ProveExecutionPayloadAgainstBlockHeaderCapella(
	blockHeader *phase0.BeaconBlockHeader,
	withdrawalBeaconBlockBody *capella.BeaconBlockBody,
) ([][32]byte, [32]byte, error) {
	// prove block body root against block header
	beaconBlockBodyAgainstBeaconBlockHeaderProof, err := ProveBlockBodyAgainstBlockHeader(blockHeader)
	if err != nil {
		return nil, [32]byte{}, err
	}

	// proof execution payload against the block body
	executionPayloadAgainstBlockHeaderProof, executionPayloadRoot, err := ProveExecutionPayloadAgainstBlockBodyCapella(withdrawalBeaconBlockBody)
	if err != nil {
		return nil, [32]byte{}, err
	}

	fullExecutionPayloadProof := append(executionPayloadAgainstBlockHeaderProof, beaconBlockBodyAgainstBeaconBlockHeaderProof...)
	return fullExecutionPayloadProof, executionPayloadRoot, nil
}

// Refer to beaconblockbody.go in go-eth2-client
// https://github.com/attestantio/go-eth2-client/blob/654ac05b4c534d96562329f988655e49e5743ff5/spec/bellatrix/beaconblockbody_encoding.go
func ProveExecutionPayloadAgainstBlockBodyCapella(beaconBlockBody *capella.BeaconBlockBody) (common.Proof, [32]byte, error) {
	beaconBlockBodyContainerRoots := make([]phase0.Root, 11)
	var err error

	hh := ssz.NewHasher()
	//Field 0: RANDAOReveal
	hh.PutBytes(beaconBlockBody.RANDAOReveal[:])
	copy(beaconBlockBodyContainerRoots[0][:], hh.Hash())
	hh.Reset()
	//Field 1: ETH1Data
	{
		if err = beaconBlockBody.ETH1Data.HashTreeRootWith(hh); err != nil {
			return nil, [32]byte{}, err
		}
		copy(beaconBlockBodyContainerRoots[1][:], hh.Hash())
	}
	//Field 2: Graffiti
	{
		hh.PutBytes(beaconBlockBody.Graffiti[:])
		copy(beaconBlockBodyContainerRoots[2][:], hh.Hash())
		hh.Reset()
	}

	//Field 3: ProposerSlashings
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.ProposerSlashings))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.ProposerSlashings {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(beaconBlockBodyContainerRoots[3][:], hh.Hash())
		hh.Reset()
	}

	//Field 4: AttesterSlashings
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.AttesterSlashings))
		if num > 2 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.AttesterSlashings {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 2)
		copy(beaconBlockBodyContainerRoots[4][:], hh.Hash())
		hh.Reset()
	}

	//Field 5: Attestations
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.Attestations))
		if num > 128 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.Attestations {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 128)
		copy(beaconBlockBodyContainerRoots[5][:], hh.Hash())
		hh.Reset()
	}

	//Field 6: Deposits
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.Deposits))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.Deposits {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(beaconBlockBodyContainerRoots[6][:], hh.Hash())
		hh.Reset()
	}

	//Field 7: VoluntaryExits
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.VoluntaryExits))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.VoluntaryExits {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(beaconBlockBodyContainerRoots[7][:], hh.Hash())
		hh.Reset()
	}

	//Field 8: SyncAggregate
	{
		if err = beaconBlockBody.SyncAggregate.HashTreeRootWith(hh); err != nil {
			return nil, [32]byte{}, err
		}
		copy(beaconBlockBodyContainerRoots[8][:], hh.Hash())
	}

	//Field 9: ExecutionPayload
	{
		if err = beaconBlockBody.ExecutionPayload.HashTreeRootWith(hh); err != nil {
			return nil, [32]byte{}, err
		}
		copy(beaconBlockBodyContainerRoots[9][:], hh.Hash())
	}

	//Field 10: BLSToExecutionChanges
	{
		subIndx := hh.Index()
		num := uint64(len(beaconBlockBody.BLSToExecutionChanges))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			return nil, [32]byte{}, err
		}
		for _, elem := range beaconBlockBody.BLSToExecutionChanges {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, [32]byte{}, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(beaconBlockBodyContainerRoots[10][:], hh.Hash())
		hh.Reset()
	}

	proof, err := common.GetProof(beaconBlockBodyContainerRoots, ExecutionPayloadIndex, BlockBodyMerkleSubtreeNumLayers)

	return proof, beaconBlockBodyContainerRoots[ExecutionPayloadIndex], err
}

func GetExecutionPayloadFieldRootsCapella(executionPayloadFields *capella.ExecutionPayload) ([]phase0.Root, error) {
	executionPayloadFieldRoots := make([]phase0.Root, 15)
	var err error

	hh := ssz.NewHasher()

	//Field 0: ParentHash
	hh.PutBytes(executionPayloadFields.ParentHash[:])
	copy(executionPayloadFieldRoots[0][:], hh.Hash())
	hh.Reset()

	//Field 1: FeeRecipient
	hh.PutBytes(executionPayloadFields.FeeRecipient[:])
	copy(executionPayloadFieldRoots[1][:], hh.Hash())
	hh.Reset()

	//Field 2: StateRoot
	hh.PutBytes(executionPayloadFields.StateRoot[:])
	copy(executionPayloadFieldRoots[2][:], hh.Hash())
	hh.Reset()

	//Field 3: ReceiptRoot
	hh.PutBytes(executionPayloadFields.ReceiptsRoot[:])
	copy(executionPayloadFieldRoots[3][:], hh.Hash())
	hh.Reset()

	//Field 4: LogsBloom
	hh.PutBytes(executionPayloadFields.LogsBloom[:])
	copy(executionPayloadFieldRoots[4][:], hh.Hash())
	hh.Reset()

	//Field 5: PrevRandao
	hh.PutBytes(executionPayloadFields.PrevRandao[:])
	copy(executionPayloadFieldRoots[5][:], hh.Hash())
	hh.Reset()

	//Field 6: BlockNumber
	hh.PutUint64(executionPayloadFields.BlockNumber)
	copy(executionPayloadFieldRoots[6][:], hh.Hash())
	hh.Reset()

	//Field 7: GasLimit
	hh.PutUint64(executionPayloadFields.GasLimit)
	copy(executionPayloadFieldRoots[7][:], hh.Hash())
	hh.Reset()

	//Field 8: GasUsed
	hh.PutUint64(executionPayloadFields.GasUsed)
	copy(executionPayloadFieldRoots[8][:], hh.Hash())
	hh.Reset()

	//Field 9: Timestamp
	hh.PutUint64(executionPayloadFields.Timestamp)
	copy(executionPayloadFieldRoots[9][:], hh.Hash())
	hh.Reset()

	//Field 10: ExtraData

	// //If the field is empty, we set it to 0
	// if len(executionPayloadFields.ExtraData) == 0 {
	// 	executionPayloadFields.ExtraData = []byte{0}
	// }

	{
		elemIndx := hh.Index()
		byteLen := uint64(len(executionPayloadFields.ExtraData))
		if byteLen > 32 {
			err = ssz.ErrIncorrectListSize
			fmt.Println(err)
		}
		hh.PutBytes(executionPayloadFields.ExtraData)
		hh.MerkleizeWithMixin(elemIndx, byteLen, (32+31)/32)
		copy(executionPayloadFieldRoots[10][:], hh.Hash())
		hh.Reset()
	}

	//Field 11: BaseFeePerGas
	hh.PutBytes(executionPayloadFields.BaseFeePerGas[:])
	copy(executionPayloadFieldRoots[11][:], hh.Hash())
	hh.Reset()

	//Field 12: BlockHash
	hh.PutBytes(executionPayloadFields.BlockHash[:])
	copy(executionPayloadFieldRoots[12][:], hh.Hash())
	hh.Reset()

	//Field 13: Transactions
	{
		subIndx := hh.Index()
		num := uint64(len(executionPayloadFields.Transactions))
		if num > 1048576 {
			err = ssz.ErrIncorrectListSize
			fmt.Println(err)
		}
		for _, elem := range executionPayloadFields.Transactions {
			{
				elemIndx := hh.Index()
				byteLen := uint64(len(elem))
				if byteLen > 1073741824 {
					err = ssz.ErrIncorrectListSize
					fmt.Println(err)
				}
				hh.AppendBytes32(elem)
				hh.MerkleizeWithMixin(elemIndx, byteLen, (1073741824+31)/32)
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 1048576)
		copy(executionPayloadFieldRoots[13][:], hh.Hash())
		hh.Reset()
	}

	//Field 14: Withdrawals
	{
		subIndx := hh.Index()
		num := uint64(len(executionPayloadFields.Withdrawals))
		if num > 16 {
			err := ssz.ErrIncorrectListSize
			return nil, err
		}
		for _, elem := range executionPayloadFields.Withdrawals {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16)
		copy(executionPayloadFieldRoots[14][:], hh.Hash())
		hh.Reset()
	}

	return executionPayloadFieldRoots, nil
}

// taken from https://github.com/attestantio/go-eth2-client/blob/654ac05b4c534d96562329f988655e49e5743ff5/spec/capella/beaconstate_ssz.go#L639
func ComputeBeaconStateTopLevelRootsCapella(b *capella.BeaconState) (*BeaconStateTopLevelRoots, error) {

	var err error
	beaconStateTopLevelRoots := &BeaconStateTopLevelRoots{}

	hh := ssz.NewHasher()

	// Field (0) 'GenesisTime'
	hh.PutUint64(b.GenesisTime)
	tmp0 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.GenesisTimeRoot = &tmp0
	hh.Reset()

	// Field (1) 'GenesisValidatorsRoot'
	if size := len(b.GenesisValidatorsRoot); size != 32 {
		err = ssz.ErrBytesLengthFn("BeaconState.GenesisValidatorsRoot", size, 32)
		return nil, err
	}
	hh.PutBytes(b.GenesisValidatorsRoot[:])
	tmp1 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.GenesisValidatorsRoot = &tmp1
	hh.Reset()

	// Field (2) 'Slot'
	hh.PutUint64(uint64(b.Slot))
	tmp2 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.SlotRoot = &tmp2
	hh.Reset()

	// Field (3) 'Fork'
	if b.Fork == nil {
		b.Fork = new(phase0.Fork)
	}
	if err = b.Fork.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp3 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.ForkRoot = &tmp3
	// copy(beaconStateTopLevelRoots.ForkRoot[:], hh.Hash())
	hh.Reset()

	// Field (4) 'LatestBlockHeader'
	if b.LatestBlockHeader == nil {
		b.LatestBlockHeader = new(phase0.BeaconBlockHeader)
	}
	if err = b.LatestBlockHeader.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp4 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.LatestBlockHeaderRoot = &tmp4
	// copy(beaconStateTopLevelRoots.LatestBlockHeaderRoot[:], hh.Hash())
	hh.Reset()

	// Field (5) 'BlockRoots'
	{
		if size := len(b.BlockRoots); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.BlockRoots", size, 8192)
			return nil, err
		}
		subIndx := hh.Index()
		for _, i := range b.BlockRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return nil, err
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp5 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.BlockRootsRoot = &tmp5
		// copy(beaconStateTopLevelRoots.BlockRootsRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (6) 'StateRoots'
	{
		if size := len(b.StateRoots); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.StateRoots", size, 8192)
			return nil, err
		}
		subIndx := hh.Index()
		for _, i := range b.StateRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return nil, err
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp6 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.StateRootsRoot = &tmp6
		// copy(beaconStateTopLevelRoots.StateRootsRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (7) 'HistoricalRoots'
	{
		if size := len(b.HistoricalRoots); size > 16777216 {
			err = ssz.ErrListTooBigFn("BeaconState.HistoricalRoots", size, 16777216)
			return nil, err
		}
		subIndx := hh.Index()
		for _, i := range b.HistoricalRoots {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return nil, err
			}
			hh.Append(i[:])
		}
		numItems := uint64(len(b.HistoricalRoots))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(16777216, numItems, 32))
		tmp7 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.HistoricalRootsRoot = &tmp7
		// copy(beaconStateTopLevelRoots.HistoricalRootsRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (8) 'ETH1Data'
	if b.ETH1Data == nil {
		b.ETH1Data = new(phase0.ETH1Data)
	}
	if err = b.ETH1Data.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp8 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.ETH1DataRoot = &tmp8
	// copy(beaconStateTopLevelRoots.ETH1DataRoot[:], hh.Hash())
	hh.Reset()

	// Field (9) 'ETH1DataVotes'
	{
		subIndx := hh.Index()
		num := uint64(len(b.ETH1DataVotes))
		if num > 2048 {
			err = ssz.ErrIncorrectListSize
			return nil, err
		}
		for _, elem := range b.ETH1DataVotes {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 2048)
		tmp9 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ETH1DataVotesRoot = &tmp9
		// copy(beaconStateTopLevelRoots.ETH1DataVotesRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (10) 'ETH1DepositIndex'
	hh.PutUint64(b.ETH1DepositIndex)
	tmp10 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.ETH1DepositIndexRoot = &tmp10
	// copy(beaconStateTopLevelRoots.ETH1DepositIndexRoot[:], hh.Hash())
	hh.Reset()

	// Field (11) 'Validators'
	{
		subIndx := hh.Index()
		num := uint64(len(b.Validators))
		if num > 1099511627776 {
			err = ssz.ErrIncorrectListSize
			return nil, err
		}
		for _, elem := range b.Validators {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 1099511627776)
		tmp11 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.ValidatorsRoot = &tmp11
		// copy(beaconStateTopLevelRoots.ValidatorsRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (12) 'Balances'
	{
		if size := len(b.Balances); size > 1099511627776 {
			err = ssz.ErrListTooBigFn("BeaconState.Balances", size, 1099511627776)
			return nil, err
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
		// copy(beaconStateTopLevelRoots.BalancesRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (13) 'RANDAOMixes'
	{
		if size := len(b.RANDAOMixes); size != 65536 {
			err = ssz.ErrVectorLengthFn("BeaconState.RANDAOMixes", size, 65536)
			return nil, err
		}
		subIndx := hh.Index()
		for _, i := range b.RANDAOMixes {
			if len(i) != 32 {
				err = ssz.ErrBytesLength
				return nil, err
			}
			hh.Append(i[:])
		}
		hh.Merkleize(subIndx)
		tmp13 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.RANDAOMixesRoot = &tmp13
		// copy(beaconStateTopLevelRoots.RANDAOMixesRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (14) 'Slashings'
	{
		if size := len(b.Slashings); size != 8192 {
			err = ssz.ErrVectorLengthFn("BeaconState.Slashings", size, 8192)
			return nil, err
		}
		subIndx := hh.Index()
		for _, i := range b.Slashings {
			hh.AppendUint64(uint64(i))
		}
		hh.Merkleize(subIndx)
		tmp14 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.SlashingsRoot = &tmp14
		// copy(beaconStateTopLevelRoots.SlashingsRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (15) 'PreviousEpochParticipation'
	{
		if size := len(b.PreviousEpochParticipation); size > 1099511627776 {
			err = ssz.ErrListTooBigFn("BeaconState.PreviousEpochParticipation", size, 1099511627776)
			return nil, err
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
		// copy(beaconStateTopLevelRoots.PreviousEpochParticipationRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (16) 'CurrentEpochParticipation'
	{
		if size := len(b.CurrentEpochParticipation); size > 1099511627776 {
			err = ssz.ErrListTooBigFn("BeaconState.CurrentEpochParticipation", size, 1099511627776)
			return nil, err
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
		// copy(beaconStateTopLevelRoots.CurrentEpochParticipationRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (17) 'JustificationBits'
	if size := len(b.JustificationBits); size != 1 {
		err = ssz.ErrBytesLengthFn("BeaconState.JustificationBits", size, 1)
		return nil, err
	}
	hh.PutBytes(b.JustificationBits)
	tmp17 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.JustificationBitsRoot = &tmp17
	hh.Reset()

	// Field (18) 'PreviousJustifiedCheckpoint'
	if b.PreviousJustifiedCheckpoint == nil {
		b.PreviousJustifiedCheckpoint = new(phase0.Checkpoint)
	}
	if err = b.PreviousJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp18 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.PreviousJustifiedCheckpointRoot = &tmp18
	// copy(beaconStateTopLevelRoots.PreviousJustifiedCheckpointRoot[:], hh.Hash())
	hh.Reset()

	// Field (19) 'CurrentJustifiedCheckpoint'
	if b.CurrentJustifiedCheckpoint == nil {
		b.CurrentJustifiedCheckpoint = new(phase0.Checkpoint)
	}
	if err = b.CurrentJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp19 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.CurrentJustifiedCheckpointRoot = &tmp19
	// copy(beaconStateTopLevelRoots.CurrentJustifiedCheckpointRoot[:], hh.Hash())
	hh.Reset()

	// Field (20) 'FinalizedCheckpoint'
	if b.FinalizedCheckpoint == nil {
		b.FinalizedCheckpoint = new(phase0.Checkpoint)
	}
	if err = b.FinalizedCheckpoint.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp20 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.FinalizedCheckpointRoot = &tmp20
	// copy(beaconStateTopLevelRoots.FinalizedCheckpointRoot[:], hh.Hash())
	hh.Reset()

	// Field (21) 'InactivityScores'
	{
		if size := len(b.InactivityScores); size > 1099511627776 {
			err = ssz.ErrListTooBigFn("BeaconState.InactivityScores", size, 1099511627776)
			return nil, err
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
		// copy(beaconStateTopLevelRoots.InactivityScoresRoot[:], hh.Hash())
		hh.Reset()
	}

	// Field (22) 'CurrentSyncCommittee'
	if b.CurrentSyncCommittee == nil {
		b.CurrentSyncCommittee = new(altair.SyncCommittee)
	}
	if err = b.CurrentSyncCommittee.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp22 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.CurrentSyncCommitteeRoot = &tmp22
	// copy(beaconStateTopLevelRoots.CurrentSyncCommitteeRoot[:], hh.Hash())
	hh.Reset()

	// Field (23) 'NextSyncCommittee'
	if b.NextSyncCommittee == nil {
		b.NextSyncCommittee = new(altair.SyncCommittee)
	}
	if err = b.NextSyncCommittee.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp23 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.NextSyncCommitteeRoot = &tmp23
	// copy(beaconStateTopLevelRoots.NextSyncCommitteeRoot[:], hh.Hash())
	hh.Reset()

	// Field (24) 'LatestExecutionPayloadHeader'
	if err = b.LatestExecutionPayloadHeader.HashTreeRootWith(hh); err != nil {
		return nil, err
	}
	tmp24 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.LatestExecutionPayloadHeaderRoot = &tmp24
	// copy(beaconStateTopLevelRoots.LatestExecutionPayloadHeaderRoot[:], hh.Hash())
	hh.Reset()

	// Field (25) 'NextWithdrawalIndex'
	hh.PutUint64(uint64(b.NextWithdrawalIndex))
	tmp25 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.NextWithdrawalIndexRoot = &tmp25
	hh.Reset()

	// Field (26) 'NextWithdrawalValidatorIndex'
	hh.PutUint64(uint64(b.NextWithdrawalValidatorIndex))
	tmp26 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
	beaconStateTopLevelRoots.NextWithdrawalValidatorIndexRoot = &tmp26
	hh.Reset()

	// Field (27) 'HistoricalSummaries'
	{
		subIndx := hh.Index()
		num := uint64(len(b.HistoricalSummaries))
		if num > 16777216 {
			err = ssz.ErrIncorrectListSize
			return nil, err
		}
		for _, elem := range b.HistoricalSummaries {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return nil, err
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 16777216)
		tmp27 := phase0.Root(common.ConvertTo32ByteArray(hh.Hash()))
		beaconStateTopLevelRoots.HistoricalSummariesRoot = &tmp27
		hh.Reset()
	}

	return beaconStateTopLevelRoots, nil
}
