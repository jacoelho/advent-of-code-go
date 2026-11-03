package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("ne,ne,ne"), Want: "3"},
		{Input: strings.NewReader("ne,ne,sw,sw"), Want: "0"},
		{Input: strings.NewReader("ne,ne,s,s"), Want: "2"},
		{Input: strings.NewReader("se,sw,se,sw,sw"), Want: "3"},
		{Input: strings.NewReader("ne,n"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 11), Want: "705"},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("ne,ne,sw,sw"), Want: "2"},
		{Input: strings.NewReader("ne,n"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 11), Want: "1469"},
	}

	aoc.AOCTest(t, day11p02, tests)
}

func TestDay11ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"11/1", day11p01},
		{"11/2", day11p02},
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

func TestDay11MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"hex direction", "ne,east", day11p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
