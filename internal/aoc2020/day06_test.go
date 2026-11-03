package aoc2020

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
		{Input: strings.NewReader(`abc

a
b
c

ab
ac

a
a
a
a

b`), Want: "11"},
		{Input: strings.NewReader(`ab
ac

a
b`), Want: "5"},
		{Input: strings.NewReader(`

ab


a

`), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 6), Want: "6504"},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abc

a
b
c

ab
ac

a
a
a
a

b`), Want: "6"},
		{Input: strings.NewReader(`ab
ac

a
b`), Want: "1"},
		{Input: strings.NewReader(`

ab


a

`), Want: "3"},
		{Input: aoc.FileInput(t, 2020, 6), Want: "3351"},
	}
	aoc.AOCTest(t, day06p02, tests)
}

func TestDay06ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"06/1", day06p01},
		{"06/2", day06p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay06MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"customs form", "aB", day06p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
