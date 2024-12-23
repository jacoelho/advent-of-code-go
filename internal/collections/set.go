package collections

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elements ...T) Set[T] {
	s := make(Set[T], len(elements))
	s.Add(elements...)
	return s
}

func (s Set[T]) Add(e ...T) {
	for _, e := range e {
		s[e] = struct{}{}
	}
}

func (s Set[T]) Remove(e T) {
	delete(s, e)
}

func (s Set[T]) Contains(e T) bool {
	_, ok := s[e]
	return ok
}

func (s Set[T]) Iter() iter.Seq[T] {
	return maps.Keys(s)
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

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

func (s Set[T]) Difference(other Set[T]) Set[T] {
	difference := NewSet[T]()
	for key := range s {
		if _, exists := other[key]; !exists {
			difference[key] = struct{}{}
		}
	}
	return difference
}

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

func (s Set[T]) Union(other Set[T]) Set[T] {
	union := NewSet[T]()
	for k := range s.Iter() {
		union.Add(k)
	}
	for k := range other.Iter() {
		union.Add(k)
	}
	return union
}
