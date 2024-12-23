package graph

import (
	"maps"
	"slices"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

// bron-kerbosch finds all maximal cliques in a graph using the Bron-Kerbosch algorithm.
func bronKerbosch[T comparable](
	graph map[T]collections.Set[T],
	currentClique, potential, excluded collections.Set[T],
	cliques *[]collections.Set[T],
) {
	if potential.IsEmpty() && excluded.IsEmpty() {
		// currentClique is a maximal clique
		maximalClique := currentClique.Clone() // Copy to avoid mutating it later
		*cliques = append(*cliques, maximalClique)
		return
	}

	// Choose a pivot to reduce the number of recursive calls
	pivot, _ := xiter.Next(potential.Iter())
	nonNeighbors := potential.Difference(graph[pivot])

	for vertex := range nonNeighbors.Iter() {
		neighbors := graph[vertex]
		// Recurse with vertex added to currentClique
		bronKerbosch(
			graph,
			currentClique.Union(collections.NewSet[T](vertex)),
			potential.Intersect(neighbors),
			excluded.Intersect(neighbors),
			cliques,
		)
		// Backtrack: remove vertex from potential and add it to excluded
		potential.Remove(vertex)
		excluded.Add(vertex)
	}
}

// MaximalCliques runs the Bron-Kerbosch algorithm on a graph.
func MaximalCliques[T comparable](graph map[T]collections.Set[T]) []collections.Set[T] {
	// Initialize the sets: currentClique (empty), potential (all vertices), and excluded (empty)
	currentClique := collections.NewSet[T]()
	potential := collections.NewSet[T](slices.Collect(maps.Keys(graph))...)
	excluded := collections.NewSet[T]()

	var cliques []collections.Set[T]
	bronKerbosch(graph, currentClique, potential, excluded, &cliques)
	return cliques
}
