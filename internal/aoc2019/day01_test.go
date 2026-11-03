package aoc2019

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
			Input: strings.NewReader(`100756`),
			Want:  "33583",
		},
		{
			Input: aoc.FileInput(t, 2019, 01),
			Want:  "3399394",
		},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`100756`),
			Want:  "50346",
		},
		{
			Input: aoc.FileInput(t, 2019, 01),
			Want:  "5096223",
		},
	}
	aoc.AOCTest(t, day01p02, tests)
}

func TestDay01ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day01p01", "12\n", day01p01},
		{"day01p02", "12\n", day01p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}
