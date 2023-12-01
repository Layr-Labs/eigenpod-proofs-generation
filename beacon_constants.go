package main

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
const TIMESTAMP_INDEX = uint64(9)

// Index of the withdrawals inside the execution payload
const WITHDRAWALS_INDEX = uint64(14)

// Index of the slot in the beacon block header
const SLOT_INDEX = uint64(0)
const STATE_ROOT_INDEX = uint64(3)

const VALIDATORS_INDEX = uint64(11)

const WITHDRAWAL_CREDENTIALS_INDEX = uint64(1)

// in the historical summary coontainer, the block root summary is at index 0
const BLOCK_SUMMARY_ROOT_INDEX = uint64(0)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

// Number of layers for various merkle subtrees
const BLOCK_HEADER_MERKLE_SUBTREE_NUM_LAYERS = uint64(3)

const BLOCK_BODY_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkelization of the Execution Payload
const EXECUTION_PAYLOAD_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkleization of the Validator container
const VALIDATOR_MERKLE_SUBTREE_NUM_LAYERS = uint64(3)

// Number of layers for the merkleixation of the Validator List in the Beacon State
const VALIDATOR_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(40)

// Number of layers for the merkleixation of the Historical Summary List in the Beacon State
const HISTORICAL_SUMMARY_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(24)

// Number of layers for the merkleixation of the Balance List in the Beacon State
const BALANCE_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(40)

// Number of layers for the merkleization of the Withdrawal List in the Exection Payload
const WITHDRAWAL_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkleization of the Beacon State
const BEACON_STATE_MERKLE_SUBTREE_NUM_LAYERS = uint64(5)

// Number of layers for the merkleization of the Block Roots in the Beacon State
const BLOCK_ROOTS_MERKLE_SUBTREE_NUM_LAYERS = uint64(13)

// **************Number of fields of various containers**************
const BEACON_STATE_NUM_FIELDS = uint64(28)
const BEACON_BLOCK_HEADER_NUM_FIELDS = uint64(5)
const BEACON_BLOCK_BODY_NUM_FIELDS = uint64(11)
const EXECUTION_PAYLOAD_NUM_FIELDS = uint64(15)
