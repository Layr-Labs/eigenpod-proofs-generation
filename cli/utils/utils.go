package utils

import (
	"math/big"

	lo "github.com/samber/lo"
)

// maximum number of proofs per txn for each of the following proof types:
const DEFAULT_BATCH_CREDENTIALS = 60
const DEFAULT_BATCH_CHECKPOINT = 80

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
