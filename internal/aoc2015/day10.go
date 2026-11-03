package aoc2015

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func lookAndSay(sequence []byte) []byte {
	next := make([]byte, 0, len(sequence)*2)
	for start := 0; start < len(sequence); {
		end := start + 1
		for end < len(sequence) && sequence[end] == sequence[start] {
			end++
		}
		next = strconv.AppendInt(next, int64(end-start), 10)
		next = append(next, sequence[start])
		start = end
	}
	return next
}

func lookAndSayLength(r io.Reader, rounds int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	sequence := bytes.TrimSpace(input)
	if len(sequence) == 0 {
		return "", fmt.Errorf("empty look-and-say sequence")
	}
	for _, digit := range sequence {
		if digit < '0' || digit > '9' {
			return "", fmt.Errorf("invalid look-and-say digit %q", digit)
		}
	}
	for range rounds {
		sequence = lookAndSay(sequence)
	}
	return strconv.Itoa(len(sequence)), nil
}

func day10p01(r io.Reader) (string, error) {
	return lookAndSayLength(r, 40)
}

func day10p02(r io.Reader) (string, error) {
	return lookAndSayLength(r, 50)
}
