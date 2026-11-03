package aoc2016

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	aoc.AOCTest(t, day07p01, []aoc.TestInput{
		{Input: strings.NewReader("abba[mnop]qrst\nabcd[bddb]xyyx\naaaa[qwer]tyui\nioxxoj[asdfgh]zxcvbn"), Want: "2"},
		{Input: aoc.FileInput(t, 2016, 7), Want: "115"},
	})
}

func Test_day07p02(t *testing.T) {
	aoc.AOCTest(t, day07p02, []aoc.TestInput{
		{Input: strings.NewReader("aba[bab]xyz\nxyx[xyx]xyx\naaa[kek]eke\nzazbz[bzb]cdb"), Want: "3"},
		{Input: aoc.FileInput(t, 2016, 7), Want: "231"},
	})
}

func TestDay07ReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"day07p01", day07p01},
		{"day07p02", day07p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got error %v, want input failure", err)
			}
		})
	}
}

func TestDay07MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"IPv7 brackets", "aba[bab", day07p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
