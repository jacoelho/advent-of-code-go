package aoc2023

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9`),
			Want: "5",
		},
		{
			Input: aoc.FileInput(t, 2023, 22),
			Want:  "527",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1,0,1~1,2,1
0,0,2~2,0,2
0,2,3~2,2,3
0,0,4~0,2,4
2,0,5~2,2,5
0,1,6~2,1,6
1,1,8~1,1,9`),
			Want: "7",
		},
		{
			Input: aoc.FileInput(t, 2023, 22),
			Want:  "100376",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestDay22ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day22p01, "0,0,1~0,0,1\n"},
		{"part2", day22p02, "0,0,1~0,0,1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, mode := range []struct{ name, prefix string }{{"immediate", ""}, {"after prefix", tt.prefix}} {
				t.Run(mode.name, func(t *testing.T) {
					input := io.MultiReader(strings.NewReader(mode.prefix), iotest.ErrReader(failure))
					_, err := tt.solve(input)
					if !errors.Is(err, failure) {
						t.Fatalf("got %v, want reader failure", err)
					}
				})
			}
		})
	}
}
