// Package collections provides generic container types: a set, a stack, a
// double-ended queue, and a binary heap.
package collections

import (
	"iter"
	"maps"
)

// Set is an unordered set of comparable values. A nil Set reads as empty but
// panics on Add; use NewSet or NewSetFromIter.
type Set[T comparable] map[T]struct{}

// NewSet returns a set containing elements.
func NewSet[T comparable](elements ...T) Set[T] {
	s := make(Set[T], len(elements))
	s.Add(elements...)
	return s
}

// NewSetFromIter returns a set containing every value yielded by iter.
func NewSetFromIter[T comparable](iter iter.Seq[T]) Set[T] {
	s := NewSet[T]()
	for e := range iter {
		s.Add(e)
	}
	return s
}

// Add inserts each of e into s.
func (s Set[T]) Add(e ...T) {
	for _, e := range e {
		s[e] = struct{}{}
	}
}

// Remove deletes e from s if present.
func (s Set[T]) Remove(e T) {
	delete(s, e)
}

// Contains reports whether e is in s.
func (s Set[T]) Contains(e T) bool {
	_, ok := s[e]
	return ok
}

// Iter returns an iterator over the elements of s in unspecified order.
func (s Set[T]) Iter() iter.Seq[T] {
	return maps.Keys(s)
}

// Len returns the number of elements in s.
func (s Set[T]) Len() int {
	return len(s)
}

// IsEmpty reports whether s has no elements.
func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Clone returns a shallow copy of s. Cloning a nil set returns nil.
func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// Intersect returns a new set of the elements in both s and other.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	if len(s) > len(other) {
		s, other = other, s
	}
	intersection := NewSet[T]()
	for key := range s {
		if _, exists := other[key]; exists {
			intersection[key] = struct{}{}
		}
	}
	return intersection
}

// Difference returns a new set of the elements in s that are not in other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	difference := NewSet[T]()
	for key := range s {
		if _, exists := other[key]; !exists {
			difference[key] = struct{}{}
		}
	}
	return difference
}

// SymmetricDifference returns a new set of the elements in exactly one of s
// and other.
func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	difference := NewSet[T]()
	for key := range s {
		if _, exists := other[key]; !exists {
			difference[key] = struct{}{}
		}
	}
	for key := range other {
		if _, exists := s[key]; !exists {
			difference[key] = struct{}{}
		}
	}

	return difference
}

// Union returns a new set of the elements in s, other, or both.
func (s Set[T]) Union(other Set[T]) Set[T] {
	union := make(Set[T], len(s)+len(other))
	maps.Copy(union, s)
	maps.Copy(union, other)
	return union
}

// Next returns an arbitrary element of s without removing it, or
// (zero, false) when s is empty.
func (s Set[T]) Next() (T, bool) {
	for e := range s {
		return e, true
	}
	var zero T
	return zero, false
}
