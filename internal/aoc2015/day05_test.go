package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	aoc.AOCTest(t, day05p01, []aoc.TestInput{
		{Input: strings.NewReader("ugknbfddgicrmopn\naaa\njchzalrnumimnmhp\nhaegwjzuvuyypxyu\ndvszwmarrgswjxmb"), Want: "2"},
		{Input: aoc.FileInput(t, 2015, 5), Want: "258"},
	})
}

func Test_day05p02(t *testing.T) {
	aoc.AOCTest(t, day05p02, []aoc.TestInput{
		{Input: strings.NewReader("qjhvhtzxzqqjkmpb\nxxyxx\nuurcxstgmygtbstg\nieodomkazucvgmuy"), Want: "2"},
		{Input: strings.NewReader("aaa"), Want: "0"},
		{Input: strings.NewReader("aaaa"), Want: "1"},
		{Input: aoc.FileInput(t, 2015, 5), Want: "53"},
	})
}

func TestDay05ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day05p01", day05p01},
		{"day05p02", day05p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}
