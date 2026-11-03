package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// BFSWithVisited marks nodes when queued, using the caller's visited set. The
// start is always yielded, even if already visited. Stopping skips expansion of
// the last yielded node but leaves queued nodes marked in visited.
func BFSWithVisited[T comparable](start T, visited collections.Set[T], neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited.Add(start)
		frontier := collections.NewDeque[T](10)
		frontier.PushBack(start)
		enqueue := func(node T) bool {
			if !visited.Contains(node) {
				visited.Add(node)
				frontier.PushBack(node)
			}
			return true
		}

		for frontier.Size() > 0 {
			node, ok := frontier.PopFront()
			if !ok || !yield(node) {
				return
			}

			neighbours(node)(enqueue)
		}
	}
}

func BFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		BFSWithVisited(start, collections.NewSet[T](), neighbours)(yield)
	}
}

// BFSTree visits every occurrence without deduplication, preserving child order
// within each depth. It yields a node before expanding it; stopping iteration
// skips expansion. An expansion failure yields (zero, err) once and ends traversal,
// discarding any children returned with the error. An empty slice makes a leaf.
//
// Child values are copied into the frontier; their referenced data remains owned
// by the caller. Each iteration has a fresh frontier, but replay also requires a
// repeatable callback. The caller must ensure finiteness or enforce search limits.
func BFSTree[T any](start T, children func(T) ([]T, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		frontier := collections.NewDeque[T](32)
		frontier.PushBack(start)
		for frontier.Size() > 0 {
			node, _ := frontier.PopFront()
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
				frontier.PushBack(child)
			}
		}
	}
}

func BFSDistanceTo[T comparable](start, target T, neighbours func(T) iter.Seq[T]) int {
	if start == target {
		return 0
	}
	distance := -1
	bfsDistances(start, neighbours, func(node T, depth int) bool {
		if node != target {
			return true
		}
		distance = depth
		return false
	})
	return distance
}

func BFSMaxDistance[T comparable](start T, neighbours func(T) iter.Seq[T]) int {
	maxDist := 0
	bfsDistances(start, neighbours, func(_ T, depth int) bool {
		maxDist = depth
		return true
	})
	return maxDist
}

func BFSDistances[T comparable](start T, neighbours func(T) iter.Seq[T]) map[T]int {
	return bfsDistances(start, neighbours, nil)
}

func bfsDistances[T comparable](start T, neighbours func(T) iter.Seq[T], visit func(T, int) bool) map[T]int {
	distances := map[T]int{start: 0}
	frontier := collections.NewDeque[T](10)
	frontier.PushBack(start)
	var depth int
	enqueue := func(neighbour T) bool {
		if _, seen := distances[neighbour]; !seen {
			distances[neighbour] = depth + 1
			frontier.PushBack(neighbour)
		}
		return true
	}
	for frontier.Size() > 0 {
		node, _ := frontier.PopFront()
		depth = distances[node]
		if visit != nil && !visit(node, depth) {
			break
		}
		neighbours(node)(enqueue)
	}
	return distances
}
