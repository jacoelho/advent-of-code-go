package aoc2021

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func stepSeaCucumbers(floor grid.Bytes, next []byte) bool {
	cells := floor.Cells()
	width, height := floor.Width(), floor.Height()
	moved := false
	for _, kind := range [...]byte{'>', 'v'} {
		copy(next, cells)
		for i, c := range cells {
			if c != kind {
				continue
			}
			x, y := i%width, i/width
			var destination int
			if kind == '>' {
				destination = y*width + (x+1)%width
			} else {
				destination = ((y+1)%height)*width + x
			}
			if cells[destination] == '.' {
				next[i], next[destination] = '.', kind
				moved = true
			}
		}
		copy(cells, next)
	}
	return moved
}

func day25p01(r io.Reader) (string, error) {
	floor, err := grid.ReadBytes(r, ".>v")
	if err != nil {
		return "", err
	}
	cells := floor.Cells()
	next := make([]byte, len(cells))
	// Brent's cycle detection: the step is deterministic, so revisiting the
	// saved state proves the sea cucumbers never stop.
	saved := slices.Clone(cells)
	for step, power := 1, 1; ; step++ {
		if !stepSeaCucumbers(floor, next) {
			return strconv.Itoa(step), nil
		}
		if bytes.Equal(saved, cells) {
			return "", fmt.Errorf("sea cucumbers repeat a moving state after %d steps", step)
		}
		if step == power {
			copy(saved, cells)
			power *= 2
		}
	}
}
