package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`),
			Want: "136",
		},
		{
			Input: aoc.FileInput(t, 2023, 14),
			Want:  "105249",
		},
	}
	aoc.AOCTest(t, day14p01, tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`O....#....
O.OO#....#
.....##...
OO.#O....O
.O.....O#.
O.#..O.#.#
..O..#O..O
.......O..
#....###..
#OO..#....`),
			Want: "64",
		},
		{
			Input: aoc.FileInput(t, 2023, 14),
			Want:  "88680",
		},
	}
	aoc.AOCTest(t, day14p02, tests)
}
