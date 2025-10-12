package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	t.Skip()
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(``),
			Want:  "",
		},
		{
			Input: aoc.FileInput(t, 2019, 25),
			Want:  "",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func Test_day25p02(t *testing.T) {
	t.Skip()
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(``),
			Want:  "",
		},
		{
			Input: aoc.FileInput(t, 2019, 25),
			Want:  "",
		},
	}
	aoc.AOCTest(t, day25p02, tests)
}
