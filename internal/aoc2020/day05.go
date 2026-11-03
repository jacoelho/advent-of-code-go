package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseBoardingPasses(r io.Reader) ([]int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int, error) {
		if len(line) != 10 {
			return 0, fmt.Errorf("boarding pass must have 10 characters: %q", line)
		}
		seat := 0
		for i, c := range line {
			low, high := byte('F'), byte('B')
			if i >= 7 {
				low, high = 'L', 'R'
			}
			if c != low && c != high {
				return 0, fmt.Errorf("invalid boarding pass character %q at position %d", c, i+1)
			}
			seat <<= 1
			if c == high {
				seat++
			}
		}
		return seat, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func day05p01(r io.Reader) (string, error) {
	seats, err := parseBoardingPasses(r)
	if err != nil {
		return "", err
	}
	if len(seats) == 0 {
		return "", fmt.Errorf("no boarding passes")
	}
	return strconv.Itoa(slices.Max(seats)), nil
}

func day05p02(r io.Reader) (string, error) {
	seats, err := parseBoardingPasses(r)
	if err != nil {
		return "", err
	}
	var occupied [1024]bool
	for _, seat := range seats {
		occupied[seat] = true
	}
	missing := -1
	for seat := 8; seat < 127*8; seat++ {
		if !occupied[seat] && occupied[seat-1] && occupied[seat+1] {
			if missing >= 0 {
				return "", fmt.Errorf("multiple missing seats with occupied neighbours")
			}
			missing = seat
		}
	}
	if missing < 0 {
		return "", fmt.Errorf("no missing seat with occupied neighbours")
	}
	return strconv.Itoa(missing), nil
}
