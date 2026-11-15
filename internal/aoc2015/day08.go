package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func hexDigit(b byte) bool {
	return '0' <= b && b <= '9' || 'a' <= b && b <= 'f' || 'A' <= b && b <= 'F'
}

func decodedLength(body string) (int, error) {
	characters := 0
	for i := 0; i < len(body); characters++ {
		switch {
		case body[i] == '"':
			return 0, fmt.Errorf("unescaped quote at %d", i+1)
		case body[i] != '\\':
			i++
		case i+1 < len(body) && (body[i+1] == '\\' || body[i+1] == '"'):
			i += 2
		case i+3 < len(body) && body[i+1] == 'x' && hexDigit(body[i+2]) && hexDigit(body[i+3]):
			i += 4
		default:
			return 0, fmt.Errorf("invalid escape at %d", i+1)
		}
	}
	return characters, nil
}

func stringLiteralOverhead(r io.Reader) (decoded, encoded int, err error) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if len(line) < 2 || line[0] != '"' || line[len(line)-1] != '"' {
			return 0, 0, fmt.Errorf("invalid string literal %q", line)
		}
		characters, err := decodedLength(line[1 : len(line)-1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid string literal %q: %w", line, err)
		}
		decoded += len(line) - characters
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
