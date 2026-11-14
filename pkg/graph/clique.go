package graph

import (
	"iter"
	"maps"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// MaximalCliques yields the inclusion-maximal cliques of an undirected graph
// using the Bron-Kerbosch algorithm. Each yielded set is independent.
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
				if !yield(current.currentClique) {
					return
				}
				continue
			}

			pivot := choosePivot(graph, current.potential, current.excluded)
			nonNeighbors := current.potential.Difference(graph[pivot])

			for vertex := range nonNeighbors.Iter() {
				neighbors := graph[vertex]

				stack.Push(state{
					currentClique: extend(current.currentClique, vertex),
					potential:     current.potential.Intersect(neighbors),
					excluded:      current.excluded.Intersect(neighbors),
				})

				// Exclude this vertex from later branches to avoid yielding the
				// same clique twice.
				current.potential.Remove(vertex)
				current.excluded.Add(vertex)
			}
		}
	}
}

func extend[T comparable](clique collections.Set[T], vertex T) collections.Set[T] {
	result := make(collections.Set[T], len(clique)+1)
	maps.Copy(result, clique)
	result.Add(vertex)
	return result
}

// choosePivot returns the vertex of potential ∪ excluded with the most
// neighbours in potential, which minimises the branches Bron–Kerbosch explores
// (Tomita, Tanaka and Takahashi, 2006).
func choosePivot[T comparable](graph map[T]collections.Set[T], potential, excluded collections.Set[T]) T {
	var pivot T
	best := -1
	for _, candidates := range [2]collections.Set[T]{potential, excluded} {
		for vertex := range candidates {
			count := 0
			for neighbour := range graph[vertex] {
				if potential.Contains(neighbour) {
					count++
				}
			}
			if count > best {
				pivot, best = vertex, count
			}
		}
	}
	return pivot
}
