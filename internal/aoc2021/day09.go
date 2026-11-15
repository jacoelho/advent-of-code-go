package aoc2021

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func day09p01(r io.Reader) (string, error) {
	heights, err := grid.ReadBytes(r, "0123456789")
	if err != nil {
		return "", err
	}
	risk := 0
	for index, height := range heights.Cells() {
		low := true
		for neighbour := range heights.Neighbours4(heights.Position(index)) {
			if heights.At(neighbour) <= height {
				low = false
				break
			}
		}
		if low {
			risk += int(height-'0') + 1
		}
	}
	return strconv.Itoa(risk), nil
}

func day09p02(r io.Reader) (string, error) {
	heights, err := grid.ReadBytes(r, "0123456789")
	if err != nil {
		return "", err
	}
	visited := make([]bool, len(heights.Cells()))
	var sizes []int
	// Puzzle basins are connected components separated by height-nine ridges.
	for start, height := range heights.Cells() {
		if height == '9' || visited[start] {
			continue
		}
		queue := []int{start}
		visited[start] = true
		for head := 0; head < len(queue); head++ {
			for neighbour := range heights.Neighbours4(heights.Position(queue[head])) {
				index := heights.Index(neighbour)
				if heights.At(neighbour) != '9' && !visited[index] {
					visited[index] = true
					queue = append(queue, index)
				}
			}
		}
		sizes = append(sizes, len(queue))
	}
	if len(sizes) < 3 {
		return "", fmt.Errorf("heightmap has fewer than three basins")
	}
	slices.Sort(sizes)
	last := len(sizes) - 1
	return strconv.Itoa(sizes[last] * sizes[last-1] * sizes[last-2]), nil
}
