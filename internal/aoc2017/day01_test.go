package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1122"), Want: "3"},
		{Input: strings.NewReader("1111"), Want: "4"},
		{Input: strings.NewReader("1234"), Want: "0"},
		{Input: strings.NewReader("91212129"), Want: "9"},
		{Input: aoc.FileInput(t, 2017, 1), Want: "1144"},
	}

	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1212"), Want: "6"},
		{Input: strings.NewReader("1221"), Want: "0"},
		{Input: strings.NewReader("123425"), Want: "4"},
		{Input: strings.NewReader("123123"), Want: "12"},
		{Input: strings.NewReader("12131415"), Want: "4"},
		{Input: aoc.FileInput(t, 2017, 1), Want: "1194"},
	}

	aoc.AOCTest(t, day01p02, tests)
}

func TestDay01ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"01/1", day01p01},
		{"01/2", day01p02},
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

func TestDay01MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"captcha digit", "1x", day01p01},
		{"odd captcha", "123", day01p02},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
