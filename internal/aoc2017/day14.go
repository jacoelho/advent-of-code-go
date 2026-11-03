package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type diskGrid [128][128]bool

func diskUsage(r io.Reader) (diskGrid, error) {
	key, err := io.ReadAll(r)
	if err != nil {
		return diskGrid{}, err
	}

	key = bytes.TrimSpace(key)
	if len(key) == 0 || bytes.ContainsAny(key, " \t\r\n") {
		return diskGrid{}, fmt.Errorf("invalid disk key")
	}

	var grid diskGrid
	for row := range grid {
		input := append(bytes.Clone(key), '-')
		input = strconv.AppendInt(input, int64(row), 10)
		hash := knotHash(input)
		for column := range grid[row] {
			grid[row][column] = hash[column/8]&(1<<(7-column%8)) != 0
		}
	}

	return grid, nil
}

func diskRegions(grid diskGrid) int {
	regions := 0
	stack := collections.NewStack[[2]int]()

	for y := range grid {
		for x := range grid[y] {
			if !grid[y][x] {
				continue
			}
			regions++
			grid[y][x] = false
			stack.Push([2]int{x, y})
			for !stack.IsEmpty() {
				point, _ := stack.Pop()
				for _, offset := range [...][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
					nx, ny := point[0]+offset[0], point[1]+offset[1]
					if nx >= 0 && nx < 128 && ny >= 0 && ny < 128 && grid[ny][nx] {
						grid[ny][nx] = false
						stack.Push([2]int{nx, ny})
					}
				}
			}
		}
	}

	return regions
}

func day14p01(r io.Reader) (string, error) {
	grid, err := diskUsage(r)
	if err != nil {
		return "", err
	}

	used := 0
	for _, row := range grid {
		for _, occupied := range row {
			if occupied {
				used++
			}
		}
	}

	return strconv.Itoa(used), nil
}

func day14p02(r io.Reader) (string, error) {
	grid, err := diskUsage(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(diskRegions(grid)), nil
}
