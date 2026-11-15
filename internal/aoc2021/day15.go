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
	base, err := grid.ReadBytes(r, "123456789")
	if err != nil {
		return 0, err
	}
	baseRisks, baseWidth, baseHeight := base.Cells(), base.Width(), base.Height()
	width, height := baseWidth*factor, baseHeight*factor
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
			if !next.Within(width, height) {
				continue
			}
			index := next.Y*width + next.X
			risk := (int(baseRisks[(next.Y%baseHeight)*baseWidth+next.X%baseWidth]-'0')+next.X/baseWidth+next.Y/baseHeight-1)%9 + 1
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
