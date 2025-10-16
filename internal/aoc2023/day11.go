package aoc2023

import (
	"bufio"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseGalaxies(r io.Reader) ([]grid.Position2D[int], collections.Set[int], collections.Set[int], error) {
	scanner := bufio.NewScanner(r)
	var galaxies []grid.Position2D[int]
	rowsWithGalaxies := collections.NewSet[int]()
	colsWithGalaxies := collections.NewSet[int]()

	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		for col, ch := range line {
			if ch == '#' {
				galaxies = append(galaxies, grid.NewPosition2D(col, row))
				rowsWithGalaxies.Add(row)
				colsWithGalaxies.Add(col)
			}
		}
		row++
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, nil, err
	}

	return galaxies, rowsWithGalaxies, colsWithGalaxies, nil
}

func countEmptyBetween(start, end int, occupied collections.Set[int]) int {
	min, max := start, end
	if start > end {
		min, max = end, start
	}
	count := 0
	for i := min + 1; i < max; i++ {
		if !occupied.Contains(i) {
			count++
		}
	}
	return count
}

func day11(expansionFactor int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		galaxies, rowsWithGalaxies, colsWithGalaxies, err := parseGalaxies(r)
		if err != nil {
			return "", err
		}

		result := xiter.Sum(xiter.Map(func(pair xslices.Pair[grid.Position2D[int], grid.Position2D[int]]) int {
			g1, g2 := pair.V1, pair.V2

			baseDistance := g1.Distance(g2)
			emptyRows := countEmptyBetween(g1.Y, g2.Y, rowsWithGalaxies)
			emptyCols := countEmptyBetween(g1.X, g2.X, colsWithGalaxies)

			return baseDistance + (expansionFactor-1)*emptyRows + (expansionFactor-1)*emptyCols
		}, xslices.Pairwise(galaxies)))

		return strconv.Itoa(result), nil
	}
}

func day11p01(r io.Reader) (string, error) {
	return day11(2)(r)
}

func day11p02(r io.Reader) (string, error) {
	return day11(1000000)(r)
}
