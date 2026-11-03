package xslices

import (
	"iter"
	"slices"
)

// Window yields overlapping, borrowed views of s. A nonempty slice shorter than
// n yields one short window. An empty slice yields none. It panics if n < 1.
func Window[Slice ~[]E, E any](n int, s Slice) iter.Seq[Slice] {
	if n < 1 {
		panic("cannot be less than 1")
	}

	return func(yield func(Slice) bool) {
		if len(s) == 0 {
			return
		}
		if n > len(s) {
			if !yield(s) {
				return
			}
		}
		for i := 0; i <= len(s)-n; i++ {
			if !yield(s[i : i+n]) {
				return
			}
		}
	}
}

type Pair[E, T any] struct {
	V1 E
	V2 T
}

// Pairwise yields every pair of elements with distinct increasing indexes,
// rather than only adjacent pairs.
func Pairwise[Slice ~[]E, E any](s Slice) iter.Seq[Pair[E, E]] {
	return func(yield func(Pair[E, E]) bool) {
		for i := 0; i < len(s); i++ {
			for j := i + 1; j < len(s); j++ {
				if !yield(Pair[E, E]{s[i], s[j]}) {
					return
				}
			}
		}
	}
}

func LastIndexFunc[Slice ~[]E, E any](f func(E) bool, s Slice) int {
	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

func HasDuplicates[Slice ~[]E, E comparable](slice Slice) bool {
	seen := make(map[E]struct{}, len(slice))
	for _, v := range slice {
		if _, exists := seen[v]; exists {
			return true
		}
		seen[v] = struct{}{}
	}
	return false
}

func HasSuffixFunc[Slice ~[]E, E any](slice Slice, suffix Slice, compare func(E, E) bool) bool {
	if len(suffix) > len(slice) {
		return false
	}
	offset := len(slice) - len(suffix)
	for i := range suffix {
		if !compare(slice[offset+i], suffix[i]) {
			return false
		}
	}
	return true
}

func HasSuffix[Slice ~[]E, E comparable](slice Slice, suffix Slice) bool {
	return HasSuffixFunc(slice, suffix, func(a, b E) bool { return a == b })
}

// SubSlices returns the prefixes s[:1] through s[:min(n, len(s))]. The prefixes
// share s's backing array. It panics if n is negative.
func SubSlices[Slice ~[]E, E any](s Slice, n int) []Slice {
	if n > len(s) {
		n = len(s)
	}

	result := make([]Slice, n)
	for i := 1; i <= n; i++ {
		result[i-1] = s[:i]
	}

	return result
}

// Combinations yields a new slice for each k-element selection by increasing
// input index. It yields one empty slice for k == 0 and no values for k < 0 or
// k > len(s). Equal input values can produce duplicate combinations.
func Combinations[Slice ~[]E, E any](s Slice, k int) iter.Seq[Slice] {
	if k > len(s) || k < 0 {
		return func(yield func(Slice) bool) {}
	}
	if k == 0 {
		return func(yield func(Slice) bool) {
			yield(Slice{})
		}
	}

	return func(yield func(Slice) bool) {
		var generate func(int, Slice) bool
		generate = func(start int, current Slice) bool {
			if len(current) == k {
				return yield(slices.Clone(current))
			}

			for i := start; i <= len(s)-k+len(current); i++ {
				if !generate(i+1, append(current, s[i])) {
					return false
				}
			}
			return true
		}
		generate(0, make(Slice, 0, k))
	}
}
