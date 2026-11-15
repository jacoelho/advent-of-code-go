package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// elfGrove keeps elves on a flat occupancy grid whose border is always at
// least one cell away from every elf, so neighbour indices stay in bounds.
type elfGrove struct {
	width, height int
	origin        grid.Position2D[int]
	occupied      []bool
	elves         []int
	proposals     []uint8
	destinations  []int
}

func parseElves(r io.Reader) (*elfGrove, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	rows := strings.Split(input, "\n")
	var positions []grid.Position2D[int]
	for y, row := range rows {
		if len(row) == 0 || len(row) != len(rows[0]) {
			return nil, fmt.Errorf("elf grove must be a nonempty rectangle")
		}
		for x := range len(row) {
			switch row[x] {
			case '#':
				positions = append(positions, grid.Position2D[int]{X: x, Y: y})
			case '.':
			default:
				return nil, fmt.Errorf("invalid elf grove cell %q", row[x])
			}
		}
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("elf grove has no elves")
	}
	grove := &elfGrove{}
	grove.place(positions)
	return grove, nil
}

func elfBounds(positions []grid.Position2D[int]) (low, high grid.Position2D[int]) {
	low, high = positions[0], positions[0]
	for _, p := range positions {
		low.X, low.Y = min(low.X, p.X), min(low.Y, p.Y)
		high.X, high.Y = max(high.X, p.X), max(high.Y, p.Y)
	}
	return low, high
}

func (g *elfGrove) place(positions []grid.Position2D[int]) {
	low, high := elfBounds(positions)
	margin := max(8, (high.X-low.X)/2, (high.Y-low.Y)/2)
	g.width = high.X - low.X + 1 + 2*margin
	g.height = high.Y - low.Y + 1 + 2*margin
	g.origin = grid.Position2D[int]{X: low.X - margin, Y: low.Y - margin}
	g.occupied = make([]bool, g.width*g.height)
	g.proposals = make([]uint8, g.width*g.height)
	g.elves = g.elves[:0]
	for _, p := range positions {
		index := (p.Y-g.origin.Y)*g.width + p.X - g.origin.X
		g.occupied[index] = true
		g.elves = append(g.elves, index)
	}
	g.destinations = make([]int, len(g.elves))
}

func (g *elfGrove) positions() []grid.Position2D[int] {
	positions := make([]grid.Position2D[int], len(g.elves))
	for i, index := range g.elves {
		positions[i] = grid.Position2D[int]{X: g.origin.X + index%g.width, Y: g.origin.Y + index/g.width}
	}
	return positions
}

func (g *elfGrove) touchesBorder() bool {
	for _, index := range g.elves {
		x, y := index%g.width, index/g.width
		if x == 0 || y == 0 || x == g.width-1 || y == g.height-1 {
			return true
		}
	}
	return false
}

func (g *elfGrove) round(round int) bool {
	if g.touchesBorder() {
		g.place(g.positions())
	}
	w := g.width
	// Each direction lists the cell moved into followed by its two diagonal flanks.
	directions := [4][3]int{
		{-w, -w - 1, -w + 1},
		{w, w - 1, w + 1},
		{-1, -w - 1, w - 1},
		{1, -w + 1, w + 1},
	}
	neighbours := [8]int{-w - 1, -w, -w + 1, -1, 1, w - 1, w, w + 1}
	for i, elf := range g.elves {
		g.destinations[i] = -1
		crowded := false
		for _, offset := range neighbours {
			if g.occupied[elf+offset] {
				crowded = true
				break
			}
		}
		if !crowded {
			continue
		}
		for k := range 4 {
			direction := directions[(round+k)%4]
			if !g.occupied[elf+direction[0]] && !g.occupied[elf+direction[1]] && !g.occupied[elf+direction[2]] {
				destination := elf + direction[0]
				g.destinations[i] = destination
				g.proposals[destination]++
				break
			}
		}
	}
	moved := false
	for i, destination := range g.destinations {
		if destination < 0 {
			continue
		}
		// Every destination was empty at the start of the round, so moves cannot block each other.
		if g.proposals[destination] == 1 {
			g.occupied[g.elves[i]] = false
			g.occupied[destination] = true
			g.elves[i] = destination
			moved = true
		}
		g.proposals[destination] = 0
	}
	return moved
}

func day23p01(r io.Reader) (string, error) {
	grove, err := parseElves(r)
	if err != nil {
		return "", err
	}
	for round := range 10 {
		grove.round(round)
	}
	positions := grove.positions()
	low, high := elfBounds(positions)
	return strconv.Itoa((high.X-low.X+1)*(high.Y-low.Y+1) - len(positions)), nil
}

func day23p02(r io.Reader) (string, error) {
	grove, err := parseElves(r)
	if err != nil {
		return "", err
	}
	round := 0
	for grove.round(round) {
		round++
	}
	return strconv.Itoa(round + 1), nil
}
