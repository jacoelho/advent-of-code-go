package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type seaFloor struct {
	cells         []byte
	width, height int
}

func parseSeaFloor(r io.Reader) (seaFloor, error) {
	var floor seaFloor
	s := bufio.NewScanner(r)
	for s.Scan() {
		row := s.Bytes()
		if floor.height == 0 {
			floor.width = len(row)
		}
		if len(row) == 0 || len(row) != floor.width {
			return floor, fmt.Errorf("sea floor must be a nonempty rectangle")
		}
		for _, c := range row {
			if c != '.' && c != '>' && c != 'v' {
				return floor, fmt.Errorf("invalid sea cucumber %q", c)
			}
		}
		floor.cells = append(floor.cells, row...)
		floor.height++
		if len(floor.cells) > 1_000_000 {
			return floor, fmt.Errorf("sea floor exceeds one million cells")
		}
	}
	if err := s.Err(); err != nil {
		return floor, err
	}
	if len(floor.cells) == 0 {
		return floor, fmt.Errorf("empty sea floor")
	}
	return floor, nil
}

func (floor *seaFloor) step(next []byte) bool {
	moved := false
	for _, kind := range [...]byte{'>', 'v'} {
		copy(next, floor.cells)
		for i, c := range floor.cells {
			if c != kind {
				continue
			}
			x, y := i%floor.width, i/floor.width
			var destination int
			if kind == '>' {
				destination = y*floor.width + (x+1)%floor.width
			} else {
				destination = ((y+1)%floor.height)*floor.width + x
			}
			if floor.cells[destination] == '.' {
				next[i], next[destination] = '.', kind
				moved = true
			}
		}
		copy(floor.cells, next)
	}
	return moved
}

func day25p01(r io.Reader) (string, error) {
	floor, err := parseSeaFloor(r)
	if err != nil {
		return "", err
	}
	next := make([]byte, len(floor.cells))
	seen := map[string]struct{}{string(floor.cells): {}}
	work := 100_000_000
	for step := 1; ; step++ {
		work -= 2 * len(floor.cells)
		if work < 0 {
			return "", fmt.Errorf("sea cucumber simulation exceeds 100 million cell visits")
		}
		if !floor.step(next) {
			return strconv.Itoa(step), nil
		}
		key := string(floor.cells)
		if _, exists := seen[key]; exists {
			return "", fmt.Errorf("sea cucumbers repeat a moving state")
		}
		if (len(seen)+1)*(len(floor.cells)+48) > 64<<20 {
			return "", fmt.Errorf("sea cucumber history exceeds 64 MiB")
		}
		seen[key] = struct{}{}
	}
}
