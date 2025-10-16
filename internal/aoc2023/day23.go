package aoc2023

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type trailSegment struct {
	toIdx    int
	distance int
}

type trailNetwork struct {
	segments [][]trailSegment
	indexMap map[grid.Position2D[int]]int
	startIdx int
	endIdx   int
}

const (
	tilePath       = '.'
	tileForest     = '#'
	tileSlopeUp    = '^'
	tileSlopeRight = '>'
	tileSlopeDown  = 'v'
	tileSlopeLeft  = '<'
)

var slopeDirections = map[rune][]grid.Position2D[int]{
	tileSlopeUp:    {{X: 0, Y: -1}},
	tileSlopeRight: {{X: 1, Y: 0}},
	tileSlopeDown:  {{X: 0, Y: 1}},
	tileSlopeLeft:  {{X: -1, Y: 0}},
}

func isSlope(tile rune) bool {
	return tile == tileSlopeUp || tile == tileSlopeRight || tile == tileSlopeDown || tile == tileSlopeLeft
}

func parseTrail(r io.Reader) (grid.Grid2D[int, rune], grid.Position2D[int], grid.Position2D[int], error) {
	lineScanner := scanner.NewScanner(r, func(bytes []byte) (string, error) {
		return string(bytes), nil
	})

	trail := make(grid.Grid2D[int, rune])
	var start, end grid.Position2D[int]
	var startFound, endFound bool

	y := 0
	for line := range lineScanner.Values() {
		for x, char := range line {
			trail[grid.NewPosition2D(x, y)] = char
		}
		y++
	}

	// Find start position (top row, single '.')
	minX, maxX, minY, maxY := trail.Dimensions()
	for x := minX; x <= maxX; x++ {
		pos := grid.NewPosition2D(x, minY)
		if trail[pos] == tilePath {
			start = pos
			startFound = true
			break
		}
	}

	// Find end position (bottom row, single '.')
	for x := minX; x <= maxX; x++ {
		pos := grid.NewPosition2D(x, maxY)
		if trail[pos] == tilePath {
			end = pos
			endFound = true
			break
		}
	}

	if !startFound || !endFound {
		return nil, grid.Position2D[int]{}, grid.Position2D[int]{}, fmt.Errorf("start or end not found")
	}

	return trail, start, end, nil
}

func dfsLongestPath(
	trail grid.Grid2D[int, rune],
	current, end grid.Position2D[int],
	visited collections.Set[grid.Position2D[int]],
	pathLength int,
) int {
	if current == end {
		return pathLength
	}

	visited.Add(current)
	maxLength := 0

	for _, neighbor := range getAccessiblePositions(trail, current, visited, false) {
		length := dfsLongestPath(trail, neighbor, end, visited, pathLength+1)
		maxLength = max(maxLength, length)
	}

	visited.Remove(current)
	return maxLength
}

func getAccessiblePositions(
	trail grid.Grid2D[int, rune],
	pos grid.Position2D[int],
	visited collections.Set[grid.Position2D[int]],
	climbableSlopes bool,
) []grid.Position2D[int] {
	tile := trail[pos]
	var directions []grid.Position2D[int]

	if tile == tilePath || (climbableSlopes && isSlope(tile)) {
		directions = grid.OffsetsNeighbours4[int]()
	} else if dirs, ok := slopeDirections[tile]; ok {
		directions = dirs
	} else {
		return nil
	}

	var result []grid.Position2D[int]
	for _, dir := range directions {
		neighbor := pos.Add(dir)
		if trail.Contains(neighbor) && trail[neighbor] != tileForest && !visited.Contains(neighbor) {
			result = append(result, neighbor)
		}
	}
	return result
}

func day23p01(r io.Reader) (string, error) {
	trail, start, end, err := parseTrail(r)
	if err != nil {
		return "", err
	}

	visited := collections.NewSet[grid.Position2D[int]]()
	longestPath := dfsLongestPath(trail, start, end, visited, 0)
	return strconv.Itoa(longestPath), nil
}

func day23p02(r io.Reader) (string, error) {
	trail, start, end, err := parseTrail(r)
	if err != nil {
		return "", err
	}

	network := buildTrailNetwork(trail, start, end)
	longestPath := findLongestHike(network)

	return strconv.Itoa(longestPath), nil
}

func buildTrailNetwork(trail grid.Grid2D[int, rune], start, end grid.Position2D[int]) trailNetwork {
	junctions := findTrailJunctions(trail, start, end)

	indexMap := make(map[grid.Position2D[int]]int)
	var startIdx, endIdx int
	idx := 0
	for junction := range junctions.Iter() {
		indexMap[junction] = idx
		if junction == start {
			startIdx = idx
		}
		if junction == end {
			endIdx = idx
		}
		idx++
	}

	segments := make([][]trailSegment, len(indexMap))
	emptyVisited := collections.NewSet[grid.Position2D[int]]()

	for junction := range junctions.Iter() {
		jIdx := indexMap[junction]
		neighbors := getAccessiblePositions(trail, junction, emptyVisited, true)

		for _, neighbor := range neighbors {
			toJunction, distance := followTrail(trail, junction, neighbor, 1)
			if junctions.Contains(toJunction) {
				toIdx := indexMap[toJunction]
				segments[jIdx] = append(segments[jIdx], trailSegment{toIdx: toIdx, distance: distance})
			}
		}
	}

	return trailNetwork{
		segments: segments,
		indexMap: indexMap,
		startIdx: startIdx,
		endIdx:   endIdx,
	}
}

func findTrailJunctions(trail grid.Grid2D[int, rune], start, end grid.Position2D[int]) collections.Set[grid.Position2D[int]] {
	junctions := collections.NewSet(start, end)
	emptyVisited := collections.NewSet[grid.Position2D[int]]()

	for pos := range trail {
		tile := trail[pos]
		if tile == tilePath || isSlope(tile) {
			neighbors := getAccessiblePositions(trail, pos, emptyVisited, true)
			if len(neighbors) != 2 && pos != start && pos != end {
				junctions.Add(pos)
			}
		}
	}

	return junctions
}

func followTrail(
	trail grid.Grid2D[int, rune],
	from, current grid.Position2D[int],
	distance int,
) (grid.Position2D[int], int) {
	emptyVisited := collections.NewSet[grid.Position2D[int]]()
	prev := from
	pos := current
	dist := distance

	for {
		neighbors := getAccessiblePositions(trail, pos, emptyVisited, true)
		nextOptions := xslices.Filter(func(n grid.Position2D[int]) bool {
			return n != prev
		}, neighbors)

		if len(nextOptions) == 0 {
			return pos, dist
		}

		nextPos := nextOptions[0]
		nextNeighbors := getAccessiblePositions(trail, nextPos, emptyVisited, true)
		if len(nextNeighbors) != 2 {
			return nextPos, dist + 1
		}

		prev = pos
		pos = nextPos
		dist++
	}
}

func findLongestHike(network trailNetwork) int {
	var visited uint64
	return exploreTrail(network, network.startIdx, network.endIdx, visited, 0)
}

func exploreTrail(network trailNetwork, current, end int, visited uint64, pathLength int) int {
	if current == end {
		return pathLength
	}

	visited |= (1 << current)
	maxLength := 0

	for _, segment := range network.segments[current] {
		if visited&(1<<segment.toIdx) == 0 {
			length := exploreTrail(network, segment.toIdx, end, visited, pathLength+segment.distance)
			maxLength = max(maxLength, length)
		}
	}

	return maxLength
}
