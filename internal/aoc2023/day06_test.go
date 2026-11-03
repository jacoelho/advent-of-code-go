package aoc2023

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day06p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Time:      7  15   30
Distance:  9  40  200`),
			Want: "288",
		},
		{
			Input: aoc.FileInput(t, 2023, 6),
			Want:  "1195150",
		},
	}
	aoc.AOCTest(t, day06p01, tests)
}

func Test_day06p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Time:      7  15   30
Distance:  9  40  200`),
			Want: "71503",
		},
		{
			Input: aoc.FileInput(t, 2023, 6),
			Want:  "42550411",
		},
	}
	aoc.AOCTest(t, day06p02, tests)
}

func TestRaceIntegerBoundaries(t *testing.T) {
	tests := []struct{ time, distance, want string }{
		{"7", "12", "0"},
		{"7", "11", "2"},
		{"6", "8", "1"},
		{"0", "0", "0"},
		{"18014398509481984", "0", "18014398509481983"},
	}
	for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
		for _, tt := range tests {
			input := fmt.Sprintf("Time: %s\nDistance: %s", tt.time, tt.distance)
			aoc.AOCTest(t, solve, []aoc.TestInput{{Input: strings.NewReader(input), Want: tt.want}})
		}
	}
	aoc.AOCTest(t, day06p01, []aoc.TestInput{{Input: strings.NewReader("Time: 9223372036854775807 3 1\nDistance: 0 0 0"), Want: "0"}})
	if _, err := day06p01(strings.NewReader("Time: 9223372036854775807 3\nDistance: 0 0")); err == nil {
		t.Fatal("expected product overflow")
	}
}

func TestRaceInputErrors(t *testing.T) {
	for _, input := range []string{
		"Time: 7 15\nDistance: 9",
		"Time: -7\nDistance: 9",
		"Time: seven\nDistance: 9",
	} {
		for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected race error for %q", input)
			}
		}
	}
	if _, err := day06p02(strings.NewReader("Time: 9223372036854775807 0\nDistance: 0 0")); err == nil {
		t.Fatal("expected combined race overflow")
	}
	failure := errors.New("race read failed")
	for _, solve := range []func(io.Reader) (string, error){day06p01, day06p02} {
		_, err := solve(io.MultiReader(strings.NewReader("Time: 7\nDistance: 9"), iotest.ErrReader(failure)))
		if !errors.Is(err, failure) {
			t.Fatalf("got %v, want reader failure", err)
		}
	}
}
