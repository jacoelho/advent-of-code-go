package aoc2017

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestKnotProductExample(t *testing.T) {
	got, err := knotProduct(strings.NewReader("3,4,1,5"), 5)
	if err != nil {
		t.Fatal(err)
	}
	if got != "12" {
		t.Fatalf("got %s, want 12", got)
	}
}

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("3,4,1,5"), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 10), Want: "52070"},
	}

	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(""), Want: "a2582a3a0e66e6e86e3812dcb672a272"},
		{Input: strings.NewReader("AoC 2017"), Want: "33efeb34ea91902bb2f59c9920caa6cd"},
		{Input: strings.NewReader("1,2,3"), Want: "3efbe78a8d82f29979031a4aa0b16a9d"},
		{Input: strings.NewReader("1,2,4"), Want: "63960835bcdc130f0b66d7ff4f6a5a8e"},
		{Input: aoc.FileInput(t, 2017, 10), Want: "7f94112db4e32e19cf6502073c66f9bb"},
	}

	aoc.AOCTest(t, day10p02, tests)
}
