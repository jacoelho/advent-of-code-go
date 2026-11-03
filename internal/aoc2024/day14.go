package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type robot struct {
	position grid.Position2D[int]
	velocity grid.Position2D[int]
}

func parseBathroomRobots(r io.Reader) ([]robot, error) {
	parsePair := func(text, prefix string) (grid.Position2D[int], error) {
		text, ok := strings.CutPrefix(text, prefix)
		if !ok {
			return grid.Position2D[int]{}, fmt.Errorf("robot field needs prefix %q", prefix)
		}
		xText, yText, ok := strings.Cut(text, ",")
		if !ok {
			return grid.Position2D[int]{}, fmt.Errorf("robot field needs two coordinates")
		}
		x, err := convert.ScanNumber[int]([]byte(xText))
		if err != nil {
			return grid.Position2D[int]{}, fmt.Errorf("robot x: %w", err)
		}
		y, err := convert.ScanNumber[int]([]byte(yText))
		if err != nil {
			return grid.Position2D[int]{}, fmt.Errorf("robot y: %w", err)
		}
		return grid.NewPosition2D(x, y), nil
	}
	count := 0
	s := scanner.NewScanner(r, func(line []byte) (robot, error) {
		count++
		if count > 100_000 {
			return robot{}, fmt.Errorf("robot input exceeds 100000 robots")
		}
		fields := strings.Fields(string(line))
		if len(fields) != 2 {
			return robot{}, fmt.Errorf("robot needs position and velocity")
		}
		position, err := parsePair(fields[0], "p=")
		if err != nil {
			return robot{}, err
		}
		velocity, err := parsePair(fields[1], "v=")
		if err != nil {
			return robot{}, err
		}
		return robot{position: position, velocity: velocity}, nil
	})
	robots := slices.Collect(s.Values())
	return robots, s.Err()
}

func positionsAfterIterations(
	robots []robot,
	width int,
	height int,
	iterations int,
) []grid.Position2D[int] {
	advance := func(position, velocity, size int) int {
		position = xmath.Modulo(position, size)
		velocity = xmath.Modulo(velocity, size)
		return int((int64(position) + int64(velocity)*int64(iterations)) % int64(size))
	}
	result := make([]grid.Position2D[int], 0, len(robots))
	for _, r := range robots {
		position := grid.Position2D[int]{
			X: advance(r.position.X, r.velocity.X, width),
			Y: advance(r.position.Y, r.velocity.Y, height),
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
		initialPositions, err := parseBathroomRobots(reader)
		if err != nil {
			return "", err
		}
		if width < 1 || height < 1 || width > 1_000_000 || height > 1_000_000 {
			return "", fmt.Errorf("bathroom dimensions must be between 1 and 1000000")
		}

		finalPositions := positionsAfterIterations(initialPositions, width, height, 100)

		count := countQuadrants(finalPositions, width, height)

		return strconv.Itoa(xiter.Product(slices.Values(count))), nil
	}
}

func uniqueRobotFrame(robots []robot, width, height int) (int, error) {
	if len(robots) > width*height {
		return 0, fmt.Errorf("no collision-free frame: more robots than cells")
	}
	candidate := -1
	for iteration := range xmath.LCM(width, height) {
		positions := positionsAfterIterations(robots, width, height, iteration)
		if xslices.HasDuplicates(positions) {
			continue
		}
		if candidate != -1 {
			return 0, fmt.Errorf("cannot identify tree: multiple collision-free frames")
		}
		candidate = iteration
	}
	if candidate == -1 {
		return 0, fmt.Errorf("cannot identify tree: no collision-free frame")
	}
	return candidate, nil
}

func day14p02(reader io.Reader) (string, error) {
	robots, err := parseBathroomRobots(reader)
	if err != nil {
		return "", err
	}
	// This specialization assumes the unique collision-free frame is the intended tree.
	iteration, err := uniqueRobotFrame(robots, 101, 103)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(iteration), nil
}
