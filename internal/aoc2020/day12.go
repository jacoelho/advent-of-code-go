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

func day12p01(r io.Reader) (string, error) {
	instructions, err := parseNavigation(r)
	if err != nil {
		return "", err
	}
	x, y, dx, dy := 0, 0, 1, 0
	for _, instruction := range instructions {
		value := instruction.value
		switch instruction.action {
		case 'N':
			y += value
		case 'S':
			y -= value
		case 'E':
			x += value
		case 'W':
			x -= value
		case 'F':
			x += value * dx
			y += value * dy
		case 'L':
			dx, dy = rotateNavigation(dx, dy, -value/90)
		case 'R':
			dx, dy = rotateNavigation(dx, dy, value/90)
		}
	}
	return strconv.Itoa(xmath.Abs(x) + xmath.Abs(y)), nil
}

func day12p02(r io.Reader) (string, error) {
	instructions, err := parseNavigation(r)
	if err != nil {
		return "", err
	}
	x, y, wx, wy := 0, 0, 10, 1
	for _, instruction := range instructions {
		value := instruction.value
		switch instruction.action {
		case 'N':
			wy += value
		case 'S':
			wy -= value
		case 'E':
			wx += value
		case 'W':
			wx -= value
		case 'F':
			x += value * wx
			y += value * wy
		case 'L':
			wx, wy = rotateNavigation(wx, wy, -value/90)
		case 'R':
			wx, wy = rotateNavigation(wx, wy, value/90)
		}
	}
	return strconv.Itoa(xmath.Abs(x) + xmath.Abs(y)), nil
}
