package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseSeatLayout(r io.Reader) ([]string, error) {
	width := 0
	s := scanner.NewScanner(r, func(line []byte) (string, error) {
		if len(line) == 0 {
			return "", fmt.Errorf("empty seat-layout row")
		}
		if width == 0 {
			width = len(line)
		}
		if len(line) != width {
			return "", fmt.Errorf("seat-layout row has width %d, want %d", len(line), width)
		}
		for _, cell := range line {
			if cell != '.' && cell != 'L' && cell != '#' {
				return "", fmt.Errorf("invalid seat-layout cell %q", cell)
			}
		}
		return string(line), nil
	})
	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty seat layout")
	}
	return rows, nil
}

func occupiedSeatsAtRest(rows []string, sightDistance, occupiedLimit int) int {
	width, height := len(rows[0]), len(rows)
	seats := make([]int, width*height)
	for i := range seats {
		seats[i] = -1
	}
	var occupied []bool
	for y, row := range rows {
		for x, cell := range row {
			if cell != '.' {
				seats[y*width+x] = len(occupied)
				occupied = append(occupied, cell == '#')
			}
		}
	}
	neighbours := make([][]int, len(occupied))
	for y, row := range rows {
		for x, cell := range row {
			if cell == '.' {
				continue
			}
			seat := seats[y*width+x]
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}
					for distance := 1; distance <= sightDistance; distance++ {
						nx, ny := x+dx*distance, y+dy*distance
						if nx < 0 || nx >= width || ny < 0 || ny >= height {
							break
						}
						if other := seats[ny*width+nx]; other >= 0 {
							neighbours[seat] = append(neighbours[seat], other)
							break
						}
					}
				}
			}
		}
	}
	next := make([]bool, len(occupied))
	for {
		changed, total := false, 0
		for seat, taken := range occupied {
			count := 0
			for _, other := range neighbours[seat] {
				if occupied[other] {
					count++
				}
			}
			next[seat] = taken
			if !taken && count == 0 {
				next[seat] = true
			}
			if taken && count >= occupiedLimit {
				next[seat] = false
			}
			if next[seat] {
				total++
			}
			changed = changed || next[seat] != taken
		}
		if !changed {
			return total
		}
		occupied, next = next, occupied
	}
}

func day11p01(r io.Reader) (string, error) {
	rows, err := parseSeatLayout(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(occupiedSeatsAtRest(rows, 1, 4)), nil
}

func day11p02(r io.Reader) (string, error) {
	rows, err := parseSeatLayout(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(occupiedSeatsAtRest(rows, max(len(rows), len(rows[0])), 5)), nil
}
