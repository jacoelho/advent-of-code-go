package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`),
			Want: "36",
		},
		{
			Input: aoc.FileInput(t, 2024, 10),
			Want:  "482",
		},
	}

	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`89010123
78121874
87430965
96549874
45678903
32019012
01329801
10456732`),
			Want: "81",
		},
		{
			Input: aoc.FileInput(t, 2024, 10),
			Want:  "1094",
		},
	}

	aoc.AOCTest(t, day10p02, tests)
}

func TestTrailRows(t *testing.T) {
	for _, test := range []struct {
		name, input, want string
	}{
		{
			name: "empty input",
			want: "0",
		},
		{
			name:  "row boundaries",
			input: "01234\n56789",
			want:  "0",
		},
		{
			name:  "ragged rows and nondigits",
			input: "0123456789\n.x",
			want:  "1",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			for i, solve := range []func(io.Reader) (string, error){day10p01, day10p02} {
				got, err := solve(strings.NewReader(test.input))
				if err != nil || got != test.want {
					t.Fatalf("part %d = %q, %v; want %s", i+1, got, err, test.want)
				}
			}
		})
	}
}

func TestDay10ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day10p01, "0123456789\n"},
		{"part2", day10p02, "0123456789\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
