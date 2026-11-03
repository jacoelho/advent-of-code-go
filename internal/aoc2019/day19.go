package aoc2019

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func checkTractorBeam(program []int, x, y int) (bool, error) {
	computer := New(program)
	computer.SetInput(x, y)
	if err := computer.Run(); err != nil {
		return false, fmt.Errorf("beam query at (%d,%d): %w", x, y, err)
	}
	output := computer.GetOutput()
	if !computer.IsHalted() || len(output) != 1 || output[0] < 0 || output[0] > 1 {
		return false, fmt.Errorf("beam query at (%d,%d) needs one zero/one reply followed by halt", x, y)
	}
	return output[0] == 1, nil
}

func day19p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	count := 0
	for y := range 50 {
		for x := range 50 {
			pulled, err := checkTractorBeam(program, x, y)
			if err != nil {
				return "", err
			}
			if pulled {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

func findLeftEdge(program []int, y int) (int, error) {
	for x := 0; x < y*2; x++ {
		pulled, err := checkTractorBeam(program, x, y)
		if err != nil {
			return 0, err
		}
		if pulled {
			return x, nil
		}
	}
	return -1, nil
}

func day19p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	const squareSize = 100

	y, err := search.BinarySearch(squareSize-1, 10000, func(y int) (bool, error) {
		leftX, err := findLeftEdge(program, y)
		if err != nil {
			return false, err
		}
		if leftX == -1 {
			return false, nil
		}

		topY := y - (squareSize - 1)
		pulled, err := checkTractorBeam(program, leftX, topY)
		if err != nil || !pulled {
			return false, err
		}
		return checkTractorBeam(program, leftX+(squareSize-1), topY)
	})
	if err != nil {
		return "", err
	}

	leftX, err := findLeftEdge(program, y)
	if err != nil {
		return "", err
	}

	topY := y - (squareSize - 1)
	return strconv.Itoa(leftX*10000 + topY), nil
}
