package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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
		return grid.ParsePosition2D[int](text)
	}
	s := scanner.NewScanner(r, func(line []byte) (robot, error) {
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
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(robots) == 0 {
		return nil, fmt.Errorf("bathroom has no robots")
	}
	return robots, nil
}

func day14CheckBounds(robots []robot, width, height int) error {
	if width < 1 || height < 1 {
		return fmt.Errorf("bathroom dimensions %dx%d must be positive", width, height)
	}
	for i, r := range robots {
		if r.position.X < 0 || r.position.X >= width || r.position.Y < 0 || r.position.Y >= height {
			return fmt.Errorf("robot %d position %v is outside the %dx%d bathroom", i+1, r.position, width, height)
		}
	}
	return nil
}

func positionsAfterIterations(
	robots []robot,
	width int,
	height int,
	iterations int,
) []grid.Position2D[int] {
	advance := func(position, velocity, size int) int {
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
		if err := day14CheckBounds(initialPositions, width, height); err != nil {
			return "", err
		}

		finalPositions := positionsAfterIterations(initialPositions, width, height, 100)

		count := countQuadrants(finalPositions, width, height)

		return strconv.Itoa(xiter.Product(slices.Values(count))), nil
	}
}

// day14LeastSpreadTime returns the time in [0, size) at which the robots'
// coordinates along one axis have the least variance.
func day14LeastSpreadTime(positions, velocities []int, size int) int {
	best, bestSpread := 0, -1
	n := len(positions)
	for t := range size {
		var sum, sumSquares int
		for i, p := range positions {
			q := (p + velocities[i]*t) % size
			sum += q
			sumSquares += q * q
		}
		// n² times the variance.
		if spread := n*sumSquares - sum*sum; bestSpread < 0 || spread < bestSpread {
			best, bestSpread = t, spread
		}
	}
	return best
}

// day14TreeTime returns the frame where the robots cluster into the tree,
// which is the frame of least spread on both axes. x positions repeat every
// width steps and y positions every height steps, independently, so each axis
// minimum is found within one period and, for coprime dimensions, the Chinese
// remainder theorem gives the unique frame below width*height matching both.
func day14TreeTime(robots []robot, width, height int) int {
	xs, ys := make([]int, len(robots)), make([]int, len(robots))
	vxs, vys := make([]int, len(robots)), make([]int, len(robots))
	for i, r := range robots {
		xs[i], ys[i] = r.position.X, r.position.Y
		vxs[i], vys[i] = xmath.Modulo(r.velocity.X, width), xmath.Modulo(r.velocity.Y, height)
	}
	tx := day14LeastSpreadTime(xs, vxs, width)
	ty := day14LeastSpreadTime(ys, vys, height)
	t := tx
	for t%height != ty {
		t += width
	}
	return t
}

func day14p02(reader io.Reader) (string, error) {
	const width, height = 101, 103
	robots, err := parseBathroomRobots(reader)
	if err != nil {
		return "", err
	}
	if err := day14CheckBounds(robots, width, height); err != nil {
		return "", err
	}
	return strconv.Itoa(day14TreeTime(robots, width, height)), nil
}
