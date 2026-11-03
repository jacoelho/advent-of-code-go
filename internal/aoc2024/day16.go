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

func reindeerEdges(maze grid.Grid2D[int, rune]) search.Edges[reindeerState] {
	return func(p reindeerState, yield func(reindeerState, int) bool) error {
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
		for i, candidate := range [...]reindeerState{forward, clockwise, counter} {
			if maze[candidate.position] == '.' || maze[candidate.position] == 'E' || maze[candidate.position] == 'S' {
				cost := 1000
				if i == 0 {
					cost = 1
				}
				if !yield(candidate, cost) {
					return nil
				}
			}
		}
		return nil
	}
}

func reindeerGoal(maze grid.Grid2D[int, rune]) func(reindeerState) bool {
	return func(state reindeerState) bool {
		return maze[state.position] == 'E'
	}
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

	score, found, err := search.AStarDistance(start, reindeerGoal(maze), reindeerEdges(maze), nil, search.Limits{})
	if err != nil {
		return "", fmt.Errorf("reindeer maze: %w", err)
	}
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

	_, states, found, err := search.DijkstraOptimalStates(start, reindeerGoal(maze), reindeerEdges(maze), search.Limits{})
	if err != nil {
		return "", fmt.Errorf("reindeer maze: %w", err)
	}
	if !found {
		return "", fmt.Errorf("reindeer maze has no route to the end")
	}

	visited := collections.NewSet[grid.Position2D[int]]()
	for _, state := range states {
		visited.Add(state.position)
	}

	return strconv.Itoa(visited.Len()), nil
}
