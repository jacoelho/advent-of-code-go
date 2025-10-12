package aoc2019

import (
	"io"
	"iter"
	"slices"
	"strconv"
	"unicode"

	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
)

// 1 space padding + 1 letter
const outerPortalThreshold = 2

type donutMaze struct {
	passages     grid.Grid2D[int, bool]
	portals      map[grid.Position2D[int]]grid.Position2D[int]
	outerPortals map[grid.Position2D[int]]bool
	start        grid.Position2D[int]
	end          grid.Position2D[int]
}

type recursiveState struct {
	pos   grid.Position2D[int]
	level int
}

func findPortalLabel(lines []string, x, y int) (string, grid.Position2D[int], bool) {
	if y >= len(lines) || x >= len(lines[y]) {
		return "", grid.Position2D[int]{}, false
	}

	ch := rune(lines[y][x])
	if !unicode.IsUpper(ch) {
		return "", grid.Position2D[int]{}, false
	}

	directions := []struct {
		dx, dy            int                    // second letter direction
		passages          []struct{ dx, dy int } // passages to check
		checkSecondLetter func() bool
		getSecondLetter   func() rune
	}{
		{
			dx: 1, dy: 0, // horizontal
			passages:          []struct{ dx, dy int }{{-1, 0}, {2, 0}},
			checkSecondLetter: func() bool { return x+1 < len(lines[y]) },
			getSecondLetter:   func() rune { return rune(lines[y][x+1]) },
		},
		{
			dx: 0, dy: 1, // vertical
			passages:          []struct{ dx, dy int }{{0, -1}, {0, 2}},
			checkSecondLetter: func() bool { return y+1 < len(lines) && x < len(lines[y+1]) },
			getSecondLetter:   func() rune { return rune(lines[y+1][x]) },
		},
	}

	for _, dir := range directions {
		if !dir.checkSecondLetter() {
			continue
		}

		ch2 := dir.getSecondLetter()
		if !unicode.IsUpper(ch2) {
			continue
		}

		label := string(ch) + string(ch2)

		for _, p := range dir.passages {
			px, py := x+p.dx, y+p.dy
			if py >= 0 && py < len(lines) && px >= 0 && px < len(lines[py]) && lines[py][px] == '.' {
				return label, grid.NewPosition2D(px, py), true
			}
		}
	}

	return "", grid.Position2D[int]{}, false
}

func parseMaze(r io.Reader) *donutMaze {
	s := scanner.NewScanner(r, func(line []byte) (string, error) {
		return string(line), nil
	})
	lines := slices.Collect(s.Values())

	maze := &donutMaze{
		passages:     make(grid.Grid2D[int, bool]),
		portals:      make(map[grid.Position2D[int]]grid.Position2D[int]),
		outerPortals: make(map[grid.Position2D[int]]bool),
	}

	// first pass: identify passages
	for y, line := range lines {
		for x, ch := range line {
			if ch == '.' {
				pos := grid.NewPosition2D(x, y)
				maze.passages[pos] = true
			}
		}
	}

	// second pass: find portal labels and map them
	portalLabels := make(map[string][]grid.Position2D[int])

	for y, line := range lines {
		for x := range line {
			label, passagePos, found := findPortalLabel(lines, x, y)
			if !found {
				continue
			}

			switch label {
			case "AA":
				maze.start = passagePos
			case "ZZ":
				maze.end = passagePos
			default:
				portalLabels[label] = append(portalLabels[label], passagePos)
			}
		}
	}

	// connect portal pairs
	for _, positions := range portalLabels {
		if len(positions) == 2 {
			maze.portals[positions[0]] = positions[1]
			maze.portals[positions[1]] = positions[0]
		}
	}

	minX, maxX, minY, maxY := maze.passages.Dimensions()

	isOuterPosition := func(pos grid.Position2D[int]) bool {
		return pos.X <= minX+outerPortalThreshold || pos.X >= maxX-outerPortalThreshold ||
			pos.Y <= minY+outerPortalThreshold || pos.Y >= maxY-outerPortalThreshold
	}

	// classify all portals (including AA/ZZ) as outer or inner
	for pos := range maze.portals {
		if isOuterPosition(pos) {
			maze.outerPortals[pos] = true
		}
	}

	// AA and ZZ are always outer portals
	maze.outerPortals[maze.start] = true
	maze.outerPortals[maze.end] = true

	return maze
}

func (m *donutMaze) yieldAdjacentPassages(pos grid.Position2D[int], yield func(grid.Position2D[int]) bool) bool {
	for next := range grid.Neighbours4(pos) {
		if m.passages.Contains(next) {
			if !yield(next) {
				return false
			}
		}
	}
	return true
}

func (m *donutMaze) neighbours(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(yield func(grid.Position2D[int]) bool) {
		if !m.yieldAdjacentPassages(pos, yield) {
			return
		}

		// portal teleportation
		if dest, hasPortal := m.portals[pos]; hasPortal {
			if !yield(dest) {
				return
			}
		}
	}
}

func day20p01(r io.Reader) (string, error) {
	maze := parseMaze(r)
	steps := search.BFSDistanceTo(maze.start, maze.end, maze.neighbours)
	return strconv.Itoa(steps), nil
}

func (m *donutMaze) recursiveNeighbours(state recursiveState) iter.Seq[recursiveState] {
	return func(yield func(recursiveState) bool) {
		yieldAtCurrentLevel := func(pos grid.Position2D[int]) bool {
			return yield(recursiveState{pos: pos, level: state.level})
		}
		if !m.yieldAdjacentPassages(state.pos, yieldAtCurrentLevel) {
			return
		}

		// portal teleportation with level changes
		dest, hasPortal := m.portals[state.pos]
		if !hasPortal {
			return
		}

		isOuter := m.outerPortals[state.pos]

		// at level 0, outer portals act as walls (can't go up from outermost level)
		if state.level == 0 && isOuter {
			return
		}

		// outer portals go up a level (decrease), inner portals go down (increase)
		newLevel := state.level
		if isOuter {
			newLevel--
		} else {
			newLevel++
		}

		yield(recursiveState{pos: dest, level: newLevel})
	}
}

func day20p02(r io.Reader) (string, error) {
	maze := parseMaze(r)
	start := recursiveState{pos: maze.start, level: 0}
	end := recursiveState{pos: maze.end, level: 0}
	steps := search.BFSDistanceTo(start, end, maze.recursiveNeighbours)
	return strconv.Itoa(steps), nil
}
