package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const (
	maxVaultCells  = 100_000
	maxVaultNodes  = 256
	maxVaultStates = 1_000_000
	maxVaultWork   = 200_000_000
)

type vault struct {
	width     int
	cells     []byte
	entrances []int
	allKeys   uint32
}

type vaultEdge struct {
	to       uint8
	distance int
}

type vaultNode struct {
	tile  byte
	edges []vaultEdge
}

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
	graph     vaultGraph
	remaining int
	memo      map[vaultState]int
}

func keyBit(key byte) uint32 {
	return 1 << (key - 'a')
}

func parseVault(r io.Reader) (vault, error) {
	var result vault
	s := scanner.NewScanner(r, func(line []byte) ([]byte, error) {
		if len(line) == 0 || result.width != 0 && len(line) != result.width {
			return nil, fmt.Errorf("empty or ragged vault")
		}
		if len(line) > maxVaultCells-len(result.cells) {
			return nil, fmt.Errorf("vault exceeds %d cells", maxVaultCells)
		}
		for _, tile := range line {
			if tile != '#' && tile != '.' && tile != '@' && !(tile >= 'a' && tile <= 'z') && !(tile >= 'A' && tile <= 'Z') {
				return nil, fmt.Errorf("invalid vault tile %q", tile)
			}
		}
		return bytes.Clone(line), nil
	})
	for row := range s.Values() {
		result.width = len(row)
		for x, tile := range row {
			switch {
			case tile == '@':
				result.entrances = append(result.entrances, len(result.cells)+x)
				if len(result.entrances) > 4 {
					return vault{}, fmt.Errorf("vault has more than four entrances")
				}
			case tile >= 'a' && tile <= 'z':
				if result.allKeys&keyBit(tile) != 0 {
					return vault{}, fmt.Errorf("duplicate vault key %q", tile)
				}
				result.allKeys |= keyBit(tile)
			}
		}
		result.cells = append(result.cells, row...)
	}
	if err := s.Err(); err != nil {
		return vault{}, err
	}
	if len(result.cells) == 0 {
		return vault{}, fmt.Errorf("empty vault")
	}
	return result, nil
}

func (v *vault) splitEntrance() error {
	center := v.entrances[0]
	x, y := center%v.width, center/v.width
	if x == 0 || x+1 == v.width || y == 0 || y+1 == len(v.cells)/v.width {
		return fmt.Errorf("vault entrance needs a surrounding 3x3 floor")
	}
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if (dx != 0 || dy != 0) && v.cells[center+dy*v.width+dx] != '.' {
				return fmt.Errorf("vault entrance needs a surrounding 3x3 floor")
			}
		}
	}
	v.entrances = nil
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			at := center + dy*v.width + dx
			v.cells[at] = '#'
			if dx != 0 && dy != 0 {
				v.cells[at] = '@'
				v.entrances = append(v.entrances, at)
			}
		}
	}
	return nil
}

func (s *vaultSearch) spend() error {
	if s.remaining == 0 {
		return fmt.Errorf("vault search exceeds %d work units", maxVaultWork)
	}
	s.remaining--
	return nil
}

func (s *vaultSearch) buildGraph(v vault) error {
	indices := make([]int, len(v.cells))
	var points []int
	for at, tile := range v.cells {
		indices[at] = -1
		if tile == '#' || tile == '.' {
			continue
		}
		if len(points) == maxVaultNodes {
			return fmt.Errorf("vault exceeds %d entrances, keys and doors", maxVaultNodes)
		}
		indices[at] = len(points)
		points = append(points, at)
		s.graph.nodes = append(s.graph.nodes, vaultNode{tile: tile})
		if tile == '@' {
			s.graph.entrances = append(s.graph.entrances, uint8(indices[at]))
		}
	}
	s.graph.allKeys = v.allKeys

	distance := make([]int, len(v.cells))
	queue := make([]int, 0, len(v.cells))
	for source, start := range points {
		for i := range distance {
			distance[i] = -1
		}
		distance[start] = 0
		queue = append(queue[:0], start)
		for head := 0; head < len(queue); head++ {
			at := queue[head]
			x, y := at%v.width, at/v.width
			for _, delta := range [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				if err := s.spend(); err != nil {
					return err
				}
				nx, ny := x+delta[0], y+delta[1]
				if nx < 0 || nx >= v.width || ny < 0 || ny >= len(v.cells)/v.width {
					continue
				}
				next := ny*v.width + nx
				if v.cells[next] == '#' || distance[next] >= 0 {
					continue
				}
				distance[next] = distance[at] + 1
				if target := indices[next]; target >= 0 {
					// Corridor interiors contain no keys or doors, so their shortest route is always valid.
					s.graph.nodes[source].edges = append(s.graph.nodes[source].edges, vaultEdge{to: uint8(target), distance: distance[next]})
					continue
				}
				queue = append(queue, next)
			}
		}
	}
	return nil
}

func (s *vaultSearch) reachableKeys(start uint8, keys uint32) ([]vaultEdge, error) {
	distance := make([]int, len(s.graph.nodes))
	for i := range distance {
		distance[i] = math.MaxInt
	}
	distance[start] = 0
	queue := collections.NewHeap(func(a, b vaultEdge) bool { return a.distance < b.distance })
	queue.Push(vaultEdge{to: start})
	var reachable []vaultEdge
	for current := range queue.PopSeq() {
		if err := s.spend(); err != nil {
			return nil, err
		}
		if current.distance != distance[current.to] {
			continue
		}
		node := s.graph.nodes[current.to]
		if node.tile >= 'a' && node.tile <= 'z' && keys&keyBit(node.tile) == 0 {
			reachable = append(reachable, current)
			continue
		}
		for _, edge := range node.edges {
			if err := s.spend(); err != nil {
				return nil, err
			}
			tile := s.graph.nodes[edge.to].tile
			if tile >= 'A' && tile <= 'Z' && keys&keyBit(tile-'A'+'a') == 0 {
				continue
			}
			cost := current.distance + edge.distance
			if cost < distance[edge.to] {
				distance[edge.to] = cost
				queue.Push(vaultEdge{to: edge.to, distance: cost})
			}
		}
	}
	return reachable, nil
}

func (s *vaultSearch) remainingDistance(current vaultState) (int, error) {
	if current.keys == s.graph.allKeys {
		return 0, nil
	}
	if distance, ok := s.memo[current]; ok {
		return distance, nil
	}
	best := math.MaxInt
	for robot := range s.graph.entrances {
		moves, err := s.reachableKeys(current.positions[robot], current.keys)
		if err != nil {
			return 0, err
		}
		for _, move := range moves {
			if err := s.spend(); err != nil {
				return 0, err
			}
			next := current
			next.positions[robot] = move.to
			next.keys |= keyBit(s.graph.nodes[move.to].tile)
			remaining, err := s.remainingDistance(next)
			if err != nil {
				return 0, err
			}
			if remaining != math.MaxInt {
				best = min(best, move.distance+remaining)
			}
		}
	}
	if len(s.memo) >= maxVaultStates {
		return 0, fmt.Errorf("vault search exceeds %d states", maxVaultStates)
	}
	s.memo[current] = best
	return best, nil
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

	search := vaultSearch{remaining: maxVaultWork, memo: make(map[vaultState]int)}
	if err := search.buildGraph(v); err != nil {
		return "", err
	}
	var start vaultState
	copy(start.positions[:], search.graph.entrances)
	distance, err := search.remainingDistance(start)
	if err != nil {
		return "", err
	}
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
