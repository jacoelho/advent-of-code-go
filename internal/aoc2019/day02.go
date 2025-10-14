package aoc2019

import (
	"fmt"
	"io"
	"strconv"
)

func runIntcodeWithInputs(program []int, noun, verb int) (int, error) {
	computer := New(program)

	if err := computer.SetMemory(1, noun); err != nil {
		return 0, err
	}
	if err := computer.SetMemory(2, verb); err != nil {
		return 0, err
	}

	if err := computer.Run(); err != nil {
		return 0, err
	}

	return computer.Memory()[0], nil
}

func day2p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	result, err := runIntcodeWithInputs(program, 12, 2)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func day2p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	const target = 19690720

	for noun := range 100 {
		for verb := range 100 {
			result, err := runIntcodeWithInputs(program, noun, verb)
			if err != nil {
				continue
			}
			if result == target {
				return strconv.Itoa(100*noun + verb), nil
			}
		}
	}

	return "", fmt.Errorf("no solution found for target %d", target)
}
