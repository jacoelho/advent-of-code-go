package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const blueprintExample = `Begin in state A.
Perform a diagnostic checksum after 6 steps.

In state A:
  If the current value is 0:
    - Write the value 1.
    - Move one slot to the right.
    - Continue with state B.
  If the current value is 1:
    - Write the value 0.
    - Move one slot to the left.
    - Continue with state B.

In state B:
  If the current value is 0:
    - Write the value 1.
    - Move one slot to the left.
    - Continue with state A.
  If the current value is 1:
    - Write the value 1.
    - Move one slot to the right.
    - Continue with state A.`

func Test_day25p01(t *testing.T) {
	aoc.AOCTest(t, day25p01, []aoc.TestInput{
		{Input: strings.NewReader(blueprintExample), Want: "3"},
		{Input: strings.NewReader(strings.Replace(blueprintExample, "6 steps", "0 steps", 1)), Want: "0"},
		{Input: aoc.FileInput(t, 2017, 25), Want: "3578"},
	})
}

func TestBlueprintRejectsInvalidTransitions(t *testing.T) {
	tests := []struct{ from, to string }{
		{"state B.", "state C."},
		{"value 1.", "value 2."},
		{"to the left.", "to the middle."},
		{"current value is 1:", "current value is 0:"},
	}
	for _, tt := range tests {
		t.Run(tt.to, func(t *testing.T) {
			input := strings.Replace(blueprintExample, tt.from, tt.to, 1)
			if _, err := parseBlueprint(strings.NewReader(input)); err == nil {
				t.Fatal("expected invalid transition error")
			}
		})
	}
}

func TestDay25ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"25/1", day25p01},
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
