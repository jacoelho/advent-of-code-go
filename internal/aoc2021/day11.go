package aoc2021

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func octopusStep(energy digitGrid) int {
	var queue []int
	for i := range energy.values {
		energy.values[i]++
		if energy.values[i] == 10 {
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		index := queue[head]
		energy.values[index] = 0
		point := grid.Position2D[int]{X: index % energy.width, Y: index / energy.width}
		for neighbour := range grid.Neighbours8(point) {
			if neighbour.X >= energy.width || neighbour.Y >= energy.height {
				continue
			}
			next := neighbour.Y*energy.width + neighbour.X
			if energy.values[next] == 0 {
				continue
			}
			energy.values[next]++
			if energy.values[next] == 10 {
				queue = append(queue, next)
			}
		}
	}
	return len(queue)
}

func day11p01(r io.Reader) (string, error) {
	energy, err := parseDigitGrid(r)
	if err != nil {
		return "", err
	}
	flashes := 0
	for range 100 {
		flashes += octopusStep(energy)
	}
	return strconv.Itoa(flashes), nil
}

func day11p02(r io.Reader) (string, error) {
	energy, err := parseDigitGrid(r)
	if err != nil {
		return "", err
	}
	seen := map[string]bool{string(energy.values): true}
	stored := len(energy.values)
	for step := 1; ; step++ {
		if octopusStep(energy) == len(energy.values) {
			return strconv.Itoa(step), nil
		}
		state := string(energy.values)
		if seen[state] {
			return "", fmt.Errorf("octopus energy cycles without synchronization")
		}
		if len(seen) == 1_000_000 || stored > 64<<20-len(state) {
			return "", fmt.Errorf("octopus synchronization exceeds state budget")
		}
		seen[state] = true
		stored += len(state)
	}
}
