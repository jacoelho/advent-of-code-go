package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`5764801
17807724`),
			Want: "14897079",
		},
		{
			Input: aoc.FileInput(t, 2020, 25),
			Want:  "4441893",
		},
	}

	aoc.AOCTest(t, day25p01, tests)
}
