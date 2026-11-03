package aoc2022

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

func treeVisibility(r io.Reader) (int, int64, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return 0, 0, err
	}
	rows := strings.Split(input, "\n")
	width, height := len(rows[0]), len(rows)
	if width == 0 || height > 1_000_000/width {
		return 0, 0, fmt.Errorf("forest exceeds one million cells")
	}
	trees := make([]byte, 0, width*height)
	for _, row := range rows {
		if len(row) != width {
			return 0, 0, fmt.Errorf("ragged forest")
		}
		for _, c := range []byte(row) {
			if c < '0' || c > '9' {
				return 0, 0, fmt.Errorf("invalid tree height")
			}
			trees = append(trees, c-'0')
		}
	}
	visible := make([]bool, len(trees))
	scores := make([]int64, len(trees))
	for i := range scores {
		scores[i] = 1
	}
	stack := make([]int, 0, max(width, height))
	scan := func(start, step, length int) {
		stack = stack[:0]
		for i := range length {
			index := start + i*step
			for len(stack) > 0 && trees[start+stack[len(stack)-1]*step] < trees[index] {
				stack = stack[:len(stack)-1]
			}
			distance := i
			if len(stack) == 0 {
				visible[index] = true
			} else {
				distance -= stack[len(stack)-1]
			}
			scores[index] *= int64(distance)
			stack = append(stack, i)
		}
	}
	for y := range height {
		scan(y*width, 1, width)
		scan((y+1)*width-1, -1, width)
	}
	for x := range width {
		scan(x, width, height)
		scan((height-1)*width+x, -width, height)
	}
	count := 0
	for _, v := range visible {
		if v {
			count++
		}
	}
	return count, slices.Max(scores), nil
}

func day08p01(r io.Reader) (string, error) {
	visible, _, err := treeVisibility(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(visible), nil
}

func day08p02(r io.Reader) (string, error) {
	_, score, err := treeVisibility(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(score, 10), nil
}
