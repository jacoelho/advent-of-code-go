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
