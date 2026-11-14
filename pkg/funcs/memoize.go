// Package funcs provides function memoization.
package funcs

// Memoize returns a function that caches the result of f for each argument.
// The cache is unbounded and the returned function is not safe for concurrent
// use.
func Memoize[In comparable, Out any](f func(In) Out) func(In) Out {
	cache := make(map[In]Out)
	return func(in In) Out {
		if result, found := cache[in]; found {
			return result
		}
		result := f(in)
		cache[in] = result
		return result
	}
}

type key[In1, In2 comparable] struct {
	V1 In1
	V2 In2
}

// Memoize2 is Memoize for functions of two arguments.
func Memoize2[In1, In2 comparable, Out any](f func(In1, In2) Out) func(In1, In2) Out {
	cache := make(map[key[In1, In2]]Out)
	return func(v1 In1, v2 In2) Out {
		k := key[In1, In2]{V1: v1, V2: v2}
		if result, found := cache[k]; found {
			return result
		}
		result := f(v1, v2)
		cache[k] = result
		return result
	}
}
