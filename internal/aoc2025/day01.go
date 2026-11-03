package aoc2025

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type dialCountMode uint8

const (
	dialEndings dialCountMode = iota
	dialClicks
)

func parseDialRotations(r io.Reader) ([]int64, error) {
	s := scanner.NewScanner(r, func(line []byte) (int64, error) {
		if len(line) < 2 || line[0] != 'L' && line[0] != 'R' {
			return 0, fmt.Errorf("dial rotation needs L or R and a distance")
		}
		if line[1] == '-' || line[1] == '+' {
			return 0, fmt.Errorf("dial distance must be unsigned decimal")
		}
		amount, err := convert.ScanNumber[int64](line[1:])
		if err != nil {
			return 0, fmt.Errorf("dial distance: %w", err)
		}
		if line[0] == 'L' {
			amount = -amount
		}
		return amount, nil
	})
	var rotations []int64
	for rotation := range s.Values() {
		if len(rotations) == 100_000 {
			return nil, fmt.Errorf("dial input exceeds 100000 rotations")
		}
		rotations = append(rotations, rotation)
	}
	return rotations, s.Err()
}

func countClicks(position, amount int64) int64 {
	distance := amount
	if distance < 0 {
		distance = -distance
	}
	clicks, remainder := distance/100, distance%100
	if amount > 0 && position+remainder >= 100 || amount < 0 && position > 0 && remainder >= position {
		clicks++
	}
	return clicks
}

func dialPassword(r io.Reader, mode dialCountMode) (string, error) {
	rotations, err := parseDialRotations(r)
	if err != nil {
		return "", err
	}
	position := int64(50)
	var count int64
	for _, amount := range rotations {
		next := xmath.Modulo(position+amount%100, 100)
		var increment int64
		if mode == dialClicks {
			increment = countClicks(position, amount)
		} else if next == 0 {
			increment = 1
		}
		if increment > math.MaxInt64-count {
			return "", fmt.Errorf("dial password overflows int64")
		}
		count += increment
		position = next
	}
	return strconv.FormatInt(count, 10), nil
}

func day01p01(r io.Reader) (string, error) {
	return dialPassword(r, dialEndings)
}

func day01p02(r io.Reader) (string, error) {
	return dialPassword(r, dialClicks)
}
