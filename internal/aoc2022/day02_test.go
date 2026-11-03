package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day02Example = `A Y
B X
C Z`

func Test_day02p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day02Example), Want: "15"},
		{Input: aoc.FileInput(t, 2022, 2), Want: "11475"},
	}
	aoc.AOCTest(t, day02p01, tests)
}

func Test_day02p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day02Example), Want: "12"},
		{Input: aoc.FileInput(t, 2022, 2), Want: "16862"},
	}
	aoc.AOCTest(t, day02p02, tests)
}

func TestStrategyInvalidRound(t *testing.T) {
	if _, err := day02p02(strings.NewReader("A Q")); err == nil {
		t.Fatal("expected invalid round error")
	}
}

func TestDay02ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day02p01},
		{"part2", day02p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
