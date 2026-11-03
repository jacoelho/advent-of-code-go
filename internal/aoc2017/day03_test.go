package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "0"},
		{Input: strings.NewReader("12"), Want: "3"},
		{Input: strings.NewReader("23"), Want: "2"},
		{Input: strings.NewReader("1024"), Want: "31"},
		{Input: aoc.FileInput(t, 2017, 3), Want: "438"},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "2"},
		{Input: strings.NewReader("747"), Want: "806"},
		{Input: aoc.FileInput(t, 2017, 3), Want: "266330"},
	}

	aoc.AOCTest(t, day03p02, tests)
}
