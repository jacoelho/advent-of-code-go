package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type trailSegment struct {
	toIdx    int
	distance int
}

type trailNetwork struct {
	segments [][]trailSegment
	startIdx int
	endIdx   int
}

type trailConditions uint8

const (
	slipperyTrail trailConditions = iota
	dryTrail
)

const (
	tilePath       = '.'
	tileForest     = '#'
	tileSlopeUp    = '^'
	tileSlopeRight = '>'
	tileSlopeDown  = 'v'
	tileSlopeLeft  = '<'

	maxTrailCells     = 1_000_000
	maxTrailJunctions = 4096
	maxHikingStates   = 50_000_000
)

func isSlope(tile rune) bool {
	return tile == tileSlopeUp || tile == tileSlopeRight || tile == tileSlopeDown || tile == tileSlopeLeft
}

func parseTrail(r io.Reader) (grid.Grid2D[int, rune], grid.Position2D[int], grid.Position2D[int], error) {
	lineScanner := scanner.NewScanner(r, func(bytes []byte) (string, error) {
		return string(bytes), nil
	})

	trail := make(grid.Grid2D[int, rune])
	var start, end grid.Position2D[int]
	y := 0
	for line := range lineScanner.Values() {
		if len(trail)+len(line) > maxTrailCells || y >= maxTrailCells {
			return nil, start, end, fmt.Errorf("trail exceeds %d cells or rows", maxTrailCells)
		}
		for x, char := range line {
			if char != tilePath && char != tileForest && !isSlope(char) {
				return nil, start, end, fmt.Errorf("invalid trail tile %q at %d,%d", char, x, y)
			}
			trail[grid.NewPosition2D(x, y)] = char
		}
		y++
	}
	if err := lineScanner.Err(); err != nil {
		return nil, start, end, err
	}

	var startFound, endFound bool
	minX, maxX, minY, maxY := trail.Dimensions()
	for x := minX; x <= maxX; x++ {
		pos := grid.NewPosition2D(x, minY)
		if trail[pos] == tilePath {
			start = pos
			startFound = true
			break
		}
	}
	for x := minX; x <= maxX; x++ {
		pos := grid.NewPosition2D(x, maxY)
		if trail[pos] == tilePath {
			end = pos
			endFound = true
			break
		}
	}
	if !startFound || !endFound {
		return nil, start, end, fmt.Errorf("trail start or end not found")
	}

	return trail, start, end, nil
}

func trailNeighbors(trail grid.Grid2D[int, rune], pos grid.Position2D[int], conditions trailConditions) []grid.Position2D[int] {
	directions := grid.OffsetsNeighbours4[int]()
	if conditions == slipperyTrail {
		switch trail[pos] {
		case tileSlopeUp:
			directions = []grid.Position2D[int]{{X: 0, Y: -1}}
		case tileSlopeRight:
			directions = []grid.Position2D[int]{{X: 1, Y: 0}}
		case tileSlopeDown:
			directions = []grid.Position2D[int]{{X: 0, Y: 1}}
		case tileSlopeLeft:
			directions = []grid.Position2D[int]{{X: -1, Y: 0}}
		}
	}

	var neighbors []grid.Position2D[int]
	for _, dir := range directions {
		neighbor := pos.Add(dir)
		if tile, exists := trail[neighbor]; exists && tile != tileForest {
			neighbors = append(neighbors, neighbor)
		}
	}
	return neighbors
}

func day23p01(r io.Reader) (string, error) {
	return longestHike(r, slipperyTrail)
}

func day23p02(r io.Reader) (string, error) {
	return longestHike(r, dryTrail)
}

func longestHike(r io.Reader, conditions trailConditions) (string, error) {
	trail, start, end, err := parseTrail(r)
	if err != nil {
		return "", err
	}

	network, err := buildTrailNetwork(trail, start, end, conditions)
	if err != nil {
		return "", err
	}
	length, err := findLongestHike(network)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(length), nil
}

func buildTrailNetwork(trail grid.Grid2D[int, rune], start, end grid.Position2D[int], conditions trailConditions) (trailNetwork, error) {
	junctions := findTrailJunctions(trail, start, end)
	if len(junctions) > maxTrailJunctions {
		return trailNetwork{}, fmt.Errorf("trail exceeds %d junctions", maxTrailJunctions)
	}

	indexMap := make(map[grid.Position2D[int]]int, len(junctions))
	for i, junction := range junctions {
		indexMap[junction] = i
	}

	network := trailNetwork{
		segments: make([][]trailSegment, len(junctions)),
		startIdx: indexMap[start],
		endIdx:   indexMap[end],
	}
	for i, junction := range junctions {
		for _, neighbor := range trailNeighbors(trail, junction, conditions) {
			segment, ok, err := followTrail(trail, indexMap, junction, neighbor, conditions)
			if err != nil {
				return trailNetwork{}, err
			}
			if ok {
				network.segments[i] = append(network.segments[i], segment)
			}
		}
	}
	return network, nil
}

func findTrailJunctions(trail grid.Grid2D[int, rune], start, end grid.Position2D[int]) []grid.Position2D[int] {
	var junctions []grid.Position2D[int]
	for pos, tile := range trail {
		if tile == tileForest {
			continue
		}
		if pos == start || pos == end || len(trailNeighbors(trail, pos, dryTrail)) != 2 {
			junctions = append(junctions, pos)
		}
	}
	slices.SortFunc(junctions, func(a, b grid.Position2D[int]) int {
		return cmp.Or(cmp.Compare(a.Y, b.Y), cmp.Compare(a.X, b.X))
	})
	return junctions
}

func followTrail(
	trail grid.Grid2D[int, rune],
	junctions map[grid.Position2D[int]]int,
	previous, current grid.Position2D[int],
	conditions trailConditions,
) (trailSegment, bool, error) {
	for distance := 1; distance <= len(trail); distance++ {
		if idx, ok := junctions[current]; ok {
			return trailSegment{toIdx: idx, distance: distance}, true, nil
		}

		// Only degree-two cells remain inside corridors; a slope can block their sole continuation.
		var next grid.Position2D[int]
		found := false
		for _, neighbor := range trailNeighbors(trail, current, conditions) {
			if neighbor != previous {
				next = neighbor
				found = true
				break
			}
		}
		if !found {
			return trailSegment{}, false, nil
		}
		previous, current = current, next
	}
	return trailSegment{}, false, fmt.Errorf("trail corridor does not reach a junction")
}

func findLongestHike(network trailNetwork) (int, error) {
	visited := make([]bool, len(network.segments))
	states := 0
	longest := -1

	var explore func(int, int) error
	explore = func(current, distance int) error {
		states++
		if states > maxHikingStates {
			return fmt.Errorf("hiking search exceeds %d states", maxHikingStates)
		}
		if current == network.endIdx {
			longest = max(longest, distance)
			return nil
		}

		visited[current] = true
		for _, segment := range network.segments[current] {
			if !visited[segment.toIdx] {
				if err := explore(segment.toIdx, distance+segment.distance); err != nil {
					return err
				}
			}
		}
		visited[current] = false
		return nil
	}
	if err := explore(network.startIdx, 0); err != nil {
		return 0, err
	}
	if longest < 0 {
		return 0, fmt.Errorf("no hiking route reaches the end")
	}
	return longest, nil
}
