package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("9 players; last marble is worth 25 points"), Want: "32"},
		{Input: strings.NewReader("9 players; last marble is worth 23 points"), Want: "32"},
		{Input: strings.NewReader("10 players; last marble is worth 1618 points"), Want: "8317"},
		{Input: aoc.FileInput(t, 2018, 9), Want: "386018"},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("9 players; last marble is worth 25 points"), Want: "22563"},
		{Input: aoc.FileInput(t, 2018, 9), Want: "3085518618"},
	}

	aoc.AOCTest(t, day09p02, tests)
}

func TestDay09ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"09/1", day09p01},
		{"09/2", day09p02},
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

func TestDay09InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "no marble players", solve: day09p01, input: "0 players; last marble is worth 25 points"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
