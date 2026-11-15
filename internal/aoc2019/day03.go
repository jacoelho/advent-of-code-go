package aoc2019

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type wireSegment struct {
	offset   grid.Position2D[int]
	distance int
}

func parseWirePath(line []byte) ([]wireSegment, error) {
	var segments []wireSegment
	for token := range strings.SplitSeq(strings.TrimSpace(string(line)), ",") {
		if len(token) < 2 || token[1] < '0' || token[1] > '9' {
			return nil, fmt.Errorf("invalid wire segment %q", token)
		}
		var offset grid.Position2D[int]
		switch token[0] {
		case 'U':
			offset = grid.Position2D[int]{Y: -1}
		case 'D':
			offset = grid.Position2D[int]{Y: 1}
		case 'L':
			offset = grid.Position2D[int]{X: -1}
		case 'R':
			offset = grid.Position2D[int]{X: 1}
		default:
			return nil, fmt.Errorf("invalid direction in wire segment %q", token)
		}
		distance, err := convert.ScanNumber[int]([]byte(token[1:]))
		if err != nil {
			return nil, fmt.Errorf("wire segment %q: %w", token, err)
		}
		segments = append(segments, wireSegment{offset: offset, distance: distance})
	}
	return segments, nil
}

// traceWire maps every point the wire visits, except the central port, to the
// steps taken to first reach it.
func traceWire(segments []wireSegment) map[grid.Position2D[int]]int {
	steps := make(map[grid.Position2D[int]]int)
	var current grid.Position2D[int]
	stepCount := 0
	for _, segment := range segments {
		for range segment.distance {
			current = current.Add(segment.offset)
			stepCount++
			if _, seen := steps[current]; !seen && current != (grid.Position2D[int]{}) {
				steps[current] = stepCount
			}
		}
	}
	return steps
}

// closestCrossing returns the lowest cost over the points both wires visit.
func closestCrossing(r io.Reader, cost func(pos grid.Position2D[int], steps1, steps2 int) int) (string, error) {
	s := scanner.NewScanner(r, parseWirePath)
	var wires [][]wireSegment
	for wire := range s.Values() {
		wires = append(wires, wire)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(wires) != 2 {
		return "", fmt.Errorf("expected 2 wires, got %d", len(wires))
	}

	steps1, steps2 := traceWire(wires[0]), traceWire(wires[1])
	best := math.MaxInt
	for pos, s1 := range steps1 {
		if s2, ok := steps2[pos]; ok {
			best = min(best, cost(pos, s1, s2))
		}
	}
	if best == math.MaxInt {
		return "", fmt.Errorf("wires do not intersect outside the central port")
	}
	return strconv.Itoa(best), nil
}

func day03p01(r io.Reader) (string, error) {
	return closestCrossing(r, func(pos grid.Position2D[int], _, _ int) int {
		return pos.Distance(grid.Position2D[int]{})
	})
}

func day03p02(r io.Reader) (string, error) {
	return closestCrossing(r, func(_ grid.Position2D[int], steps1, steps2 int) int {
		return steps1 + steps2
	})
}
