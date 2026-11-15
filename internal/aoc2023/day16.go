package aoc2023

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// Beam headings in clockwise order: right, down, left, up.
var (
	beamDX = [4]int{1, 0, -1, 0}
	beamDY = [4]int{0, 1, 0, -1}
)

// contraption reuses its visit marks across beams: a mark counts only when it
// equals the current generation.
type contraption struct {
	width, height int
	tiles         []byte
	seenBeam      []uint32
	seenTile      []uint32
	generation    uint32
	stack         []int
}

func parseContraption(r io.Reader) (*contraption, error) {
	g, err := grid.ReadBytes(r, `.|-/\`)
	if err != nil {
		return nil, err
	}
	c := &contraption{width: g.Width(), height: g.Height(), tiles: g.Cells()}
	c.seenBeam = make([]uint32, 4*len(c.tiles))
	c.seenTile = make([]uint32, len(c.tiles))
	return c, nil
}

// energise returns how many tiles a beam entering x,y with heading lights.
func (c *contraption) energise(x, y, heading int) int {
	c.generation++
	energised := 0
	c.stack = c.stack[:0]
	visit := func(x, y, heading int) {
		if x < 0 || x >= c.width || y < 0 || y >= c.height {
			return
		}
		beam := 4*(y*c.width+x) + heading
		if c.seenBeam[beam] == c.generation {
			return
		}
		c.seenBeam[beam] = c.generation
		c.stack = append(c.stack, beam)
	}
	visit(x, y, heading)

	for len(c.stack) > 0 {
		beam := c.stack[len(c.stack)-1]
		c.stack = c.stack[:len(c.stack)-1]
		tile, heading := beam/4, beam%4
		x, y := tile%c.width, tile/c.width
		if c.seenTile[tile] != c.generation {
			c.seenTile[tile] = c.generation
			energised++
		}
		next := [2]int{heading, -1}
		switch c.tiles[tile] {
		case '/':
			next[0] = heading ^ 3
		case '\\':
			next[0] = heading ^ 1
		case '|':
			if heading%2 == 0 {
				next = [2]int{1, 3}
			}
		case '-':
			if heading%2 == 1 {
				next = [2]int{0, 2}
			}
		}
		for _, d := range next {
			if d >= 0 {
				visit(x+beamDX[d], y+beamDY[d], d)
			}
		}
	}
	return energised
}

func day16p01(r io.Reader) (string, error) {
	c, err := parseContraption(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(c.energise(0, 0, 0)), nil
}

func day16p02(r io.Reader) (string, error) {
	c, err := parseContraption(r)
	if err != nil {
		return "", err
	}
	best := 0
	for x := range c.width {
		best = max(best, c.energise(x, 0, 1), c.energise(x, c.height-1, 3))
	}
	for y := range c.height {
		best = max(best, c.energise(0, y, 0), c.energise(c.width-1, y, 2))
	}
	return strconv.Itoa(best), nil
}
