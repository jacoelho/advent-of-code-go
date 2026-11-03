package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`), Want: "2"},
		{Input: strings.NewReader(`1-3 a: aba
1-3 a: aab`), Want: "2"},
		{Input: aoc.FileInput(t, 2020, 2), Want: "493"},
	}
	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`), Want: "1"},
		{Input: strings.NewReader(`1-3 a: aba
1-3 a: aab`), Want: "1"},
		{Input: aoc.FileInput(t, 2020, 2), Want: "593"},
	}
	aoc.AOCTest(t, day02p02, tests)
}
