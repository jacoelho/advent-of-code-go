package aoc2019

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const (
	erisSize   = 5
	erisCells  = erisSize * erisSize
	erisCenter = erisCells / 2
)

// erisNeighbours holds, for each cell, bit masks of its neighbours on the same
// level, on the enclosing level and on the level nested in the centre.
type erisNeighbours struct {
	same, outer, inner [erisCells]uint32
}

var erisAdjacency = func() erisNeighbours {
	var n erisNeighbours
	bit := func(x, y int) uint32 { return 1 << (y*erisSize + x) }
	edge := func(x0, y0, dx, dy int) uint32 {
		var mask uint32
		for i := range erisSize {
			mask |= bit(x0+i*dx, y0+i*dy)
		}
		return mask
	}
	const c = erisSize / 2
	for y := range erisSize {
		for x := range erisSize {
			cell := y*erisSize + x
			for _, d := range [][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
				nx, ny := x+d[0], y+d[1]
				switch {
				case nx < 0 || nx >= erisSize || ny < 0 || ny >= erisSize:
					n.outer[cell] |= bit(c+d[0], c+d[1])
				case nx == c && ny == c:
					n.same[cell] |= bit(nx, ny)
					switch d {
					case [2]int{0, -1}:
						n.inner[cell] |= edge(0, erisSize-1, 1, 0)
					case [2]int{1, 0}:
						n.inner[cell] |= edge(0, 0, 0, 1)
					case [2]int{0, 1}:
						n.inner[cell] |= edge(0, 0, 1, 0)
					case [2]int{-1, 0}:
						n.inner[cell] |= edge(erisSize-1, 0, 0, 1)
					}
				default:
					n.same[cell] |= bit(nx, ny)
				}
			}
		}
	}
	return n
}()

func parseEris(r io.Reader) (uint32, error) {
	var state uint32
	rows := 0
	s := scanner.NewScanner(r, func(line []byte) ([]byte, error) { return line, nil })
	for line := range s.Values() {
		if rows == erisSize || len(line) != erisSize {
			return 0, fmt.Errorf("eris scan must be %dx%d", erisSize, erisSize)
		}
		for x, c := range line {
			cell := rows*erisSize + x
			switch {
			case c == '#':
				state |= 1 << cell
			case c == '.', c == '?' && cell == erisCenter:
			default:
				return 0, fmt.Errorf("invalid eris cell %q at (%d,%d)", c, x, rows)
			}
		}
		rows++
	}
	if err := s.Err(); err != nil {
		return 0, err
	}
	if rows != erisSize {
		return 0, fmt.Errorf("eris scan must be %dx%d", erisSize, erisSize)
	}
	return state, nil
}

func bugSurvives(hasBug bool, adjacent int) bool {
	return adjacent == 1 || !hasBug && adjacent == 2
}

func day24p01(r io.Reader) (string, error) {
	state, err := parseEris(r)
	if err != nil {
		return "", err
	}
	seen := map[uint32]bool{}
	for !seen[state] {
		seen[state] = true
		var next uint32
		for cell := range erisCells {
			adjacent := bits.OnesCount32(state & erisAdjacency.same[cell])
			if bugSurvives(state&(1<<cell) != 0, adjacent) {
				next |= 1 << cell
			}
		}
		state = next
	}
	return strconv.FormatUint(uint64(state), 10), nil
}

func simulateRecursive(r io.Reader, minutes int) (string, error) {
	state, err := parseEris(r)
	if err != nil {
		return "", err
	}
	// Bugs spread at most one level per minute, so the levels within minutes
	// of the start, plus an empty border on each side, hold every bug.
	levels := make([]uint32, 2*minutes+3)
	levels[minutes+1] = state &^ (1 << erisCenter)
	next := make([]uint32, len(levels))
	for range minutes {
		for level := 1; level < len(levels)-1; level++ {
			var grid uint32
			for cell := range erisCells {
				if cell == erisCenter {
					continue
				}
				adjacent := bits.OnesCount32(levels[level]&erisAdjacency.same[cell]&^(1<<erisCenter)) +
					bits.OnesCount32(levels[level-1]&erisAdjacency.outer[cell]) +
					bits.OnesCount32(levels[level+1]&erisAdjacency.inner[cell])
				if bugSurvives(levels[level]&(1<<cell) != 0, adjacent) {
					grid |= 1 << cell
				}
			}
			next[level] = grid
		}
		levels, next = next, levels
	}

	total := 0
	for _, grid := range levels {
		total += bits.OnesCount32(grid)
	}
	return strconv.Itoa(total), nil
}

func day24p02(r io.Reader) (string, error) {
	return simulateRecursive(r, 200)
}
