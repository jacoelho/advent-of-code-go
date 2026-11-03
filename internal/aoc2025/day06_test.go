package aoc2025

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +  `),
			Want: "4277556",
		},
		{
			Input: aoc.FileInput(t, 2025, 06),
			Want:  "3261038365331",
		},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +`),
			Want: "3263827",
		},
		{
			Input: aoc.FileInput(t, 2025, 06),
			Want:  "8342588849093",
		},
	}
	aoc.AOCTest(t, day06p02, tests)
}

func TestHomeworkSpacing(t *testing.T) {
	for _, tt := range []struct {
		input string
		want  [2]string
	}{
		{"12\n3\n+", [2]string{"15", "15"}},
		{"12  \n3\n+ ", [2]string{"15", "15"}},
		{"12  4\n 3  5\n+   *", [2]string{"35", "69"}},
	} {
		for _, input := range []string{tt.input, tt.input + "\n", strings.ReplaceAll(tt.input+"\n", "\n", "\r\n")} {
			for part, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
				got, err := solve(strings.NewReader(input))
				if err != nil || got != tt.want[part] {
					t.Fatalf("part %d input %q = %q, %v; want %s", part+1, input, got, err, tt.want[part])
				}
			}
		}
	}
}

func TestHomeworkInputErrors(t *testing.T) {
	for _, input := range []string{"", "  \n  ", "12\n3\n ", "12\n3\n-", "1x\n3\n+", "12\n3\n++"} {
		for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
}

func TestDay06ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day06p01, "12\n3\n+\n"},
		{"part2", day06p02, "12\n3\n+\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(input); !errors.Is(err, failure) {
					t.Fatalf("got %v; want reader failure", err)
				}
			}
		})
	}
}
