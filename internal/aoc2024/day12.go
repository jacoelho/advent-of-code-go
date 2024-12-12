package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
)

func parseGardenPlots(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		res := make([]rune, len(b))
		for i, v := range b {
			res[i] = rune(v)
		}
		return res, nil
	})
	return grid.NewGrid2D[int, rune](slices.Collect(s.Values())), s.Err()
}

func day12neighbours(plots grid.Grid2D[int, rune], plant rune) func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(p grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for neighbour := range grid.Neighbours4(p) {
				if plots[neighbour] == plant {
					if !yield(neighbour) {
						return
					}
				}
			}
		}
	}
}

func perimeter(region collections.Set[grid.Position2D[int]]) int {
	result := region.Len() * 4
	for el := range region.Iter() {
		for n := range grid.Neighbours4(el) {
			if region.Contains(n) {
				result -= 1
			}
		}
	}
	return result
}

func day12p01(r io.Reader) (string, error) {
	plots := aoc.Must(parseGardenPlots(r))

	visited := collections.NewSet[grid.Position2D[int]]()

	var regions []collections.Set[grid.Position2D[int]]
	for position, plant := range plots {
		neighbours := day12neighbours(plots, plant)

		if !visited.Contains(position) {
			region := slices.Collect(search.BFS(position, neighbours))
			regions = append(regions, collections.NewSet(region...))
			visited.Add(region...)
		}
	}

	var total int
	for _, region := range regions {
		total += perimeter(region) * region.Len()
	}

	return strconv.Itoa(total), nil
}

func corners(region collections.Set[grid.Position2D[int]]) int {
	var count int

	for _, offset := range grid.OffsetsNeighbours4[int]() {
		sides := collections.NewSet[grid.Position2D[int]]()

		for p := range region.Iter() {
			n := p.Add(offset)
			if !region.Contains(n) {
				sides.Add(n)
			}
		}

		remove := collections.NewSet[grid.Position2D[int]]()
		for side := range sides.Iter() {
			p := side.Add(grid.Position2D[int]{X: offset.Y, Y: offset.X})
			for sides.Contains(p) {
				remove.Add(p)
				p = p.Add(grid.Position2D[int]{X: offset.Y, Y: offset.X})
			}
		}
		count += sides.Len() - remove.Len()
	}

	return count
}

func day12p02(r io.Reader) (string, error) {
	plots := aoc.Must(parseGardenPlots(r))

	visited := collections.NewSet[grid.Position2D[int]]()

	var regions []collections.Set[grid.Position2D[int]]
	for position, plant := range plots {
		neighbours := day12neighbours(plots, plant)

		if !visited.Contains(position) {
			region := slices.Collect(search.BFS(position, neighbours))
			regions = append(regions, collections.NewSet(region...))
			visited.Add(region...)
		}
	}

	var total int
	for _, region := range regions {
		total += corners(region) * region.Len()
	}

	return strconv.Itoa(total), nil
}
