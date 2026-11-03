package aoc2024

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`),
			Want: "7036",
		},
		{
			Input: strings.NewReader(`#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`),
			Want: "11048",
		},
		{
			Input: aoc.FileInput(t, 2024, 16),
			Want:  "143564",
		},
	}
	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`###############
#.......#....E#
#.#.###.#.###.#
#.....#.#...#.#
#.###.#####.#.#
#.#.#.......#.#
#.#.#####.###.#
#...........#.#
###.#.#####.#.#
#...#.....#.#.#
#.#.#.###.#.#.#
#.....#...#.#.#
#.###.#.#.#.#.#
#S..#.....#...#
###############`),
			Want: "45",
		},
		{
			Input: strings.NewReader(`#################
#...#...#...#..E#
#.#.#.#.#.#.#.#.#
#.#.#.#...#...#.#
#.#.#.#.###.#.#.#
#...#.#.#.....#.#
#.#.#.#.#.#####.#
#.#...#.#.#.....#
#.#.#####.#.###.#
#.#.#.......#...#
#.#.###.#####.###
#.#.#...#.....#.#
#.#.#.#####.###.#
#.#.#.........#.#
#.#.#.#########.#
#S#.............#
#################`),
			Want: "64",
		},
		{
			Input: aoc.FileInput(t, 2024, 16),
			Want:  "593",
		},
	}
	aoc.AOCTest(t, day16p02, tests)
}

func TestReindeerRoutes(t *testing.T) {
	for part, solve := range []func(io.Reader) (string, error){day16p01, day16p02} {
		for _, input := range []string{"S#E", ".#E"} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("part %d input %q: expected route error", part+1, input)
			}
		}
		want := []string{"2001", "2"}[part]
		got, err := solve(strings.NewReader("ES"))
		if err != nil || got != want {
			t.Fatalf("part %d reversed route = %q, %v; want %s", part+1, got, err, want)
		}
	}
}
