package aoc2018

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const stepsExample = `Step C must be finished before step A can begin.
Step C must be finished before step F can begin.
Step A must be finished before step B can begin.
Step A must be finished before step D can begin.
Step B must be finished before step E can begin.
Step D must be finished before step E can begin.
Step F must be finished before step E can begin.`

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(stepsExample), Want: "CABDFE"},
		{Input: aoc.FileInput(t, 2018, 7), Want: "GLMVWXZDKOUCEJRHFAPITSBQNY"},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(stepsExample), Want: "253"},
		{Input: aoc.FileInput(t, 2018, 7), Want: "1105"},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func TestAssemblyTime(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "example", input: stepsExample, want: 15},
		{name: "simultaneous completions", input: `Step A must be finished before step D can begin.
Step C must be finished before step B can begin.
Step D must be finished before step H can begin.
Step B must be finished before step E can begin.
Step B must be finished before step F can begin.`, want: 18},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assemblyTime(strings.NewReader(tt.input), 2, 0)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
