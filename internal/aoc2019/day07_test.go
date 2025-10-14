package aoc2019

import (
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day7p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "437860",
		},
	}
	aoc.AOCTest(t, day7p01, tests)
}

func Test_day7p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "49810599",
		},
	}
	aoc.AOCTest(t, day7p02, tests)
}
