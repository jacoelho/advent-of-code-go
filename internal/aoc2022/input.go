package aoc2022

import (
	"fmt"
	"io"
	"strings"
)

func puzzleInput(r io.Reader) (string, error) {
	const limit = 16 << 20
	data, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return "", err
	}
	if len(data) > limit {
		return "", fmt.Errorf("input exceeds 16 MiB")
	}
	input := strings.ReplaceAll(string(data), "\r\n", "\n")
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		return "", fmt.Errorf("empty input")
	}
	return input, nil
}
