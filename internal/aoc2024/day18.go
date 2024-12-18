package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

func parsePushDownAutomatonMemory(r io.Reader) ([]grid.Position2D[int], error) {
	s := scanner.NewScanner(r, func(bytes []byte) (grid.Position2D[int], error) {
		digits := convert.ExtractDigits[int](string(bytes))
		return grid.Position2D[int]{
			X: digits[0],
			Y: digits[1],
		}, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func day18Neighbours(memory map[grid.Position2D[int]]bool, dimensions int) func(grid.Position2D[int]) []grid.Position2D[int] {
	return func(p grid.Position2D[int]) []grid.Position2D[int] {
		return slices.Collect(xiter.Filter(func(v grid.Position2D[int]) bool {
			return v.X >= 0 && v.X <= dimensions && v.Y >= 0 && v.Y <= dimensions && !memory[v]
		}, grid.Neighbours4(p)))
	}
}

func day18Heuristic(dimensions int) func(p grid.Position2D[int]) int {
	return func(p grid.Position2D[int]) int {
		if p.X == dimensions && p.Y == dimensions {
			return 0
		}
		return 1
	}
}

func day18search(
	corrupted []grid.Position2D[int],
	dimensions int,
	start grid.Position2D[int],
	steps int,
) (int, bool) {
	memory := make(map[grid.Position2D[int]]bool, steps)
	for _, v := range corrupted[:steps] {
		memory[v] = true
	}
	cost, _, found := search.AStar(
		start,
		day18Neighbours(memory, dimensions),
		day18Heuristic(dimensions),
		search.ConstantStepCost,
	)
	return cost, found
}

func day18p01(dimensions, steps int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted := aoc.Must(parsePushDownAutomatonMemory(r))

		start := grid.Position2D[int]{X: 0, Y: 0}

		cost, _ := day18search(corrupted, dimensions, start, steps)

		return strconv.Itoa(cost), nil
	}
}

func day18p02(dimensions int, steps int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted := aoc.Must(parsePushDownAutomatonMemory(r))
		start := grid.Position2D[int]{X: 0, Y: 0}

		low, high := steps, len(corrupted)-1
		for low < high {
			mid := (low + high) / 2
			if _, found := day18search(corrupted, dimensions, start, mid+1); found {
				low = mid + 1
			} else {
				high = mid
			}
		}

		p := corrupted[low]
		return fmt.Sprintf("%d,%d", p.X, p.Y), nil
	}
}
