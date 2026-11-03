package aoc2020

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2020, 9), Want: "1309761972"},
	}
	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2020, 9), Want: "177989832"},
	}
	aoc.AOCTest(t, day09p02, tests)
}

func TestXMASExample(t *testing.T) {
	numbers := []int64{35, 20, 15, 25, 47, 40, 62, 55, 65, 95, 102, 117, 150, 182, 127, 219, 299, 277, 309, 576}
	invalid, err := invalidXMASNumber(numbers, 5)
	if err != nil || invalid != 127 {
		t.Fatalf("got invalid number %d, %v; want 127", invalid, err)
	}
	weakness, err := xmasWeakness(numbers, 127)
	if err != nil || weakness != 62 {
		t.Fatalf("got weakness %d, %v; want 62", weakness, err)
	}
}

func TestXMASPairMustHaveDifferentValues(t *testing.T) {
	for _, tt := range []struct {
		numbers []int64
		want    int64
	}{
		{[]int64{10, 10, 20}, 20},
		{[]int64{10, 11, 20}, 20},
		{[]int64{10, 11, 21, 32, 54}, 54},
	} {
		got, err := invalidXMASNumber(tt.numbers, 2)
		if err != nil || got != tt.want {
			t.Fatalf("%v: got %d, %v; want %d", tt.numbers, got, err, tt.want)
		}
	}
	if _, err := xmasWeakness([]int64{127, 1, 2}, 127); err == nil {
		t.Fatal("single number is not a contiguous range of two")
	}
}

func TestDay09ReadErrors(t *testing.T) {
	failure := errors.New("read failure")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"09/1", day09p01},
		{"09/2", day09p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(iotest.ErrReader(failure)); !errors.Is(err, failure) {
				t.Fatalf("got %v, want input failure", err)
			}
		})
	}
}

func TestDay09MalformedInput(t *testing.T) {
	for _, tt := range []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"XMAS", "not a number", day09p01},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected malformed-input error")
			}
		})
	}
}
