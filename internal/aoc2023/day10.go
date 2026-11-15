package aoc2023

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type pipeDirection struct {
	dx, dy int
}

var (
	pipeNorth = pipeDirection{0, -1}
	pipeSouth = pipeDirection{0, 1}
	pipeEast  = pipeDirection{1, 0}
	pipeWest  = pipeDirection{-1, 0}
)

var pipeEnds = map[byte][2]pipeDirection{
	'|': {pipeNorth, pipeSouth},
	'-': {pipeEast, pipeWest},
	'L': {pipeNorth, pipeEast},
	'J': {pipeNorth, pipeWest},
	'7': {pipeSouth, pipeWest},
	'F': {pipeSouth, pipeEast},
}

func (d pipeDirection) reverse() pipeDirection {
	return pipeDirection{-d.dx, -d.dy}
}

func pipeOpensTo(rows [][]byte, x, y int, d pipeDirection) bool {
	if y < 0 || y >= len(rows) || x < 0 || x >= len(rows[y]) {
		return false
	}
	ends, ok := pipeEnds[rows[y][x]]
	return ok && (ends[0] == d || ends[1] == d)
}

// pipeLoop returns the number of tiles in the loop through S and the area the
// loop's tile centres enclose.
func pipeLoop(r io.Reader) (length, area int, err error) {
	g, err := grid.ReadBytes(r, "|-LJ7F.S")
	if err != nil {
		return 0, 0, err
	}
	start, err := g.Find('S')
	if err != nil {
		return 0, 0, err
	}
	rows, startX, startY := g.Rows(), start.X, start.Y

	var exits []pipeDirection
	for _, d := range []pipeDirection{pipeNorth, pipeSouth, pipeEast, pipeWest} {
		if pipeOpensTo(rows, startX+d.dx, startY+d.dy, d.reverse()) {
			exits = append(exits, d)
		}
	}
	if len(exits) != 2 {
		return 0, 0, fmt.Errorf("start at %d,%d connects to %d pipes, want 2", startX, startY, len(exits))
	}

	x, y, heading := startX, startY, exits[0]
	twiceArea := 0
	for {
		nextX, nextY := x+heading.dx, y+heading.dy
		twiceArea += x*nextY - nextX*y
		x, y = nextX, nextY
		length++
		if x == startX && y == startY {
			break
		}
		if !pipeOpensTo(rows, x, y, heading.reverse()) {
			return 0, 0, fmt.Errorf("loop breaks at %d,%d", x, y)
		}
		ends := pipeEnds[rows[y][x]]
		if ends[0] == heading.reverse() {
			heading = ends[1]
		} else {
			heading = ends[0]
		}
	}
	return length, max(twiceArea, -twiceArea) / 2, nil
}

func day10p01(r io.Reader) (string, error) {
	length, _, err := pipeLoop(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(length / 2), nil
}

func day10p02(r io.Reader) (string, error) {
	length, area, err := pipeLoop(r)
	if err != nil {
		return "", err
	}
	// Pick's theorem: interior points = area - boundary/2 + 1.
	return strconv.Itoa(area - length/2 + 1), nil
}
