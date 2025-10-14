package aoc2019

import (
	"cmp"
	"io"
	"slices"
	"strconv"
)

func day13p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	if err := computer.Run(); err != nil {
		return "", err
	}

	output := computer.GetOutput()

	blockCount := 0
	for tile := range slices.Chunk(output, 3) {
		if tile[2] == 2 {
			blockCount++
		}
	}

	return strconv.Itoa(blockCount), nil
}

func day13p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	if err := computer.SetMemory(0, 2); err != nil {
		return "", err
	}

	var score, ballX, paddleX int

	for {
		if err := computer.Run(); err != nil {
			return "", err
		}

		output := computer.GetOutput()

		for tile := range slices.Chunk(output, 3) {
			x, y, tileID := tile[0], tile[1], tile[2]

			if x == -1 && y == 0 {
				score = tileID
			} else {
				switch tileID {
				case 3:
					paddleX = x
				case 4:
					ballX = x
				}
			}
		}

		if computer.IsHalted() {
			break
		}

		direction := cmp.Compare(ballX, paddleX)
		computer.AddInput(direction)
	}

	return strconv.Itoa(score), nil
}
