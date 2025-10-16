package aoc2023

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseSequence(line []byte) ([]int, error) {
	return convert.ExtractDigits[int](string(line)), nil
}

func differences(seq []int) []int {
	if len(seq) <= 1 {
		return []int{}
	}
	result := make([]int, len(seq)-1)
	for i := 0; i < len(seq)-1; i++ {
		result[i] = seq[i+1] - seq[i]
	}
	return result
}

func allZeros(seq []int) bool {
	return xslices.Every(func(v int) bool { return v == 0 }, seq)
}

func extrapolateForward(seq []int) int {
	if allZeros(seq) {
		return 0
	}
	diffs := differences(seq)
	return seq[len(seq)-1] + extrapolateForward(diffs)
}

func extrapolateBackward(seq []int) int {
	if allZeros(seq) {
		return 0
	}
	diffs := differences(seq)
	return seq[0] - extrapolateBackward(diffs)
}

func day09p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseSequence)
	total := xiter.Sum(xiter.Map(extrapolateForward, s.Values()))
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day09p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseSequence)
	total := xiter.Sum(xiter.Map(extrapolateBackward, s.Values()))
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}
