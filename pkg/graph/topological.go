package graph

import (
	"fmt"
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

// TopologicalOrder orders every node before its successors, choosing the earliest
// node in nodes whenever several are ready. Nodes must be distinct and contain
// every edge endpoint. Duplicate edges do not change the ordering.
//
// Successors is called once per node; inputs are not modified. A cycle or invalid
// endpoint returns an error and no partial order. Time is O(V log V + E), with
// O(V + E) storage.
func TopologicalOrder[T comparable](nodes []T, successors func(T) iter.Seq[T]) ([]T, error) {
	indices := make(map[T]int, len(nodes))
	for i, node := range nodes {
		if _, exists := indices[node]; exists {
			return nil, fmt.Errorf("duplicate graph node %v", node)
		}
		indices[node] = i
	}

	var edges []int
	offsets := make([]int, len(nodes)+1)
	incoming := make([]int, len(nodes))
	var source T
	var edgeErr error
	appendEdge := func(next T) bool {
		if edgeErr != nil {
			return false
		}
		j, exists := indices[next]
		if !exists {
			edgeErr = fmt.Errorf("edge from %v to unknown node %v", source, next)
			return false
		}
		edges = append(edges, j)
		incoming[j]++
		return true
	}
	for i, node := range nodes {
		source = node
		offsets[i] = len(edges)
		successors(node)(appendEdge)
		if edgeErr != nil {
			return nil, edgeErr
		}
	}
	offsets[len(nodes)] = len(edges)

	ready := collections.NewHeap(func(a, b int) bool { return a < b })
	for i, count := range incoming {
		if count == 0 {
			ready.Push(i)
		}
	}
	order := make([]T, 0, len(nodes))
	for i := range ready.PopSeq() {
		order = append(order, nodes[i])
		for _, j := range edges[offsets[i]:offsets[i+1]] {
			incoming[j]--
			if incoming[j] == 0 {
				ready.Push(j)
			}
		}
	}
	if len(order) != len(nodes) {
		return nil, fmt.Errorf("graph contains a dependency cycle")
	}
	return order, nil
}
