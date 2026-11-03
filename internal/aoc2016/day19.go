package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func elfCount(r io.Reader) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	n, err := convert.ScanNumber[int64](bytes.TrimSpace(data))
	if err != nil {
		return 0, fmt.Errorf("elf count: %w", err)
	}
	if n < 1 {
		return 0, fmt.Errorf("elf count must be positive")
	}
	return n, nil
}

func day19p01(r io.Reader) (string, error) {
	n, err := elfCount(r)
	if err != nil {
		return "", err
	}
	power := int64(1)
	for power <= n/2 {
		power *= 2
	}
	return strconv.FormatInt(2*(n-power)+1, 10), nil
}

func day19p02(r io.Reader) (string, error) {
	n, err := elfCount(r)
	if err != nil {
		return "", err
	}
	power := int64(1)
	for power <= n/3 {
		power *= 3
	}
	var winner int64
	switch {
	case n == power:
		winner = n
	case n-power <= power:
		winner = n - power
	default:
		winner = (n - power) + (n - power - power)
	}
	return strconv.FormatInt(winner, 10), nil
}
