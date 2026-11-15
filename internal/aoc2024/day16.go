package aoc2024

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

// parseReindeerMaze returns the maze and its unique start tile.
func parseReindeerMaze(r io.Reader) (grid.Bytes, grid.Position2D[int], error) {
	maze, err := grid.ReadBytes(r, "#.SE")
	if err != nil {
		return grid.Bytes{}, grid.Position2D[int]{}, fmt.Errorf("reindeer maze: %w", err)
	}
	start, err := maze.Find('S')
	if err != nil {
		return grid.Bytes{}, grid.Position2D[int]{}, fmt.Errorf("reindeer maze: %w", err)
	}
	if _, err := maze.Find('E'); err != nil {
		return grid.Bytes{}, grid.Position2D[int]{}, fmt.Errorf("reindeer maze: %w", err)
	}
	return maze, start, nil
}

type reindeerState struct {
	position  grid.Position2D[int]
	direction grid.Position2D[int]
}

func reindeerEdges(maze grid.Bytes) search.Edges[reindeerState] {
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
			if maze.Contains(candidate.position) && maze.At(candidate.position) != '#' {
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

func reindeerGoal(maze grid.Bytes) func(reindeerState) bool {
	return func(state reindeerState) bool {
		return maze.At(state.position) == 'E'
	}
}

func day16p01(r io.Reader) (string, error) {
	maze, position, err := parseReindeerMaze(r)
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
	maze, position, err := parseReindeerMaze(r)
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
