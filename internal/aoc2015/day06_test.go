package aoc2015

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	aoc.AOCTest(t, day06p01, []aoc.TestInput{
		{Input: strings.NewReader("turn on 0,0 through 999,999"), Want: "1000000"},
		{Input: strings.NewReader("toggle 0,0 through 999,0"), Want: "1000"},
		{Input: strings.NewReader("turn on 0,0 through 999,999\nturn off 499,499 through 500,500"), Want: "999996"},
		{Input: aoc.FileInput(t, 2015, 6), Want: "377891"},
	})
}

func Test_day06p02(t *testing.T) {
	aoc.AOCTest(t, day06p02, []aoc.TestInput{
		{Input: strings.NewReader("turn on 0,0 through 0,0"), Want: "1"},
		{Input: strings.NewReader("toggle 0,0 through 999,999"), Want: "2000000"},
		{Input: strings.NewReader("turn off 0,0 through 0,0\nturn on 999,999 through 999,999\nturn off 0,0 through 0,0"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 6), Want: "14110788"},
	})
}
