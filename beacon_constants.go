package main

var VALIDATOR_REGISTRY_LIMIT = uint64(1099511627776)
var SLOTS_PER_HISTORICAL_ROOT = uint64(8192)

// Maximum number of withdrawals in an execution payload
var MAX_WITHDRAWALS_PER_PAYLOAD = uint64(16)

// **************Indexes of relevant containers**************
var BLOCK_ROOTS_INDEX = uint64(5)

// Index of the historical summaries in the beacon state
var HISTORICAL_SUMMARY_INDEX = uint64(27)

// Index of validator list in beacon state
var VALIDATOR_LIST_INDEX = uint64(11)
var BALANCE_LIST_INDEX = uint64(12)

var BEACON_STATE_SLOT_INDEX = uint64(2)

var LATEST_BLOCK_HEADER_INDEX = uint64(4)

// Index of the beacon body root inside the beacon body header
var BEACON_BLOCK_BODY_ROOT_INDEX = uint64(4)

// Index of the execution payload in the BeaconBlockBody container
var EXECUTION_PAYLOAD_INDEX = uint64(9)

// Index of the block number inside the execution payload
var BLOCK_NUMBER_INDEX = uint64(6)

// Index of the timestamp inside the execution payload
var TIMESTAMP_INDEX = uint64(9)

// Index of the withdrawals inside the execution payload
var WITHDRAWALS_INDEX = uint64(14)

// Index of the slot in the beacon block header
var SLOT_INDEX = uint64(0)
var STATE_ROOT_INDEX = uint64(3)

var VALIDATORS_INDEX = uint64(11)

var WITHDRAWAL_CREDENTIALS_INDEX = uint64(1)

// in the historical summary coontainer, the block root summary is at index 0
var BLOCK_SUMMARY_ROOT_INDEX = uint64(0)

//
//
// **************Number of Layers in Various Subtrees**************
//
//

// Number of layers for various merkle subtrees
var BLOCK_HEADER_MERKLE_SUBTREE_NUM_LAYERS = uint64(3)

var BLOCK_BODY_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkelization of the Execution Payload
var EXECUTION_PAYLOAD_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkleization of the Validator container
var VALIDATOR_MERKLE_SUBTREE_NUM_LAYERS = uint64(3)

// Number of layers for the merkleixation of the Validator List in the Beacon State
var VALIDATOR_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(40)

// Number of layers for the merkleixation of the Historical Summary List in the Beacon State
var HISTORICAL_SUMMARY_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(24)

// Number of layers for the merkleixation of the Balance List in the Beacon State
var BALANCE_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(40)

// Number of layers for the merkleization of the Withdrawal List in the Exection Payload
var WITHDRAWAL_LIST_MERKLE_SUBTREE_NUM_LAYERS = uint64(4)

// Number of layers for the merkleization of the Beacon State
var BEACON_STATE_MERKLE_SUBTREE_NUM_LAYERS = uint64(5)

// Number of layers for the merkleization of the Block Roots in the Beacon State
var BLOCK_ROOTS_MERKLE_SUBTREE_NUM_LAYERS = uint64(13)

// **************Number of fields of various containers**************
var BEACON_STATE_NUM_FIELDS = uint64(28)
var BEACON_BLOCK_HEADER_NUM_FIELDS = uint64(5)
var BEACON_BLOCK_BODY_NUM_FIELDS = uint64(11)
var EXECUTION_PAYLOAD_NUM_FIELDS = uint64(15)
