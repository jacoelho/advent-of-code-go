package aoc2019

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func checkTractorBeam(program []int, x, y int) int {
	computer := New(program)
	computer.SetInput(x, y)
	if err := computer.Run(); err != nil {
		panic(err)
	}
	output, err := computer.LastOutput()
	if err != nil {
		panic(err)
	}
	return output
}

func day19p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	count := 0
	for y := range 50 {
		for x := range 50 {
			if checkTractorBeam(program, x, y) == 1 {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

func findLeftEdge(program []int, y int) int {
	for x := 0; x < y*2; x++ {
		if checkTractorBeam(program, x, y) == 1 {
			return x
		}
	}
	return -1
}

func day19p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	const squareSize = 100

	y := search.BinarySearch(squareSize, 10000, func(y int) bool {
		leftX := findLeftEdge(program, y)
		if leftX == -1 {
			return false
		}

		topY := y - (squareSize - 1)
		if topY < 0 {
			return false
		}

		return checkTractorBeam(program, leftX, topY) == 1 &&
			checkTractorBeam(program, leftX+(squareSize-1), topY) == 1
	})

	if y > 10000 {
		return "not found", nil
	}

	leftX := findLeftEdge(program, y)

	topY := y - (squareSize - 1)
	return strconv.Itoa(leftX*10000 + topY), nil
}
