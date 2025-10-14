package aoc2019

import (
	"fmt"
	"io"
	"strconv"
)

func day5p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	computer.SetInput(1)

	if err := computer.Run(); err != nil {
		return "", err
	}

	output := computer.GetOutput()
	if len(output) == 0 {
		return "", fmt.Errorf("no output produced")
	}

	// All outputs except the last should be 0 (diagnostic tests passing)
	for i := 0; i < len(output)-1; i++ {
		if output[i] != 0 {
			return "", fmt.Errorf("diagnostic test failed at position %d: got %d, expected 0", i, output[i])
		}
	}

	return strconv.Itoa(output[len(output)-1]), nil
}

func day5p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	computer.SetInput(5)

	if err := computer.Run(); err != nil {
		return "", err
	}

	output := computer.GetOutput()
	if len(output) == 0 {
		return "", fmt.Errorf("no output produced")
	}

	return strconv.Itoa(output[len(output)-1]), nil
}
