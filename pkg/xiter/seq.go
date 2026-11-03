package xiter

import (
	"iter"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/xconstraints"
)

// Seq describes a lazy, synchronous traversal. Its zero value is empty.
// Adapters do not cache values or make a single-use source replayable. Repeated
// traversals also repeat callbacks; their state and the source determine reuse.
type Seq[T any] struct {
	seq iter.Seq[T]
}

// From wraps seq without consuming it. A nil seq is empty.
func From[T any](seq iter.Seq[T]) Seq[T] {
	return Seq[T]{seq: seq}
}

// FromSlice borrows the slice. Changes to its elements are visible to traversal.
func FromSlice[T any](values []T) Seq[T] {
	return From(slices.Values(values))
}

// All exposes the traversal to range loops and standard-library adapters.
func (s Seq[T]) All() iter.Seq[T] {
	if s.seq == nil {
		return func(func(T) bool) {}
	}
	return s.seq
}

// Map applies f lazily, preserving source order.
func (s Seq[T]) Map[U any](f func(T) U) Seq[U] {
	return From(func(yield func(U) bool) {
		for v := range s.All() {
			if !yield(f(v)) {
				return
			}
		}
	})
}

// Filter yields the values satisfying predicate, preserving source order.
func (s Seq[T]) Filter(predicate func(T) bool) Seq[T] {
	return From(func(yield func(T) bool) {
		for v := range s.All() {
			if predicate(v) && !yield(v) {
				return
			}
		}
	})
}

// Take yields at most n values. If n <= 0, it does not invoke the source.
func (s Seq[T]) Take(n int) Seq[T] {
	if n <= 0 {
		return Seq[T]{}
	}
	return From(func(yield func(T) bool) {
		remaining := n
		for v := range s.All() {
			remaining--
			if !yield(v) || remaining == 0 {
				return
			}
		}
	})
}

// Skip omits the first n values. If n <= 0, it returns s unchanged.
func (s Seq[T]) Skip(n int) Seq[T] {
	if n <= 0 {
		return s
	}
	return From(func(yield func(T) bool) {
		remaining := n
		for v := range s.All() {
			if remaining > 0 {
				remaining--
				continue
			}
			if !yield(v) {
				return
			}
		}
	})
}

// Reduce folds the source in order, returning initial for an empty source.
func (s Seq[T]) Reduce[A any](initial A, f func(A, T) A) A {
	for v := range s.All() {
		initial = f(initial, v)
	}
	return initial
}

// Collect returns a new slice, or nil for an empty source. Elements are copied
// shallowly; any storage referenced by an element keeps its original ownership.
func (s Seq[T]) Collect() []T {
	return slices.Collect(s.All())
}

// Count consumes the source and returns its number of values.
func (s Seq[T]) Count() int {
	var count int
	for range s.All() {
		count++
	}
	return count
}

// Any stops at the first value satisfying predicate. It is false when empty.
func (s Seq[T]) Any(predicate func(T) bool) bool {
	for v := range s.All() {
		if predicate(v) {
			return true
		}
	}
	return false
}

// Every stops at the first value failing predicate. It is true when empty.
func (s Seq[T]) Every(predicate func(T) bool) bool {
	for v := range s.All() {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// First consumes at most one value, returning (zero, false) when empty.
func (s Seq[T]) First() (T, bool) {
	for v := range s.All() {
		return v, true
	}
	var zero T
	return zero, false
}

// SumBy sums f(v) in source order, returning zero when empty.
func (s Seq[T]) SumBy[N xconstraints.Number](f func(T) N) N {
	return Sum(s.Map(f).All())
}
