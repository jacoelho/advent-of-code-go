package xiter

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/xconstraints"
)

func Unique[T comparable](seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		found := make(map[T]struct{})
		for v := range seq {
			if _, ok := found[v]; ok {
				continue
			}
			found[v] = struct{}{}
			if !yield(v) {
				return
			}
		}
	}
}

func Sum[V xconstraints.Number](seq iter.Seq[V]) V {
	var total V
	for v := range seq {
		total += v
	}
	return total
}

func Length[V any](seq iter.Seq[V]) int {
	var total int
	for range seq {
		total++
	}
	return total
}
