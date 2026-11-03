package aoc2021

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day12Example = `start-A
start-b
A-c
A-b
b-d
A-end
b-end`

func Test_day12p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day12Example), Want: "10"},
		{Input: strings.NewReader(`start-end`), Want: "1"},
		{Input: strings.NewReader(`start-A
A-B
start-end`), Want: "1"},
		{Input: strings.NewReader(`start-a
a-A
A-B
a-end`), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 12), Want: "4754"},
	}

	aoc.AOCTest(t, day12p01, tests)
}

func Test_day12p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day12Example), Want: "36"},
		{Input: strings.NewReader(`start-end`), Want: "1"},
		{Input: strings.NewReader(`start-A
A-B
start-end`), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 12), Want: "143562"},
	}

	aoc.AOCTest(t, day12p02, tests)
}

func TestDay12UsefulCyclePartOne(t *testing.T) {
	if _, err := day12p01(strings.NewReader(`start-A
A-B
B-end`)); err == nil {
		t.Fatal("expected useful cycle part one error")
	}
}

func TestDay12UsefulCyclePartTwo(t *testing.T) {
	if _, err := day12p02(strings.NewReader(`start-A
A-B
B-end`)); err == nil {
		t.Fatal("expected useful cycle part two error")
	}
}

func TestDay12StateDependentUsefulCycle(t *testing.T) {
	if _, err := day12p02(strings.NewReader(`start-a
a-A
A-B
a-end`)); err == nil {
		t.Fatal("expected state-dependent useful cycle error")
	}
}

func TestDay12ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day12p01},
		{"part2", day12p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
