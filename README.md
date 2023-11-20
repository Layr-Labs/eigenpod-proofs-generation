# Introduction
This repository allows users to generate the proofs necessary to prove consensus layer state on EigenLayer, specifically the EigenPods system.  These proofs are used for verifying 1) that withdrawal credentials of a validator are pointed to an EigenPod, 2) that changes in balance of a validator due to slashing, etc can be propagated to EigenLayer's smart contracts 3) Prove a validator's withdrawal on the consensus layer, which then allows a staker to withdraw their validator's ETH from their EigenPod. Specifically, these proofs are passed as inputs to the verifyWithdrawalCredentials(), verifyBalanceUpdates() and verifyAndProcessWithdrawals() functions, see [here](https://github.com/Layr-Labs/eigenlayer-contracts/blob/master/src/contracts/interfaces/IEigenPod.sol) for the exact function interface definitions.

An important note is that this CLI is designed to be used with inputs that can be retrieved from a consensus layer client, [here](https://ethereum.github.io/beacon-APIs/) is the relevant API specification.


# GIT Large File Storage
Please install git LFS (large file storage) before using this REPO in order for the CI to function.  We use LFS to store our large state files that are needed to generate the proofs and verify their integrity via the CI.

To install LFS: 
```
# On macOS using Homebrew
brew install git-lfs

# On Debian-based Linux
sudo apt-get install git-lfs

# On RPM-based Linux
sudo yum install git-lfs
```

Then to install it on this repo, navigate to apg/ and run:
```
git lfs install
```

If you have the repo cloned already, run: 
```
git lfs pull
```
If this doesn't work, it may require you to re-clone the repo from scratch.  



# How to Generate the Proofs with the Proof Generation library
This package allows you to generate withdrawal credential proofs, withdrawal proofs and balance update proofs. To generate the proofs using this library, run the following commands:

## Build the Executable

```
$go build
```

### Generate Validator Withdrawal Credential Proof
Here is the command:
```
$ ./proof-gen -command ValidatorFieldsProof -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] -stateFile [STATE_FILE_PATH]-validatorIndex [VALIDATOR_INDEX] -outputFile [OUTPUT_FILE_PATH] -chainID [CHAIN_ID]
```
Here is an example of running this command with the sample state/block files in the `/data` folder
```
./proof-gen -command ValidatorFieldsProof -oracleBlockHeaderFile "data/goerli_block_header_6399998.json" -stateFile "data/goerli_slot_6399998.json" -validatorIndex 302913 -outputFile "withdrawal_credential_proof_302913.json" -chainID 5
```
### Generate Withdrawal Proof
Here is the command:
```
$ ./proofGeneration -command WithdrawalFieldsProof -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] -stateFile [STATE_FILE_PATH]-validatorIndex [VALIDATOR_INDEX] -outputFile [OUTPUT_FILE_PATH] -chainID [CHAIN_ID] -historicalSummariesIndex [HISTORICAL_SUMMARIES_INDEX] -blockHeaderIndex [BLOCK_HEADER_INDEX] -historicalSummaryStateFile [HISTORICAL_SUMMARY_STATE_FILE_PATH] -blockHeaderFile [BLOCK_HEADER_FILE_PATH] -blockBodyFile [BLOCK_BODY_FILE_PATH] -withdrawalIndex [WITHDRAWAL_INDEX]
```
Here is an example of running this command with the sample state/block files in the `/data` folder

```
./proofGeneration -command WithdrawalFieldsProof -oracleBlockHeaderFile "data/goerli_block_header_6399998.json" -stateFile "data/goerli_slot_6399998.json" -validatorIndex 200240 -outputFile "withdrawal_proof_302913.json" -chainID 5 -historicalSummariesIndex 146 -blockHeaderIndex 8092 -historicalSummaryStateFile "data/goerli_slot_6397852.json" -blockHeaderFile  "data/goerli_block_header_6397852.json" -blockBodyFile "data/goerli_block_6397852.json" -withdrawalIndex 0
```

### Generate a Balance Update Proof.  
```
$ ./proofGeneration "BalanceUpdateProof"  -oracleBlockHeaderFile [ORACLE_BLOCK_HEADER_FILE_PATH] -stateFile [STATE_FILE_PATH]-validatorIndex [VALIDATOR_INDEX] -outputFile [OUTPUT_FILE_PATH] -chainID [CHAIN_ID]
```
Here is an example of running this command with the sample state/block files in the `/data` folder:
```
./proof-gen -command BalanceUpdateProof -oracleBlockHeaderFile "data/goerli_block_header_6399998.json" -stateFile "data/goerli_slot_6399998.json" -validatorIndex 302913 -outputFile "withdrawal_credential_proof_302913.json" -chainID 5
```

# Proof Generation Input Glossary
- `oracleBlockHeaderFile` is the block header of the oracle block heade root being used to make the proof
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
Historical summary proofs allow us to prove withdrawals from more than 27 hours after they are included in the beacon chain. Currently we require the submission of a historical summary proof which forces any withdrawers to wait at the most 27 hours before the historical_summaries container is updated with the latest historical summary that includes the block root with the withdrawal in question.  Refer [here](https://github.com/ethereum/consensus-specs/blob/dev/specs/capella/beacon-chain.md#historicalsummary) for the beacon chain specs for historical summaries. 





