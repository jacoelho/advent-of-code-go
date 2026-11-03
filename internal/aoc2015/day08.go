package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func stringLiteralOverhead(r io.Reader) (decoded, encoded int, err error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if len(line) < 2 || line[0] != '"' || line[len(line)-1] != '"' {
			return 0, 0, fmt.Errorf("invalid string literal %q", line)
		}
		value, err := strconv.Unquote(line)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid string literal %q: %w", line, err)
		}
		decoded += len(line) - len(value)
		encoded += 2
		for i := range len(line) {
			if line[i] == '"' || line[i] == '\\' {
				encoded++
			}
		}
	}
	return decoded, encoded, s.Err()
}

func day08p01(r io.Reader) (string, error) {
	decoded, _, err := stringLiteralOverhead(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(decoded), nil
}

func day08p02(r io.Reader) (string, error) {
	_, encoded, err := stringLiteralOverhead(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(encoded), nil
}
