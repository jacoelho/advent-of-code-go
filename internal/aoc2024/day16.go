package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseReindeerMaze(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		res := make([]rune, len(b))
		for i, v := range b {
			res[i] = rune(v)
		}
		return res, nil
	})
	return grid.NewGrid2D[int, rune](slices.Collect(s.Values())), s.Err()
}

func mazeStartPosition(g grid.Grid2D[int, rune]) (grid.Position2D[int], error) {
	for position, value := range g {
		if value == 'S' {
			return position, nil
		}
	}
	return grid.Position2D[int]{}, fmt.Errorf("reindeer maze has no start")
}

type reindeerState struct {
	position  grid.Position2D[int]
	direction grid.Position2D[int]
}

func reindeerNeighbours(maze grid.Grid2D[int, rune]) func(p reindeerState) []reindeerState {
	return func(p reindeerState) []reindeerState {
		forward := reindeerState{
			position:  p.position.Add(p.direction),
			direction: p.direction,
		}
		clockwise := reindeerState{
			position:  p.position,
			direction: p.direction.TurnRight(),
		}
		counter := reindeerState{
			position:  p.position,
			direction: p.direction.TurnLeft(),
		}
		var result []reindeerState
		for _, candidate := range []reindeerState{forward, clockwise, counter} {
			if maze[candidate.position] == '.' || maze[candidate.position] == 'E' || maze[candidate.position] == 'S' {
				result = append(result, candidate)
			}
		}
		return result
	}
}

func reindeerHeuristic(maze grid.Grid2D[int, rune]) func(pair reindeerState) int {
	return func(pair reindeerState) int {
		if maze[pair.position] == 'E' {
			return 0
		}
		return 1
	}
}

func reindeerStepCost(p1, p2 reindeerState) int {
	if p1.position.Add(p1.direction) == p2.position {
		return 1
	}
	return 1000
}

func day16p01(r io.Reader) (string, error) {
	maze, err := parseReindeerMaze(r)
	if err != nil {
		return "", err
	}
	position, err := mazeStartPosition(maze)
	if err != nil {
		return "", err
	}
	start := reindeerState{
		position:  position,
		direction: grid.Position2D[int]{X: 1, Y: 0},
	}

	score, _, found := search.AStar(start, reindeerNeighbours(maze), reindeerHeuristic(maze), reindeerStepCost)
	if !found {
		return "", fmt.Errorf("reindeer maze has no route to the end")
	}

	return strconv.Itoa(score), nil
}

func day16p02(r io.Reader) (string, error) {
	maze, err := parseReindeerMaze(r)
	if err != nil {
		return "", err
	}
	position, err := mazeStartPosition(maze)
	if err != nil {
		return "", err
	}
	start := reindeerState{
		position:  position,
		direction: grid.Position2D[int]{X: 1, Y: 0},
	}

	_, paths, found := search.AStarBag(start, reindeerNeighbours(maze), reindeerHeuristic(maze), reindeerStepCost)
	if !found {
		return "", fmt.Errorf("reindeer maze has no route to the end")
	}

	visited := collections.NewSet[grid.Position2D[int]]()
	for _, path := range paths {
		for _, element := range path {
			visited.Add(element.position)
		}
	}

	return strconv.Itoa(visited.Len()), nil
}
