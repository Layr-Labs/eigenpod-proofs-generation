WARNING: Please note that this repository is a work in progress and is currently unaudited. Use it with caution as it may contain unfinished features or bugs."
# Introduction
This repository allows users to generate the proofs necessary to prove consensus layer state on EigenLayer, specifically the EigenPods system.  These proofs are used for verifying 1) that withdrawal credentials of a validator are pointed to an EigenPod, 2) that changes in balance of a validator due to slashing, etc can be propagated to EigenLayer's smart contracts 3) Prove a validator's withdrawal on the consensus layer, which then allows a staker to withdraw their validator's ETH from their EigenPod. Specifically, these proofs are passed as inputs to the verifyWithdrawalCredentials(), verifyBalanceUpdates() and verifyAndProcessWithdrawals() functions, see [here](https://github.com/Layr-Labs/eigenlayer-contracts/blob/master/src/contracts/interfaces/IEigenPod.sol) for the exact function interface definitions.


## How to Retrieve Data

An important note is that this CLI is designed to be used with inputs that can be retrieved from a consensus layer full node, [here](https://ethereum.github.io/beacon-APIs/) is the relevant API specification.  These are the api ednpoints that are required to retrieve the 3 consensus layer object required to generate proofs with this CLI:

### Beacon State
[This](https://ethereum.github.io/beacon-APIs/#/Debug/getStateV2) is the entire consensus layer [state](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate) object at a given slot.  The following endpoint returns this object:
```
/beacon/eth/v2/debug/beacon/states/[SLOT_NUMBER]
```
### Beacon Block
[This](https://ethereum.github.io/beacon-APIs/#/Beacon/getBlockV2) is the [beacon block](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconstate) object.  The following endpoint returns this object:
```
beacon/eth/v2/beacon/blocks/[SLOT_NUMBER]
```
### Beacon Block Header
[This](https://ethereum.github.io/beacon-APIs/#/Beacon/getBlockHeader) is the [block header](https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#beaconblockheader) for a beacon block.  The following endpoint returns this object:
```
/beacon/eth/v1/beacon/headers/[SLOT_NUMBER]
```


# How to Generate the Proofs with the Proof Generation library
This package allows you to generate withdrawal credential proofs, withdrawal proofs and balance update proofs. Please note that in order to run the sample commands, you must unzip the state files included in the data repo.  To generate the proofs using this library, run the following commands:

## Build the Executable

```bash
$ cd generation
$ go build
$ cd ..
```

### Generate Validator Withdrawal Credential Proof
Here is the command:
```bash
$ ./generation/generation \
    -command ValidatorFieldsProof \
    -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] \
    -stateFile [STATE_FILE_PATH] \
    -validatorIndex [VALIDATOR_INDEX] \
    -outputFile [OUTPUT_FILE_PATH] \
    -chainID [CHAIN_ID]
```
Here is a breakdown of the inputs here:
- “Command” aka the type of proof being generated
- “oracleBlockHeaderFile” is the path to the oracle block header file, that we are proving all of this against
- “stateFile” is the consensus state from that slot, containing the validator information
- “validatorIndex” is the index of the validator being proven inside state.validators
- “outputFile” - setting this will write the proofs to a json file
- “chainID” this parameter allows certain constants to be set depending on whether the proof is being generated for a goerli or mainnet state.

Here is an example of running this command with the sample state/block files in the `/data` folder
```bash
./generation/generation \
  -command ValidatorFieldsProof \
  -oracleBlockHeaderFile "./data/deneb_goerli_block_header_7431952.json" \
  -stateFile "./data/deneb_goerli_slot_7431952.json" \
  -validatorIndex 302913 \
  -outputFile "withdrawal_credential_proof_302913.json" \
  -chainID 5
```
### Generate Withdrawal Proof
Here is the command:
```bash
$ ./generation/generation \
  -command WithdrawalFieldsProof \
  -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] \
  -stateFile [STATE_FILE_PATH] \
  -validatorIndex [VALIDATOR_INDEX] \
  -outputFile [OUTPUT_FILE_PATH] \
  -chainID [CHAIN_ID] \
  -historicalSummariesIndex [HISTORICAL_SUMMARIES_INDEX] \
  -blockHeaderIndex [BLOCK_HEADER_INDEX] \
  -historicalSummaryStateFile [HISTORICAL_SUMMARY_STATE_FILE_PATH] \
  -blockHeaderFile [BLOCK_HEADER_FILE_PATH] \
  -blockBodyFile [BLOCK_BODY_FILE_PATH] \
  -withdrawalIndex [WITHDRAWAL_INDEX]
```
Here is an example of running this command with the sample state/block files in the `/data` folder
```bash
./generation/generation \
  -command WithdrawalFieldsProof \
  -oracleBlockHeaderFile ./data/deneb_goerli_block_header_7431952.json \
  -stateFile ./data/deneb_goerli_slot_7431952.json \
  -validatorIndex 627559 \
  -outputFile full_withdrawal_proof_627559.json \
  -chainID 5 \
  -historicalSummariesIndex 271 \
  -blockHeaderIndex 8191 \
  -historicalSummaryStateFile ./data/deneb_goerli_slot_7421952.json \
  -blockHeaderFile  data/deneb_goerli_block_header_7421951.json \
  -blockBodyFile data/deneb_goerli_block_7421951.json \
  -withdrawalIndex 0
```



### Generate a Balance Update Proof.  
```bash
$ ./generation/generation \
  - command BalanceUpdateProof \
  -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] \
  -stateFile [STATE_FILE_PATH] \
  -validatorIndex [VALIDATOR_INDEX] \
  -outputFile [OUTPUT_FILE_PATH] \
  -chainID [CHAIN_ID]
```
Here is a breakdown of the inputs here:
- “Command” aka the type of proof being generated
- “oracleBlockHeaderFile” is the path to the oracle block header file, that we are proving all of this against
- “stateFile” is the consensus state from that slot, containing the validator information
- “validatorIndex” is the index of the validator being proven inside state.validators
- “outputFile” - setting this will write the proofs to a json file
- “chainID” this parameter allows certain constants to be set depending on whether the proof is being generated for a goerli or mainnet state.

Here is an example of running this command with the sample state/block files in the `/data` folder:
```bash
./generation/generation \
  -command BalanceUpdateProof \
  -oracleBlockHeaderFile "./data/deneb_goerli_block_header_7431952.json" \
  -stateFile "./data/deneb_goerli_slot_7431952.json" \
  -validatorIndex 302913 \
  -outputFile "withdrawal_credential_proof_302913.json" \
  -chainID 5
```

# Proof Generation Input Glossary
- `oracleBlockHeaderFile` is the block header of the oracle block header root being used to make the proof
- `stateFile` is the associated state file of the oracle block being used
- `validatorIndex` is the index of the validator being proven for in the consensus layer
- `outputFile` is the location where the generated proof will be written to
- `chainID` is the chainID (either goerli = 5 or mainnet = 1) being generated for.
- `historicalSummariesIndex` refer to *What Are Historical Summary Proofs?* secion.  This is the index of the historical summary we're gonna use to prove the withdrawal
- `historicalSummaryStateFile` state file corresponding to the `state_summary_root` stored in the historical summary we're gonna use.
- `blockHeaderIndex` index of the block header that contains the withdrawal being proven
- `blockHeaderFile` file containing the block header that contains the withdrawal being proven
- `blockBodyFile` file containing the block body that contains the withdrawal being proven
- `withdrawalIndex` index of the withdrawal being proven within the block (there are 16 withdrawals per block).



# What Are Historical Summary Proofs?
Every block contains 16 withdrawals and any given beacon state stores the last 8192 block roots.  Thus if a withdrawal was within the last 8192 blocks, we can prove any withdrawal against one of the block roots, and then prove that block root against the state root.  However, what happens when we need to prove something from further in the past than 8192 blocks? That is where historical summaries come in. 
	Every 8192 blocks, the state transition function takes a “snapshot” of the state.block_roots that are stored in the beacon state by taking the hash tree root of the state.block_roots, and adding that root to state.historical_summaries.  Then, state.block_roots is cleared and the next 8192 block_roots will be added to state.block_roots.  Thus, to prove an old withdrawal, we need to take the extra step of retrieving the state at the slot at which the snapshot that contains the root of the block when the withdrawal was included. Refer [here](https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#historicalsummary) for the beacon chain specs for historical summaries. 





