package search

import (
	"iter"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
)

func DFSWithVisited[T comparable](start T, visited collections.Set[T], neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		visited.Add(start)
		frontier := collections.NewStack[T](start)

		for frontier.Len() > 0 {
			node, ok := frontier.Pop()
			if !ok || !yield(node) {
				return
			}

			for el := range neighbours(node) {
				if visited.Contains(el) {
					continue
				}
				visited.Add(el)
				frontier.Push(el)
			}
		}
	}
}

func DFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return DFSWithVisited(start, collections.NewSet[T](), neighbours)
}
