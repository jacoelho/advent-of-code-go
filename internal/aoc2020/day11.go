package aoc2020

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func occupiedSeatsAtRest(layout grid.Bytes, sightDistance, occupiedLimit int) int {
	width, height := layout.Width(), layout.Height()
	cells := layout.Cells()
	seats := make([]int, len(cells))
	var occupied []bool
	for i, cell := range cells {
		seats[i] = -1
		if cell != '.' {
			seats[i] = len(occupied)
			occupied = append(occupied, cell == '#')
		}
	}
	neighbours := make([][]int, len(occupied))
	for i, seat := range seats {
		if seat < 0 {
			continue
		}
		x, y := i%width, i/width
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
	layout, err := grid.ReadBytes(r, ".L#")
	if err != nil {
		return "", err
	}
	return strconv.Itoa(occupiedSeatsAtRest(layout, 1, 4)), nil
}

func day11p02(r io.Reader) (string, error) {
	layout, err := grid.ReadBytes(r, ".L#")
	if err != nil {
		return "", err
	}
	return strconv.Itoa(occupiedSeatsAtRest(layout, max(layout.Width(), layout.Height()), 5)), nil
}
