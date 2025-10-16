package aoc2023

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type position = grid.Position2D[int]

var (
	north = position{X: 0, Y: -1}
	south = position{X: 0, Y: 1}
	east  = position{X: 1, Y: 0}
	west  = position{X: -1, Y: 0}
)

func parseGrid(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	return grid.NewGrid2D[int, rune](rows), nil
}

func findStart(g grid.Grid2D[int, rune]) (position, bool) {
	for pos, ch := range g {
		if ch == 'S' {
			return pos, true
		}
	}
	return position{}, false
}

func connects(g grid.Grid2D[int, rune], pos position, dir position) bool {
	ch, exists := g[pos]
	if !exists {
		return false
	}

	switch dir {
	case north:
		return ch == '|' || ch == 'L' || ch == 'J' || ch == 'S'
	case south:
		return ch == '|' || ch == '7' || ch == 'F' || ch == 'S'
	case east:
		return ch == '-' || ch == 'L' || ch == 'F' || ch == 'S'
	case west:
		return ch == '-' || ch == 'J' || ch == '7' || ch == 'S'
	default:
		return false
	}
}

func determineStartPipe(g grid.Grid2D[int, rune], start position) rune {
	n := connects(g, start.Add(north), south)
	s := connects(g, start.Add(south), north)
	e := connects(g, start.Add(east), west)
	w := connects(g, start.Add(west), east)

	switch {
	case n && s:
		return '|'
	case e && w:
		return '-'
	case n && e:
		return 'L'
	case n && w:
		return 'J'
	case s && w:
		return '7'
	case s && e:
		return 'F'
	default:
		return 'S'
	}
}

var pipeConnections = map[rune][]position{
	'|': {north, south},
	'-': {east, west},
	'L': {north, east},
	'J': {north, west},
	'7': {south, west},
	'F': {south, east},
}

var oppositeDir = map[position]position{
	north: south,
	south: north,
	east:  west,
	west:  east,
}

func neighbors(g grid.Grid2D[int, rune], pos position) iter.Seq[position] {
	return func(yield func(position) bool) {
		ch := g[pos]
		dirs, ok := pipeConnections[ch]
		if !ok {
			return
		}

		for _, dir := range dirs {
			neighbor := pos.Add(dir)
			if connects(g, neighbor, oppositeDir[dir]) {
				if !yield(neighbor) {
					return
				}
			}
		}
	}
}

func day10p01(r io.Reader) (string, error) {
	g, err := parseGrid(r)
	if err != nil {
		return "", err
	}

	start, found := findStart(g)
	if !found {
		return "", fmt.Errorf("start position not found")
	}

	g[start] = determineStartPipe(g, start)

	maxDist := search.BFSMaxDistance(start, func(pos position) iter.Seq[position] {
		return neighbors(g, pos)
	})

	return strconv.Itoa(maxDist), nil
}

func getLoopPositions(g grid.Grid2D[int, rune], start position) collections.Set[position] {
	bfs := search.BFS(start, func(pos position) iter.Seq[position] {
		return neighbors(g, pos)
	})

	loopSet := collections.NewSet[position]()
	for pos := range bfs {
		loopSet.Add(pos)
	}
	return loopSet
}

func day10p02(r io.Reader) (string, error) {
	g, err := parseGrid(r)
	if err != nil {
		return "", err
	}

	start, found := findStart(g)
	if !found {
		return "", fmt.Errorf("start position not found")
	}

	g[start] = determineStartPipe(g, start)

	loopPositions := getLoopPositions(g, start)

	minX, maxX, minY, maxY := g.Dimensions()

	enclosed := 0
	for y := minY; y <= maxY; y++ {
		inside := false
		entryBend := ' '

		for x := minX; x <= maxX; x++ {
			pos := position{X: x, Y: y}

			if loopPositions.Contains(pos) {
				ch := g[pos]

				switch ch {
				case '|':
					inside = !inside
				case 'F', 'L':
					entryBend = ch
				case '7', 'J':
					// end of a horizontal segment
					// F---7: both point down, doesn't cross (same side)
					// L---J: both point up, doesn't cross (same side)
					// F---J: F points down, J points up, crosses
					// L---7: L points up, 7 points down, crosses
					if (entryBend == 'F' && ch == 'J') || (entryBend == 'L' && ch == '7') {
						inside = !inside
					}
					entryBend = ' '
				}
			} else if inside {
				// not on loop and inside - count it
				enclosed++
			}
		}
	}

	return strconv.Itoa(enclosed), nil
}
