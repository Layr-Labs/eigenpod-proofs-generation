package utils

import "math/big"

// maximum number of proofs per txn for each of the following proof types:
const DEFAULT_BATCH_CREDENTIALS = 60
const DEFAULT_BATCH_CHECKPOINT = 80

// imagine if golang had a standard library
func Map[A any, B any](coll []A, mapper func(i A, index uint64) B) []B {
	out := make([]B, len(coll))
	for i, item := range coll {
		out[i] = mapper(item, uint64(i))
	}
	return out
}

type Addable interface {
	~int | ~float64 | ~int64 | ~float32 | ~uint64
}

// A generic Sum function that sums up all elements in a list.
func Sum[T Addable](list []T) T {
	var sum T
	for _, item := range list {
		sum += item
	}
	return sum
}

func BigSum(list []*big.Int) *big.Int {
	return Reduce(list, func(sum *big.Int, cur *big.Int) *big.Int {
		return sum.Add(sum, cur)
	}, big.NewInt(0))
}

func Filter[A any](coll []A, criteria func(i A) bool) []A {
	out := []A{}
	for _, item := range coll {
		if criteria(item) {
			out = append(out, item)
		}
	}
	return out
}

func FilterI[A any](coll []A, criteria func(i A, index uint64) bool) []A {
	out := []A{}
	i := uint64(0)
	for _, item := range coll {
		if criteria(item, i) {
			out = append(out, item)
		}
		i++
	}
	return out
}

func Reduce[A any, B any](coll []A, processor func(accum B, next A) B, initialState B) B {
	val := initialState
	for _, item := range coll {
		val = processor(val, item)
	}
	return val
}

func Flatten[A any](coll [][]A) []A {
	out := []A{}
	for _, arr := range coll {
		out = append(out, arr...)
	}
	return out
}

func ShortenHex(publicKey string) string {
	return publicKey[0:6] + ".." + publicKey[len(publicKey)-4:]
}
