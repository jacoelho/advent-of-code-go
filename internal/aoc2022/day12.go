package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func hillDistance(r io.Reader, anyStart bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	rows := strings.Split(input, "\n")
	width, height := len(rows[0]), len(rows)
	if width == 0 || height > 1_000_000/width {
		return "", fmt.Errorf("heightmap exceeds one million cells")
	}
	heights := make([]byte, 0, width*height)
	start, end := -1, -1
	for _, row := range rows {
		if len(row) != width {
			return "", fmt.Errorf("ragged heightmap")
		}
		for _, c := range []byte(row) {
			switch c {
			case 'S':
				if start >= 0 {
					return "", fmt.Errorf("multiple starts")
				}
				start = len(heights)
				c = 'a'
			case 'E':
				if end >= 0 {
					return "", fmt.Errorf("multiple ends")
				}
				end = len(heights)
				c = 'z'
			}
			if c < 'a' || c > 'z' {
				return "", fmt.Errorf("invalid height")
			}
			heights = append(heights, c-'a')
		}
	}
	if start < 0 || end < 0 {
		return "", fmt.Errorf("missing start or end")
	}
	distances := make([]int, len(heights))
	for i := range distances {
		distances[i] = -1
	}
	distances[end] = 0
	queue := collections.NewDeque[int](64)
	queue.PushBack(end)
	for queue.Size() > 0 {
		at, _ := queue.PopFront()
		if at == start || anyStart && heights[at] == 0 {
			return strconv.Itoa(distances[at]), nil
		}
		position := grid.Position2D[int]{X: at % width, Y: at / width}
		for next := range grid.Neighbours4(position) {
			if next.X >= width || next.Y >= height {
				continue
			}
			index := next.Y*width + next.X
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
