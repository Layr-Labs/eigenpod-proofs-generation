package beacon

// various constants used in the beacon package
const (
	SlotsPerHistoricalRoot = uint64(8192)

	HistoricalSummaryListIndex = uint64(27)

	ValidatorListIndex         = uint64(11)
	ValidatorBalancesListIndex = uint64(12)

	BeaconBlockBodyRootIndex = uint64(4)

	ExecutionPayloadIndex = uint64(9)

	// Index of the timestamp inside the execution payload
	TimestampIndex = uint64(9)

	// Index of the withdrawals inside the execution payload
	WithdrawalsIndex = uint64(14)

	// Index of the slot in the beacon block header
	SlotIndex      = uint64(0)
	StateRootIndex = uint64(3)

	// in the historical summary coontainer, the block root summary is at index 0
	BlockSummaryRootIndex = uint64(0)
)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

const (
	BlockHeaderMerkleSubtreeNumLayers = uint64(3)

	BlockBodyMerkleSubtreeNumLayers = uint64(4)

	ValidatorListMerkleSubtreeNumLayers = uint64(40)

	HistoricalSummaryListMerkleSubtreeNumLayers = uint64(24)

	WithdrawalListMerkleSubtreeNumLayers = uint64(4)

	beaconStateMerkleSubtreeNumLayers = uint64(5)

	BlockRootsMerkleSubtreeNumLayers = uint64(13)

	BeaconBlockHeaderNumFields = uint64(5)
)
