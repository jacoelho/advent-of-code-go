package aoc2024

import (
	"bufio"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type frequencyMap struct {
	maxX, maxY int
	antennas   map[rune][]grid.Position2D[int]
}

func (m frequencyMap) isWithinBounds(p grid.Position2D[int]) bool {
	return p.X >= 0 && p.X < m.maxX && p.Y >= 0 && p.Y < m.maxY
}

func parseFrequencyMap(r io.Reader) (frequencyMap, error) {
	result := frequencyMap{
		antennas: make(map[rune][]grid.Position2D[int]),
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		for x, v := range s.Text() {
			if v == '.' {
				continue
			}
			result.antennas[v] = append(
				result.antennas[v],
				grid.Position2D[int]{X: x, Y: result.maxY},
			)

			result.maxX = max(result.maxX, len(s.Text()))
		}
		result.maxY++
	}
	return result, s.Err()
}

func day08p01(r io.Reader) (string, error) {
	m := aoc.Must(parseFrequencyMap(r))

	antiNodes := collections.NewSet[grid.Position2D[int]]()
	for _, antenna := range m.antennas {
		for pair := range xslices.Pairwise(antenna) {
			a1 := grid.Position2D[int]{
				// https://www.andreaminini.net/math/central-symmetry
				X: 2*pair.V1.X - pair.V2.X,
				Y: 2*pair.V1.Y - pair.V2.Y,
			}
			if m.isWithinBounds(a1) {
				antiNodes.Add(a1)
			}

			a2 := grid.Position2D[int]{
				X: 2*pair.V2.X - pair.V1.X,
				Y: 2*pair.V2.Y - pair.V1.Y,
			}
			if m.isWithinBounds(a2) {
				antiNodes.Add(a2)
			}
		}
	}

	return strconv.Itoa(antiNodes.Len()), nil
}

func day08p02(r io.Reader) (string, error) {
	m := aoc.Must(parseFrequencyMap(r))

	antiNodes := collections.NewSet[grid.Position2D[int]]()
	for _, antenna := range m.antennas {
		for pair := range xslices.Pairwise(antenna) {
			dX := pair.V2.X - pair.V1.X
			dY := pair.V2.Y - pair.V1.Y

			p := pair.V1
			for m.isWithinBounds(p) {
				antiNodes.Add(p)
				p.X += dX
				p.Y += dY
			}
			p = pair.V1
			for m.isWithinBounds(p) {
				antiNodes.Add(p)
				p.X -= dX
				p.Y -= dY
			}
		}
	}

	return strconv.Itoa(antiNodes.Len()), nil
}
