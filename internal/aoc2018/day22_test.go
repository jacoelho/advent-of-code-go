package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestCaveRequiresTorchAtTarget(t *testing.T) {
	if _, err := day22p02(strings.NewReader("depth: 1\ntarget: 1,1")); err == nil {
		t.Fatal("expected invalid torch region error")
	}
}

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("depth: 510\ntarget: 10,10"),
			Want:  "114",
		},
		{
			Input: strings.NewReader("depth: 510\ntarget: 0,0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 22),
			Want:  "5622",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("depth: 510\ntarget: 10,10"),
			Want:  "45",
		},
		{
			Input: strings.NewReader("depth: 510\ntarget: 0,0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2018, 22),
			Want:  "1089",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestDay22ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"22/1", day22p01},
		{"22/2", day22p02},
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
