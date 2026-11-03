package aoc2023

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type beam struct {
	pos position
	dir position
}

var (
	right = position{X: 1, Y: 0}
	left  = position{X: -1, Y: 0}
	up    = position{X: 0, Y: -1}
	down  = position{X: 0, Y: 1}
)

func parseContraption(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	return grid.NewGrid2D[int](rows), nil
}

func reflect(dir position, mirror rune) position {
	if mirror == '/' {
		// right -> up, up -> right, left -> down, down -> left
		return position{X: -dir.Y, Y: -dir.X}
	}
	// mirror == '\\'
	// right -> down, down -> right, left -> up, up -> left
	return position{X: dir.Y, Y: dir.X}
}

func nextBeam(pos, dir position) beam {
	return beam{pos: pos.Add(dir), dir: dir}
}

func simulate(g grid.Grid2D[int, rune], start beam) int {
	children := func(current beam) iter.Seq[beam] {
		return func(yield func(beam) bool) {
			switch tile := g[current.pos]; tile {
			case '.':
				yield(nextBeam(current.pos, current.dir))
			case '/', '\\':
				yield(nextBeam(current.pos, reflect(current.dir, tile)))
			case '|':
				if current.dir == left || current.dir == right {
					if yield(nextBeam(current.pos, up)) {
						yield(nextBeam(current.pos, down))
					}
				} else {
					yield(nextBeam(current.pos, current.dir))
				}
			case '-':
				if current.dir == up || current.dir == down {
					if yield(nextBeam(current.pos, left)) {
						yield(nextBeam(current.pos, right))
					}
				} else {
					yield(nextBeam(current.pos, current.dir))
				}
			}
		}
	}
	energized := collections.NewSet[position]()
	for current := range search.DFS(start, children) {
		if g.Contains(current.pos) {
			energized.Add(current.pos)
		}
	}
	return energized.Len()
}

func edgeBeams(g grid.Grid2D[int, rune]) iter.Seq[beam] {
	return func(yield func(beam) bool) {
		minX, maxX, minY, maxY := g.Dimensions()

		for x := minX; x <= maxX; x++ {
			if !yield(beam{pos: position{X: x, Y: minY}, dir: down}) {
				return
			}
		}

		for x := minX; x <= maxX; x++ {
			if !yield(beam{pos: position{X: x, Y: maxY}, dir: up}) {
				return
			}
		}

		for y := minY; y <= maxY; y++ {
			if !yield(beam{pos: position{X: minX, Y: y}, dir: right}) {
				return
			}
		}

		for y := minY; y <= maxY; y++ {
			if !yield(beam{pos: position{X: maxX, Y: y}, dir: left}) {
				return
			}
		}
	}
}

func day16p01(r io.Reader) (string, error) {
	g, err := parseContraption(r)
	if err != nil {
		return "", err
	}

	start := beam{pos: position{X: 0, Y: 0}, dir: right}
	count := simulate(g, start)

	return strconv.Itoa(count), nil
}

func day16p02(r io.Reader) (string, error) {
	g, err := parseContraption(r)
	if err != nil {
		return "", err
	}

	maxEnergized := xiter.From(edgeBeams(g)).Map(func(b beam) int { return simulate(g, b) }).Reduce(0, func(maxSoFar, count int) int { return max(maxSoFar, count) })

	return strconv.Itoa(maxEnergized), nil
}
