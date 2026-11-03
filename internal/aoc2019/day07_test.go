package aoc2019

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("3,0,3,1,104,-1,99"), Want: "-1"},
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "437860",
		},
	}
	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("3,0,3,1,104,-1,99"), Want: "-1"},
		{
			Input: aoc.FileInput(t, 2019, 7),
			Want:  "49810599",
		},
	}
	aoc.AOCTest(t, day07p02, tests)
}
