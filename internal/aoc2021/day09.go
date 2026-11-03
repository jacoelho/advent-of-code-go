package aoc2021

import (
	"fmt"
	"io"
	"slices"
	"strconv"
)

func day09p01(r io.Reader) (string, error) {
	heights, err := parseDigitGrid(r)
	if err != nil {
		return "", err
	}
	risk := 0
	for index, height := range heights.values {
		low := true
		for neighbour := range heights.neighbours(index) {
			if heights.values[neighbour] <= height {
				low = false
				break
			}
		}
		if low {
			risk += int(height) + 1
		}
	}
	return strconv.Itoa(risk), nil
}

func day09p02(r io.Reader) (string, error) {
	heights, err := parseDigitGrid(r)
	if err != nil {
		return "", err
	}
	visited := make([]bool, len(heights.values))
	var sizes []int
	// Puzzle basins are connected components separated by height-nine ridges.
	for start, height := range heights.values {
		if height == 9 || visited[start] {
			continue
		}
		queue := []int{start}
		visited[start] = true
		for head := 0; head < len(queue); head++ {
			for neighbour := range heights.neighbours(queue[head]) {
				if heights.values[neighbour] != 9 && !visited[neighbour] {
					visited[neighbour] = true
					queue = append(queue, neighbour)
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
