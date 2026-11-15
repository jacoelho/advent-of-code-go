package aoc2021

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func depthIncreases(r io.Reader, window int) (int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int64, error) {
		depth, err := convert.ScanNumber[int64](line)
		if err == nil && depth < 0 {
			err = fmt.Errorf("negative depth %d", depth)
		}
		return depth, err
	})
	previous := make([]int64, window)
	position, filled, increases := 0, 0, 0
	for depth := range s.Values() {
		if filled == window {
			if depth > previous[position] {
				increases++
			}
		} else {
			filled++
		}
		previous[position] = depth
		position = (position + 1) % window
	}
	return increases, s.Err()
}

func day01p01(r io.Reader) (string, error) {
	increases, err := depthIncreases(r, 1)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(increases), nil
}

func day01p02(r io.Reader) (string, error) {
	increases, err := depthIncreases(r, 3)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(increases), nil
}
