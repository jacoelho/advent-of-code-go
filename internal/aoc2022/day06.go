package aoc2022

import (
	"fmt"
	"io"
	"strconv"
)

func markerPosition(r io.Reader, width int) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	for i, c := range []byte(input) {
		if c < 'a' || c > 'z' {
			return "", fmt.Errorf("invalid datastream character %q at %d", c, i)
		}
	}
	var counts [256]int
	duplicates := 0
	for i, c := range []byte(input) {
		counts[c]++
		if counts[c] == 2 {
			duplicates++
		}
		if i >= width {
			old := input[i-width]
			if counts[old] == 2 {
				duplicates--
			}
			counts[old]--
		}
		if i+1 >= width && duplicates == 0 {
			return strconv.Itoa(i + 1), nil
		}
	}
	return "", fmt.Errorf("datastream has no marker of length %d", width)
}

func day06p01(r io.Reader) (string, error) {
	return markerPosition(r, 4)
}

func day06p02(r io.Reader) (string, error) {
	return markerPosition(r, 14)
}
