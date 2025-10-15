package xmaps

type Pair[V1, V2 any] struct {
	K V1
	V V2
}

func Any[M ~map[K]V, K comparable, V any](predicate func(k K, v V) bool, m M) bool {
	for k, v := range m {
		if predicate(k, v) {
			return true
		}
	}
	return false
}

func Find[M ~map[K]V, K comparable, V any](m M, predicate func(k K, v V) bool) (Pair[K, V], bool) {
	for k, v := range m {
		if predicate(k, v) {
			return Pair[K, V]{K: k, V: v}, true
		}
	}
	var empty Pair[K, V]
	return empty, false
}

func Filter[M ~map[K]V, K comparable, V any](predicate func(k K, v V) bool, m M) []Pair[K, V] {
	if len(m) == 0 {
		return nil
	}

	var result []Pair[K, V]
	for k, v := range m {
		if predicate(k, v) {
			result = append(result, Pair[K, V]{K: k, V: v})
		}
	}

	return result
}
