package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func day25p01(r io.Reader) (string, error) {
	type point [4]int32
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
			result[i] = value
		}
		return result, nil
	})
	points := collections.NewSet[point]()
	for p := range s.Values() {
		points.Add(p)
		if len(points) > 5000 {
			return "", fmt.Errorf("constellation search supports at most 5000 points")
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}

	neighbours := func(from point) iter.Seq[point] {
		return func(yield func(point) bool) {
			for to := range points {
				var distance int64
				for axis := range 4 {
					distance += xmath.Abs(int64(from[axis]) - int64(to[axis]))
				}
				if distance <= 3 && !yield(to) {
					return
				}
			}
		}
	}
	visited := collections.NewSet[point]()
	count := 0
	for p := range points {
		if visited.Contains(p) {
			continue
		}
		for range search.BFSWithVisited(p, visited, neighbours) {
		}
		count++
	}

	return strconv.Itoa(count), nil
}
