package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// DFS yields each node reachable from start once, expanding the most recently
// discovered node first.
func DFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited := collections.NewSet(start)
		frontier := collections.NewStack(start)
		push := func(node T) bool {
			if !visited.Contains(node) {
				visited.Add(node)
				frontier.Push(node)
			}
			return true
		}
		for frontier.Len() > 0 {
			node, _ := frontier.Pop()
			if !yield(node) {
				return
			}
			neighbours(node)(push)
		}
	}
}

// DFSTree visits every occurrence without deduplication, exploring the last child
// first. It yields a node before expanding it; stopping iteration skips expansion.
// An expansion failure yields (zero, err) once and ends traversal, discarding any
// children returned with the error. An empty child slice makes a node a leaf.
//
// Child values are copied into the frontier; their referenced data remains owned
// by the caller. Each iteration has a fresh frontier, but replay also requires a
// repeatable callback. The caller must ensure finiteness or enforce search limits.
func DFSTree[T any](start T, children func(T) ([]T, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		frontier := collections.NewStack(start)
		for frontier.Len() > 0 {
			node, _ := frontier.Pop()
			if !yield(node, nil) {
				return
			}
			next, err := children(node)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, child := range next {
				frontier.Push(child)
			}
		}
	}
}
