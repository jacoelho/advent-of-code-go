package aoc2023

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const maxHeatGridCells = 250_000

type heatGrid struct {
	width, height int
	cells         []byte
}

type crucibleState struct {
	index int
	heat  int
}

func parseHeatGrid(r io.Reader) (heatGrid, error) {
	var result heatGrid
	s := scanner.NewScanner(r, func(line []byte) ([]byte, error) {
		if len(line) == 0 || result.height > 0 && len(line) != result.width {
			return nil, fmt.Errorf("empty or ragged heat grid at row %d", result.height+1)
		}
		if len(line) > maxHeatGridCells-len(result.cells) {
			return nil, fmt.Errorf("heat grid exceeds %d cells", maxHeatGridCells)
		}

		row := make([]byte, len(line))
		for i, ch := range line {
			if ch < '0' || ch > '9' {
				return nil, fmt.Errorf("invalid heat digit %q at row %d", ch, result.height+1)
			}
			row[i] = ch - '0'
		}
		return row, nil
	})
	for row := range s.Values() {
		result.width = len(row)
		result.height++
		result.cells = append(result.cells, row...)
	}
	if err := s.Err(); err != nil {
		return heatGrid{}, err
	}
	if result.height == 0 {
		return heatGrid{}, fmt.Errorf("empty heat grid")
	}
	return result, nil
}

func minimumHeat(heatMap heatGrid, minSteps, maxSteps int) (int, error) {
	distance := make([]int, 2*len(heatMap.cells))
	for i := range distance {
		distance[i] = math.MaxInt
	}
	queue := collections.NewHeap(func(a, b crucibleState) bool { return a.heat < b.heat })
	for axis := range 2 {
		distance[axis] = 0
		queue.Push(crucibleState{index: axis})
	}

	for current := range queue.PopSeq() {
		if current.heat != distance[current.index] {
			continue
		}
		cell := current.index / 2
		if cell == len(heatMap.cells)-1 && (cell != 0 || minSteps == 1) {
			return current.heat, nil
		}

		// A state ends a legal stride; its parity selects the next, perpendicular axis.
		dx, dy := 1, 0
		if current.index%2 == 1 {
			dx, dy = 0, 1
		}
		nextAxis := 1 - current.index%2
		for _, sign := range []int{-1, 1} {
			x, y := cell%heatMap.width, cell/heatMap.width
			heat := current.heat
			for length := 1; length <= maxSteps; length++ {
				x += sign * dx
				y += sign * dy
				if x < 0 || x >= heatMap.width || y < 0 || y >= heatMap.height {
					break
				}
				next := y*heatMap.width + x
				heat += int(heatMap.cells[next])
				if length < minSteps {
					continue
				}
				index := 2*next + nextAxis
				if heat < distance[index] {
					distance[index] = heat
					queue.Push(crucibleState{index: index, heat: heat})
				}
			}
		}
	}
	return 0, fmt.Errorf("heat-grid destination is unreachable with strides %d..%d", minSteps, maxSteps)
}

func crucibleHeat(r io.Reader, minSteps, maxSteps int) (string, error) {
	heatMap, err := parseHeatGrid(r)
	if err != nil {
		return "", err
	}
	heat, err := minimumHeat(heatMap, minSteps, maxSteps)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(heat), nil
}

func day17p01(r io.Reader) (string, error) {
	return crucibleHeat(r, 1, 3)
}

func day17p02(r io.Reader) (string, error) {
	return crucibleHeat(r, 4, 10)
}
