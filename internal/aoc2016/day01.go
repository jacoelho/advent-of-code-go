package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func taxicabDistances(r io.Reader) (final, repeated int, err error) {
	type instruction struct {
		turn   byte
		blocks int
	}
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte{','}), func(token []byte) (instruction, error) {
		token = bytes.TrimSpace(token)
		if len(token) < 2 || (token[0] != 'L' && token[0] != 'R') {
			return instruction{}, fmt.Errorf("invalid taxicab instruction %q", token)
		}
		blocks, err := convert.ScanNumber[int](token[1:])
		if err != nil {
			return instruction{}, fmt.Errorf("invalid block count: %w", err)
		}
		if blocks < 0 {
			return instruction{}, fmt.Errorf("negative block count %d", blocks)
		}
		return instruction{token[0], blocks}, nil
	})
	origin := grid.Position2D[int]{}
	position, direction := origin, grid.Position2D[int]{Y: -1}
	visited := collections.NewSet(origin)
	repeated = -1
	for instruction := range s.Values() {
		if instruction.turn == 'R' {
			direction = direction.TurnRight()
		} else {
			direction = direction.TurnLeft()
		}
		for range instruction.blocks {
			position = position.Add(direction)
			if repeated < 0 && visited.Contains(position) {
				repeated = position.Distance(origin)
			}
			visited.Add(position)
		}
	}
	return position.Distance(origin), repeated, s.Err()
}

func day01p01(r io.Reader) (string, error) {
	distance, _, err := taxicabDistances(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(distance), nil
}

func day01p02(r io.Reader) (string, error) {
	_, distance, err := taxicabDistances(r)
	if err != nil {
		return "", err
	}
	if distance < 0 {
		return "", fmt.Errorf("route never revisits a location")
	}
	return strconv.Itoa(distance), nil
}
