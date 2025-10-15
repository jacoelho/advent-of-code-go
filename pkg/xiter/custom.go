package xiter

import (
	"cmp"
	"iter"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/xconstraints"
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

func CountBy[V any](predicate func(V) bool, seq iter.Seq[V]) int {
	var total int
	for v := range seq {
		if predicate(v) {
			total++
		}
	}
	return total
}

func MaxBy[V any](comparator func(V, V) bool, seq iter.Seq[V]) V {
	currentMax, ok := Next(seq)
	if !ok {
		panic("invalid sequence")
	}
	for v := range seq {
		if comparator(currentMax, v) {
			currentMax = v
		}
	}
	return currentMax
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

// Nth returns the nth element from the sequence (0-indexed).
func Nth[T any](seq iter.Seq[T], n int) (T, bool) {
	var result T
	found := false
	seq(func(v T) bool {
		if n > 0 {
			n--
			return true
		}
		result = v
		found = true
		return false
	})
	return result, found
}

// Permutations generates all permutations of the input slice using Heap's algorithm.
func Permutations[T any](slice []T) iter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if len(slice) == 0 {
			return
		}

		a := slices.Clone(slice)
		if !yield(slices.Clone(a)) {
			return
		}

		c := make([]int, len(slice))
		for i := 0; i < len(slice); {
			if c[i] < i {
				if i&1 == 0 {
					a[0], a[i] = a[i], a[0]
				} else {
					a[c[i]], a[i] = a[i], a[c[i]]
				}
				if !yield(slices.Clone(a)) {
					return
				}
				c[i]++
				i = 0
			} else {
				c[i] = 0
				i++
			}
		}
	}
}

// DotProduct computes the dot product of two sequences (element-wise multiply and sum).
// If the sequences have different lengths, computation stops at the shorter sequence.
func DotProduct[V xconstraints.Number](x, y iter.Seq[V]) V {
	var sum V
	for z := range Zip(x, y) {
		if z.Ok1 && z.Ok2 {
			sum += z.V1 * z.V2
		}
	}
	return sum
}
