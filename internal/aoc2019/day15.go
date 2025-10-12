package aoc2019

import (
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/search"
)

type status int

const (
	statusWall   status = 0
	statusOxygen status = 2
)

type droidDirection int

const (
	north droidDirection = 1
	south droidDirection = 2
	west  droidDirection = 3
	east  droidDirection = 4
)

var allDirections = []droidDirection{north, south, west, east}

type tile int

const (
	tileWall   tile = 1
	tileOpen   tile = 2
	tileOxygen tile = 3
	tileStart  tile = 4
)

type explorationState struct {
	pos      grid.Position2D[int]
	computer *IntcodeComputer
}

func directionToMovement(dir droidDirection) grid.Position2D[int] {
	switch dir {
	case north:
		return grid.NewPosition2D(0, -1)
	case south:
		return grid.NewPosition2D(0, 1)
	case west:
		return grid.NewPosition2D(-1, 0)
	case east:
		return grid.NewPosition2D(1, 0)
	default:
		return grid.NewPosition2D(0, 0)
	}
}

func exploreWithDroid(program []int) (grid.Grid2D[int, tile], grid.Position2D[int]) {
	maze := make(grid.Grid2D[int, tile])
	start := grid.NewPosition2D(0, 0)
	maze[start] = tileStart

	queue := collections.NewDeque[explorationState](100)
	queue.PushBack(explorationState{
		pos:      start,
		computer: New(program),
	})

	visited := collections.NewSet[grid.Position2D[int]]()
	visited.Add(start)

	var oxygenPos grid.Position2D[int]
	for state, ok := queue.PopFront(); ok; state, ok = queue.PopFront() {
		for _, dir := range allDirections {
			movement := directionToMovement(dir)
			nextPos := state.pos.Add(movement)

			if visited.Contains(nextPos) {
				continue
			}

			computer := New(state.computer.Memory())
			computer.AddInput(int(dir))
			if err := computer.Run(); err != nil {
				panic(err)
			}

			statusCode, err := computer.LastOutput()
			if err != nil {
				panic(err)
			}

			s := status(statusCode)
			visited.Add(nextPos)

			switch s {
			case statusWall:
				maze[nextPos] = tileWall
			case statusOxygen:
				maze[nextPos] = tileOxygen
				oxygenPos = nextPos
			default:
				maze[nextPos] = tileOpen
			}

			if s != statusWall {
				queue.PushBack(explorationState{
					pos:      nextPos,
					computer: computer,
				})
			}
		}
	}

	return maze, oxygenPos
}

func mazeNeighbours(maze grid.Grid2D[int, tile]) func(grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
	return func(pos grid.Position2D[int]) iter.Seq[grid.Position2D[int]] {
		return func(yield func(grid.Position2D[int]) bool) {
			for _, dir := range allDirections {
				movement := directionToMovement(dir)
				nextPos := pos.Add(movement)

				t, exists := maze[nextPos]
				if !exists || t == tileWall {
					continue
				}

				if !yield(nextPos) {
					return
				}
			}
		}
	}
}

func findShortestPath(maze grid.Grid2D[int, tile], start, target grid.Position2D[int]) int {
	return search.BFSDistanceTo(start, target, mazeNeighbours(maze))
}

func day15p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	maze, oxygenPos := exploreWithDroid(program)
	start := grid.NewPosition2D(0, 0)
	steps := findShortestPath(maze, start, oxygenPos)

	return strconv.Itoa(steps), nil
}

func findMaxDistanceFromOxygen(maze grid.Grid2D[int, tile], oxygenPos grid.Position2D[int]) int {
	return search.BFSMaxDistance(oxygenPos, mazeNeighbours(maze))
}

func day15p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	maze, oxygenPos := exploreWithDroid(program)
	minutes := findMaxDistanceFromOxygen(maze, oxygenPos)

	return strconv.Itoa(minutes), nil
}
