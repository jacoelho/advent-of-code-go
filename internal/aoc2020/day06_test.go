package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abc

a
b
c

ab
ac

a
a
a
a

b`), Want: "11"},
		{Input: strings.NewReader(`ab
ac

a
b`), Want: "5"},
		{Input: strings.NewReader(`

ab


a

`), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 6), Want: "6504"},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abc

a
b
c

ab
ac

a
a
a
a

b`), Want: "6"},
		{Input: strings.NewReader(`ab
ac

a
b`), Want: "1"},
		{Input: strings.NewReader(`

ab


a

`), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 6), Want: "3351"},
	}
	aoc.AOCTest(t, day06p02, tests)
}
