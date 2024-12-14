package aoc2024

import (
	"bufio"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

type robot struct {
	position grid.Position2D[int]
	velocity grid.Position2D[int]
}

func parseRobot(line string) robot {
	digits := convert.ExtractDigits[int](line)
	return robot{
		position: grid.Position2D[int]{X: digits[0], Y: digits[1]},
		velocity: grid.Position2D[int]{X: digits[2], Y: digits[3]},
	}
}

func parseBathroomRobotsPositions(r io.Reader) []robot {
	var result []robot
	s := bufio.NewScanner(r)
	for s.Scan() {
		result = append(result, parseRobot(s.Text()))
	}
	if s.Err() != nil {
		panic(s.Err())
	}
	return result
}

func positionsAfterIterations(
	robots []robot,
	width int,
	height int,
	iterations int,
) []grid.Position2D[int] {
	result := make([]grid.Position2D[int], 0, len(robots))
	for _, r := range robots {
		position := grid.Position2D[int]{
			X: xmath.Modulo(r.position.X+r.velocity.X*iterations, width),
			Y: xmath.Modulo(r.position.Y+r.velocity.Y*iterations, height),
		}
		result = append(result, position)
	}
	return result
}

func countQuadrants(
	positions []grid.Position2D[int],
	width int,
	height int,
) []int {
	verticalSplit := height / 2
	horizontalSplit := width / 2

	result := make([]int, 4)
	for _, position := range positions {
		if position.X == horizontalSplit || position.Y == verticalSplit {
			continue
		}
		quadrant := 0
		if position.X >= horizontalSplit {
			quadrant += 2
		}
		if position.Y >= verticalSplit {
			quadrant += 1
		}
		result[quadrant]++
	}
	return result
}

func day14p01(width, height int) func(io.Reader) (string, error) {
	return func(reader io.Reader) (string, error) {
		initialPositions := parseBathroomRobotsPositions(reader)

		finalPositions := positionsAfterIterations(initialPositions, width, height, 100)

		count := countQuadrants(finalPositions, width, height)

		return strconv.Itoa(xslices.Product(count)), nil
	}
}

func day14p02(reader io.Reader) (string, error) {
	const (
		width  = 101
		height = 103
	)

	initialPositions := parseBathroomRobotsPositions(reader)

	for iteration := 1; ; iteration++ {
		finalPositions := positionsAfterIterations(initialPositions, width, height, iteration)

		// check if robots no longer overlap
		if !xslices.HasDuplicates(finalPositions) {
			return strconv.Itoa(iteration), nil
		}
	}

	panic("unreachable")
}
