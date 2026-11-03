package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("0,0,0,0\n3,0,0,0\n0,3,0,0\n0,0,3,0\n0,0,0,3\n0,0,0,6\n9,0,0,0\n12,0,0,0"),
			Want:  "2",
		},
		{
			Input: strings.NewReader("0,0,0,0\n0,0,0,0\n3,0,0,0\n6,0,0,0\n10,0,0,0"),
			Want:  "2",
		},
		{
			Input: aoc.FileInput(t, 2018, 25),
			Want:  "383",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func TestDay25ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"25/1", day25p01},
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
