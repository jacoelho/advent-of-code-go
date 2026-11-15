package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func day25p01(r io.Reader) (string, error) {
	type point [4]int64
	s := scanner.NewScanner(r, func(line []byte) (point, error) {
		fields := bytes.Split(line, []byte(","))
		if len(fields) != 4 {
			return point{}, fmt.Errorf("expected four constellation coordinates, got %q", line)
		}
		var result point
		for i, field := range fields {
			value, err := convert.ScanNumber[int32](bytes.TrimSpace(field))
			if err != nil {
				return point{}, err
			}
			result[i] = int64(value)
		}
		return result, nil
	})
	var points []point
	for p := range s.Values() {
		points = append(points, p)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(points) == 0 {
		return "", fmt.Errorf("no constellation points")
	}

	parent := make([]int, len(points))
	for i := range parent {
		parent[i] = i
	}
	var root func(int) int
	root = func(i int) int {
		if parent[i] != i {
			parent[i] = root(parent[i])
		}
		return parent[i]
	}
	constellations := len(points)
	for i, from := range points {
		for j, to := range points[:i] {
			var distance int64
			for axis := range from {
				distance += xmath.Abs(from[axis] - to[axis])
			}
			if distance > 3 {
				continue
			}
			if a, b := root(i), root(j); a != b {
				parent[a] = b
				constellations--
			}
		}
	}

	return strconv.Itoa(constellations), nil
}
