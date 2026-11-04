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
		{"unique noninteger solution", 2, 0, 0, 2, 1, 2, "0"},
		{"prize off line", 2, 2, 1, 1, 4, 5, "0"},
		{"zero button", 0, 0, 2, 3, 4, 6, "2"},
		{"zero button B", 2, 3, 0, 0, 4, 6, "6"},
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

func TestClawPrizeOffsetOverflow(t *testing.T) {
	input := "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=9223372036854775807, Y=9223372036854775807"
	got, err := day13p02(strings.NewReader(input))
	if err != nil || got != "36893528147419103228" {
		t.Fatalf("got %q, %v; want 36893528147419103228", got, err)
	}
}

func TestClawInputFormatting(t *testing.T) {
	const valid = "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=1, Y=0"
	for _, tt := range []struct {
		name, input string
	}{
		{"edge whitespace, CRLF and signs", "\u2003\r\nButton A: X++0001, Y+-0\r\nButton B: X+0, Y+1\r\nPrize: X=+0001, Y=-0\r\n\u00a0"},
		{"long numeric field", strings.Replace(valid, "X+1", "X+"+strings.Repeat("0", 70_000)+"1", 1)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day13p01(strings.NewReader(tt.input))
			if err != nil || got != "3" {
				t.Fatalf("got %q, %v; want 3", got, err)
			}
		})
	}
	for _, tt := range []struct {
		name, input string
		cause       error
	}{
		{"internal whitespace", strings.Replace(valid, "Button B", "\u2003Button B", 1), nil},
		{"missing separator", valid + "\n" + valid, nil},
		{"extra separator", valid + "\n\n\n" + valid, nil},
		{"numeric syntax", strings.Replace(valid, "X+1", "X+1x", 1), strconv.ErrSyntax},
		{"negative coordinate", strings.Replace(valid, "X+1", "X+-1", 1), nil},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := day13p01(strings.NewReader(tt.input)); err == nil || tt.cause != nil && !errors.Is(err, tt.cause) {
				t.Fatalf("got %v; want rejection with cause %v", err, tt.cause)
			}
		})
	}
}

func TestClawFrameErrors(t *testing.T) {
	const (
		invalidA = "Button A: X+1x, Y+0"
		buttonB  = "Button B: X+0, Y+1"
		prize    = "Prize: X=1, Y=0"
		valid    = "Button A: X+1, Y+0\n" + buttonB + "\n" + prize
		shape    = "needs two buttons and a prize"
	)
	for _, tt := range []struct {
		name, input, wantContext string
	}{
		{"empty input", "", shape},
		{"missing prize before numeric syntax", invalidA + "\n" + buttonB, shape},
		{"empty button B before numeric syntax", invalidA + "\n\n" + prize, shape},
		{"empty prize before numeric syntax", invalidA + "\n" + buttonB + "\n\n\n" + valid, shape},
		{"fourth line", valid + "\nextra", shape},
		{"empty button A in three-line frame", valid + "\n\n\n" + buttonB + "\n" + prize, "invalid claw coordinate line"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day13p01(strings.NewReader(tt.input))
			if got != "" || err == nil {
				t.Fatalf("got %q, %v; want empty result and %q error", got, err, tt.wantContext)
			}
			if !strings.Contains(err.Error(), tt.wantContext) || errors.Is(err, strconv.ErrSyntax) {
				t.Errorf("got %v; want %q error without numeric syntax cause", err, tt.wantContext)
			}
			if tt.wantContext != shape && strings.Contains(err.Error(), shape) {
				t.Errorf("got %v; want coordinate error after accepting frame shape", err)
			}
		})
	}
}

func TestClawDecimalBoundaries(t *testing.T) {
	const (
		buttonB = "Button B: X+0, Y+1"
		prize   = "Prize: X=1, Y=0"
		valid   = "Button A: X+1, Y+0\n" + buttonB + "\n" + prize
	)
	t.Run("eighteen digit coordinates", func(t *testing.T) {
		input := "Button A: X+1, Y+0\n" + buttonB + "\nPrize: X=999999999999999999, Y=999999999999999999"
		got, err := day13p02(strings.NewReader(input))
		if err != nil || got != "4000039999999999996" {
			t.Fatalf("got %q, %v; want 4000039999999999996", got, err)
		}
	})
	for _, tt := range []struct {
		name, value string
		cause       error
	}{
		{"empty Y", "", strconv.ErrSyntax},
		{"trailing junk in Y", "999999999999999999x", strconv.ErrSyntax},
		{"overflowing Y", "9223372036854775808", strconv.ErrRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			input := valid + "\n\nButton A: X+1, Y+" + tt.value + "\n" + buttonB + "\n" + prize
			got, err := day13p01(strings.NewReader(input))
			if got != "" || err == nil {
				t.Fatalf("got %q, %v; want empty result and numeric error", got, err)
			}
			if !strings.Contains(err.Error(), "claw y coordinate") || !errors.Is(err, tt.cause) {
				t.Errorf("got %v; want claw y coordinate context and %v", err, tt.cause)
			}
			number, ok := errors.AsType[*strconv.NumError](err)
			if !ok {
				t.Fatalf("got %v; want *strconv.NumError", err)
			}
			if number.Func != "ParseInt" || number.Num != tt.value || !errors.Is(number.Err, tt.cause) {
				t.Errorf("got %#v; want ParseInt(%q) with cause %v", number, tt.value, tt.cause)
			}
		})
	}
}

func TestClawCoordinateSeparatorErrors(t *testing.T) {
	const (
		buttonA = "Button A: X+1, Y+0"
		buttonB = "Button B: X+0, Y+1"
		prize   = "Prize: X=1, Y=0"
		valid   = buttonA + "\n" + buttonB + "\n" + prize
	)
	for _, tt := range []struct {
		name, invalid, wantContext, wantNumber string
	}{
		{"later valid separator", "Button A: X+1, Z+garbage, Y+2\n" + buttonB + "\n" + prize, "claw x coordinate", "1, Z+garbage"},
		{"wrong prize separator", buttonA + "\n" + buttonB + "\nPrize: X=1, Y+0", "invalid claw coordinates", ""},
		{"truncated separator", "Button A: X+1,\n" + buttonB + "\n" + prize, "invalid claw coordinates", ""},
		{"missing separator", "Button A: X+1\n" + buttonB + "\n" + prize, "invalid claw coordinates", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := day13p01(strings.NewReader(valid + "\n\n" + tt.invalid))
			if got != "" || err == nil {
				t.Fatalf("got %q, %v; want empty result and %q error", got, err, tt.wantContext)
			}
			if !strings.Contains(err.Error(), tt.wantContext) {
				t.Errorf("got %v; want %q context", err, tt.wantContext)
			}
			number, ok := errors.AsType[*strconv.NumError](err)
			if tt.wantNumber == "" {
				if ok || errors.Is(err, strconv.ErrSyntax) {
					t.Errorf("got %v; want coordinate error without numeric syntax cause", err)
				}
				return
			}
			if !ok {
				t.Fatalf("got %v; want numeric syntax error", err)
			}
			if number.Func != "ParseInt" || number.Num != tt.wantNumber ||
				!errors.Is(number.Err, strconv.ErrSyntax) || !errors.Is(err, strconv.ErrSyntax) {
				t.Errorf("got %#v wrapped by %v; want ParseInt(%q) with syntax cause", number, err, tt.wantNumber)
			}
		})
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

func TestClawInputLimits(t *testing.T) {
	const machine = "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=1, Y=2"
	t.Run("machine count", func(t *testing.T) {
		input := strings.Repeat(machine+"\n\n", 10_000)
		got, err := day13p01(strings.NewReader(input))
		if err != nil || got != "50000" {
			t.Fatalf("at limit: got %q, %v; want 50000", got, err)
		}
		if _, err := day13p01(strings.NewReader(input + machine)); err == nil {
			t.Fatal("accepted more than 10000 machines")
		}
	})
	t.Run("original byte count", func(t *testing.T) {
		const limit = 4 << 20
		padding := limit - len(machine)
		input := machine + strings.Repeat("\r\n", padding/2) + strings.Repeat(" ", padding%2)
		got, err := day13p01(strings.NewReader(input))
		if err != nil || got != "5" {
			t.Fatalf("at limit: got %q, %v; want 5", got, err)
		}
		if _, err := day13p01(strings.NewReader(input + "\r")); err == nil {
			t.Fatal("accepted more than 4 MiB before normalization")
		}
	})
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
