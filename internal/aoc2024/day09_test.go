package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2333133121414131402`),
			Want:  "1928",
		},
		{
			Input: aoc.FileInput(t, 2024, 9),
			Want:  "6332189866718",
		},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2333133121414131402`),
			Want:  "2858",
		},
		{
			Input: aoc.FileInput(t, 2024, 9),
			Want:  "6353648390778",
		},
	}

	aoc.AOCTest(t, day09p02, tests)
}
