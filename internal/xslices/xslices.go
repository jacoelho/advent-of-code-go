package xslices

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/xconstraints"
)

func Frequencies[Slice ~[]E, E comparable](s Slice) map[E]int {
	count := make(map[E]int)
	for _, item := range s {
		count[item]++
	}
	return count
}

func Sum[Slice ~[]E, E xconstraints.Number](s Slice) E {
	var total E
	for _, v := range s {
		total += v
	}
	return total
}

func Window[Slice ~[]E, E any](s Slice, n int) iter.Seq[Slice] {
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

func Every[Slice ~[]E, E any](s Slice, predicate func(E) bool) bool {
	for _, v := range s {
		if !predicate(v) {
			return false
		}
	}
	return true
}
