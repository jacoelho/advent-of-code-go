package aoc2020

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`), Want: "37"},
		{Input: strings.NewReader(`#...L...#`), Want: "3"},
		{Input: strings.NewReader(`...`), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 11), Want: "2453"},
	}
	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`L.LL.LL.LL
LLLLLLL.LL
L.L.L..L..
LLLL.LL.LL
L.LL.LL.LL
L.LLLLL.LL
..L.L.....
LLLLLLLLLL
L.LLLLLL.L
L.LLLLL.LL`), Want: "26"},
		{Input: strings.NewReader(`#...L...#`), Want: "2"},
		{Input: strings.NewReader(`...`), Want: "0"},
		{Input: aoc.FileInput(t, 2020, 11), Want: "2159"},
	}
	aoc.AOCTest(t, day11p02, tests)
}
