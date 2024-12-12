package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
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
