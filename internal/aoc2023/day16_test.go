package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`),
			Want: "46",
		},
		{
			Input: aoc.FileInput(t, 2023, 16),
			Want:  "6605",
		},
	}
	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....`),
			Want: "51",
		},
		{
			Input: aoc.FileInput(t, 2023, 16),
			Want:  "6766",
		},
	}
	aoc.AOCTest(t, day16p02, tests)
}
