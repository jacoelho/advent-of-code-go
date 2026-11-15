package aoc2020

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func customsAnswers(r io.Reader) (int, int, error) {
	anyone, everyone := 0, 0
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for group := range xiter.Blocks(s.Values()) {
		if len(group) == 0 {
			continue
		}
		union, intersection := uint32(0), uint32(1<<26-1)
		for _, line := range group {
			var answers uint32
			for _, c := range line {
				if c < 'a' || c > 'z' {
					return 0, 0, fmt.Errorf("invalid customs answer %q", c)
				}
				answers |= 1 << (c - 'a')
			}
			union |= answers
			intersection &= answers
		}
		anyone += bits.OnesCount32(union)
		everyone += bits.OnesCount32(intersection)
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
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
