package aoc2017

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
		{Input: strings.NewReader("{}"), Want: "1"},
		{Input: strings.NewReader("{{{}}}"), Want: "6"},
		{Input: strings.NewReader("{{},{}}"), Want: "5"},
		{Input: strings.NewReader("{{{},{},{{}}}}"), Want: "16"},
		{Input: strings.NewReader("{<a>,<a>,<a>,<a>}"), Want: "1"},
		{Input: strings.NewReader("{{<ab>},{<ab>},{<ab>},{<ab>}}"), Want: "9"},
		{Input: strings.NewReader("{{<!!>},{<!!>},{<!!>},{<!!>}}"), Want: "9"},
		{Input: strings.NewReader("{{<a!>},{<a!>},{<a!>},{<ab>}}"), Want: "3"},
		{Input: aoc.FileInput(t, 2017, 9), Want: "7616"},
	}

	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("<>"), Want: "0"},
		{Input: strings.NewReader("<random characters>"), Want: "17"},
		{Input: strings.NewReader("<<<<>"), Want: "3"},
		{Input: strings.NewReader("<{!>}>"), Want: "2"},
		{Input: strings.NewReader("<!!>"), Want: "0"},
		{Input: strings.NewReader("<!!!>>"), Want: "0"},
		{Input: strings.NewReader("<{o\"i!a,<{i<a>"), Want: "10"},
		{Input: aoc.FileInput(t, 2017, 9), Want: "3838"},
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

func TestDay09MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"unfinished garbage", "{<abc!", day09p01},
		{"group close", "}", day09p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
