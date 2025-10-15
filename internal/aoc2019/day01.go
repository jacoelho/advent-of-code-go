package aoc2019

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseFuel(r io.Reader) ([]int, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	return slices.Collect(s.Values()), s.Err()
}

func calculateFuel(mass int) int { return mass/3 - 2 }

func calculateFuelRecursive(mass int) int {
	fuel := calculateFuel(mass)
	if fuel <= 0 {
		return 0
	}
	return fuel + calculateFuelRecursive(fuel)
}

func day01p01(r io.Reader) (string, error) {
	masses, err := parseFuel(r)
	if err != nil {
		return "", err
	}

	total := xslices.Reduce(func(sum int, mass int) int {
		return sum + calculateFuel(mass)
	}, 0, masses)

	return strconv.Itoa(total), nil
}

func day01p02(r io.Reader) (string, error) {
	masses, err := parseFuel(r)
	if err != nil {
		return "", err
	}

	total := xslices.Reduce(func(sum int, mass int) int {
		return sum + calculateFuelRecursive(mass)
	}, 0, masses)

	return strconv.Itoa(total), nil
}
