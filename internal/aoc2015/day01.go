package aoc2015

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func santaFloors(r io.Reader) (floor, basement int, err error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}
	for i, direction := range bytes.TrimSpace(input) {
		switch direction {
		case '(':
			floor++
		case ')':
			floor--
		default:
			return 0, 0, fmt.Errorf("invalid floor direction %q", direction)
		}
		if floor < 0 && basement == 0 {
			basement = i + 1
		}
	}
	return floor, basement, nil
}

func day01p01(r io.Reader) (string, error) {
	floor, _, err := santaFloors(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(floor), nil
}

func day01p02(r io.Reader) (string, error) {
	_, basement, err := santaFloors(r)
	if err != nil {
		return "", err
	}
	if basement == 0 {
		return "", fmt.Errorf("directions never enter the basement")
	}
	return strconv.Itoa(basement), nil
}
