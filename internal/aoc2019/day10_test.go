package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.#..#
.....
#####
....#
...##`),
			Want: "8",
		},
		{
			Input: strings.NewReader(`......#.#.
#..#.#....
..#######.
.#.#.###..
.#..#.....
..#....#.#
#..#....#.
.##.#..###
##...#..#.
.#....####`),
			Want: "33",
		},
		{
			Input: strings.NewReader(`#.#...#.#.
.###....#.
.#....#...
##.#.#.#.#
....#.#.#.
.##..###.#
..#...##..
..##....##
......#...
.####.###.`),
			Want: "35",
		},
		{
			Input: strings.NewReader(`.#..#..###
####.###.#
....###.#.
..###.##.#
##.##.#.#.
....###..#
..#.#..#.#
#..#.#.###
.##...##.#
.....#.#..`),
			Want: "41",
		},
		{
			Input: strings.NewReader(`.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`),
			Want: "210",
		},
		{
			Input: aoc.FileInput(t, 2019, 10),
			Want:  "267",
		},
	}
	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`.#..##.###...#######
##.############..##.
.#.######.########.#
.###.#######.####.#.
#####.##.#.##.###.##
..#####..#.#########
####################
#.####....###.#.#.##
##.#################
#####.##.###..####..
..######..##.#######
####.##.####...##..#
.#####..#.######.###
##...#.##########...
#.##########.#######
.####.#.###.###.#.##
....##.##.###..#####
.#.#.###########.###
#.#.#.#####.####.###
###.##.####.##.#..##`),
			Want: "802",
		},
		{
			Input: aoc.FileInput(t, 2019, 10),
			Want:  "1309",
		},
	}
	aoc.AOCTest(t, day10p02, tests)
}

func TestAsteroidMapErrors(t *testing.T) {
	for _, input := range []string{"", "...\n...", "#.\n#", "#x"} {
		for _, solve := range []func(io.Reader) (string, error){day10p01, day10p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected asteroid map error for %q", input)
			}
		}
	}
	if _, err := day10p02(strings.NewReader("###\n###")); err == nil {
		t.Fatal("expected error with fewer than 200 asteroids to vaporize")
	}
}
