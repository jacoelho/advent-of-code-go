package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`100756`),
			Want:  "33583",
		},
		{
			Input: aoc.FileInput(t, 2019, 01),
			Want:  "3399394",
		},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`100756`),
			Want:  "50346",
		},
		{
			Input: aoc.FileInput(t, 2019, 01),
			Want:  "5096223",
		},
	}
	aoc.AOCTest(t, day01p02, tests)
}
