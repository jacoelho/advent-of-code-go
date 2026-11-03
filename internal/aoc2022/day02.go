package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func strategyScore(r io.Reader, outcomes bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	score := 0
	for line := range strings.SplitSeq(input, "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 || len(fields[0]) != 1 || len(fields[1]) != 1 || fields[0][0] < 'A' || fields[0][0] > 'C' || fields[1][0] < 'X' || fields[1][0] > 'Z' {
			return "", fmt.Errorf("invalid round %q", line)
		}
		opponent, second := int(fields[0][0]-'A'), int(fields[1][0]-'X')
		shape := second
		if outcomes {
			shape = (opponent + second + 2) % 3
		}
		score += shape + 1 + ((shape-opponent+4)%3)*3
	}
	return strconv.Itoa(score), nil
}

func day02p01(r io.Reader) (string, error) {
	return strategyScore(r, false)
}

func day02p02(r io.Reader) (string, error) {
	return strategyScore(r, true)
}
