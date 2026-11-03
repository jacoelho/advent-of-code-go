package aoc2016

import (
	"io"
	"strconv"
)

func day12p01(r io.Reader) (string, error) {
	program, err := parseAssembunny(r)
	if err != nil {
		return "", err
	}
	registers, err := runAssembunny(program, [4]int64{})
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(registers[0], 10), nil
}

func day12p02(r io.Reader) (string, error) {
	program, err := parseAssembunny(r)
	if err != nil {
		return "", err
	}
	registers, err := runAssembunny(program, [4]int64{2: 1})
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(registers[0], 10), nil
}
