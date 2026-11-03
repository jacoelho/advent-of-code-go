package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func readPolymer(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	polymer := bytes.TrimSpace(data)
	for _, unit := range polymer {
		if unit < 'A' || unit > 'Z' && unit < 'a' || unit > 'z' {
			return nil, fmt.Errorf("invalid polymer unit %q", unit)
		}
	}
	return polymer, nil
}

func reactedLength(polymer []byte, omitted byte) int {
	stack := make([]byte, 0, len(polymer))
	for _, unit := range polymer {
		if unit|32 == omitted {
			continue
		}
		if len(stack) > 0 && stack[len(stack)-1]^unit == 32 {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, unit)
		}
	}
	return len(stack)
}

func day05p01(r io.Reader) (string, error) {
	polymer, err := readPolymer(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(reactedLength(polymer, 0)), nil
}

func day05p02(r io.Reader) (string, error) {
	polymer, err := readPolymer(r)
	if err != nil {
		return "", err
	}

	shortest := len(polymer)
	for unit := byte('a'); unit <= 'z'; unit++ {
		shortest = min(shortest, reactedLength(polymer, unit))
	}
	return strconv.Itoa(shortest), nil
}
