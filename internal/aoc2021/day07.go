package aoc2021

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func crabPositions(r io.Reader) ([]int64, error) {
	positions, err := commaSeparatedNumbers(r)
	if err != nil {
		return nil, err
	}
	for _, position := range positions {
		if position < 0 {
			return nil, fmt.Errorf("negative crab position")
		}
	}
	return positions, nil
}

func day07p01(r io.Reader) (string, error) {
	positions, err := crabPositions(r)
	if err != nil {
		return "", err
	}
	slices.Sort(positions)
	target := positions[len(positions)/2]
	var fuel int64
	for _, position := range positions {
		distance := xmath.Abs(position - target)
		if fuel > math.MaxInt64-distance {
			return "", fmt.Errorf("crab fuel overflows int64")
		}
		fuel += distance
	}
	return strconv.FormatInt(fuel, 10), nil
}

func triangularFuel(positions []int64, target int64) (int64, error) {
	var fuel int64
	for _, position := range positions {
		distance := xmath.Abs(position - target)
		if distance == math.MaxInt64 {
			return 0, fmt.Errorf("crab fuel overflows int64")
		}
		a, b := distance, distance+1
		if a%2 == 0 {
			a /= 2
		} else {
			b /= 2
		}
		if a > math.MaxInt64/b || fuel > math.MaxInt64-a*b {
			return 0, fmt.Errorf("crab fuel overflows int64")
		}
		fuel += a * b
	}
	return fuel, nil
}

func day07p02(r io.Reader) (string, error) {
	positions, err := crabPositions(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, position := range positions {
		if total > math.MaxInt64-position {
			return "", fmt.Errorf("crab position sum overflows int64")
		}
		total += position
	}
	target := total / int64(len(positions))
	lower, err := triangularFuel(positions, target)
	if err != nil {
		return "", err
	}
	if total%int64(len(positions)) == 0 {
		return strconv.FormatInt(lower, 10), nil
	}
	upper, err := triangularFuel(positions, target+1)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(min(lower, upper), 10), nil
}
