package aoc2019

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func calculateFuel(mass int) int { return mass/3 - 2 }

func calculateFuelRecursive(mass int) int {
	fuel := calculateFuel(mass)
	if fuel <= 0 {
		return 0
	}
	return fuel + calculateFuelRecursive(fuel)
}

func day01p01(r io.Reader) (string, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	total := xiter.From(s.Values()).SumBy(calculateFuel)
	if err := s.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(total), nil
}

func day01p02(r io.Reader) (string, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	total := xiter.From(s.Values()).SumBy(calculateFuelRecursive)
	if err := s.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(total), nil
}
