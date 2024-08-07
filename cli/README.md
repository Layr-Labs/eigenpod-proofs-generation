# EigenPod CLI

# Quickstart

## Dependencies

- Golang >= 1.21
- URL of an execution ETH node
- URL of a beacon ETH node

## Building from source

>> `go build`
>> `./cli`

## Help

Any of the below commands has a built in help dialogue which further explains how to structure the command properly. Simply add `--help` to the end of any command below for more information. For example:
```bash
./cli credentials --help
```


## Key Management and EigenPod Proof Submitter

EigenLayer Native Restaking requires submitting proofs to EigenLayer contracts to prove the amount of validator ETH is active and its withdrawal address is pointing to the EigenPod. Be sure to use the most secure key management solution available for your EigenPod generation key (aka the EigenPod "owner"), such as a hardware wallet or cold wallet solution.

For users who do not wish to include the EigenPod Owner Private Key in their proof generation commands, you may identify another wallet and delegate its privilege to submit proofs on its behalf using the assign_submitter command. This is a **one time process** to assign a submitter for proofs. At any point in the future the `sender` of the proof can be the assigned submitter.

We recommend using a **different key** for the Proof Submitter vs the EigenPod owner. The Proof Submitter is any other address that is approved to submit proofs on behalf of the EigenPod owner. This allows the EigenPod owner key to remain used less frequently and remain more secure.

Use the following command to assign a submitter for your EigenPod:
```bash
/cli assign-submitter --execNode $NODE_ETH --podAddress $EIGENPOD_ADDRESS --sender $EIGENPOD_OWNER_PK
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

