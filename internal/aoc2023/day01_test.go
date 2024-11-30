package aoc2023

import (
	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"strings"
	"testing"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1abc2
pqr3stu8vwx
a1b2c3d4e5f
treb7uchet`),
			Want: "142",
		},
		{
			Input: aoc.FileInput(t, 2023, 1),
			Want:  "54708",
		},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`two1nine
eightwothree
abcone2threexyz
xtwone3four
4nineeightseven2
zoneight234
7pqrstsixteen`),
			Want: "281",
		},
		{
			Input: aoc.FileInput(t, 2023, 1),
			Want:  "54087",
		},
	}
	aoc.AOCTest(t, day01p02, tests)
}
