package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`FBFBBFFRLR
BFFFBBFRRR
FFFBBBFRRR
BBFFBBFRLL`), Want: "820"},
		{Input: strings.NewReader(`FFFFFBFRLL
FFFFFBFRLR
FFFFFBFRRR`), Want: "23"},
		{Input: aoc.FileInput(t, 2020, 5), Want: "974"},
	}
	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`FFFFFBFRLL
FFFFFBFRLR
FFFFFBFRRR`), Want: "22"},
		{Input: aoc.FileInput(t, 2020, 5), Want: "646"},
	}
	aoc.AOCTest(t, day05p02, tests)
}

func TestBoardingPassSeatIDs(t *testing.T) {
	for _, tt := range []struct{ pass, want string }{
		{"FBFBBFFRLR", "357"}, {"BFFFBBFRRR", "567"}, {"FFFBBBFRRR", "119"}, {"BBFFBBFRLL", "820"},
	} {
		t.Run(tt.pass, func(t *testing.T) {
			aoc.AOCTest(t, day05p01, []aoc.TestInput{{Input: strings.NewReader(tt.pass), Want: tt.want}})
		})
	}
	if _, err := day05p02(strings.NewReader("FFFFFBFRLL\nFFFFFBFRLR")); err == nil {
		t.Fatal("expected no-gap error")
	}
}

func TestDay05ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"05/1", day05p01},
		{"05/2", day05p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay05MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"boarding pass", "FBFBBFFRLF", day05p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
