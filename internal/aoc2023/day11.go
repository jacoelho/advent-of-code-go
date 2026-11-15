package aoc2023

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// expandedDistanceSum returns the sum of pairwise distances along one axis,
// where counts[i] galaxies sit at line i and every empty line counts as
// expansion lines.
func expandedDistanceSum(counts []int, expansion int) int {
	sum, seen, coordinateSum, coordinate := 0, 0, 0, 0
	for _, count := range counts {
		if count == 0 {
			coordinate += expansion
			continue
		}
		sum += count * (seen*coordinate - coordinateSum)
		seen += count
		coordinateSum += count * coordinate
		coordinate++
	}
	return sum
}

func galaxyDistances(r io.Reader, expansion int) (string, error) {
	g, err := grid.ReadBytes(r, ".#")
	if err != nil {
		return "", err
	}
	rowCounts := make([]int, g.Height())
	columnCounts := make([]int, g.Width())
	for i, c := range g.Cells() {
		if c == '#' {
			p := g.Position(i)
			rowCounts[p.Y]++
			columnCounts[p.X]++
		}
	}
	total := expandedDistanceSum(rowCounts, expansion) + expandedDistanceSum(columnCounts, expansion)
	return strconv.Itoa(total), nil
}

func day11p01(r io.Reader) (string, error) {
	return galaxyDistances(r, 2)
}

func day11p02(r io.Reader) (string, error) {
	return galaxyDistances(r, 1000000)
}
