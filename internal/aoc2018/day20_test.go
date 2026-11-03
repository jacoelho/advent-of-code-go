package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("^WNE$"),
			Want:  "3",
		},
		{
			Input: strings.NewReader("^ENWWW(NEEE|SSE(EE|N))$"),
			Want:  "10",
		},
		{
			Input: strings.NewReader("^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$"),
			Want:  "18",
		},
		{
			Input: strings.NewReader("^N(E|W)N$"),
			Want:  "3",
		},
		{
			Input: strings.NewReader("^$"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 20),
			Want:  "4308",
		},
	}
	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("^WNE$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^ENWWW(NEEE|SSE(EE|N))$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^N(E|W)N$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^$"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 20),
			Want:  "8528",
		},
	}
	aoc.AOCTest(t, day20p02, tests)
}
