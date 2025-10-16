package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`broadcaster -> a, b, c
%a -> b
%b -> c
%c -> inv
&inv -> a`),
			Want: "32000000",
		},
		{
			Input: strings.NewReader(`broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> output`),
			Want: "11687500",
		},
		{
			Input: aoc.FileInput(t, 2023, 20),
			Want:  "787056720",
		},
	}
	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`broadcaster -> a
%a -> inv, con
&inv -> b
%b -> con
&con -> rx`),
			Want: "1",
		},
		{
			Input: aoc.FileInput(t, 2023, 20),
			Want:  "212986464842911",
		},
	}
	aoc.AOCTest(t, day20p02, tests)
}
