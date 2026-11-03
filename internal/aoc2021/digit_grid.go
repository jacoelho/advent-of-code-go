package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"iter"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const maxDigitGridCells = 1_000_000

type digitGrid struct {
	width, height int
	values        []byte
}

func parseDigitGrid(r io.Reader) (digitGrid, error) {
	var digits digitGrid
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 1024), maxDigitGridCells+1)
	for s.Scan() {
		row := s.Bytes()
		if len(row) == 0 || digits.height > 0 && len(row) != digits.width {
			return digits, fmt.Errorf("empty or ragged digit grid")
		}
		if len(row) > maxDigitGridCells-len(digits.values) {
			return digits, fmt.Errorf("digit grid exceeds %d cells", maxDigitGridCells)
		}
		digits.width = len(row)
		digits.height++
		for _, value := range row {
			if value < '0' || value > '9' {
				return digits, fmt.Errorf("invalid grid digit %q", value)
			}
			digits.values = append(digits.values, value-'0')
		}
	}
	if err := s.Err(); err != nil {
		return digits, err
	}
	if digits.height == 0 {
		return digits, fmt.Errorf("empty digit grid")
	}
	return digits, nil
}

func (digits digitGrid) neighbours(index int) iter.Seq[int] {
	return func(yield func(int) bool) {
		point := grid.Position2D[int]{X: index % digits.width, Y: index / digits.width}
		for neighbour := range grid.Neighbours4(point) {
			if neighbour.X < digits.width && neighbour.Y < digits.height && !yield(neighbour.Y*digits.width+neighbour.X) {
				return
			}
		}
	}
}
