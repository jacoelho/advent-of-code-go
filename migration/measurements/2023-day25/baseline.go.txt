package aoc2023

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type edge struct {
	a, b string
}

type wiringLine struct {
	node      string
	neighbors []string
}

func (e edge) normalized() edge {
	if e.a > e.b {
		return edge{e.b, e.a}
	}
	return e
}

func parseWiringLine(line []byte) (wiringLine, error) {
	parts := strings.Split(string(line), ": ")
	if len(parts) != 2 {
		return wiringLine{}, fmt.Errorf("invalid line format: %s", line)
	}

	return wiringLine{
		node:      parts[0],
		neighbors: strings.Fields(parts[1]),
	}, nil
}

func parseWiringDiagram(r io.Reader) (map[string]collections.Set[string], error) {
	graph := make(map[string]collections.Set[string])

	addEdge := func(a, b string) {
		if _, exists := graph[a]; !exists {
			graph[a] = collections.NewSet[string]()
		}
		graph[a].Add(b)
	}

	s := scanner.NewScanner(r, parseWiringLine)
	for wl := range s.Values() {
		for _, neighbor := range wl.neighbors {
			addEdge(wl.node, neighbor)
			addEdge(neighbor, wl.node)
		}
	}

	return graph, s.Err()
}

func findEdgeUsageFrequency(graph map[string]collections.Set[string], samples int) map[edge]int {
	edgeCount := make(map[edge]int)
	nodes := make([]string, 0, len(graph))
	for node := range graph {
		nodes = append(nodes, node)
	}

	neighbours := func(node string) iter.Seq[string] {
		return graph[node].Iter()
	}

	for i := 0; i < samples && i < len(nodes); i++ {
		start := nodes[i]
		distances := search.BFSDistances(start, neighbours)

		parent := make(map[string]string)
		for node, dist := range distances {
			if dist == 0 {
				continue
			}

			for neighbor := range graph[node].Iter() {
				if neighborDist, exists := distances[neighbor]; exists && neighborDist == dist-1 {
					parent[node] = neighbor
					break
				}
			}
		}

		for end := range distances {
			if end == start {
				continue
			}

			current := end
			for {
				prev, exists := parent[current]
				if !exists {
					break
				}
				e := edge{prev, current}.normalized()
				edgeCount[e]++
				current = prev
			}
		}
	}

	return edgeCount
}

func removeEdges(graph map[string]collections.Set[string], edges []edge) map[string]collections.Set[string] {
	newGraph := make(map[string]collections.Set[string])

	for node, neighbors := range graph {
		newGraph[node] = neighbors.Clone()
	}

	for _, e := range edges {
		if set, exists := newGraph[e.a]; exists {
			set.Remove(e.b)
		}
		if set, exists := newGraph[e.b]; exists {
			set.Remove(e.a)
		}
	}

	return newGraph
}

func findComponentSizes(graph map[string]collections.Set[string]) []int {
	visited := collections.NewSet[string]()
	sizes := []int{}

	neighbours := func(node string) iter.Seq[string] {
		return graph[node].Iter()
	}

	for start := range graph {
		if visited.Contains(start) {
			continue
		}

		componentSize := xiter.Len(search.BFSWithVisited(start, visited, neighbours))
		sizes = append(sizes, componentSize)
	}

	return sizes
}

func day25p01(r io.Reader) (string, error) {
	graph, err := parseWiringDiagram(r)
	if err != nil {
		return "", err
	}

	edgeFreq := findEdgeUsageFrequency(graph, len(graph))
	allEdges := make([]edge, 0, len(edgeFreq))
	for e := range edgeFreq {
		allEdges = append(allEdges, e)
	}

	slices.SortFunc(allEdges, func(a, b edge) int {
		return edgeFreq[b] - edgeFreq[a]
	})

	topN := min(20, len(allEdges))
	topEdges := allEdges[:topN]

	for testEdges := range xslices.Combinations(topEdges, 3) {
		testGraph := removeEdges(graph, testEdges)

		sizes := findComponentSizes(testGraph)

		if len(sizes) == 2 {
			return strconv.Itoa(xslices.Product(sizes)), nil
		}
	}

	return "", fmt.Errorf("no valid partition found")
}
