package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type wiringLine struct {
	node      string
	neighbors []string
}

type wiringDiagram struct {
	adjacent    [][]int
	destination []int
}

func parseWiringLine(line []byte) (wiringLine, error) {
	node, neighbors, ok := strings.Cut(string(line), ":")
	validName := func(name string) bool {
		return len(name) > 0 && len(name) <= 128 && !strings.ContainsRune(name, ':') && !strings.ContainsFunc(name, unicode.IsSpace)
	}
	if !ok || !validName(node) {
		return wiringLine{}, fmt.Errorf("invalid wiring declaration %q", line)
	}

	result := wiringLine{node: node, neighbors: strings.Fields(neighbors)}
	for _, neighbor := range result.neighbors {
		if !validName(neighbor) || neighbor == node {
			return wiringLine{}, fmt.Errorf("invalid wire %q to %q", node, neighbor)
		}
	}
	return result, nil
}

func parseWiringDiagram(r io.Reader) (wiringDiagram, error) {
	const inputLimit = 4 << 20
	limited := &io.LimitedReader{R: r, N: inputLimit + 1}
	nodes := make(map[string]int)
	wires := make(map[[2]string]struct{})
	s := scanner.NewScanner(limited, parseWiringLine)
	for line := range s.Values() {
		nodes[line.node] = 0
		for _, neighbor := range line.neighbors {
			nodes[neighbor] = 0
			wire := [2]string{line.node, neighbor}
			if wire[0] > wire[1] {
				wire[0], wire[1] = wire[1], wire[0]
			}
			wires[wire] = struct{}{}
			if len(nodes) > 4096 || len(wires) > 20_000 {
				return wiringDiagram{}, fmt.Errorf("wiring exceeds 4096 vertices or 20000 wires")
			}
		}
		if len(nodes) > 4096 {
			return wiringDiagram{}, fmt.Errorf("wiring exceeds 4096 vertices")
		}
	}
	if err := s.Err(); err != nil {
		return wiringDiagram{}, err
	}
	if limited.N == 0 {
		return wiringDiagram{}, fmt.Errorf("wiring input exceeds 4 MiB")
	}
	if len(nodes) < 2 {
		return wiringDiagram{}, fmt.Errorf("wiring needs at least two vertices")
	}

	names := make([]string, 0, len(nodes))
	for name := range nodes {
		names = append(names, name)
	}
	slices.Sort(names)
	for i, name := range names {
		nodes[name] = i
	}
	edges := make([][2]int, 0, len(wires))
	for wire := range wires {
		edges = append(edges, [2]int{nodes[wire[0]], nodes[wire[1]]})
	}
	slices.SortFunc(edges, func(a, b [2]int) int {
		if order := cmp.Compare(a[0], b[0]); order != 0 {
			return order
		}
		return cmp.Compare(a[1], b[1])
	})

	graph := wiringDiagram{
		adjacent:    make([][]int, len(nodes)),
		destination: make([]int, 0, 2*len(edges)),
	}
	for _, edge := range edges {
		arc := len(graph.destination)
		graph.destination = append(graph.destination, edge[1], edge[0])
		graph.adjacent[edge[0]] = append(graph.adjacent[edge[0]], arc)
		graph.adjacent[edge[1]] = append(graph.adjacent[edge[1]], arc+1)
	}
	return graph, nil
}

func (graph wiringDiagram) threeEdgeCut() ([]bool, error) {
	residual := make([]uint8, len(graph.destination))
	previous := make([]int, len(graph.adjacent))
	queue := make([]int, 0, len(graph.adjacent))
	remaining := 200_000_000
	var cut []bool

	// Every global cut separates vertex zero from at least one target.
	for target := 1; target < len(graph.adjacent); target++ {
		for i := range residual {
			residual[i] = 1
		}
		for flow := range 4 {
			for i := range previous {
				previous[i] = -1
			}
			previous[0] = -2
			queue = append(queue[:0], 0)
			for head := 0; head < len(queue) && previous[target] == -1; head++ {
				for _, arc := range graph.adjacent[queue[head]] {
					if remaining == 0 {
						return nil, fmt.Errorf("wiring search exceeds 200 million arc visits")
					}
					remaining--
					next := graph.destination[arc]
					if residual[arc] == 0 || previous[next] != -1 {
						continue
					}
					previous[next] = arc
					queue = append(queue, next)
					if next == target {
						break
					}
				}
			}
			if previous[target] == -1 {
				if flow < 3 {
					return nil, fmt.Errorf("global wiring cut is smaller than three")
				}
				if cut == nil {
					cut = make([]bool, len(graph.adjacent))
					for _, node := range queue {
						cut[node] = true
					}
				}
				break
			}
			for node := target; node != 0; {
				arc := previous[node]
				residual[arc]--
				// A reverse residual can reach two: cancellation must remain possible.
				residual[arc^1]++
				node = graph.destination[arc^1]
			}
		}
	}
	if cut == nil {
		return nil, fmt.Errorf("global wiring cut exceeds three")
	}
	return cut, nil
}

func (graph wiringDiagram) partitionProduct(cut []bool) (int, error) {
	crossing := 0
	for arc := 0; arc < len(graph.destination); arc += 2 {
		if cut[graph.destination[arc]] != cut[graph.destination[arc+1]] {
			crossing++
		}
	}
	if crossing != 3 {
		return 0, fmt.Errorf("partition does not cross exactly three wires")
	}

	visited := make([]bool, len(graph.adjacent))
	queue := make([]int, 0, len(graph.adjacent))
	var sizes [2]int
	components := 0
	for start := range graph.adjacent {
		if visited[start] {
			continue
		}
		if components == len(sizes) {
			return 0, fmt.Errorf("partition has more than two components")
		}
		visited[start] = true
		queue = append(queue[:0], start)
		for head := 0; head < len(queue); head++ {
			at := queue[head]
			for _, arc := range graph.adjacent[at] {
				next := graph.destination[arc]
				if visited[next] || cut[at] != cut[next] {
					continue
				}
				visited[next] = true
				queue = append(queue, next)
			}
		}
		sizes[components] = len(queue)
		components++
	}
	if components != 2 {
		return 0, fmt.Errorf("partition needs two nonempty components")
	}
	return sizes[0] * sizes[1], nil
}

func day25p01(r io.Reader) (string, error) {
	graph, err := parseWiringDiagram(r)
	if err != nil {
		return "", err
	}
	cut, err := graph.threeEdgeCut()
	if err != nil {
		return "", err
	}
	product, err := graph.partitionProduct(cut)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(product), nil
}
