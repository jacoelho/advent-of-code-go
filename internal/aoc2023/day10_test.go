package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.....
.S-7.
.|.|.
.L-J.
.....`),
			Want: "4",
		},
		{
			Input: strings.NewReader(`..F7.
.FJ|.
SJ.L7
|F--J
LJ...`),
			Want: "8",
		},
		{
			Input: aoc.FileInput(t, 2023, 10),
			Want:  "6815",
		},
	}
	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`...........
.S-------7.
.|F-----7|.
.||.....||.
.||.....||.
.|L-7.F-J|.
.|..|.|..|.
.L--J.L--J.
...........`),
			Want: "4",
		},
		{
			Input: strings.NewReader(`..........
.S------7.
.|F----7|.
.||....||.
.||....||.
.|L-7F-J|.
.|..||..|.
.L--JL--J.
..........`),
			Want: "4",
		},
		{
			Input: strings.NewReader(`.F----7F7F7F7F-7....
.|F--7||||||||FJ....
.||.FJ||||||||L7....
FJL7L7LJLJ||LJ.L-7..
L--J.L7...LJS7F-7L7.
....F-J..F7FJ|L7L7L7
....L7.F7||L7|.L7L7|
.....|FJLJ|FJ|F7|.LJ
....FJL-7.||.||||...
....L---J.LJ.LJLJ...`),
			Want: "8",
		},
		{
			Input: strings.NewReader(`FF7FSF7F7F7F7F7F---7
L|LJ||||||||||||F--J
FL-7LJLJ||||||LJL-77
F--JF--7||LJLJ7F7FJ-
L---JF-JLJ.||-FJLJJ7
|F|F-JF---7F7-L7L|7|
|FFJF7L7F-JF7|JL---7
7-L-JL7||F7|L7F-7F7|
L.L7LFJ|||||FJL7||LJ
L7JLJL-JLJLJL--JLJ.L`),
			Want: "10",
		},
		{
			Input: aoc.FileInput(t, 2023, 10),
			Want:  "269",
		},
	}
	aoc.AOCTest(t, day10p02, tests)
}
