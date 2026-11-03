package aoc2021

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type riskLocation struct {
	index, cost int
}

func lowestRisk(r io.Reader, factor int) (int, error) {
	base, err := parseDigitGrid(r)
	if err != nil {
		return 0, err
	}
	for _, risk := range base.values {
		if risk == 0 {
			return 0, fmt.Errorf("risk values must be between one and nine")
		}
	}
	if factor < 1 || len(base.values) > maxDigitGridCells/factor/factor {
		return 0, fmt.Errorf("expanded risk grid exceeds %d cells", maxDigitGridCells)
	}
	width, height := base.width*factor, base.height*factor
	distances := make([]int, width*height)
	for i := range distances {
		distances[i] = math.MaxInt
	}
	distances[0] = 0
	queue := collections.NewHeap(func(a, b riskLocation) bool { return a.cost < b.cost })
	queue.Push(riskLocation{})
	for current := range queue.PopSeq() {
		if current.cost != distances[current.index] {
			continue
		}
		if current.index == len(distances)-1 {
			return current.cost, nil
		}
		point := grid.Position2D[int]{X: current.index % width, Y: current.index / width}
		for next := range grid.Neighbours4(point) {
			if next.X >= width || next.Y >= height {
				continue
			}
			index := next.Y*width + next.X
			risk := (int(base.values[(next.Y%base.height)*base.width+next.X%base.width])+next.X/base.width+next.Y/base.height-1)%9 + 1
			cost := current.cost + risk
			if cost < distances[index] {
				distances[index] = cost
				queue.Push(riskLocation{index: index, cost: cost})
			}
		}
	}
	return 0, fmt.Errorf("risk destination is unreachable")
}

func day15p01(r io.Reader) (string, error) {
	risk, err := lowestRisk(r, 1)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(risk), nil
}

func day15p02(r io.Reader) (string, error) {
	risk, err := lowestRisk(r, 5)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(risk), nil
}
