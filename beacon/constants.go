package beacon

const (
	SlotsPerHistoricalRoot = uint64(8192)

	// Index of the historical summaries in the beacon state
	HistoricalSummaryListIndex = uint64(27)

	// Index of validator list in beacon state
	ValidatorListIndex = uint64(11)
	balanceListIndex   = uint64(12)

	// Index of the beacon body root inside the beacon body header
	BeaconBlockBodyRootIndex = uint64(4)

	// Index of the execution payload in the BeaconBlockBody container
	ExecutionPayloadIndex = uint64(9)

	// Index of the timestamp inside the execution payload
	TimestampIndex = uint64(9)

	// Index of the withdrawals inside the execution payload
	WithdrawalsIndex = uint64(14)

	// Index of the slot in the beacon block header
	SlotIndex      = uint64(0)
	stateRootIndex = uint64(3)

	// in the historical summary coontainer, the block root summary is at index 0
	BlockSummaryRootIndex = uint64(0)
)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

const (
	// Number of layers for various merkle subtrees
	blockHeaderMerkleSubtreeNumLayers = uint64(3)

	BlockBodyMerkleSubtreeNumLayers = uint64(4)

	// TODO unused; remove
	// Number of layers for the merkelization of the Execution Payload
	executionPayloadMerkleSubtreeNumLayersDeneb = uint64(5)

	executionPayloadMerkleSubtreeNumLayersCapella = uint64(4)

	// Number of layers for the merkleixation of the Validator List in the Beacon State
	ValidatorListMerkleSubtreeNumLayers = uint64(40)

	// Number of layers for the merkleixation of the Historical Summary List in the Beacon State
	HistoricalSummaryListMerkleSubtreeNumLayers = uint64(24)

	// Number of layers for the merkleization of the Withdrawal List in the Execution Payload
	WithdrawalListMerkleSubtreeNumLayers = uint64(4)

	// Number of layers for the merkleization of the Beacon State
	beaconStateMerkleSubtreeNumLayers = uint64(5)

	// Number of layers for the merkleization of the Block Roots in the Beacon State
	BlockRootsMerkleSubtreeNumLayers = uint64(13)

	// **************Number of fields of various containers**************
	beaconBlockHeaderNumFields = uint64(5)
)
