package aoc2016

import (
	"io"
	"strconv"
)

func safeCode(r io.Reader, eggs int64) (string, error) {
	program, err := parseAssembunny(r)
	if err != nil {
		return "", err
	}
	registers, err := runAssembunny(program, [4]int64{eggs, 0, 0, 0})
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(registers[0], 10), nil
}

func day23p01(r io.Reader) (string, error) {
	return safeCode(r, 7)
}

func day23p02(r io.Reader) (string, error) {
	return safeCode(r, 12)
}
