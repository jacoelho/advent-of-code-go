package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`/->-\
|   |  /----\
| /-+--+-\  |
| | |  | v  |
\-+-/  \-+--/
  \------/`),
			Want: "7,3",
		},
		{
			Input: aoc.FileInput(t, 2018, 13),
			Want:  "57,104",
		},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`/>-<\
|   |
| /<+-\
| | | v
\>+</ |
  |   ^
  \<->/`),
			Want: "6,4",
		},
		{
			Input: strings.NewReader(">><"),
			Want:  "1,0",
		},
		{
			Input: aoc.FileInput(t, 2018, 13),
			Want:  "67,74",
		},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func TestMineCartErrors(t *testing.T) {
	for _, input := range []string{"", "><", "?"} {
		if _, err := day13p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay13ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"13/1", day13p01},
		{"13/2", day13p02},
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
