package aoc2019

import (
	"slices"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func TestApplyFFTPhase(t *testing.T) {
	signal := []int{1, 2, 3, 4, 5, 6, 7, 8}
	for phase, want := range [][]int{
		{4, 8, 2, 2, 6, 1, 5, 8},
		{3, 4, 0, 4, 0, 4, 3, 8},
		{0, 3, 4, 1, 5, 5, 1, 8},
		{0, 1, 0, 2, 9, 4, 9, 8},
	} {
		signal = applyFFTPhase(signal)
		if !slices.Equal(signal, want) {
			t.Fatalf("phase %d: got %v, want %v", phase+1, signal, want)
		}
	}

	for _, tt := range []struct {
		name        string
		input, want []int
	}{
		{"odd length", []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{5, 8, 3, 2, 5, 0, 4, 7, 9}},
		{"singleton", []int{9}, []int{9}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyFFTPhase(tt.input); !slices.Equal(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_day16p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("80871224585914546619083218645595"),
			Want:  "24176176",
		},
		{
			Input: strings.NewReader("19617804207202209144916044189917"),
			Want:  "73745418",
		},
		{
			Input: strings.NewReader("69317163492948606335995924319873"),
			Want:  "52432133",
		},
		{
			Input: aoc.FileInput(t, 2019, 16),
			Want:  "40921727",
		},
	}
	aoc.AOCTest(t, day16p01, tests)
}

func Test_day16p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("03036732577212944063491565474664"),
			Want:  "84462026",
		},
		{
			Input: strings.NewReader("02935109699940807407585447034323"),
			Want:  "78725270",
		},
		{
			Input: strings.NewReader("03081770884921959731165446850517"),
			Want:  "53553731",
		},
		{
			Input: aoc.FileInput(t, 2019, 16),
			Want:  "89950138",
		},
	}
	aoc.AOCTest(t, day16p02, tests)
}
