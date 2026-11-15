package aoc2022

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func hillDistance(r io.Reader, anyStart bool) (string, error) {
	hill, err := grid.ReadBytes(r, "SEabcdefghijklmnopqrstuvwxyz")
	if err != nil {
		return "", err
	}
	startPosition, err := hill.Find('S')
	if err != nil {
		return "", err
	}
	endPosition, err := hill.Find('E')
	if err != nil {
		return "", err
	}
	heights := hill.Cells()
	start, end := hill.Index(startPosition), hill.Index(endPosition)
	heights[start], heights[end] = 'a', 'z'
	distances := make([]int, len(heights))
	for i := range distances {
		distances[i] = -1
	}
	distances[end] = 0
	queue := collections.NewDeque[int](64)
	queue.PushBack(end)
	for queue.Size() > 0 {
		at, _ := queue.PopFront()
		if at == start || anyStart && heights[at] == 'a' {
			return strconv.Itoa(distances[at]), nil
		}
		for next := range hill.Neighbours4(hill.Position(at)) {
			index := hill.Index(next)
			if distances[index] >= 0 || int(heights[at])-int(heights[index]) > 1 {
				continue
			}
			distances[index] = distances[at] + 1
			queue.PushBack(index)
		}
	}
	return "", fmt.Errorf("hill destination is unreachable")
}

func day12p01(r io.Reader) (string, error) {
	return hillDistance(r, false)
}

func day12p02(r io.Reader) (string, error) {
	return hillDistance(r, true)
}
