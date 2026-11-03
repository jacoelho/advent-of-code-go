package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`), Want: "2"},
		{Input: strings.NewReader(`1-3 a: aba
1-3 a: aab`), Want: "2"},
		{Input: aoc.FileInput(t, 2020, 2), Want: "493"},
	}
	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1-3 a: abcde
1-3 b: cdefg
2-9 c: ccccccccc`), Want: "1"},
		{Input: strings.NewReader(`1-3 a: aba
1-3 a: aab`), Want: "1"},
		{Input: aoc.FileInput(t, 2020, 2), Want: "593"},
	}
	aoc.AOCTest(t, day02p02, tests)
}

func TestDay02ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"02/1", day02p01},
		{"02/2", day02p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay02MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"policy", "0-2 a: ab", day02p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
