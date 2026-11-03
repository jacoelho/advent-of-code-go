package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`), Want: "7"},
		{Input: strings.NewReader(`...
...
...`), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 3), Want: "225"},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`..##.......
#...#...#..
.#....#..#.
..#.#...#.#
.#...##..#.
..#.##.....
.#.#.#....#
.#........#
#.##...#...
#...##....#
.#..#...#.#`), Want: "336"},
		{Input: strings.NewReader(`...
...
...`), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 3), Want: "1115775000"},
	}
	aoc.AOCTest(t, day03p02, tests)
}
