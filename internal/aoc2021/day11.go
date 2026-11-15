package aoc2021

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func octopusEnergy(r io.Reader) (grid.Bytes, error) {
	energy, err := grid.ReadBytes(r, "0123456789")
	if err != nil {
		return grid.Bytes{}, err
	}
	cells := energy.Cells()
	for i, c := range cells {
		cells[i] = c - '0'
	}
	return energy, nil
}

func octopusStep(energy grid.Bytes) int {
	levels := energy.Cells()
	var queue []int
	for i := range levels {
		levels[i]++
		if levels[i] == 10 {
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		index := queue[head]
		levels[index] = 0
		for neighbour := range energy.Neighbours8(energy.Position(index)) {
			next := energy.Index(neighbour)
			if levels[next] == 0 {
				continue
			}
			levels[next]++
			if levels[next] == 10 {
				queue = append(queue, next)
			}
		}
	}
	return len(queue)
}

func day11p01(r io.Reader) (string, error) {
	energy, err := octopusEnergy(r)
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
	energy, err := octopusEnergy(r)
	if err != nil {
		return "", err
	}
	// Brent's cycle detection: the step is deterministic, so revisiting the
	// saved state proves the octopuses never synchronize.
	saved := slices.Clone(energy.Cells())
	for step, power := 1, 1; ; step++ {
		if octopusStep(energy) == len(saved) {
			return strconv.Itoa(step), nil
		}
		if bytes.Equal(saved, energy.Cells()) {
			return "", fmt.Errorf("octopus energy cycles without synchronization after %d steps", step)
		}
		if step == power {
			copy(saved, energy.Cells())
			power *= 2
		}
	}
}
