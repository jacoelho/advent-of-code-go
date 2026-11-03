package aoc2020

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type pocketCube [4]int

func parsePocketDimension(r io.Reader) (map[pocketCube]struct{}, error) {
	active := make(map[pocketCube]struct{})
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	width, y := 0, 0
	for line := range s.Values() {
		if y == 0 {
			width = len(line)
		}
		if len(line) == 0 || len(line) != width {
			return nil, fmt.Errorf("invalid pocket dimension row width at row %d", y+1)
		}
		for x, c := range line {
			switch c {
			case '#':
				active[pocketCube{x, y, 0, 0}] = struct{}{}
			case '.':
			default:
				return nil, fmt.Errorf("invalid pocket dimension cell %q", c)
			}
		}
		y++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if y == 0 {
		return nil, fmt.Errorf("empty pocket dimension")
	}
	return active, nil
}

func simulatePocketDimension(r io.Reader, dimensions int) (string, error) {
	active, err := parsePocketDimension(r)
	if err != nil {
		return "", err
	}
	wLimit := dimensions - 3
	var offsets []pocketCube
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			for z := -1; z <= 1; z++ {
				for w := -wLimit; w <= wLimit; w++ {
					offset := pocketCube{x, y, z, w}
					if offset != (pocketCube{}) {
						offsets = append(offsets, offset)
					}
				}
			}
		}
	}
	for range 6 {
		neighbours := make(map[pocketCube]int)
		for cube := range active {
			for _, offset := range offsets {
				neighbour := pocketCube{cube[0] + offset[0], cube[1] + offset[1], cube[2] + offset[2], cube[3] + offset[3]}
				neighbours[neighbour]++
			}
		}
		next := make(map[pocketCube]struct{})
		for cube, count := range neighbours {
			_, occupied := active[cube]
			if count == 3 || occupied && count == 2 {
				next[cube] = struct{}{}
			}
		}
		active = next
	}
	return strconv.Itoa(len(active)), nil
}

func day17p01(r io.Reader) (string, error) {
	return simulatePocketDimension(r, 3)
}

func day17p02(r io.Reader) (string, error) {
	return simulatePocketDimension(r, 4)
}
