package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("3"),
			Want:  "638",
		},
		{
			Input: strings.NewReader("0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2017, 17),
			Want:  "777",
		},
	}

	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	if testing.Short() {
		t.Skip("full 50-million-insertion spinlock fixture")
	}

	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2017, 17),
			Want:  "39289581",
		},
	}

	aoc.AOCTest(t, day17p02, tests)
}

func TestSpinlockAfterZero(t *testing.T) {
	tests := []struct {
		insertions, want int
	}{
		{insertions: 1, want: 1},
		{insertions: 2, want: 2},
		{insertions: 3, want: 2},
		{insertions: 5, want: 5},
		{insertions: 9, want: 9},
	}

	for _, tt := range tests {
		if got := spinlockAfterZero(3, tt.insertions); got != tt.want {
			t.Errorf("insertions=%d: got %d, want %d", tt.insertions, got, tt.want)
		}
	}
}

func TestSpinlockErrors(t *testing.T) {
	for _, input := range []string{"", "-1", "one"} {
		if _, err := day17p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
}

func TestDay17ReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"17/1", day17p01},
		{"17/2", day17p02},
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
