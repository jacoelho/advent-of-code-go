package aoc2025

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day01p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`),
			Want: "3",
		},
		{
			Input: aoc.FileInput(t, 2025, 01),
			Want:  "1055",
		},
	}
	aoc.AOCTest(t, day01p01, tests)
}

func Test_day01p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`),
			Want: "6",
		},
		{
			Input: aoc.FileInput(t, 2025, 01),
			Want:  "6386",
		},
	}
	aoc.AOCTest(t, day01p02, tests)
}

func TestDialEndingsAndClicks(t *testing.T) {
	tests := []struct {
		input string
		want  [2]string
	}{
		{"R50", [2]string{"1", "1"}},
		{"R50\nR1", [2]string{"1", "1"}},
		{"R50\nR0", [2]string{"2", "1"}},
		{"R9223372036854775807", [2]string{"0", "92233720368547758"}},
		{"L9223372036854775807", [2]string{"0", "92233720368547758"}},
	}
	for _, tt := range tests {
		for part, solve := range []func(io.Reader) (string, error){day01p01, day01p02} {
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != tt.want[part] {
				t.Errorf("part %d input %q: got %q, %v; want %q", part+1, tt.input, got, err, tt.want[part])
			}
		}
	}
}

func TestDialClicksAgainstLiteralRotations(t *testing.T) {
	for position := int64(0); position < 100; position++ {
		for amount := int64(-300); amount <= 300; amount++ {
			steps, direction := amount, int64(1)
			if steps < 0 {
				steps, direction = -steps, -1
			}
			next, want := position, int64(0)
			for range steps {
				next += direction
				if next == 100 {
					next = 0
				} else if next == -1 {
					next = 99
				}
				if next == 0 {
					want++
				}
			}
			if got := countClicks(position, amount); got != want {
				t.Fatalf("position %d, amount %d: got %d, want %d", position, amount, got, want)
			}
		}
	}
}

func TestDialPasswordOverflow(t *testing.T) {
	input := strings.Repeat("R9223372036854775807\n", 100)
	if got, err := day01p02(strings.NewReader(input)); err != nil || got != "9223372036854775807" {
		t.Errorf("got %q, %v; want MaxInt64", got, err)
	}
	if got, err := day01p02(strings.NewReader(input + "R9223372036854775807\n")); err == nil {
		t.Errorf("got %q, want overflow error", got)
	}
}

func TestDialInvalidInput(t *testing.T) {
	for _, input := range []string{"X1", "L-1", "R+1", "R9223372036854775808", "R1\n\nR2"} {
		for _, solve := range []func(io.Reader) (string, error){day01p01, day01p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
	}
}
