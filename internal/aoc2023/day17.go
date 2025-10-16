package aoc2023

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type crucibleState struct {
	position  grid.Position2D[int]
	direction grid.Position2D[int]
	steps     int
}

func crucibleNeighbours(heatMap grid.Grid2D[int, int], minSteps, maxSteps int) func(crucibleState) []crucibleState {
	return func(s crucibleState) []crucibleState {
		var neighbours []crucibleState

		if s.steps < maxSteps {
			nextPos := s.position.Add(s.direction)
			if heatMap.Contains(nextPos) {
				neighbours = append(neighbours, crucibleState{
					position:  nextPos,
					direction: s.direction,
					steps:     s.steps + 1,
				})
			}
		}

		if s.steps >= minSteps {
			for _, turnDir := range []grid.Position2D[int]{s.direction.TurnLeft(), s.direction.TurnRight()} {
				nextPos := s.position.Add(turnDir)
				if heatMap.Contains(nextPos) {
					neighbours = append(neighbours, crucibleState{
						position:  nextPos,
						direction: turnDir,
						steps:     1,
					})
				}
			}
		}

		return neighbours
	}
}

func crucibleHeuristic(target grid.Position2D[int], minSteps int) func(crucibleState) int {
	return func(s crucibleState) int {
		if s.position == target && s.steps >= minSteps {
			return 0
		}
		return 1
	}
}

func crucibleStepCost(heatMap grid.Grid2D[int, int]) func(crucibleState, crucibleState) int {
	return func(from, to crucibleState) int {
		return heatMap[to.position]
	}
}

func parseHeatGrid(r io.Reader) (grid.Grid2D[int, int], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]int, error) {
		row := make([]int, len(line))
		for i, ch := range line {
			row[i] = int(ch - '0')
		}
		return row, nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	return grid.NewGrid2D[int](rows), nil
}

func day17p01(r io.Reader) (string, error) {
	heatMap, err := parseHeatGrid(r)
	if err != nil {
		return "", err
	}

	_, maxX, _, maxY := heatMap.Dimensions()
	target := grid.Position2D[int]{X: maxX, Y: maxY}

	minHeat := -1
	for _, startDir := range []grid.Position2D[int]{{X: 1, Y: 0}, {X: 0, Y: 1}} {
		start := crucibleState{
			position:  grid.Position2D[int]{X: 0, Y: 0},
			direction: startDir,
			steps:     0,
		}

		heat, _, found := search.AStar(
			start,
			crucibleNeighbours(heatMap, 0, 3),
			crucibleHeuristic(target, 0),
			crucibleStepCost(heatMap),
		)

		if found && (minHeat == -1 || heat < minHeat) {
			minHeat = heat
		}
	}

	return strconv.Itoa(minHeat), nil
}

func day17p02(r io.Reader) (string, error) {
	heatMap, err := parseHeatGrid(r)
	if err != nil {
		return "", err
	}

	_, maxX, _, maxY := heatMap.Dimensions()
	target := grid.Position2D[int]{X: maxX, Y: maxY}

	minHeat := -1
	for _, startDir := range []grid.Position2D[int]{{X: 1, Y: 0}, {X: 0, Y: 1}} {
		start := crucibleState{
			position:  grid.Position2D[int]{X: 0, Y: 0},
			direction: startDir,
			steps:     0,
		}

		heat, _, found := search.AStar(
			start,
			crucibleNeighbours(heatMap, 4, 10),
			crucibleHeuristic(target, 4),
			crucibleStepCost(heatMap),
		)

		if found && (minHeat == -1 || heat < minHeat) {
			minHeat = heat
		}
	}

	return strconv.Itoa(minHeat), nil
}
