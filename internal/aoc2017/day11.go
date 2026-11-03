package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func hexDistances(r io.Reader) (final, farthest int, err error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte{','}), func(token []byte) (string, error) { return string(bytes.TrimSpace(token)), nil })
	x, y, z := 0, 0, 0
	for direction := range s.Values() {
		switch direction {
		case "n":
			y++
			z--
		case "ne":
			x++
			z--
		case "se":
			x++
			y--
		case "s":
			y--
			z++
		case "sw":
			x--
			z++
		case "nw":
			x--
			y++
		default:
			return 0, 0, fmt.Errorf("invalid hex direction %q", direction)
		}
		final = max(xmath.Abs(x), xmath.Abs(y), xmath.Abs(z))
		farthest = max(farthest, final)
	}

	return final, farthest, s.Err()
}

func day11p01(r io.Reader) (string, error) {
	final, _, err := hexDistances(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(final), nil
}

func day11p02(r io.Reader) (string, error) {
	_, farthest, err := hexDistances(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(farthest), nil
}
