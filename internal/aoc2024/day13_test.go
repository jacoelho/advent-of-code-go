package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`),
			Want: "480",
		},
		{
			Input: aoc.FileInput(t, 2024, 13),
			Want:  "28138",
		},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`),
			Want: "875318608908",
		},
		{
			Input: aoc.FileInput(t, 2024, 13),
			Want:  "108394825772874",
		},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func TestClawFeasibilityAndMinimum(t *testing.T) {
	tests := []struct {
		name       string
		ax, ay     int64
		bx, by     int64
		px, py     int64
		wantTokens string
	}{
		{"negative solution", 1, 1, 1, 2, 1, 3, "0"},
		{"above press limit", 1, 1, 1, 2, 102, 103, "0"},
		{"at press limit", 1, 1, 1, 2, 101, 102, "301"},
		{"dependent cheaper B", 2, 2, 1, 1, 4, 4, "4"},
		{"dependent cheaper A", 4, 4, 1, 1, 8, 8, "6"},
		{"dependent bounded minimum", 2, 2, 1, 1, 250, 250, "325"},
		{"noninteger solution", 4, 4, 6, 6, 5, 5, "0"},
		{"prize off line", 2, 2, 1, 1, 4, 5, "0"},
		{"zero button", 0, 0, 2, 3, 4, 6, "2"},
		{"both zero buttons", 0, 0, 0, 0, 1, 1, "0"},
		{"vertical buttons", 0, 2, 0, 1, 0, 4, "4"},
		{"large determinant", 4_000_000_000, 1, 1, 4_000_000_000, 8_000_000_003, 12_000_000_002, "9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf("Button A: X+%d, Y+%d\nButton B: X+%d, Y+%d\nPrize: X=%d, Y=%d", tt.ax, tt.ay, tt.bx, tt.by, tt.px, tt.py)
			got, err := day13p01(strings.NewReader(input))
			if err != nil || got != tt.wantTokens {
				t.Fatalf("got %q, %v; want %q", got, err, tt.wantTokens)
			}
		})
	}
}

func TestClawDependentPrizeOffset(t *testing.T) {
	input := "Button A: X+2, Y+2\nButton B: X+1, Y+1\nPrize: X=4, Y=4"
	got, err := day13p02(strings.NewReader(input))
	if err != nil || got != "10000000000004" {
		t.Fatalf("got %q, %v; want 10000000000004", got, err)
	}
}

func TestClawInvalidInput(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day13p01, day13p02} {
		if got, err := solve(strings.NewReader("Button A: X+1, Y+2")); err == nil {
			t.Errorf("got %q, want incomplete machine error", got)
		}
		input := "Button A: X+9223372036854775808, Y+1\nButton B: X+1, Y+2\nPrize: X=3, Y=3"
		if _, err := solve(strings.NewReader(input)); !errors.Is(err, strconv.ErrRange) {
			t.Errorf("got %v, want numeric range error", err)
		}
	}
}

func TestDay13ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day13p01, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=3\n\n"},
		{"part2", day13p02, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=3\n\n"},
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
