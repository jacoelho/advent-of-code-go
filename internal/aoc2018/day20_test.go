package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("^WNE$"),
			Want:  "3",
		},
		{
			Input: strings.NewReader("^ENWWW(NEEE|SSE(EE|N))$"),
			Want:  "10",
		},
		{
			Input: strings.NewReader("^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$"),
			Want:  "18",
		},
		{
			Input: strings.NewReader("^N(E|W)N$"),
			Want:  "3",
		},
		{
			Input: strings.NewReader("^$"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 20),
			Want:  "4308",
		},
	}
	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("^WNE$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^ENWWW(NEEE|SSE(EE|N))$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^N(E|W)N$"),
			Want:  "0",
		},
		{
			Input: strings.NewReader("^$"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 20),
			Want:  "8528",
		},
	}
	aoc.AOCTest(t, day20p02, tests)
}

func TestDay20ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"20/1", day20p01},
		{"20/2", day20p02},
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
