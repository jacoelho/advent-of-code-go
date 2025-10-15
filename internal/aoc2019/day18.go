package aoc2019

import (
	"bufio"
	"io"
	"iter"
	"maps"
	"math/bits"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type vault struct {
	passages  map[grid.Position2D[int]]rune
	entrances []grid.Position2D[int]
	keys      map[rune]grid.Position2D[int]
}

func parseVault(r io.Reader) *vault {
	v := &vault{
		passages: make(map[grid.Position2D[int]]rune),
		keys:     make(map[rune]grid.Position2D[int]),
	}

	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, ch := range line {
			if ch == '#' {
				continue
			}
			pos := grid.NewPosition2D(x, y)
			v.passages[pos] = ch

			if ch == '@' {
				v.entrances = append(v.entrances, pos)
			} else if ch >= 'a' && ch <= 'z' {
				v.keys[ch] = pos
			}
		}
		y++
	}

	return v
}

type pathInfo struct {
	distance     int
	requiredKeys uint32
}

func keyBit(k rune) uint32 {
	return 1 << (k - 'a')
}

func (v *vault) neighbours(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(yield func(grid.Position2D[int]) bool) {
		for next := range grid.Neighbours4(p) {
			if _, ok := v.passages[next]; ok {
				if !yield(next) {
					return
				}
			}
		}
	}
}

func (v *vault) allKeysMask() uint32 {
	var mask uint32
	for k := range v.keys {
		mask |= keyBit(k)
	}
	return mask
}

func countRemainingKeys(keysCollected, allKeysMask uint32) int {
	return bits.OnesCount32((allKeysMask ^ keysCollected) & allKeysMask)
}

func canReachKey(path pathInfo, keyAtPos rune, keysCollected uint32) bool {
	if (path.requiredKeys & keysCollected) != path.requiredKeys {
		return false
	}
	if keyAtPos < 'a' || keyAtPos > 'z' {
		return false
	}
	keyMask := keyBit(keyAtPos)
	return (keysCollected & keyMask) == 0
}

func (v *vault) buildDistanceGraph(startPositions []grid.Position2D[int]) map[grid.Position2D[int]]map[grid.Position2D[int]]pathInfo {
	graph := make(map[grid.Position2D[int]]map[grid.Position2D[int]]pathInfo)

	pointsOfInterest := append([]grid.Position2D[int]{}, startPositions...)
	pointsOfInterest = append(pointsOfInterest, slices.Collect(maps.Values(v.keys))...)

	for _, start := range pointsOfInterest {
		distances := make(map[grid.Position2D[int]]int)
		requiredKeys := make(map[grid.Position2D[int]]uint32)
		distances[start] = 0
		requiredKeys[start] = 0

		for node := range search.BFS(start, v.neighbours) {
			currentDist := distances[node]
			currentKeys := requiredKeys[node]

			for neighbour := range v.neighbours(node) {
				if _, seen := distances[neighbour]; !seen {
					distances[neighbour] = currentDist + 1
					newKeys := currentKeys
					ch := v.passages[neighbour]
					if ch >= 'A' && ch <= 'Z' {
						newKeys |= keyBit(rune(ch - 'A' + 'a'))
					}
					requiredKeys[neighbour] = newKeys
				}
			}
		}

		graph[start] = make(map[grid.Position2D[int]]pathInfo)
		for _, targetPos := range pointsOfInterest {
			if targetPos == start {
				continue
			}
			if dist, ok := distances[targetPos]; ok {
				graph[start][targetPos] = pathInfo{
					distance:     dist,
					requiredKeys: requiredKeys[targetPos],
				}
			}
		}
	}

	return graph
}

type state struct {
	pos           grid.Position2D[int]
	keysCollected uint32
}

func (v *vault) findShortestPath() int {
	graph := v.buildDistanceGraph(v.entrances[:1])
	allKeysMask := v.allKeysMask()

	startState := state{pos: v.entrances[0], keysCollected: 0}

	neighbours := func(current state) []state {
		var result []state
		for nextPos, path := range graph[current.pos] {
			keyAtPos := v.passages[nextPos]
			if !canReachKey(path, keyAtPos, current.keysCollected) {
				continue
			}

			newState := state{
				pos:           nextPos,
				keysCollected: current.keysCollected | keyBit(keyAtPos),
			}
			result = append(result, newState)
		}
		return result
	}

	heuristic := func(s state) int {
		return countRemainingKeys(s.keysCollected, allKeysMask)
	}

	stepCost := func(from, to state) int {
		return graph[from.pos][to.pos].distance
	}

	distance, _, _ := search.AStar(startState, neighbours, heuristic, stepCost)
	return distance
}

func day18p01(r io.Reader) (string, error) {
	v := parseVault(r)
	result := v.findShortestPath()
	return strconv.Itoa(result), nil
}

func (v *vault) transformToMultiEntrance() {
	center := v.entrances[0]

	delete(v.passages, center)
	for _, offset := range []grid.Position2D[int]{
		{X: -1, Y: 0}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: 0, Y: 1},
	} {
		delete(v.passages, center.Add(offset))
	}

	v.entrances = []grid.Position2D[int]{
		center.Add(grid.Position2D[int]{X: -1, Y: -1}),
		center.Add(grid.Position2D[int]{X: 1, Y: -1}),
		center.Add(grid.Position2D[int]{X: -1, Y: 1}),
		center.Add(grid.Position2D[int]{X: 1, Y: 1}),
	}

	for _, entrance := range v.entrances {
		v.passages[entrance] = '@'
	}
}

type multiState struct {
	positions     [4]grid.Position2D[int]
	keysCollected uint32
}

func (v *vault) findShortestPathMultiRobot() int {
	graph := v.buildDistanceGraph(v.entrances)
	allKeysMask := v.allKeysMask()

	var startPositions [4]grid.Position2D[int]
	copy(startPositions[:], v.entrances)
	startState := multiState{
		positions:     startPositions,
		keysCollected: 0,
	}

	neighbours := func(current multiState) []multiState {
		var result []multiState

		for robotIdx := range 4 {
			robotPos := current.positions[robotIdx]

			for nextPos, path := range graph[robotPos] {
				keyAtPos := v.passages[nextPos]
				if !canReachKey(path, keyAtPos, current.keysCollected) {
					continue
				}

				newState := multiState{
					positions:     current.positions,
					keysCollected: current.keysCollected | keyBit(keyAtPos),
				}
				newState.positions[robotIdx] = nextPos
				result = append(result, newState)
			}
		}
		return result
	}

	heuristic := func(s multiState) int {
		return countRemainingKeys(s.keysCollected, allKeysMask)
	}

	stepCost := func(from, to multiState) int {
		for i := range 4 {
			if from.positions[i] != to.positions[i] {
				return graph[from.positions[i]][to.positions[i]].distance
			}
		}
		return 0
	}

	distance, _, _ := search.AStar(startState, neighbours, heuristic, stepCost)
	return distance
}

func day18p02(r io.Reader) (string, error) {
	v := parseVault(r)
	if len(v.entrances) == 1 {
		v.transformToMultiEntrance()
	}
	result := v.findShortestPathMultiRobot()
	return strconv.Itoa(result), nil
}
