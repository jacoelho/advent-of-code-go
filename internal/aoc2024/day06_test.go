package aoc2024

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`),
			Want: "41",
		},
		{
			Input: aoc.FileInput(t, 2024, 6),
			Want:  "5208",
		},
	}

	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#.....
.........#
..........
..#.......
.......#..
..........
.#..^.....
........#.
#.........
......#...`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2024, 6),
			Want:  "1972",
		},
	}

	aoc.AOCTest(t, day06p02, tests)
}

func TestGuardMapErrors(t *testing.T) {
	for _, input := range []string{".", "^^"} {
		for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: got %q; want guard error", input, got)
			}
		}
	}
}

func TestEnclosedGuard(t *testing.T) {
	const input = "###\n#^#\n###"
	aoc.AOCTest(t, day06p01, []aoc.TestInput{{Input: strings.NewReader(input), Want: "1"}})
	aoc.AOCTest(t, day06p02, []aoc.TestInput{{Input: strings.NewReader(input), Want: "0"}})
}
