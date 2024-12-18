package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

func parsePushDownAutomatonMemory(r io.Reader) (iter.Seq[grid.Position2D[int]], error) {
	s := scanner.NewScanner(r, func(bytes []byte) (grid.Position2D[int], error) {
		digits := convert.ExtractDigits[int](string(bytes))
		return grid.Position2D[int]{
			X: digits[0],
			Y: digits[1],
		}, nil
	})
	return s.Values(), s.Err()
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

func day18p01(dimensions, steps int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted := aoc.Must(parsePushDownAutomatonMemory(r))
		start := grid.Position2D[int]{X: 0, Y: 0}

		memory := make(map[grid.Position2D[int]]bool)
		for p := range xiter.Take(corrupted, steps) {
			memory[p] = true
		}

		cost, _, _ := search.AStar(start, day18Neighbours(
			memory,
			dimensions,
		), day18Heuristic(dimensions), search.ConstantStepCost)

		return strconv.Itoa(cost), nil
	}
}

func day18p02(dimensions int, steps int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		corrupted := aoc.Must(parsePushDownAutomatonMemory(r))
		start := grid.Position2D[int]{X: 0, Y: 0}

		memory := make(map[grid.Position2D[int]]bool)
		for p := range xiter.Take(corrupted, steps) {
			memory[p] = true
		}

		for p := range corrupted {
			memory[p] = true

			if _, _, found := search.AStar(
				start,
				day18Neighbours(
					memory,
					dimensions,
				),
				day18Heuristic(dimensions),
				search.ConstantStepCost,
			); !found {
				return fmt.Sprintf("%d,%d", p.X, p.Y), nil
			}
		}
		return "", errors.New("not found")
	}
}
