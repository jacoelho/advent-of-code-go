package xslices

import (
	"cmp"
	"iter"
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/xconstraints"
)

func Frequencies[Slice ~[]E, E comparable](s Slice) map[E]int {
	count := make(map[E]int)
	for _, item := range s {
		count[item]++
	}
	return count
}

func Max[Slice ~[]E, E cmp.Ordered](s Slice) E {
	if len(s) == 0 {
		panic("empty slice")
	}
	return Reduce(func(m E, v1 E) E {
		return max(m, v1)
	}, s[0], s[1:])
}

func MaxBy[Slice ~[]E, E any](comparator func(a, current E) bool, s Slice) E {
	if len(s) == 0 {
		panic("empty slice")
	}
	return Reduce(func(m E, v1 E) E {
		if comparator(m, v1) {
			return v1
		}
		return m
	}, s[0], s[1:])
}

func Sum[Slice ~[]E, E xconstraints.Number](s Slice) E {
	var total E
	for _, v := range s {
		total += v
	}
	return total
}

func Product[Slice ~[]E, E xconstraints.Number](s Slice) E {
	var total E = 1
	for _, v := range s {
		total *= v
	}
	return total
}

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

func Every[Slice ~[]E, E any](predicate func(E) bool, s Slice) bool {
	for _, v := range s {
		if !predicate(v) {
			return false
		}
	}
	return true
}

func Any[Slice ~[]E, E any](predicate func(E) bool, s Slice) bool {
	return slices.ContainsFunc(s, predicate)
}

func Map[In, Out any](f func(In) Out, in []In) []Out {
	result := make([]Out, len(in))
	for i, v := range in {
		result[i] = f(v)
	}
	return result
}

func Reduce[Sum any, Slice ~[]E, E any](f func(Sum, E) Sum, sum Sum, s Slice) Sum {
	for _, v := range s {
		sum = f(sum, v)
	}
	return sum
}

func Filter[Slice ~[]E, E any](predicate func(E) bool, s Slice) Slice {
	var result Slice
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func CountFunc[Slice ~[]E, E any](predicate func(E) bool, s Slice) int {
	var result int
	for _, v := range s {
		if predicate(v) {
			result++
		}
	}
	return result
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

// SubSlices generates all subslices from the beginning of the input slice up to the given length.
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
