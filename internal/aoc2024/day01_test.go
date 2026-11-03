package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3   4
4   3
2   5
1   3
3   9
3   3`),
			Want: "11",
		},
		{
			Input: aoc.FileInput(t, 2024, 1),
			Want:  "1110981",
		},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`3   4
4   3
2   5
1   3
3   9
3   3`),
			Want: "31",
		},
		{
			Input: aoc.FileInput(t, 2024, 1),
			Want:  "24869388",
		},
	}

	aoc.AOCTest(t, day01p02, tests)
}

func TestLocationRowErrors(t *testing.T) {
	for _, input := range []string{"1", "1 2 3"} {
		t.Run(input, func(t *testing.T) {
			for _, solve := range []func(io.Reader) (string, error){day01p01, day01p02} {
				if got, err := solve(strings.NewReader(input)); err == nil {
					t.Fatalf("got %q; want malformed-row error", got)
				}
			}
		})
	}
}

func TestDay01ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day01p01, "1 2\n"},
		{"part2", day01p02, "1 2\n"},
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
