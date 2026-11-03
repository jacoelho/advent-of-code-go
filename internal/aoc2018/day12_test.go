package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const plantExample = `initial state: #..#.#..##......###...###

...## => #
..#.. => #
.#... => #
.#.#. => #
.#.## => #
.##.. => #
.#### => #
#.#.# => #
#.### => #
##.#. => #
##.## => #
###.. => #
###.# => #
####. => #`

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(plantExample),
			Want:  "325",
		},
		{
			Input: strings.NewReader("initial state: #\n"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("initial state: #\n\n.#... => #"),
			Want:  "20",
		},
		{
			Input: aoc.FileInput(t, 2018, 12),
			Want:  "2911",
		},
	}

	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("initial state: #\n"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("initial state: #\n\n.#... => #"),
			Want:  "50000000000",
		},
		{
			Input: aoc.FileInput(t, 2018, 12),
			Want:  "2500000000695",
		},
	}

	aoc.AOCTest(t, day12p02, tests)
}

func TestPlantErrors(t *testing.T) {
	for _, input := range []string{"", "initial state: x", "initial state: #\n..... => #", "initial state: #\n## => ."} {
		if _, err := day12p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay12ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"12/1", day12p01},
		{"12/2", day12p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
