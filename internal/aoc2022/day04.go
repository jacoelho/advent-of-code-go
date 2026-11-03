package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func assignmentPairs(r io.Reader, overlap bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	count := 0
	for line := range strings.SplitSeq(input, "\n") {
		pairs := strings.Split(line, ",")
		if len(pairs) != 2 {
			return "", fmt.Errorf("invalid assignment %q", line)
		}
		var bounds [4]int64
		for i, pair := range pairs {
			ends := strings.Split(pair, "-")
			if len(ends) != 2 {
				return "", fmt.Errorf("invalid range %q", pair)
			}
			for j, end := range ends {
				value, err := convert.ScanNumber[int64]([]byte(end))
				if err != nil {
					return "", err
				}
				if value < 0 {
					return "", fmt.Errorf("negative assignment")
				}
				bounds[2*i+j] = value
			}
			if bounds[2*i] > bounds[2*i+1] {
				return "", fmt.Errorf("reversed range")
			}
		}
		a, b, c, d := bounds[0], bounds[1], bounds[2], bounds[3]
		match := a <= c && b >= d || c <= a && d >= b
		if overlap {
			match = a <= d && c <= b
		}
		if match {
			count++
		}
	}
	return strconv.Itoa(count), nil
}

func day04p01(r io.Reader) (string, error) {
	return assignmentPairs(r, false)
}

func day04p02(r io.Reader) (string, error) {
	return assignmentPairs(r, true)
}
