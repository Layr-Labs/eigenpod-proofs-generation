package beacon

const (
	slotsPerHistoricalRoot = uint64(8192)

	// Index of the historical summaries in the beacon state
	historicalSummaryListIndex = uint64(27)

	// Index of validator list in beacon state
	validatorListIndex = uint64(11)
	balanceListIndex   = uint64(12)

	// Index of the beacon body root inside the beacon body header
	beaconBlockBodyRootIndex = uint64(4)

	// Index of the execution payload in the BeaconBlockBody container
	executionPayloadIndex = uint64(9)

	// Index of the timestamp inside the execution payload
	timestampIndex = uint64(9)

	// Index of the withdrawals inside the execution payload
	withdrawalsIndex = uint64(14)

	// Index of the slot in the beacon block header
	slotIndex      = uint64(0)
	stateRootIndex = uint64(3)

	// in the historical summary coontainer, the block root summary is at index 0
	blockSummaryRootIndex = uint64(0)
)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

const (
	// Number of layers for various merkle subtrees
	blockHeaderMerkleSubtreeNumLayers = uint64(3)

	blockBodyMerkleSubtreeNumLayers = uint64(4)

	// Number of layers for the merkelization of the Execution Payload
	executionPayloadMerkleSubtreeNumLayersDeneb = uint64(5)

	executionPayloadMerkleSubtreeNumLayersCapella = uint64(4)

	// Number of layers for the merkleixation of the Validator List in the Beacon State
	validatorListMerkleSubtreeNumLayers = uint64(40)

	// Number of layers for the merkleixation of the Historical Summary List in the Beacon State
	historicalSummaryListMerkleSubtreeNumLayers = uint64(24)

	// Number of layers for the merkleization of the Withdrawal List in the Exection Payload
	withdrawalListMerkleSubtreeNumLayers = uint64(4)

	// Number of layers for the merkleization of the Beacon State
	beaconStateMerkleSubtreeNumLayers = uint64(5)

	// Number of layers for the merkleization of the Block Roots in the Beacon State
	blockRootsMerkleSubtreeNumLayers = uint64(13)

	// **************Number of fields of various containers**************
	beaconBlockHeaderNumFields = uint64(5)
)
