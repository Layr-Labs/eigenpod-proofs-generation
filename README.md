> [!IMPORTANT] 
> This repository is a work in progress and is currently unaudited. Use it with caution as it may contain unfinished features or bugs.

# Supported Versions

At all times, refer to the `mainnet`, `testnet`, and `preprod` git tags to understand which versions of this codebase are compatible with the corresponding EigenLayer environment.


| Environment    |     Version   |
| -------------- | ------------- |
| Mainnet                 |       [v1.6.1](https://github.com/Layr-Labs/eigenpod-proofs-generation/releases/tag/v1.6.1) |
| Testnet(Hoodi & Holesky)|       [v1.6.1](https://github.com/Layr-Labs/eigenpod-proofs-generation/releases/tag/v1.6.1) |
| Preprod                 |       [v1.6.1](https://github.com/Layr-Labs/eigenpod-proofs-generation/releases/tag/v1.6.1) |

# Introduction

PEPE Changes how we prove balances to EigenLayer. For more information, check out some of the links below.

## Links

- [More about PEPE](https://hackmd.io/U36dE9lnQha3tbf7D0GtKw?view)
- [Contract Documentation](https://github.com/Layr-Labs/eigenlayer-contracts/blob/feat/partial-withdrawal-batching/docs/core/EigenPod.md)

# Usage

- If you want to produce and submit proofs onchain -- either immediately, or by writing to a file to submit later -- check out our [CLI](./cli/README.md). The CLI can produce both credential and checkpoint proofs, and submit them onchain if given a private key.

- If you want to produce proofs from within Golang, please use `cli/core:GenerateValidatorProof` or `cli/core:GenerateCheckpointProof` for our high-level APIs. These will handle downloading beacon state, interfacing with an eth node, and generating the relevant proofs. Lower level APIs are available in `prove_validator.go`.

## Questions

For any questions, feel free to;

- Open a Github Issue
- Ask in [Discord](https://discord.com/invite/eigenlayer)
