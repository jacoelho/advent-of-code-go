package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
)

func BFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited := make(map[T]struct{})
		frontier := collections.NewDeque[T](10)

		visited[start] = struct{}{}
		frontier.PushBack(start)

		for frontier.Size() > 0 {
			node, ok := frontier.PopFront()
			if !ok || !yield(node) {
				return
			}

			for el := range neighbours(node) {
				if _, found := visited[el]; found {
					continue
				}
				visited[el] = struct{}{}
				frontier.PushBack(el)
			}
		}
		return
	}
}
