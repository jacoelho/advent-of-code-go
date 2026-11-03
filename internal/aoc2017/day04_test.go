package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`aa bb cc dd ee
aa bb cc dd aa
aa bb cc dd aaa`), Want: "2"},
		{Input: aoc.FileInput(t, 2017, 4), Want: "337"},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abcde fghij
abcde xyz ecdab
a ab abc abd abf abj
iiii oiii ooii oooi oooo
oiii ioii iioi iiio`), Want: "3"},
		{Input: strings.NewReader("ab abb"), Want: "1"},
		{Input: aoc.FileInput(t, 2017, 4), Want: "231"},
	}

	aoc.AOCTest(t, day04p02, tests)
}

func TestDay04ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"04/1", day04p01},
		{"04/2", day04p02},
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
