package aoc2019

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func calculateFuel(mass int) int { return mass/3 - 2 }

func calculateFuelRecursive(mass int) int {
	total := 0
	for fuel := calculateFuel(mass); fuel > 0; fuel = calculateFuel(fuel) {
		total += fuel
	}
	return total
}

func parseMass(line []byte) (int, error) {
	mass, err := convert.ScanNumber[int](line)
	if err != nil {
		return 0, err
	}
	if mass < 0 {
		return 0, fmt.Errorf("negative module mass %d", mass)
	}
	return mass, nil
}

func totalFuel(r io.Reader, fuel func(int) int) (string, error) {
	s := scanner.NewScanner(r, parseMass)
	total := 0
	for mass := range s.Values() {
		total += fuel(mass)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day01p01(r io.Reader) (string, error) {
	return totalFuel(r, calculateFuel)
}

func day01p02(r io.Reader) (string, error) {
	return totalFuel(r, calculateFuelRecursive)
}
