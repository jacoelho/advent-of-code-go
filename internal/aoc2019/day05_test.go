package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 5),
			Want:  "8332629",
		},
	}
	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 5),
			Want:  "8805067",
		},
	}
	aoc.AOCTest(t, day05p02, tests)
}

func TestDiagnosticRejectsFailedCheck(t *testing.T) {
	if _, err := day05p01(strings.NewReader("104,3,104,7,99")); err == nil {
		t.Fatal("expected failed diagnostic check error")
	}
}
