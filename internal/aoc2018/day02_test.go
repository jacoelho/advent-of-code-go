package aoc2018

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
		{Input: strings.NewReader(`abcdef
bababc
abbcde
abcccd
aabcdd
abcdee
ababab`), Want: "12"},
		{Input: aoc.FileInput(t, 2018, 2), Want: "6972"},
	}

	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`abcde
fghij
klmno
pqrst
fguij
axcye
wvxyz`), Want: "fgij"},
		{Input: aoc.FileInput(t, 2018, 2), Want: "aixwcbzrmdvpsjfgllthdyoqe"},
	}

	aoc.AOCTest(t, day02p02, tests)
}

func TestDay02ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"02/1", day02p01},
		{"02/2", day02p02},
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

func TestDay02InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "unequal box IDs", solve: day02p02, input: "ab\nacd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
