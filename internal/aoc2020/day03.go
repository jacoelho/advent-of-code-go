package aoc2020

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func treesOnSlope(treeMap grid.Bytes, right, down int) int64 {
	var trees int64
	for x, y := right, down; y < treeMap.Height(); x, y = x+right, y+down {
		if treeMap.At(grid.Position2D[int]{X: x % treeMap.Width(), Y: y}) == '#' {
			trees++
		}
	}
	return trees
}

func day03p01(r io.Reader) (string, error) {
	treeMap, err := grid.ReadBytes(r, ".#")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(treesOnSlope(treeMap, 3, 1), 10), nil
}

func day03p02(r io.Reader) (string, error) {
	treeMap, err := grid.ReadBytes(r, ".#")
	if err != nil {
		return "", err
	}
	product := int64(1)
	for _, slope := range [][2]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}} {
		product *= treesOnSlope(treeMap, slope[0], slope[1])
	}
	return strconv.FormatInt(product, 10), nil
}
