package graph

import (
	"iter"
	"maps"
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

// MaximalCliques runs the Bron-Kerbosch algorithm on a graph.
func MaximalCliques[T comparable](graph map[T]collections.Set[T]) iter.Seq[collections.Set[T]] {
	return func(yield func(collections.Set[T]) bool) {
		type state struct {
			currentClique collections.Set[T]
			potential     collections.Set[T]
			excluded      collections.Set[T]
		}

		stack := collections.NewStack[state](state{
			currentClique: collections.NewSet[T](),
			potential:     collections.NewSet[T](slices.Collect(maps.Keys(graph))...),
			excluded:      collections.NewSet[T](),
		})

		for !stack.IsEmpty() {
			current, _ := stack.Pop()

			// if potential and excluded are both empty, currentClique is maximal
			if current.potential.IsEmpty() && current.excluded.IsEmpty() {
				if !yield(current.currentClique.Clone()) {
					return
				}
				continue
			}

			// choose a pivot to reduce the number of iterations
			pivot, _ := xiter.Next(current.potential.Iter())
			nonNeighbors := current.potential.Difference(graph[pivot])

			for vertex := range nonNeighbors.Iter() {
				neighbors := graph[vertex]

				stack.Push(state{
					currentClique: current.currentClique.Union(collections.NewSet[T](vertex)),
					potential:     current.potential.Intersect(neighbors),
					excluded:      current.excluded.Intersect(neighbors),
				})

				// update potential and excluded sets for backtracking
				current.potential.Remove(vertex)
				current.excluded.Add(vertex)
			}
		}
	}
}
