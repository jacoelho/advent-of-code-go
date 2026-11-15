package aoc2024

import (
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func day10neighbours(heights grid.Bytes) func(grid.Position2D[int], func(grid.Position2D[int]) bool) {
	return func(p grid.Position2D[int], yield func(grid.Position2D[int]) bool) {
		uphill := heights.At(p) + 1
		for v := range heights.Neighbours4(p) {
			if heights.At(v) == uphill && !yield(v) {
				return
			}
		}
	}
}

func trailheads(heights grid.Bytes) []grid.Position2D[int] {
	var trailheads []grid.Position2D[int]
	for i, height := range heights.Cells() {
		if height == '0' {
			trailheads = append(trailheads, heights.Position(i))
		}
	}
	return trailheads
}

func day10p01(r io.Reader) (string, error) {
	heights, err := grid.ReadBytes(r, "0123456789")
	if err != nil {
		return "", err
	}

	visit := day10neighbours(heights)
	neighbours := func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) { visit(p, yield) }
	}
	var summits int
	for _, trail := range trailheads(heights) {
		for position := range search.BFS(trail, neighbours) {
			if heights.At(position) == '9' {
				summits++
			}
		}
	}

	return strconv.Itoa(summits), nil
}

func day10p02(r io.Reader) (string, error) {
	heights, err := grid.ReadBytes(r, "0123456789")
	if err != nil {
		return "", err
	}

	summits, err := graph.CountDAGPaths(trailheads(heights), day10neighbours(heights), func(position grid.Position2D[int]) (int, bool) {
		return 1, heights.At(position) == '9'
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(summits), nil
}
