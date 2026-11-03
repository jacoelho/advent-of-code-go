package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseTriangles(r io.Reader) ([][3]int, error) {
	s := scanner.NewScanner(r, func(line []byte) ([3]int, error) {
		var sides [3]int
		fields := bytes.Fields(line)
		if len(fields) != 3 {
			return sides, fmt.Errorf("expected three triangle sides in %q", line)
		}
		for i, field := range fields {
			side, err := convert.ScanNumber[int](field)
			if err != nil {
				return sides, fmt.Errorf("invalid triangle side: %w", err)
			}
			if side <= 0 {
				return sides, fmt.Errorf("nonpositive triangle side %d", side)
			}
			sides[i] = side
		}
		return sides, nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func validTriangle(sides [3]int) bool {
	slices.Sort(sides[:])
	return sides[0] > sides[2]-sides[1]
}

func day03p01(r io.Reader) (string, error) {
	triangles, err := parseTriangles(r)
	if err != nil {
		return "", err
	}
	count := 0
	for _, triangle := range triangles {
		if validTriangle(triangle) {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func day03p02(r io.Reader) (string, error) {
	rows, err := parseTriangles(r)
	if err != nil {
		return "", err
	}
	if len(rows)%3 != 0 {
		return "", fmt.Errorf("vertical triangles require groups of three rows")
	}
	count := 0
	for start := 0; start < len(rows); start += 3 {
		for column := range 3 {
			if validTriangle([3]int{rows[start][column], rows[start+1][column], rows[start+2][column]}) {
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}
