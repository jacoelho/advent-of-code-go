package aoc2019

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type vault struct {
	tiles     grid.Bytes
	entrances []int
	allKeys   uint32
}

type vaultEdge struct {
	to       int
	distance int
}

type vaultNode struct {
	tile  byte
	edges []vaultEdge
}

// vaultGraph links entrances, keys and doors by corridor distance. Entrances
// and keys come first in nodes, so robot positions fit in a byte.
type vaultGraph struct {
	nodes     []vaultNode
	entrances []uint8
	allKeys   uint32
}

type vaultState struct {
	positions [4]uint8
	keys      uint32
}

type vaultSearch struct {
	graph vaultGraph
	memo  map[vaultState]int
}

func keyBit(key byte) uint32 {
	return 1 << (key - 'a')
}

func isKey(tile byte) bool  { return tile >= 'a' && tile <= 'z' }
func isDoor(tile byte) bool { return tile >= 'A' && tile <= 'Z' }

func parseVault(r io.Reader) (vault, error) {
	g, err := grid.ReadBytes(r, "#.@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if err != nil {
		return vault{}, err
	}
	result := vault{tiles: g}
	for at, tile := range g.Cells() {
		switch {
		case tile == '@':
			result.entrances = append(result.entrances, at)
			if len(result.entrances) > 4 {
				return vault{}, fmt.Errorf("vault has more than four entrances")
			}
		case isKey(tile):
			if result.allKeys&keyBit(tile) != 0 {
				return vault{}, fmt.Errorf("duplicate vault key %q", tile)
			}
			result.allKeys |= keyBit(tile)
		}
	}
	return result, nil
}

func (v *vault) splitEntrance() error {
	center := v.entrances[0]
	width, cells := v.tiles.Width(), v.tiles.Cells()
	x, y := center%width, center/width
	if x == 0 || x+1 == width || y == 0 || y+1 == v.tiles.Height() {
		return fmt.Errorf("vault entrance needs a surrounding 3x3 floor")
	}
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if (dx != 0 || dy != 0) && cells[center+dy*width+dx] != '.' {
				return fmt.Errorf("vault entrance needs a surrounding 3x3 floor")
			}
		}
	}
	v.entrances = nil
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			at := center + dy*width + dx
			cells[at] = '#'
			if dx != 0 && dy != 0 {
				cells[at] = '@'
				v.entrances = append(v.entrances, at)
			}
		}
	}
	return nil
}

func buildVaultGraph(v vault) vaultGraph {
	var graph vaultGraph
	cells := v.tiles.Cells()
	indices := make([]int, len(cells))
	var points []int
	addPoints := func(include func(byte) bool) {
		for at, tile := range cells {
			if include(tile) {
				indices[at] = len(points)
				points = append(points, at)
				graph.nodes = append(graph.nodes, vaultNode{tile: tile})
			}
		}
	}
	for i := range indices {
		indices[i] = -1
	}
	addPoints(func(tile byte) bool { return tile == '@' || isKey(tile) })
	addPoints(isDoor)
	for i, node := range graph.nodes {
		if node.tile == '@' {
			graph.entrances = append(graph.entrances, uint8(i))
		}
	}
	graph.allKeys = v.allKeys

	distance := make([]int, len(cells))
	queue := make([]int, 0, len(cells))
	for source, start := range points {
		for i := range distance {
			distance[i] = -1
		}
		distance[start] = 0
		queue = append(queue[:0], start)
		for head := 0; head < len(queue); head++ {
			at := queue[head]
			for neighbour := range v.tiles.Neighbours4(v.tiles.Position(at)) {
				next := v.tiles.Index(neighbour)
				if cells[next] == '#' || distance[next] >= 0 {
					continue
				}
				distance[next] = distance[at] + 1
				if target := indices[next]; target >= 0 {
					// Corridor interiors contain no keys or doors, so their shortest route is always valid.
					graph.nodes[source].edges = append(graph.nodes[source].edges, vaultEdge{to: target, distance: distance[next]})
					continue
				}
				queue = append(queue, next)
			}
		}
	}
	return graph
}

func (s *vaultSearch) reachableKeys(start uint8, keys uint32) []vaultEdge {
	distance := make([]int, len(s.graph.nodes))
	for i := range distance {
		distance[i] = math.MaxInt
	}
	distance[start] = 0
	queue := collections.NewHeap(func(a, b vaultEdge) bool { return a.distance < b.distance })
	queue.Push(vaultEdge{to: int(start)})
	var reachable []vaultEdge
	for current := range queue.PopSeq() {
		if current.distance != distance[current.to] {
			continue
		}
		node := s.graph.nodes[current.to]
		if isKey(node.tile) && keys&keyBit(node.tile) == 0 {
			reachable = append(reachable, current)
			continue
		}
		for _, edge := range node.edges {
			tile := s.graph.nodes[edge.to].tile
			if isDoor(tile) && keys&keyBit(tile-'A'+'a') == 0 {
				continue
			}
			cost := current.distance + edge.distance
			if cost < distance[edge.to] {
				distance[edge.to] = cost
				queue.Push(vaultEdge{to: edge.to, distance: cost})
			}
		}
	}
	return reachable
}

func (s *vaultSearch) remainingDistance(current vaultState) int {
	if current.keys == s.graph.allKeys {
		return 0
	}
	if distance, ok := s.memo[current]; ok {
		return distance
	}
	best := math.MaxInt
	for robot := range s.graph.entrances {
		for _, move := range s.reachableKeys(current.positions[robot], current.keys) {
			next := current
			next.positions[robot] = uint8(move.to)
			next.keys |= keyBit(s.graph.nodes[move.to].tile)
			if remaining := s.remainingDistance(next); remaining != math.MaxInt {
				best = min(best, move.distance+remaining)
			}
		}
	}
	s.memo[current] = best
	return best
}

func vaultDistance(r io.Reader, robots int) (string, error) {
	v, err := parseVault(r)
	if err != nil {
		return "", err
	}
	if robots == 4 && len(v.entrances) == 1 {
		if err := v.splitEntrance(); err != nil {
			return "", err
		}
	}
	if len(v.entrances) != robots {
		return "", fmt.Errorf("vault needs %d entrances, got %d", robots, len(v.entrances))
	}

	search := vaultSearch{graph: buildVaultGraph(v), memo: make(map[vaultState]int)}
	var start vaultState
	copy(start.positions[:], search.graph.entrances)
	distance := search.remainingDistance(start)
	if distance == math.MaxInt {
		return "", fmt.Errorf("vault keys are unreachable")
	}
	return strconv.Itoa(distance), nil
}

func day18p01(r io.Reader) (string, error) {
	return vaultDistance(r, 1)
}

func day18p02(r io.Reader) (string, error) {
	return vaultDistance(r, 4)
}
