package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

func BFSWithVisited[T comparable](start T, visited collections.Set[T], neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited.Add(start)
		frontier := collections.NewDeque[T](10)
		frontier.PushBack(start)

		for frontier.Size() > 0 {
			node, ok := frontier.PopFront()
			if !ok || !yield(node) {
				return
			}

			for el := range neighbours(node) {
				if visited.Contains(el) {
					continue
				}
				visited.Add(el)
				frontier.PushBack(el)
			}
		}
	}
}

func BFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return BFSWithVisited(start, collections.NewSet[T](), neighbours)
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

	distances := make(map[T]int)
	distances[start] = 0

	for node := range BFS(start, neighbours) {
		currentDist := distances[node]

		if node == target {
			return currentDist
		}

		for neighbour := range neighbours(node) {
			if _, seen := distances[neighbour]; !seen {
				distances[neighbour] = currentDist + 1
			}
		}
	}

	return -1
}

func BFSMaxDistance[T comparable](start T, neighbours func(T) iter.Seq[T]) int {
	distances := make(map[T]int)
	distances[start] = 0

	maxDist := 0
	for node := range BFS(start, neighbours) {
		currentDist := distances[node]
		maxDist = max(maxDist, currentDist)

		for neighbour := range neighbours(node) {
			if _, seen := distances[neighbour]; !seen {
				distances[neighbour] = currentDist + 1
			}
		}
	}

	return maxDist
}

func BFSDistances[T comparable](start T, neighbours func(T) iter.Seq[T]) map[T]int {
	distances := make(map[T]int)
	distances[start] = 0

	for node := range BFS(start, neighbours) {
		currentDist := distances[node]

		for neighbour := range neighbours(node) {
			if _, seen := distances[neighbour]; !seen {
				distances[neighbour] = currentDist + 1
			}
		}
	}

	return distances
}
