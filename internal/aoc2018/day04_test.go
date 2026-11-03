package aoc2018

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const guardExample = `[1518-11-01 00:00] Guard #10 begins shift
[1518-11-01 00:05] falls asleep
[1518-11-01 00:25] wakes up
[1518-11-01 00:30] falls asleep
[1518-11-01 00:55] wakes up
[1518-11-01 23:58] Guard #99 begins shift
[1518-11-02 00:40] falls asleep
[1518-11-02 00:50] wakes up
[1518-11-03 00:05] Guard #10 begins shift
[1518-11-03 00:24] falls asleep
[1518-11-03 00:29] wakes up
[1518-11-04 00:02] Guard #99 begins shift
[1518-11-04 00:36] falls asleep
[1518-11-04 00:46] wakes up
[1518-11-05 00:03] Guard #99 begins shift
[1518-11-05 00:45] falls asleep
[1518-11-05 00:55] wakes up`

func Test_day04p01(t *testing.T) {
	records := strings.Split(guardExample, "\n")
	slices.Reverse(records)
	tests := []aoc.TestInput{
		{Input: strings.NewReader(strings.Join(records, "\n")), Want: "240"},
		{Input: aoc.FileInput(t, 2018, 4), Want: "48680"},
	}

	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(guardExample), Want: "4455"},
		{Input: aoc.FileInput(t, 2018, 4), Want: "94826"},
	}

	aoc.AOCTest(t, day04p02, tests)
}

func TestDay04ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"04/1", day04p01},
		{"04/2", day04p02},
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

func TestDay04InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{name: "unpaired waking", solve: day04p01, input: "[1518-11-01 00:05] wakes up"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected an input error")
			}
		})
	}
}
