package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1721
979
366
299
675
1456`), Want: "514579"},
		{Input: aoc.FileInput(t, 2020, 1), Want: "842016"},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(`1721
979
366
299
675
1456`), Want: "241861950"},
		{Input: aoc.FileInput(t, 2020, 1), Want: "9199664"},
	}
	aoc.AOCTest(t, day01p02, tests)
}

func TestExpenseEntriesAreDistinct(t *testing.T) {
	aoc.AOCTest(t, day01p01, []aoc.TestInput{{Input: strings.NewReader("1010\n1010"), Want: "1020100"}})
	aoc.AOCTest(t, day01p02, []aoc.TestInput{{Input: strings.NewReader("673\n673\n674"), Want: "305274146"}})
	if _, err := day01p01(strings.NewReader("1010")); err == nil {
		t.Fatal("single entry cannot form a pair")
	}
	if _, err := day01p02(strings.NewReader("673\n674")); err == nil {
		t.Fatal("two entries cannot form a triple")
	}
}

func TestDay01ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"01/1", day01p01},
		{"01/2", day01p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay01MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"expense", "not a number", day01p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid-input error")
			}
		})
	}
}
