package merklization

const slotsPerHistoricalRoot = uint64(8192)

// Index of the historical summaries in the beacon state
const historicalSummaryListIndex = uint64(27)

// Index of validator list in beacon state
const validatorListIndex = uint64(11)
const balanceListIndex = uint64(12)

// Index of the beacon body root inside the beacon body header
const beaconBlockBodyRootIndex = uint64(4)

// Index of the execution payload in the BeaconBlockBody container
const executionPayloadIndex = uint64(9)

// Index of the timestamp inside the execution payload
const timestampIndex = uint64(9)

// Index of the withdrawals inside the execution payload
const withdrawalsIndex = uint64(14)

// Index of the slot in the beacon block header
const slotIndex = uint64(0)
const stateRootIndex = uint64(3)

// in the historical summary coontainer, the block root summary is at index 0
const blockSummaryRootIndex = uint64(0)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

// Number of layers for various merkle subtrees
const blockHeaderMerkleSubtreeNumLayers = uint64(3)

const blockBodyMerkleSubtreeNumLayers = uint64(4)

// Number of layers for the merkelization of the Execution Payload
const executionPayloadMerkleSubtreeNumLayersDeneb = uint64(5)

const executionPayloadMerkleSubtreeNumLayersCapella = uint64(4)

// Number of layers for the merkleixation of the Validator List in the Beacon State
const validatorListMerkleSubtreeNumLayers = uint64(40)

// Number of layers for the merkleixation of the Historical Summary List in the Beacon State
const historicalSummaryListMerkleSubtreeNumLayers = uint64(24)

// Number of layers for the merkleization of the Withdrawal List in the Exection Payload
const withdrawalListMerkleSubtreeNumLayers = uint64(4)

// Number of layers for the merkleization of the Beacon State
const beaconStateMerkleSubtreeNumLayers = uint64(5)

// Number of layers for the merkleization of the Block Roots in the Beacon State
const blockRootsMerkleSubtreeNumLayers = uint64(13)

// **************Number of fields of various containers**************
const beaconBlockHeaderNumFields = uint64(5)
