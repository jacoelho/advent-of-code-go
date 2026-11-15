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
		{"dependent coprime steps", 3, 3, 5, 5, 31, 31, "11"},
		{"noninteger solution", 4, 4, 6, 6, 5, 5, "0"},
		{"unique noninteger solution", 2, 0, 0, 2, 1, 2, "0"},
		{"prize off line", 2, 2, 1, 1, 4, 5, "0"},
		{"zero button", 0, 0, 2, 3, 4, 6, "2"},
		{"zero button B", 2, 3, 0, 0, 4, 6, "6"},
		{"both zero buttons", 0, 0, 0, 0, 1, 1, "0"},
		{"vertical buttons", 0, 2, 0, 1, 0, 4, "4"},
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

func TestClawEighteenDigitPrizeWithOffset(t *testing.T) {
	input := "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=999999999999999999, Y=999999999999999999"
	got, err := day13p02(strings.NewReader(input))
	if err != nil || got != "4000039999999999996" {
		t.Fatalf("got %q, %v; want 4000039999999999996", got, err)
	}
}

func TestClawArithmeticOverflowIsRejected(t *testing.T) {
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
		input string
	}{
		{"prize offset", day13p02, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=9223372036854775807, Y=9223372036854775807"},
		{"determinant", day13p01, "Button A: X+4000000000, Y+1\nButton B: X+1, Y+4000000000\nPrize: X=8000000003, Y=12000000002"},
		{"token cost", day13p02, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3074457345618258602, Y=0"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.solve(strings.NewReader(tt.input)); !errors.Is(err, errDay13Overflow) {
				t.Fatalf("got %q, %v; want overflow error", got, err)
			}
		})
	}
}

func TestClawAcceptsCRLF(t *testing.T) {
	input := "Button A: X+1, Y+0\r\nButton B: X+0, Y+1\r\nPrize: X=1, Y=0\r\n\r\nButton A: X+1, Y+0\r\nButton B: X+0, Y+1\r\nPrize: X=0, Y=1\r\n"
	got, err := day13p01(strings.NewReader(input))
	if err != nil || got != "4" {
		t.Fatalf("got %q, %v; want 4", got, err)
	}
}

func TestClawRejectsMalformedInput(t *testing.T) {
	const (
		buttonA = "Button A: X+1, Y+0"
		buttonB = "Button B: X+0, Y+1"
		prize   = "Prize: X=1, Y=0"
		valid   = buttonA + "\n" + buttonB + "\n" + prize
	)
	for _, tt := range []struct {
		name, input, wantContext string
		cause                    error
	}{
		{"empty input", "", "two buttons and a prize", nil},
		{"missing prize", buttonA + "\n" + buttonB, "two buttons and a prize", nil},
		{"fourth line", valid + "\nextra", "two buttons and a prize", nil},
		{"missing separator", valid + "\n" + valid, "two buttons and a prize", nil},
		{"extra separator", valid + "\n\n\n" + valid, "two buttons and a prize", nil},
		{"text in separator line", valid + "\nx\n" + valid, "line 4", nil},
		{"leading whitespace", " " + valid, "line 1", nil},
		{"swapped buttons", buttonB + "\n" + buttonA + "\n" + prize, "line 1", nil},
		{"wrong prize separator", buttonA + "\n" + buttonB + "\nPrize: X=1, Y+0", "line 3", nil},
		{"signed coordinate", strings.Replace(valid, "X+1", "X++1", 1), "claw x coordinate", nil},
		{"negative coordinate", strings.Replace(valid, "X+1", "X+-1", 1), "claw x coordinate", nil},
		{"empty coordinate", strings.Replace(valid, "Y+0", "Y+", 1), "claw y coordinate", nil},
		{"trailing junk", strings.Replace(valid, "Y+0", "Y+0x", 1), "claw y coordinate", strconv.ErrSyntax},
		{"out of range", strings.Replace(valid, "Y+0", "Y+9223372036854775808", 1), "claw y coordinate", strconv.ErrRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day13p01(strings.NewReader(tt.input))
			if err == nil {
				t.Fatalf("got %q; want %q error", got, tt.wantContext)
			}
			if !strings.Contains(err.Error(), tt.wantContext) || tt.cause != nil && !errors.Is(err, tt.cause) {
				t.Errorf("got %v; want %q context with cause %v", err, tt.wantContext, tt.cause)
			}
		})
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
			for _, prefix := range []string{"", "malformed prefix", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
