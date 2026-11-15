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
	key := 0
	for _, change := range []int{-6, 3, -1, 7} {
		key = key*19 + change + 9
	}
	if got := day22BananaTotals([]int{1})[key]; got != 9 {
		t.Fatalf("final price = %d; want 9", got)
	}
}

func TestSecretNumberInputError(t *testing.T) {
	for _, input := range []string{"invalid", "-1"} {
		for _, solve := range []func(io.Reader) (string, error){day22p01, day22p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected invalid secret number error", input)
			}
		}
	}
}

func TestSecretNumbersEmptyHasNoBestSequence(t *testing.T) {
	if got, err := day22p02(strings.NewReader("")); err == nil {
		t.Fatalf("got %q, want no secret numbers error", got)
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
