package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type navigationInstruction struct {
	action byte
	value  int
}

func parseNavigation(r io.Reader) ([]navigationInstruction, error) {
	s := scanner.NewScanner(r, func(line []byte) (navigationInstruction, error) {
		if len(line) < 2 {
			return navigationInstruction{}, fmt.Errorf("invalid navigation instruction %q", line)
		}
		value, err := strconv.Atoi(string(line[1:]))
		if err != nil || value < 0 {
			return navigationInstruction{}, fmt.Errorf("invalid navigation value in %q", line)
		}
		switch line[0] {
		case 'N', 'S', 'E', 'W', 'F':
		case 'L', 'R':
			if value%90 != 0 {
				return navigationInstruction{}, fmt.Errorf("rotation must be a multiple of 90 degrees: %q", line)
			}
		default:
			return navigationInstruction{}, fmt.Errorf("unknown navigation action %q", line[0])
		}
		return navigationInstruction{line[0], value}, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func rotateNavigation(x, y, quarterTurns int) (int, int) {
	for range (quarterTurns%4 + 4) % 4 {
		x, y = y, -x
	}
	return x, y
}

// navigate moves the ship along the waypoint (wx, wy) on F. Cardinal actions
// move the waypoint when moveWaypoint is set, and the ship otherwise.
func navigate(r io.Reader, wx, wy int, moveWaypoint bool) (string, error) {
	instructions, err := parseNavigation(r)
	if err != nil {
		return "", err
	}
	x, y := 0, 0
	for _, instruction := range instructions {
		value := instruction.value
		dx, dy := 0, 0
		switch instruction.action {
		case 'N':
			dy = value
		case 'S':
			dy = -value
		case 'E':
			dx = value
		case 'W':
			dx = -value
		case 'F':
			x += value * wx
			y += value * wy
		case 'L':
			wx, wy = rotateNavigation(wx, wy, -value/90)
		case 'R':
			wx, wy = rotateNavigation(wx, wy, value/90)
		}
		if moveWaypoint {
			wx, wy = wx+dx, wy+dy
		} else {
			x, y = x+dx, y+dy
		}
	}
	return strconv.Itoa(xmath.Abs(x) + xmath.Abs(y)), nil
}

func day12p01(r io.Reader) (string, error) {
	return navigate(r, 1, 0, false)
}

func day12p02(r io.Reader) (string, error) {
	return navigate(r, 10, 1, true)
}
