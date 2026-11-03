package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestLookAndSayPhases(t *testing.T) {
	sequence := []byte("1")
	for _, want := range []string{"11", "21", "1211", "111221", "312211"} {
		sequence = lookAndSay(sequence)
		if string(sequence) != want {
			t.Fatalf("got %s, want %s", sequence, want)
		}
	}
}

func Test_day10p01(t *testing.T) {
	aoc.AOCTest(t, day10p01, []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "82350"},
		{Input: aoc.FileInput(t, 2015, 10), Want: "492982"},
	})
}

func Test_day10p02(t *testing.T) {
	aoc.AOCTest(t, day10p02, []aoc.TestInput{
		{Input: strings.NewReader("1"), Want: "1166642"},
		{Input: aoc.FileInput(t, 2015, 10), Want: "6989950"},
	})
}

func TestDay10ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day10p01", day10p01},
		{"day10p02", day10p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay10MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"look-and-say digit", "12a", day10p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
