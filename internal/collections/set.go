package collections

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elements ...T) Set[T] {
	s := make(Set[T], len(elements))
	for _, element := range elements {
		s.Add(element)
	}
	return s
}

func (s Set[T]) Add(e T) {
	s[e] = struct{}{}
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
