package aoc2018

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
		{Input: strings.NewReader("dabAcCaCBAcCcaDA"), Want: "10"},
		{Input: strings.NewReader("abBA"), Want: "0"},
		{Input: strings.NewReader("abAB"), Want: "4"},
		{Input: strings.NewReader(""), Want: "0"},
		{Input: aoc.FileInput(t, 2018, 5), Want: "11636"},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("dabAcCaCBAcCcaDA"), Want: "4"},
		{Input: strings.NewReader(""), Want: "0"},
		{Input: aoc.FileInput(t, 2018, 5), Want: "5302"},
	}

	aoc.AOCTest(t, day05p02, tests)
}

func TestDay05ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"05/1", day05p01},
		{"05/2", day05p02},
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

func TestDay05InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "invalid polymer", solve: day05p01, input: "a1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
