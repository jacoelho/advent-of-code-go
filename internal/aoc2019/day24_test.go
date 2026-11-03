package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`....#
#..#.
#..##
..#..
#....`),
			Want: "2129920",
		},
		{
			Input: aoc.FileInput(t, 2019, 24),
			Want:  "18401265",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02_example(t *testing.T) {
	result, err := simulateRecursive(strings.NewReader(`....#
#..#.
#.?##
..#..
#....`), 10)
	if err != nil {
		t.Fatal(err)
	}
	if result != "99" {
		t.Errorf("got %s, want 99", result)
	}
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 24),
			Want:  "2078",
		},
	}
	aoc.AOCTest(t, day24p02, tests)
}

func TestDay24ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day24p01", ".....\n.....\n.....\n.....\n.....\n", day24p01},
		{"day24p02", ".....\n.....\n.....\n.....\n.....\n", day24p02},
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
