package xiter

import (
	"cmp"
	"iter"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/xconstraints"
)

// Unique keeps the first occurrence of each value, retaining O(distinct values)
// state for each traversal.
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

// Sum adds values in source order, returning zero when empty.
func Sum[V xconstraints.Number](seq iter.Seq[V]) V {
	var total V
	for v := range seq {
		total += v
	}
	return total
}

// Product multiplies the values in source order, returning one when empty.
func Product[V xconstraints.Number](seq iter.Seq[V]) V {
	var total V = 1
	for v := range seq {
		total *= v
	}
	return total
}

// MaxBy returns the first greatest value according to less, or (zero, false)
// when empty. It consumes the source once.
func MaxBy[V any](less func(V, V) bool, seq iter.Seq[V]) (V, bool) {
	var currentMax V
	found := false
	for v := range seq {
		if !found || less(currentMax, v) {
			currentMax = v
			found = true
		}
	}
	return currentMax, found
}

// Frequencies consumes seq and returns occurrence counts using O(distinct values)
// storage. The result is non-nil, including for an empty sequence.
func Frequencies[T comparable](seq iter.Seq[T]) map[T]int {
	count := make(map[T]int)
	for item := range seq {
		count[item]++
	}
	return count
}

// Max returns the greatest value, or (zero, false) when empty. It consumes the
// source once and propagates NaN like the built-in max.
func Max[E cmp.Ordered](seq iter.Seq[E]) (E, bool) {
	var m E
	found := false
	for v := range seq {
		if !found {
			m, found = v, true
		} else {
			m = max(m, v)
		}
	}
	return m, found
}

// Apply yields start, f(start), f(f(start)), and so on until iteration stops.
// Each traversal starts again at start and invokes f again.
func Apply[E any](start E, f func(E) E) iter.Seq[E] {
	return func(yield func(E) bool) {
		for value := start; yield(value); value = f(value) {
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

// Nth returns the zero-based nth value, or (zero, false) if the source ends
// first. A negative index returns (zero, false) without invoking the source.
func Nth[T any](seq iter.Seq[T], n int) (T, bool) {
	var result T
	if n < 0 {
		return result, false
	}
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

// Permutations yields a fresh slice for each permutation using Heap's
// algorithm, without modifying slice. Empty input yields no values; equal input
// values can produce duplicate permutations.
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
// It may consume one unmatched value from x when y ends first.
func DotProduct[V xconstraints.Number](x, y iter.Seq[V]) V {
	next, stop := iter.Pull(y)
	defer stop()
	var sum V
	for a := range x {
		b, ok := next()
		if !ok {
			break
		}
		sum += a * b
	}
	return sum
}
