package aoc2024

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseWordSearch(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		res := make([]rune, len(b))
		for i, v := range b {
			res[i] = rune(v)
		}
		return res, nil
	})
	return grid.NewGrid2D[int, rune](slices.Collect(s.Values())), s.Err()
}

func day04p01(r io.Reader) (string, error) {
	m, err := parseWordSearch(r)
	if err != nil {
		return "", err
	}

	var count int
	for p, v := range m {
		if v != 'X' {
			continue
		}
		for _, offset := range grid.OffsetsNeighbours8[int]() {
			mPosition := p.Add(offset)
			aPosition := p.Add(grid.Position2D[int]{X: offset.X * 2, Y: offset.Y * 2})
			sPosition := p.Add(grid.Position2D[int]{X: offset.X * 3, Y: offset.Y * 3})

			if m[mPosition] == 'M' && m[aPosition] == 'A' && m[sPosition] == 'S' {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

func day04p02(r io.Reader) (string, error) {
	m, err := parseWordSearch(r)
	if err != nil {
		return "", err
	}

	// Keep opposite corners two entries apart for the pattern comparisons below.
	corners := [4]grid.Position2D[int]{
		{X: -1, Y: -1},
		{X: -1, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: -1},
	}

	var count int
	for p, v := range m {
		if v != 'A' {
			continue
		}

		var values [4]rune
		for i, offset := range corners {
			values[i] = m[p.Add(offset)]
		}
		neighbours := string(values[:])

		if neighbours == "MMSS" ||
			neighbours == "MSSM" ||
			neighbours == "SSMM" ||
			neighbours == "SMMS" {
			count++
		}
	}

	return strconv.Itoa(count), nil
}
