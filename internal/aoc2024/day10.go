package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
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

func day10neighbours(m grid.Grid2D[int, int]) func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			slope := m[p]

			for v := range grid.Neighbours4(p) {
				if m[v] == slope+1 {
					if !yield(v) {
						return
					}
				}
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

	var summits int
	for _, trail := range trailheads(m) {
		for position := range search.BFS(trail, day10neighbours(m)) {
			if m[position] == 9 {
				summits++
			}
		}
	}

	return strconv.Itoa(summits), nil
}

func multiVisitBFS[T comparable](start T, neighbours func(T) iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		frontier := collections.NewDeque[T](10)
		frontier.PushBack(start)

		for frontier.Size() > 0 {
			node, ok := frontier.PopFront()
			if !ok || !yield(node) {
				return
			}

			for el := range neighbours(node) {
				frontier.PushBack(el)
			}
		}
	}
}

func day10p02(r io.Reader) (string, error) {
	m, err := parseTopographicMap(r)
	if err != nil {
		return "", err
	}

	var summits int
	for _, trail := range trailheads(m) {
		for position := range multiVisitBFS(trail, day10neighbours(m)) {
			if m[position] == 9 {
				summits++
			}
		}
	}

	return strconv.Itoa(summits), nil
}
