package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseTopographicMap(r io.Reader) (grid.Grid2D[int, int], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]int, error) {
		res := make([]int, len(b))
		for i, v := range b {
			res[i] = int(v - '0')
		}
		return res, nil
	})
	return grid.NewGrid2D[int, int](slices.Collect(s.Values())), s.Err()
}

func day10neighbours(m grid.Grid2D[int, int]) func(grid.Position2D[int], func(grid.Position2D[int]) bool) {
	return func(p grid.Position2D[int], yield func(grid.Position2D[int]) bool) {
		slope := m[p]
		for _, offset := range grid.OffsetsNeighbours4[int]() {
			v := p.Add(offset)
			if height, exists := m[v]; exists && height == slope+1 && !yield(v) {
				return
			}
		}
	}
}

func trailheads(m grid.Grid2D[int, int]) []grid.Position2D[int] {
	var trailheads []grid.Position2D[int]
	for k, slope := range m {
		if slope == 0 {
			trailheads = append(trailheads, k)
		}
	}
	return trailheads
}

func day10p01(r io.Reader) (string, error) {
	m, err := parseTopographicMap(r)
	if err != nil {
		return "", err
	}

	visit := day10neighbours(m)
	neighbours := func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) { visit(p, yield) }
	}
	var summits int
	for _, trail := range trailheads(m) {
		for position := range search.BFS(trail, neighbours) {
			if m[position] == 9 {
				summits++
			}
		}
	}

	return strconv.Itoa(summits), nil
}

func day10p02(r io.Reader) (string, error) {
	m, err := parseTopographicMap(r)
	if err != nil {
		return "", err
	}

	summits, err := graph.CountDAGPaths(trailheads(m), day10neighbours(m), func(position grid.Position2D[int]) (int, bool) {
		return 1, m[position] == 9
	})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(summits), nil
}
