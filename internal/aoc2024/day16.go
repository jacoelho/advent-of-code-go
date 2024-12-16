package aoc2024

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
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

func mazeStartPosition(g grid.Grid2D[int, rune]) grid.Position2D[int] {
	element, ok := xmaps.Find(g, func(p grid.Position2D[int], v rune) bool { return v == 'S' })
	if !ok {
		panic("not found")
	}
	return element.K
}

type day16Pair struct {
	position  grid.Position2D[int]
	direction grid.Position2D[int]
}

func day16Neighbours(maze grid.Grid2D[int, rune]) func(p day16Pair) []day16Pair {
	return func(p day16Pair) []day16Pair {
		forward := day16Pair{
			position:  p.position.Add(p.direction),
			direction: p.direction,
		}
		clockwise := day16Pair{
			position:  p.position,
			direction: p.direction.TurnRight(),
		}
		counter := day16Pair{
			position:  p.position,
			direction: p.direction.TurnLeft(),
		}
		var result []day16Pair
		for _, candidate := range []day16Pair{forward, clockwise, counter} {
			if maze[candidate.position] == '.' || maze[candidate.position] == 'E' || maze[candidate.position] == 'S' {
				result = append(result, candidate)
			}
		}
		return result
	}
}

func day16Heuristic(maze grid.Grid2D[int, rune]) func(pair day16Pair) int {
	return func(pair day16Pair) int {
		if maze[pair.position] == 'E' {
			return 0
		}
		return 1
	}
}

func day16Cost(p1, p2 day16Pair) int {
	if p1.position.Add(p1.direction) == p2.position {
		return 1
	}
	return 1000
}

func day16p01(r io.Reader) (string, error) {
	maze := aoc.Must(parseReindeerMaze(r))
	start := day16Pair{
		position:  mazeStartPosition(maze),
		direction: grid.Position2D[int]{1, 0},
	}

	score, _, _ := search.AStar(start, day16Neighbours(maze), day16Heuristic(maze), day16Cost)

	return strconv.Itoa(score), nil
}

func day16p02(r io.Reader) (string, error) {
	maze := aoc.Must(parseReindeerMaze(r))
	start := day16Pair{
		position:  mazeStartPosition(maze),
		direction: grid.Position2D[int]{1, 0},
	}

	_, paths, _ := search.AStarBag(start, day16Neighbours(maze), day16Heuristic(maze), day16Cost)

	visited := collections.NewSet[grid.Position2D[int]]()
	for _, path := range paths {
		for _, element := range path {
			visited.Add(element.position)
		}
	}

	return strconv.Itoa(visited.Len()), nil
}
