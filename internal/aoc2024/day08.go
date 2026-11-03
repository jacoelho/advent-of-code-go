package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

const maxAntinodeWork = 1_000_000

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
		line := s.Text()
		if result.maxY == 0 {
			result.maxX = len(line)
		}
		if len(line) == 0 || len(line) != result.maxX {
			return frequencyMap{}, fmt.Errorf("antenna map must have nonempty rows of equal width")
		}
		if (result.maxY+1)*result.maxX > 1_000_000 {
			return frequencyMap{}, fmt.Errorf("antenna map exceeds one million cells")
		}
		for x, v := range line {
			if v == '.' {
				continue
			}
			if !(v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' || v >= '0' && v <= '9') {
				return frequencyMap{}, fmt.Errorf("invalid antenna frequency %q", v)
			}
			result.antennas[v] = append(
				result.antennas[v],
				grid.Position2D[int]{X: x, Y: result.maxY},
			)
		}
		result.maxY++
	}
	if err := s.Err(); err != nil {
		return frequencyMap{}, err
	}
	if result.maxY == 0 {
		return frequencyMap{}, fmt.Errorf("empty antenna map")
	}
	return result, nil
}

func day08p01(r io.Reader) (string, error) {
	m, err := parseFrequencyMap(r)
	if err != nil {
		return "", err
	}

	antiNodes := collections.NewSet[grid.Position2D[int]]()
	pairs := 0
	for _, antenna := range m.antennas {
		for pair := range xslices.Pairwise(antenna) {
			pairs++
			if pairs > maxAntinodeWork {
				return "", fmt.Errorf("antenna search work limit exceeded")
			}
			a1 := grid.Position2D[int]{
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
	m, err := parseFrequencyMap(r)
	if err != nil {
		return "", err
	}

	antiNodes := collections.NewSet[grid.Position2D[int]]()
	work := 0
	for _, antenna := range m.antennas {
		for pair := range xslices.Pairwise(antenna) {
			work++
			if work > maxAntinodeWork {
				return "", fmt.Errorf("antenna search work limit exceeded")
			}
			dX := pair.V2.X - pair.V1.X
			dY := pair.V2.Y - pair.V1.Y
			divisor := xmath.GCD(dX, dY)
			dX, dY = dX/divisor, dY/divisor

			for _, direction := range [2]int{-1, 1} {
				p := pair.V1
				for m.isWithinBounds(p) {
					work++
					if work > maxAntinodeWork {
						return "", fmt.Errorf("antenna search work limit exceeded")
					}
					antiNodes.Add(p)
					p.X += direction * dX
					p.Y += direction * dY
				}
			}
		}
	}

	return strconv.Itoa(antiNodes.Len()), nil
}
