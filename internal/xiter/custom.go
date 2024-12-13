package xiter

import (
	"cmp"
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

func Frequencies[T comparable](seq iter.Seq[T]) map[T]int {
	count := make(map[T]int)
	for item := range seq {
		count[item]++
	}
	return count
}

func Max[E cmp.Ordered](seq iter.Seq[E]) E {
	next, stop := iter.Pull(seq)
	defer stop()
	m, ok := next()
	if !ok {
		panic("empty seq")
	}
	for v := range seq {
		m = max(m, v)
	}
	return m
}
