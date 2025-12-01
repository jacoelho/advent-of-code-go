package aoc2025

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const dialStartPosition = 50

func rotate(pos int, amount int) int {
	return xmath.Modulo(pos+amount, 100)
}

func countClicks(pos int, amount int) int {
	unwrapped := pos + amount
	clicks := xmath.Abs(unwrapped / 100)

	if unwrapped == 0 || (pos > 0 && unwrapped < 0) {
		clicks++
	}

	return clicks
}

func parseInstructions(r io.Reader) ([]int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int, error) {
		if len(line) == 0 {
			return 0, fmt.Errorf("empty line")
		}

		num, err := convert.ScanNumber[int](line[1:])
		if err != nil {
			return 0, fmt.Errorf("invalid number: %w", err)
		}

		switch line[0] {
		case 'L':
			return -num, nil
		case 'R':
			return num, nil
		default:
			return 0, fmt.Errorf("invalid direction: %c", line[0])
		}
	})
	return slices.Collect(s.Values()), s.Err()
}

func day01p01(r io.Reader) (string, error) {
	instructions, err := parseInstructions(r)
	if err != nil {
		return "", err
	}

	pos := dialStartPosition
	var count int
	for _, amount := range instructions {
		if pos == 0 {
			count++
		}
		pos = rotate(pos, amount)
	}
	return strconv.Itoa(count), nil
}

func day01p02(r io.Reader) (string, error) {
	instructions, err := parseInstructions(r)
	if err != nil {
		return "", err
	}

	pos := dialStartPosition
	var count int
	for _, amount := range instructions {
		count += countClicks(pos, amount)
		pos = rotate(pos, amount)
	}
	return strconv.Itoa(count), nil
}
