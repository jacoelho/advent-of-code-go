package aoc2019

import (
	"bufio"
	"io"
	"math/bits"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

const (
	gridSize  = 5
	centerPos = 12
	centerX   = 2
	centerY   = 2
)

type gridMode int

const (
	gridModeFlat gridMode = iota
	gridModeRecursive
)

func gridPosition(x, y int) int {
	return y*gridSize + x
}

func countBugsInPositions(grid uint32, positions []int) int {
	return xiter.CountBy(func(pos int) bool {
		return isBug(grid, pos)
	}, slices.Values(positions))
}

func shouldHaveBug(hasBug bool, adjacent int) bool {
	if hasBug {
		return adjacent == 1
	}
	return adjacent == 1 || adjacent == 2
}

func parseGrid(r io.Reader, mode gridMode) uint32 {
	var state uint32
	scanner := bufio.NewScanner(r)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, ch := range line {
			pos := gridPosition(x, y)
			if ch == '?' {
				continue // examples use ? to denote empty spaces
			}
			if mode == gridModeRecursive && pos == centerPos {
				continue // center is always empty in recursive mode
			}
			if ch == '#' {
				state |= 1 << pos
			}
		}
		y++
	}
	return state
}

func isBug(state uint32, pos int) bool { return state&(1<<pos) != 0 }

func countAdjacent(state uint32, x, y int) int {
	pos := grid.NewPosition2D(x, y)
	return xiter.CountBy(func(neighbor grid.Position2D[int]) bool {
		if neighbor.X < 0 || neighbor.X >= gridSize || neighbor.Y < 0 || neighbor.Y >= gridSize {
			return false
		}
		neighborPos := gridPosition(neighbor.X, neighbor.Y)
		return isBug(state, neighborPos)
	}, grid.Neighbours4(pos))
}

func step(state uint32) uint32 {
	var next uint32
	for y := range gridSize {
		for x := range gridSize {
			pos := gridPosition(x, y)
			adjacent := countAdjacent(state, x, y)

			if shouldHaveBug(isBug(state, pos), adjacent) {
				next |= 1 << pos
			}
		}
	}
	return next
}

func day24p01(r io.Reader) (string, error) {
	state := parseGrid(r, gridModeFlat)
	seen := collections.NewSet[uint32]()

	for {
		if seen.Contains(state) {
			return strconv.FormatUint(uint64(state), 10), nil
		}
		seen.Add(state)
		state = step(state)
	}
}

// check out-of-bounds cases that connect to outer level
func countOuterLevel(grids map[int]uint32, level, x, y int) int {
	var pos int
	switch {
	case x < 0:
		pos = gridPosition(centerX-1, centerY)
	case x >= gridSize:
		pos = gridPosition(centerX+1, centerY)
	case y < 0:
		pos = gridPosition(centerX, centerY-1)
	case y >= gridSize:
		pos = gridPosition(centerX, centerY+1)
	default:
		return 0
	}
	if isBug(grids[level-1], pos) {
		return 1
	}
	return 0
}

// check center adjacency cases that connect to inner level
func countInnerLevel(grids map[int]uint32, level, x, y int) int {
	var positions []int
	switch {
	case x == centerX-1 && y == centerY:
		// left of center - count left column of inner level
		for iy := range gridSize {
			positions = append(positions, gridPosition(0, iy))
		}
	case x == centerX+1 && y == centerY:
		// right of center - count right column of inner level
		for iy := range gridSize {
			positions = append(positions, gridPosition(gridSize-1, iy))
		}
	case x == centerX && y == centerY-1:
		// above center - count top row of inner level
		for ix := range gridSize {
			positions = append(positions, gridPosition(ix, 0))
		}
	case x == centerX && y == centerY+1:
		// below center - count bottom row of inner level
		for ix := range gridSize {
			positions = append(positions, gridPosition(ix, gridSize-1))
		}
	default:
		return 0
	}
	return countBugsInPositions(grids[level+1], positions)
}

func countAdjacentRecursive(grids map[int]uint32, level, x, y int) int {
	count := 0
	neighbors := [][2]int{{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}}

	for _, neighbor := range neighbors {
		nx, ny := neighbor[0], neighbor[1]

		// 1: same level neighbor (in bounds, not center)
		if nx >= 0 && nx < gridSize && ny >= 0 && ny < gridSize && !(nx == centerX && ny == centerY) {
			neighborPos := gridPosition(nx, ny)
			if isBug(grids[level], neighborPos) {
				count++
			}
		} else if nx < 0 || nx >= gridSize || ny < 0 || ny >= gridSize {
			// 2-5: out of bounds - connect to outer level
			count += countOuterLevel(grids, level, nx, ny)
		} else if nx == centerX && ny == centerY {
			// 6: adjacent to center - count edge of inner level
			count += countInnerLevel(grids, level, x, y)
		}
	}
	return count
}

func stepRecursive(grids map[int]uint32, minutes int) map[int]uint32 {
	for range minutes {
		next := make(map[int]uint32)

		// determine depth range to consider
		minLevel, maxLevel := 0, 0
		for level := range grids {
			minLevel = min(minLevel, level)
			maxLevel = max(maxLevel, level)
		}

		// expand range by 1 in each direction for new bugs
		for level := minLevel - 1; level <= maxLevel+1; level++ {
			var nextState uint32
			for y := range gridSize {
				for x := range gridSize {
					if x == centerX && y == centerY {
						continue
					}
					pos := gridPosition(x, y)
					adjacent := countAdjacentRecursive(grids, level, x, y)
					hasBug := isBug(grids[level], pos)

					if shouldHaveBug(hasBug, adjacent) {
						nextState |= 1 << pos
					}
				}
			}
			// center is never set
			nextState &= ^uint32(1 << centerPos)
			if nextState != 0 {
				next[level] = nextState
			}
		}
		grids = next
	}
	return grids
}

func simulateRecursive(r io.Reader, minutes int) (string, error) {
	state := parseGrid(r, gridModeRecursive)
	grids := map[int]uint32{0: state}
	grids = stepRecursive(grids, minutes)

	totalBugs := 0
	for _, grid := range grids {
		gridWithoutCenter := grid & ^uint32(1<<centerPos)
		totalBugs += bits.OnesCount32(gridWithoutCenter)
	}

	return strconv.Itoa(totalBugs), nil
}

func day24p02(r io.Reader) (string, error) {
	return simulateRecursive(r, 200)
}
