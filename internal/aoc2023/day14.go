package aoc2023

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type platform struct {
	width, height int
	cells         []byte
}

func parsePlatform(r io.Reader) (*platform, error) {
	g, err := grid.ReadBytes(r, "O#.")
	if err != nil {
		return nil, err
	}
	return &platform{width: g.Width(), height: g.Height(), cells: g.Cells()}, nil
}

// roll moves every round rock towards the start of its lane until it meets a
// rock. Lane k starts at cell first+k*laneStep and its cells follow by
// cellStep.
func (p *platform) roll(lanes, cells, first, laneStep, cellStep int) {
	for lane := range lanes {
		i := first + lane*laneStep
		free := i
		for range cells {
			switch p.cells[i] {
			case '#':
				free = i + cellStep
			case 'O':
				p.cells[i] = '.'
				p.cells[free] = 'O'
				free += cellStep
			}
			i += cellStep
		}
	}
}

func (p *platform) tiltNorth() {
	p.roll(p.width, p.height, 0, 1, p.width)
}

func (p *platform) spinCycle() {
	w, h := p.width, p.height
	p.tiltNorth()
	p.roll(h, w, 0, w, 1)
	p.roll(w, h, (h-1)*w, 1, -w)
	p.roll(h, w, w-1, w, -1)
}

func (p *platform) northLoad() int {
	load := 0
	for i, c := range p.cells {
		if c == 'O' {
			load += p.height - i/p.width
		}
	}
	return load
}

func day14p01(r io.Reader) (string, error) {
	p, err := parsePlatform(r)
	if err != nil {
		return "", err
	}
	p.tiltNorth()
	return strconv.Itoa(p.northLoad()), nil
}

func day14p02(r io.Reader) (string, error) {
	p, err := parsePlatform(r)
	if err != nil {
		return "", err
	}
	const cycles = 1_000_000_000
	seen := make(map[string]int)
	for i := 0; i < cycles; i++ {
		p.spinCycle()
		key := string(p.cells)
		if previous, found := seen[key]; found {
			for range (cycles - i - 1) % (i - previous) {
				p.spinCycle()
			}
			break
		}
		seen[key] = i
	}
	return strconv.Itoa(p.northLoad()), nil
}
