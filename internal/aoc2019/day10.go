package aoc2019

import (
	"cmp"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func parseAsteroidMap(r io.Reader) ([]grid.Position2D[int], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]byte, error) { return line, nil })
	var asteroids []grid.Position2D[int]
	width, height := -1, 0
	for line := range s.Values() {
		if width >= 0 && len(line) != width {
			return nil, fmt.Errorf("asteroid map row %d has width %d, want %d", height, len(line), width)
		}
		width = len(line)
		for x, c := range line {
			switch c {
			case '#':
				asteroids = append(asteroids, grid.Position2D[int]{X: x, Y: height})
			case '.':
			default:
				return nil, fmt.Errorf("invalid asteroid map cell %q at (%d,%d)", c, x, height)
			}
		}
		height++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(asteroids) == 0 {
		return nil, fmt.Errorf("asteroid map has no asteroids")
	}
	return asteroids, nil
}

func lineOfSight(station, asteroid grid.Position2D[int]) grid.Position2D[int] {
	delta := asteroid.Sub(station)
	gcd := xmath.GCD(delta.X, delta.Y)
	return grid.Position2D[int]{X: delta.X / gcd, Y: delta.Y / gcd}
}

// bestStation returns the first asteroid, in row order, that sees the most
// others, and how many it sees.
func bestStation(asteroids []grid.Position2D[int]) (grid.Position2D[int], int) {
	var best grid.Position2D[int]
	bestVisible := -1
	directions := make(map[grid.Position2D[int]]struct{})
	for _, station := range asteroids {
		clear(directions)
		for _, asteroid := range asteroids {
			if asteroid != station {
				directions[lineOfSight(station, asteroid)] = struct{}{}
			}
		}
		if len(directions) > bestVisible {
			best, bestVisible = station, len(directions)
		}
	}
	return best, bestVisible
}

func day10p01(r io.Reader) (string, error) {
	asteroids, err := parseAsteroidMap(r)
	if err != nil {
		return "", err
	}
	_, visible := bestStation(asteroids)
	return strconv.Itoa(visible), nil
}

// clockwiseAngle measures direction clockwise from straight up, in [0, 2π).
func clockwiseAngle(direction grid.Position2D[int]) float64 {
	angle := math.Atan2(float64(direction.X), float64(-direction.Y))
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

// vaporizationOrder returns the asteroids other than station in the order a
// clockwise-rotating laser, starting upward, destroys them: the n-th asteroid
// along a line of sight goes on the n-th rotation.
func vaporizationOrder(station grid.Position2D[int], asteroids []grid.Position2D[int]) []grid.Position2D[int] {
	type target struct {
		pos      grid.Position2D[int]
		rotation int
		angle    float64
	}
	byDirection := make(map[grid.Position2D[int]][]grid.Position2D[int])
	for _, asteroid := range asteroids {
		if asteroid != station {
			direction := lineOfSight(station, asteroid)
			byDirection[direction] = append(byDirection[direction], asteroid)
		}
	}
	var targets []target
	for direction, line := range byDirection {
		slices.SortFunc(line, func(a, b grid.Position2D[int]) int {
			return cmp.Compare(station.Distance(a), station.Distance(b))
		})
		angle := clockwiseAngle(direction)
		for rotation, pos := range line {
			targets = append(targets, target{pos: pos, rotation: rotation, angle: angle})
		}
	}
	slices.SortFunc(targets, func(a, b target) int {
		return cmp.Or(cmp.Compare(a.rotation, b.rotation), cmp.Compare(a.angle, b.angle))
	})
	order := make([]grid.Position2D[int], len(targets))
	for i, t := range targets {
		order[i] = t.pos
	}
	return order
}

func day10p02(r io.Reader) (string, error) {
	asteroids, err := parseAsteroidMap(r)
	if err != nil {
		return "", err
	}

	const target = 200
	station, _ := bestStation(asteroids)
	order := vaporizationOrder(station, asteroids)
	if len(order) < target {
		return "", fmt.Errorf("only %d asteroids to vaporize, need %d", len(order), target)
	}
	pos := order[target-1]
	return strconv.Itoa(pos.X*100 + pos.Y), nil
}
