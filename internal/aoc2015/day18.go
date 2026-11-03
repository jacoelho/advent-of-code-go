package aoc2015

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type lightCornerMode uint8

const (
	normalLightCorners lightCornerMode = iota
	stuckLightCorners
)

func animateLights(r io.Reader, steps int, corners lightCornerMode) (string, error) {
	s := scanner.NewScanner(r, func(line []byte) ([]bool, error) {
		row := make([]bool, len(line))
		for x, cell := range line {
			if cell != '#' && cell != '.' {
				return nil, fmt.Errorf("invalid light %q", cell)
			}
			row[x] = cell == '#'
		}
		return row, nil
	})
	var grid [][]bool
	width := 0
	for row := range s.Values() {
		if len(grid) == 0 {
			width = len(row)
		}
		if width == 0 || len(row) != width {
			return "", fmt.Errorf("invalid light grid row width")
		}
		grid = append(grid, row)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(grid) == 0 {
		return "", fmt.Errorf("empty light grid")
	}
	height := len(grid)
	forceCorners := func(grid [][]bool) {
		if corners == stuckLightCorners {
			grid[0][0], grid[0][width-1] = true, true
			grid[height-1][0], grid[height-1][width-1] = true, true
		}
	}
	forceCorners(grid)
	next := make([][]bool, height)
	for y := range next {
		next[y] = make([]bool, width)
	}
	for range steps {
		for y := range height {
			for x := range width {
				neighbours := 0
				for ny := max(0, y-1); ny <= min(height-1, y+1); ny++ {
					for nx := max(0, x-1); nx <= min(width-1, x+1); nx++ {
						if (nx != x || ny != y) && grid[ny][nx] {
							neighbours++
						}
					}
				}
				next[y][x] = neighbours == 3 || grid[y][x] && neighbours == 2
			}
		}
		forceCorners(next)
		grid, next = next, grid
	}
	on := 0
	for _, row := range grid {
		for _, light := range row {
			if light {
				on++
			}
		}
	}
	return strconv.Itoa(on), nil
}

func day18p01(r io.Reader) (string, error) {
	return animateLights(r, 100, normalLightCorners)
}

func day18p02(r io.Reader) (string, error) {
	return animateLights(r, 100, stuckLightCorners)
}
