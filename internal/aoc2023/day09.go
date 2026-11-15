package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseSequence(line []byte) ([]int, error) {
	fields := strings.Fields(string(line))
	if len(fields) == 0 {
		return nil, fmt.Errorf("empty sequence")
	}
	seq := make([]int, len(fields))
	for i, field := range fields {
		value, err := strconv.Atoi(field)
		if err != nil {
			return nil, fmt.Errorf("sequence %q: %w", line, err)
		}
		seq[i] = value
	}
	return seq, nil
}

// extrapolateForward returns the next value of seq; it overwrites seq with its
// difference rows.
func extrapolateForward(seq []int) int {
	next := 0
	for n := len(seq); n > 0; n-- {
		next += seq[n-1]
		allZero := true
		for i := range n - 1 {
			seq[i] = seq[i+1] - seq[i]
			allZero = allZero && seq[i] == 0
		}
		if allZero {
			break
		}
	}
	return next
}

func extrapolate(r io.Reader, backward bool) (string, error) {
	s := scanner.NewScanner(r, parseSequence)
	total := xiter.From(s.Values()).SumBy(func(seq []int) int {
		if backward {
			slices.Reverse(seq)
		}
		return extrapolateForward(seq)
	})
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day09p01(r io.Reader) (string, error) {
	return extrapolate(r, false)
}

func day09p02(r io.Reader) (string, error) {
	return extrapolate(r, true)
}
