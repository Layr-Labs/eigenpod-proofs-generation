package utils

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

func Filter[A any](coll []A, criteria func(i A) bool) []A {
	out := []A{}
	for _, item := range coll {
		if criteria(item) {
			out = append(out, item)
		}
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
