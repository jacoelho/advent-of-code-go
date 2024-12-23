package graph

import (
	"maps"
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

func bronKerbosch[T comparable](graph map[T]collections.Set[T]) []collections.Set[T] {
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

	var cliques []collections.Set[T]

	for !stack.IsEmpty() {
		current, _ := stack.Pop()

		// if potential and excluded are both empty, currentClique is maximal
		if current.potential.IsEmpty() && current.excluded.IsEmpty() {
			cliques = append(cliques, current.currentClique.Clone())
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

	return cliques
}

// MaximalCliques runs the Bron-Kerbosch algorithm on a graph.
func MaximalCliques[T comparable](graph map[T]collections.Set[T]) []collections.Set[T] {
	return bronKerbosch(graph)
}
