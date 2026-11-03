package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("+1\n-2\n+3\n+1"), Want: "3"},
		{Input: aoc.FileInput(t, 2018, 1), Want: "513"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("+1\n-1"), Want: "0"},
		{Input: strings.NewReader("+3\n+3\n+4\n-2\n-4"), Want: "10"},
		{Input: strings.NewReader("-6\n+3\n+8\n+5\n-6"), Want: "5"},
		{Input: strings.NewReader("+7\n+7\n-2\n-7\n-4"), Want: "14"},
		{Input: strings.NewReader("-1\n+2\n-3"), Want: "-1"},
		{Input: strings.NewReader("0"), Want: "0"},
		{Input: strings.NewReader("9223372036854775807\n-9223372036854775808\n-1\n-1"), Want: "-2"},
		{Input: aoc.FileInput(t, 2018, 1), Want: "287"},
	}

	aoc.AOCTest(t, day01p02, tests)
}
