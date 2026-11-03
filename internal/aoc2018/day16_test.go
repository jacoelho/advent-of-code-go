package aoc2018

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const opcodeExample = `Before: [3, 2, 1, 1]
9 2 1 2
After:  [3, 2, 2, 1]`

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(opcodeExample), Want: "1"},
		{Input: strings.NewReader(`Before: [9223372036854775807, 1, 0, 0]
0 0 1 2
After:  [9223372036854775807, 1, 9223372036854775807, 0]`), Want: "1"},
		{Input: strings.NewReader(`Before: [0, 0, 0, 0]
0 0 5 0
After:  [0, 0, 0, 0]`), Want: "1"},
		{Input: aoc.FileInput(t, 2018, 16), Want: "624"},
	}

	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2018, 16), Want: "584"},
	}

	aoc.AOCTest(t, day16p02, tests)
}

func TestOpcodeMappingErrors(t *testing.T) {
	tests := map[string]string{
		"ambiguous": opcodeExample,
		"inconsistent": `Before: [0, 0, 0, 0]
0 0 0 0
After: [0, 1, 0, 0]`,
	}

	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := day16p02(strings.NewReader(input)); err == nil {
				t.Fatal("expected opcode mapping error")
			}
		})
	}
}

func TestDay16ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"16/1", day16p01},
		{"16/2", day16p02},
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
