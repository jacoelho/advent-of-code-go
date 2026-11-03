package aoc2016

import (
	"bufio"
	"fmt"
	"io"
)

type messageFrequency int

const (
	mostCommon messageFrequency = iota
	leastCommon
)

func recoverMessage(r io.Reader, frequency messageFrequency) (string, error) {
	s := bufio.NewScanner(r)
	var columns [][26]int
	for s.Scan() {
		line := s.Bytes()
		if len(columns) == 0 {
			columns = make([][26]int, len(line))
		}
		if len(line) != len(columns) || len(line) == 0 {
			return "", fmt.Errorf("inconsistent message width")
		}
		for i, letter := range line {
			if letter < 'a' || letter > 'z' {
				return "", fmt.Errorf("invalid message letter %q", letter)
			}
			columns[i][letter-'a']++
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	message := make([]byte, len(columns))
	for i, counts := range columns {
		best := -1
		for letter, count := range counts {
			if count == 0 {
				continue
			}
			if best < 0 || (frequency == mostCommon && count > counts[best]) || (frequency == leastCommon && count < counts[best]) {
				best = letter
			}
		}
		message[i] = byte(best) + 'a'
	}
	return string(message), nil
}

func day06p01(r io.Reader) (string, error) {
	return recoverMessage(r, mostCommon)
}

func day06p02(r io.Reader) (string, error) {
	return recoverMessage(r, leastCommon)
}
