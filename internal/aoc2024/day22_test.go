package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1
10
100
2024`),
			Want: "37327623",
		},
		{
			Input: aoc.FileInput(t, 2024, 22),
			Want:  "12759339434",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1
2
3
2024`),
			Want: "23",
		},
		{
			Input: aoc.FileInput(t, 2024, 22),
			Want:  "1405",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}

func TestSecretFinalPrice(t *testing.T) {
	prices := sellSequences(1)
	if got, ok := prices[[4]int{-6, 3, -1, 7}]; !ok || got != 9 {
		t.Fatalf("final price = %d, present %t; want 9", got, ok)
	}
}

func TestSecretNumberInputError(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day22p01, day22p02} {
		if _, err := solve(strings.NewReader("invalid")); err == nil {
			t.Fatal("expected invalid secret number error")
		}
	}
}

func TestDay22ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day22p01, "1\n"},
		{"part2", day22p02, "1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
