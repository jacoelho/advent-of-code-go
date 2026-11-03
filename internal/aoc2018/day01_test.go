package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("+1\n-2\n+3\n+1"), Want: "3"},
		{Input: aoc.FileInput(t, 2018, 1), Want: "513"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("+1\n-1"), Want: "0"},
		{Input: strings.NewReader("+3\n+3\n+4\n-2\n-4"), Want: "10"},
		{Input: strings.NewReader("-6\n+3\n+8\n+5\n-6"), Want: "5"},
		{Input: strings.NewReader("+7\n+7\n-2\n-7\n-4"), Want: "14"},
		{Input: strings.NewReader("-1\n+2\n-3"), Want: "-1"},
		{Input: strings.NewReader("0"), Want: "0"},
		{Input: strings.NewReader("9223372036854775807\n-9223372036854775808\n-1\n-1"), Want: "-2"},
		{Input: aoc.FileInput(t, 2018, 1), Want: "287"},
	}

	aoc.AOCTest(t, day01p02, tests)
}

func TestDay01ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"01/1", day01p01},
		{"01/2", day01p02},
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

func TestDay01InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "frequency never repeats", solve: day01p02, input: "+1"},
		{name: "invalid frequency", solve: day01p01, input: "+"},
		{name: "frequency overflow part one", solve: day01p01, input: "9223372036854775807\n1"},
		{name: "frequency overflow part two", solve: day01p02, input: "9223372036854775807\n1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
