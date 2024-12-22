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

func Apply[E any](start E, f func(E) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		if !yield(start) {
			return
		}
		for {
			start = f(start)
			if !yield(start) {
				break
			}
		}
	}
}

func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := -1
		seq(func(v T) bool {
			i++
			return yield(i, v)
		})
	}
}

func Skip[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if n > 0 {
				n--
				return true
			}
			return yield(v)
		})
	}
}

func Take[T any](n int, seq iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if !yield(v) {
				return false
			}
			n--
			return n > 0
		})
	}
}

func Window[T any](n int, seq iter.Seq[T]) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		window := make([]T, 0, n)

		seq(func(v T) bool {
			if len(window) < n-1 {
				window = append(window, v)
				return true
			}
			if len(window) < n {
				window = append(window, v)
				return yield(window)
			}

			copy(window, window[1:])
			window[len(window)-1] = v
			return yield(window)
		})
		if len(window) < n {
			yield(window)
		}
	}
}

func Next[T any](seq iter.Seq[T]) (T, bool) {
	next, stop := iter.Pull(seq)
	defer stop()
	return next()
}
