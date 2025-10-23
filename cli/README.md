# EigenPod CLI

# Quickstart

## Dependencies

- Golang >= 1.21
- URL of an execution ETH node
- URL of a beacon ETH node

## Building from source

>> `go build`
>> `./cli`

# Proof Generation

The CLI produces two kinds of proofs, each corresponding to a different action you can take with your eigenpod. The CLI takes an additional `--sender $EIGENPOD_OWNER_PK` argument; if supplied, the CLI will submit proofs and act onchain for you.

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

## Consolidation Requests

#### How Does Consolidation Work?

Consolidation allows you to combine multiple validator indices together, allowing your total balance to occupy fewer validator indices and saving tons of gas on checkpoint proofs. 

Consolidations have a _source_ and a _target_. The _source_ is the validator index that will be "consumed", and the _target_ is where the source's balance will go once the consolidation request has been processed on the beacon chain.

Consolidations are initiated through your EigenPod which forwards requests to the consolidation predeploy. In order to successfully consolidate a source to a target, there are two primary requirements to keep in mind:

1. The _target_ validator must have 0x02 withdrawal credentials (and must not be exiting/exited) in order for the beacon chain to successfully process the consolidation. This is NOT checked by either the CLI or your pod.
2. The consolidation predeploy requires each request to be sent with a "request fee", which fluctuates depending on whether more requests are being added than removed. This fee is only updated at the end of each block, so if you're sending a bunch of requests in a single transaction, the current consolidation fee applies to each of the individual requests.

For more technical details and a walkthrough of how to perform a consolidation, see the [MOOCOW HackMD](https://hackmd.io/uijo9RSnSMOmejK1aKH0vw#Technical-Details).

#### Required Flags

All consolidation requests require the following flags in addition to command-specific flags:

```
-p podAddress
-b beaconNodeRPC
-e execNodeRPC
--sender senderPrivateKey
```

#### Switch to 0x02 credentials

Pass in a list of validator indices. This will initiate switch requests to change each validator's withdrawal prefix from 0x01 to 0x02. In order to be a consolidation _target_, a validator must have 0x02 credentials.

```
./cli consolidate switch --validators 425303,123444,555333
```

#### Source to Target

Pass in a target index and a list of source indices. This will initiate consolidations from each source to the specified target.

```
./cli consolidate source-to-target --target 425303 --sources 123444,555333
```

## Withdrawal Requests

#### How Do Withdrawal Requests Work?

Withdrawal requests allow your pod to initiate partial/full withdrawals on behalf of its validators. There are two kinds of withdrawal requests:
1. Full exits. Fully exit a validator from the beacon chain, withdrawing its entire balance to your EigenPod. This works just like a standard beacon-chain-initiated full exit.
2. Partial withdrawals. Withdraw a portion of your validator's balance, _down to 32 ETH_. The beacon chain will still process withdrawals that would bring a validator under 32 ETH

Withdrawal requests are initiated through your EigenPod which forwards requests to the withdrawal request predeploy. For more technical details, see the [MOOCOW HackMD](https://hackmd.io/uijo9RSnSMOmejK1aKH0vw#Technical-Details).

#### Required Flags

All withdrawal requests require the following flags in addition to command-specific flags:

```
-p podAddress
-b beaconNodeRPC
-e execNodeRPC
--sender senderPrivateKey
```

#### Full Exit

*(0x01 AND 0x02 validators)*

Pass in a list of validator indices. This will initiate full exits from the beacon chain.

```
./cli request-withdrawal  full-exit --validators 425303,123444,555333
```

#### Partial Withdrawal

*(0x02 validators only)*

Pass in a list of validator indices and an equally-sized list of amounts in gwei. This will initiate partial withdrawals from the beacon chain. Note that this method will NOT allow `amountGwei == 0`, as that is a full exit.

```
./cli request-withdrawal partial --validators 425303,123444,555333 --amounts 1000000000,2000000000,3000000000
```
