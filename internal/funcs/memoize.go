package funcs

func Memoize[In comparable, Out any](f func(In) Out) func(In) Out {
	cache := make(map[In]Out)
	return func(v In) Out {
		if result, found := cache[v]; found {
			return result
		}
		result := f(v)
		cache[v] = result
		return result
	}
}

type key[In1, In2 comparable] struct {
	V1 In1
	V2 In2
}

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
