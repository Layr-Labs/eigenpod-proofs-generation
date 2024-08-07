# EigenPod CLI

# Quickstart

## Dependencies

- Golang >= 1.21
- `abigen` (`go install github.com/ethereum/go-ethereum/cmd/abigen@latest`)
- `jq` 
- `solc` (`npm install -g solc`)
- [foundry](https://book.getfoundry.sh/getting-started/installation) (`curl -L https://foundry.paradigm.xyz | bash`)


- URL of an execution ETH node
- URL of a beacon ETH node

## Building from source

```bash
# Build CLI binary
make

# Test invoking the binary
./cli
```

# Proof Generation

The CLI produces two kinds of proofs, each corresponding to a different action you can take with your eigenpod. The CLI takes an additional `--sender $EIGENPOD_OWNER_PK` argument; if supplied, the CLI will submit proofs and act onchain for you.

Note that this is testnet software -- we aim to be addressing any bugs communicated with the team in a timely manner. We appreciate your understanding :) 

## Credential Proofs

- Credential proofs are a way of proving that a validator belongs to a given EigenPod. Later on (via a checkpointProof) you'll then prove
the _balance_ of that Validator, to represent your full staked balance within EigenLayer. **The proofs that this CLI generates are the main glue that syncs state between the beacon chain and execution chain.**.

To generate and submit a credential proof,

`./cli credentials --beaconNode $NODE_BEACON --podAddress $EIGENPOD_ADDRESS --execNode $NODE_ETH --sender $EIGENPOD_OWNER_PK`

If this is your first time, the CLI will post a transaction onchain linking your validator and eigenpod.

NOTE: If you've already linked your Validator to your EigenPod, you will see: `You have no inactive validators to verify. Everything up-to-date.`.

Once this is done, running the status command should show an "active" validator:
`./cli status --beaconNode $NODE_BEACON --podAddress $EIGENPOD_ADDRESS --execNode $NODE_ETH`


## Checkpoint Proofs

- Checkpoint proofs are a means of snapshotting your beacon state onto the execution chain via a Zero-Knowledge Proof. 
- Once you've brought your balance "up to date" on the execution chain, you can trigger a withdrawal of staking rewards and principal 
from your eigenpod.
- Checkpoints are started via the `EigenPod.StartCheckpoint()` contract function. They conclude when all proofs are submitted via `EigenPod.VerifyCheckpointProofs(...)`.
- Checkpoint proofs should be completed quickly after submission. The proofs must be generated via the becaon state that occurred at the 
time your checkpoint started. 
    - **If you wait too long to submit your checkpoint proof, you may need to use a full archival beacon node to 
re-generate the proofs.** 
    - You will not be able to start another checkpoint until the current checkpoint completes.

Note that you can also use:
    - `checkpoint --output <proof.json>` to write your proofs to a file, and 
    - `checkpoint --proof <proof.json>` to read and submit proofs that were previously written to a file.

Proofs are submitted to networks in batches by default. You can adjust the batch size with `--batch <batchSize>`. Our recommended batch sizes should provide optimal gas utilization.

- Once a checkpoint is completed, verify with the status command:

`./cli status --beaconNode $NODE_BEACON --podAddress $EIGENPOD_ADDRESS --execNode $NODE_ETH`

Congrats! Your pod balance is up-to-date.

