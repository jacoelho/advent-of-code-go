package aoc2017

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type jumpPolicy int

const (
	increaseJump jumpPolicy = iota
	decreaseLargeJump
)

func jumpSteps(r io.Reader, policy jumpPolicy) (string, error) {
	s := scanner.NewScanner(r, convert.ScanNumber[int])
	offsets := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}

	steps := 0
	for ip := 0; ip >= 0 && ip < len(offsets); steps++ {
		offset := offsets[ip]
		if offset >= len(offsets)-ip || offset < -ip {
			return strconv.Itoa(steps + 1), nil
		}
		if policy == decreaseLargeJump && offset >= 3 {
			offsets[ip]--
		} else {
			offsets[ip]++
		}
		ip += offset
	}

	return strconv.Itoa(steps), nil
}

func day05p01(r io.Reader) (string, error) {
	return jumpSteps(r, increaseJump)
}

func day05p02(r io.Reader) (string, error) {
	return jumpSteps(r, decreaseLargeJump)
}
