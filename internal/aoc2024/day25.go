package aoc2024

import (
	"bytes"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseSchematic(b []byte) (collections.Set[grid.Position2D[int]], error) {
	result := collections.NewSet[grid.Position2D[int]]()
	for y, line := range bytes.Split(b, []byte("\n")) {
		for x, v := range line {
			if v == '#' {
				result.Add(grid.NewPosition2D(x, y))
			}
		}
	}
	return result, nil
}

func parseSchematics(r io.Reader) ([]collections.Set[grid.Position2D[int]], error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte{'\n', '\n'}), parseSchematic)
	return slices.Collect(s.Values()), s.Err()
}

func day25p01(r io.Reader) (string, error) {
	var total int
	for pair := range xslices.Pairwise(aoc.Must(parseSchematics(r))) {
		if overlap := pair.V1.Intersect(pair.V2); overlap.IsEmpty() {
			total++
		}
	}
	return strconv.Itoa(total), nil
}
