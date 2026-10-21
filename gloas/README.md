# Gloas EigenPod proofs

Generates credential, checkpoint, and unsettled-checkpoint cancellation proofs
for `EigenPodGloas`. A credential proof for an eligible slashed validator can also
be used with `verifyStaleBalance`.

Requires Go 1.25 or newer. This separate module pins the Gloas SSZ dependencies;
the root module and its CLI remain unchanged. Gloas-only, mainnet preset
(including Hoodi), with SSZ beacon-state input.

## Usage

Provide a complete **post-block** beacon state in SSZ and its matching,
unwrapped `BeaconBlockHeader` JSON. Obtain SSZ from the beacon debug-state
endpoint using `Accept: application/octet-stream`.

```sh
cd gloas
CGO_ENABLED=0 go run ./cmd/gloasproof \
  -mode checkpoint -state state.ssz -header header.json \
  -indices 123,456 \
  -expected-root <saved-checkpoint-beaconBlockRoot> \
  -execution-parent <checkpointExecutionParentHash>
```

The generator authenticates the header against `-expected-root`, then the full
state against the header. Read the trusted root from the pod's saved checkpoint
or, for credentials, the EIP-4788 oracle at the proof timestamp. Checkpoint mode
also requires the execution parent saved when that checkpoint started.

- `-mode credentials`: arguments for `verifyWithdrawalCredentials`; omit
  `-execution-parent`.
- `-mode checkpoint`: balance proofs and the authenticated execution hash guard.
- `-mode cancel`: committed execution hash and proof for
  `cancelUnsettledCheckpoint`; omit indices and execution parent.

For an unsettled checkpoint, generate cancellation arguments, cancel on chain,
then start a checkpoint with a delivered parent payload and regenerate proofs.
The contract checks cancellation eligibility. For a checkpoint started before
the Gloas implementation upgrade, call `startCheckpoint` again to restart it
with a fresh root; historical proofs are unnecessary.

The CLI reads local files and prints JSON arguments; it does not submit
transactions. Memory use grows with validator count because it builds a tree
for the complete state. Run `CGO_ENABLED=0 go test ./...` to check the module.

Tests include an archived Lighthouse Gloas SSZ state, its client root and SHA-256
checksum in `testdata/`. Export synthetic and real-client Solidity proof vectors
with `GLOAS_FIXTURE_DIR=/absolute/output CGO_ENABLED=0 go test -run 'TestExportFixtures|TestRealClientProofs'`.
Synthetic fault vectors deliberately exercise states rejected by consensus.

A measured synthetic two-million-validator full-state run used 3,164,205,056
bytes peak RSS and took 8.414 s on Linux arm64 with 2 CPU / 4 GiB limits.
That is a sizing reference, not a live-mainnet measurement. Benchmark tree
construction with `CGO_ENABLED=0 go test -run '^$' -bench 'BenchmarkGenerator/2000000$' -benchtime=1x -benchmem`.
