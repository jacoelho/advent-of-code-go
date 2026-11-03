package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// DFSWithVisited marks nodes when pushed, using the caller's visited set. The
// start is always yielded, even if already visited. Stopping skips expansion of
// the last yielded node but leaves stacked nodes marked in visited.
func DFSWithVisited[T comparable](start T, visited collections.Set[T], neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited.Add(start)
		frontier := collections.NewStack[T](start)
		push := func(node T) bool {
			if !visited.Contains(node) {
				visited.Add(node)
				frontier.Push(node)
			}
			return true
		}

		for frontier.Len() > 0 {
			node, ok := frontier.Pop()
			if !ok || !yield(node) {
				return
			}

			neighbours(node)(push)
		}
	}
}

func DFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		DFSWithVisited(start, collections.NewSet[T](), neighbours)(yield)
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
