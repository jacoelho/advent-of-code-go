package aoc2019

import (
	"io"
	"strconv"
)

func day9p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	computer.SetInput(1)
	if err := computer.Run(); err != nil {
		return "", err
	}

	result, err := computer.LastOutput()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}

func day9p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	computer.SetInput(2)
	if err := computer.Run(); err != nil {
		return "", err
	}

	result, err := computer.LastOutput()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(result), nil
}
