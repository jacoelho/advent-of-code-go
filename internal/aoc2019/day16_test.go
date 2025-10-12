package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("80871224585914546619083218645595"),
			Want:  "24176176",
		},
		{
			Input: strings.NewReader("19617804207202209144916044189917"),
			Want:  "73745418",
		},
		{
			Input: strings.NewReader("69317163492948606335995924319873"),
			Want:  "52432133",
		},
		{
			Input: aoc.FileInput(t, 2019, 16),
			Want:  "40921727",
		},
	}
	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("03036732577212944063491565474664"),
			Want:  "84462026",
		},
		{
			Input: strings.NewReader("02935109699940807407585447034323"),
			Want:  "78725270",
		},
		{
			Input: strings.NewReader("03081770884921959731165446850517"),
			Want:  "53553731",
		},
		{
			Input: aoc.FileInput(t, 2019, 16),
			Want:  "89950138",
		},
	}
	aoc.AOCTest(t, day16p02, tests)
}
