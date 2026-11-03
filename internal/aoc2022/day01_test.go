package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day01Example = `1000
2000
3000

4000

5000
6000

7000
8000
9000

10000`

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "24000"},
		{Input: aoc.FileInput(t, 2022, 1), Want: "69177"},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day01Example), Want: "45000"},
		{Input: aoc.FileInput(t, 2022, 1), Want: "207456"},
	}
	aoc.AOCTest(t, day01p02, tests)
}

func TestCaloriesInvalidInput(t *testing.T) {
	if _, err := day01p01(strings.NewReader("-1")); err == nil {
		t.Fatal("expected negative calorie error")
	}
}

func TestCalorieOverflow(t *testing.T) {
	if _, err := day01p01(strings.NewReader("9223372036854775807\n1")); err == nil {
		t.Fatal("expected calorie group overflow")
	}
	if _, err := day01p02(strings.NewReader("9223372036854775807\n\n1")); err == nil {
		t.Fatal("expected calorie total overflow")
	}
}

func TestDay01ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day01p01},
		{"part2", day01p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
