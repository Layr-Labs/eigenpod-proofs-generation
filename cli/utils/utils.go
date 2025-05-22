package utils

import (
	"math/big"

	lo "github.com/samber/lo"
)

// These values are based on their corresponding Solidity types vs Ethereum's max txn size
// (See https://github.com/ethereum/go-ethereum/blob/master/core/txpool/legacypool/legacypool.go#L47-L59)
//
// txMaxSize: 131,072 bytes
// txMaxWords: txMaxSize / 32 == 4096
//
// The calculations for each theoretical max are included below.
// The default sizes given are reduced 'just in case'.
//
// TODO: Probably the best solution here is to automatically create batches when proof generation reaches
// the max transaction size, but this is the quick and dirty solution.

// input: (
//
//	uint64 beaconTimestamp,
//	BeaconChainProofs.StateRootProof stateRootProof,
//	uint40[] validatorIndices,
//	bytes[] validatorFieldsProofs,
//	bytes32[][] validatorFields
//
// )
//   - encoding: [beaconTimestamp][stateRootProof_offset][validatorIndices_offset][validatorFieldsProofs_offset][validatorFields_offset]
//   - overhead: 5 + overheads[stateRootProof, validatorIndices, validatorFieldsProofs, validatorFields]
//
// input: StateRootProof (bytes32 beaconStateRoot, bytes proof) (proof length is 3 words)
//   - encoding: [beaconStateRoot][proof_offset]
//     [proof_len][proof_0][proof_1][proof_2]
//   - overhead: 6 words
//
// input: (uint40[] validatorIndices)
//   - encoding: [arr_len][idx_0][idx_1][idx_...][idx_len-1]
//   - overhead: 1 + num_validators
//
// input: (bytes[] validatorFieldsProofs) (each proof is 47 words)
//   - encoding: [arr_len][proof0_offset][proof1_offset][...][prooflen-1_offset]
//     [proof0_len]<47 word proof>
//   - overhead: 1 + num_validators + (num_validators * 48)
//     == 1 + (49 * num_validators)
//
// input: (bytes32[][] validatorFields) (each validatorFields is 8 words)
//   - encoding: [arr_len][fields0_offset][fields1_offset][...][fieldslen-1_offset]
//     [fields0_len]<8 word validator fields>
//   - overhead: 1 + num_validators + (num_validators * 9)
//     == 1 + (10 * num_validators)
//
// TOTAL OVERHEAD (words): 5 + 6 + 3 + (60 * num_validators)
//
//	== 14 + (60 * num_validators)
//
// MAX CREDENTIAL PROOFS: 68
const DEFAULT_BATCH_CREDENTIALS = 60

// input: (
//
// BeaconChainProofs.BalanceContainerProof balanceContainerProof,
// BeaconChainProofs.BalanceProof[] proofs
//
// )
//   - encoding: [containerProof_offset][proofsArr_offset]
//   - overhead: 2 + overheads[balanceContainerProof, proofs]
//
// input: BalanceContainerProof (bytes32 balanceContainerRoot, bytes proof) (proof is 8 words)
//   - encoding: [balanceContainerRoot][proof_offset]
//     [proof_len]<8 word proof>
//   - overhead: 11 words
//
// input: BalanceProof[] (bytes32 pubkeyHash, bytes32 balanceRoot, bytes proof)[] (proof is 41 words)
//   - encoding: [arr_len][proof0_offset][proof1_offset][...][prooflen-1_offset]
//     [pubkeyHash0][balanceRoot0][proof_offset][proof_len]<41 word proof>
//   - overhead: 1 + num_validators + (45 * num_validators)
//     == 1 + (46 * num_validators)
//
// TOTAL OVERHEAD (words): 13 + (46 * num_validators)
//
// MAX CHECKPOINT PROOFS: 88
const DEFAULT_BATCH_CHECKPOINT = 80

// input: ConsolidationRequest[]
//   - encoding: [data_offset][length][0_ptr][1_ptr][...][len-1_ptr]
//   - overhead: 2 + arr_len words + (arr_len * overhead(ConsolidationRequest))
//
// input: ConsolidationRequest (bytes src, bytes target) (src and target length is 48 bytes, aka 2 words)
//   - encoding: [src_offset][target_offset]
//     [src_len][src_0][src_1]
//     [target_len][target_0][target_1]
//   - overhead: 8 words
//
// TOTAL OVERHEAD (words): 2 + (requests.length * 9)
//
// MAX CONSOLIDATE REQUESTS: 454
const DEFAULT_BATCH_CONSOLIDATE = 400

// input: WithdrawalRequest[]
//   - encoding: [data_offset][length][0_ptr][1_ptr][...][len-1_ptr]
//   - overhead: 2 + arr_len words + (arr_len * overhead(WithdrawalRequest))
//
// input: WithdrawalRequest (bytes pubkey, uint64 amountGwei) (pubkey length is 48, aka 2 words)
//   - encoding: [pubkey_offset][amountGwei]
//     [pubkey_len][pubkey_0][pubkey_1]
//   - overhead: 5 words
//
// TOTAL OVERHEAD (words): 2 + (requests.length * 5)
//
// MAX WITHDRAWAL REQUESTS: 818
const DEFAULT_BATCH_WITHDRAWREQUEST = 700

func BigSum(list []*big.Int) *big.Int {
	return lo.Reduce(list, func(sum *big.Int, cur *big.Int, index int) *big.Int {
		return sum.Add(sum, cur)
	}, big.NewInt(0))
}

func FilterDuplicates[A comparable](coll []A) []A {
	values := map[A]bool{}
	return lo.Filter(coll, func(item A, index int) bool {
		isSet := values[item]
		values[item] = true
		return !isSet
	})
}

func ShortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}
