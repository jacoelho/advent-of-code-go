package aoc2017

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// infectionState values follow the evolved cycle, so adding state+3 to a
// clockwise direction index turns left, keeps course, turns right or reverses.
type infectionState byte

const (
	clean infectionState = iota
	weakened
	infected
	flagged
)

// virusMode is how far one burst advances a node through the infection cycle.
type virusMode byte

const (
	evolvedVirus virusMode = 1
	simpleVirus  virusMode = 2
)

// infectionGrid is a square of nodes that grows around its centre when the carrier leaves it.
type infectionGrid struct {
	side  int
	nodes []infectionState
}

func (g *infectionGrid) grow() {
	side := 3 * g.side
	nodes := make([]infectionState, side*side)
	for y := range g.side {
		copy(nodes[(y+g.side)*side+g.side:], g.nodes[y*g.side:(y+1)*g.side])
	}
	g.side, g.nodes = side, nodes
}

func parseInfections(r io.Reader) (infectionGrid, error) {
	cells, err := grid.ReadBytes(r, "#.")
	if err != nil {
		return infectionGrid{}, err
	}
	width, height := cells.Width(), cells.Height()
	if width%2 != 1 || height%2 != 1 {
		return infectionGrid{}, fmt.Errorf("infection grid must have odd width and height, got %dx%d", width, height)
	}

	g := infectionGrid{side: max(width, height)}
	g.nodes = make([]infectionState, g.side*g.side)
	left, top := (g.side-width)/2, (g.side-height)/2
	for i, cell := range cells.Cells() {
		if cell == '#' {
			p := cells.Position(i)
			g.nodes[(top+p.Y)*g.side+left+p.X] = infected
		}
	}

	return g, nil
}

func virusBursts(r io.Reader, mode virusMode, bursts int) (int, error) {
	g, err := parseInfections(r)
	if err != nil {
		return 0, err
	}

	offsets := [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	x, y, dir := g.side/2, g.side/2, 0
	infections := 0
	for range bursts {
		if x < 0 || y < 0 || x >= g.side || y >= g.side {
			x, y = x+g.side, y+g.side
			g.grow()
		}
		node := &g.nodes[y*g.side+x]
		dir = (dir + int(*node) + 3) & 3
		*node = (*node + infectionState(mode)) & 3
		if *node == infected {
			infections++
		}
		x, y = x+offsets[dir][0], y+offsets[dir][1]
	}

	return infections, nil
}

func day22p01(r io.Reader) (string, error) {
	count, err := virusBursts(r, simpleVirus, 10_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}

func day22p02(r io.Reader) (string, error) {
	count, err := virusBursts(r, evolvedVirus, 10_000_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}
