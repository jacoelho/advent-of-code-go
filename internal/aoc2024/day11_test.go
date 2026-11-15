package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`125 17`),
			Want:  "55312",
		},
		{
			Input: aoc.FileInput(t, 2024, 11),
			Want:  "189092",
		},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`125 17`),
			Want:  "65601038650482",
		},
		{
			Input: aoc.FileInput(t, 2024, 11),
			Want:  "224869647102559",
		},
	}

	aoc.AOCTest(t, day11p02, tests)
}

func TestStoneArrangementRejectsMalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
	}{
		{"non-numeric stone", "125 nope"},
		{"negative stone", "125 -1"},
		{"empty input", " \n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for i, solve := range []func(io.Reader) (string, error){day11p01, day11p02} {
				if got, err := solve(strings.NewReader(tt.input)); err == nil {
					t.Fatalf("part %d = %q; want error", i+1, got)
				}
			}
		})
	}
}

func TestStoneMultiplicationOverflowIsRejected(t *testing.T) {
	if got, err := day11p01(strings.NewReader("9223372036854775807")); err == nil {
		t.Fatalf("got %q; want overflow error", got)
	}
}

func TestDay11ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day11p01, "125 17\n"},
		{"part2", day11p02, "125 17\n"},
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
