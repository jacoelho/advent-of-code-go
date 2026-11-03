package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"iter"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func spiralTarget(r io.Reader) (int, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}

	target, err := convert.ScanNumber[int](bytes.TrimSpace(input))
	if err != nil {
		return 0, err
	}
	if target < 1 {
		return 0, fmt.Errorf("spiral target must be positive")
	}

	return target, nil
}

func spiralPositions() iter.Seq[grid.Position2D[int]] {
	return func(yield func(grid.Position2D[int]) bool) {
		position, direction := grid.Position2D[int]{}, grid.Position2D[int]{X: 1}
		if !yield(position) {
			return
		}
		for length := 1; ; length++ {
			for range 2 {
				for range length {
					position = position.Add(direction)
					if !yield(position) {
						return
					}
				}
				direction = direction.TurnLeft()
			}
		}
	}
}

func day03p01(r io.Reader) (string, error) {
	target, err := spiralTarget(r)
	if err != nil {
		return "", err
	}

	index := 0
	for position := range spiralPositions() {
		index++
		if index == target {
			return strconv.Itoa(position.Distance(grid.Position2D[int]{})), nil
		}
	}

	return "", fmt.Errorf("spiral ended before target")
}

func day03p02(r io.Reader) (string, error) {
	target, err := spiralTarget(r)
	if err != nil {
		return "", err
	}

	values := make(map[grid.Position2D[int]]int)
	offsets := grid.OffsetsNeighbours8[int]()
	for position := range spiralPositions() {
		total := 0
		for _, offset := range offsets {
			value := values[position.Add(offset)]
			if total > math.MaxInt-value {
				return "", fmt.Errorf("spiral stress value exceeds int")
			}
			total += value
		}
		total = max(total, 1)
		if total > target {
			return strconv.Itoa(total), nil
		}
		values[position] = total
	}

	return "", fmt.Errorf("spiral ended before target")
}
