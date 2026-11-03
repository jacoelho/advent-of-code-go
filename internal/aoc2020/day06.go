package aoc2020

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func customsAnswers(r io.Reader) (int, int, error) {
	const alphabet = uint32(1<<26 - 1)
	union, intersection := uint32(0), alphabet
	anyone, everyone, people := 0, 0, 0
	finishGroup := func() {
		if people > 0 {
			anyone += bits.OnesCount32(union)
			everyone += bits.OnesCount32(intersection)
		}
		union, intersection, people = 0, alphabet, 0
	}
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		if line == "" {
			finishGroup()
			continue
		}
		var answers uint32
		for _, c := range line {
			if c < 'a' || c > 'z' {
				return 0, 0, fmt.Errorf("invalid customs answer %q", c)
			}
			answers |= 1 << (c - 'a')
		}
		union |= answers
		intersection &= answers
		people++
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	finishGroup()
	return anyone, everyone, nil
}

func day06p01(r io.Reader) (string, error) {
	anyone, _, err := customsAnswers(r)
	return strconv.Itoa(anyone), err
}

func day06p02(r io.Reader) (string, error) {
	_, everyone, err := customsAnswers(r)
	return strconv.Itoa(everyone), err
}
