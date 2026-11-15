package aoc2022

import (
	"fmt"
	"io"
	"strings"
)

func puzzleInput(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	input := strings.ReplaceAll(string(data), "\r\n", "\n")
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		return "", fmt.Errorf("empty input")
	}
	return input, nil
}
