package aoc2025

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("12"), Want: "12"},
		{
			Input: strings.NewReader(`987654321111111
811111111111119
234234234234278
818181911112111`),
			Want: "357",
		},
		{
			Input: aoc.FileInput(t, 2025, 03),
			Want:  "17445",
		},
	}
	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("123456789012"), Want: "123456789012"},
		{
			Input: strings.NewReader(`987654321111111
811111111111119
234234234234278
818181911112111`),
			Want: "3121910778619",
		},
		{
			Input: aoc.FileInput(t, 2025, 03),
			Want:  "173229689350551",
		},
	}
	aoc.AOCTest(t, day03p02, tests)
}

func TestBatteryBankValidation(t *testing.T) {
	for _, part := range []struct {
		name  string
		solve func(io.Reader) (string, error)
		short string
	}{
		{"part1", day03p01, "1\n"},
		{"part2", day03p02, "12345678901\n"},
	} {
		t.Run(part.name, func(t *testing.T) {
			for _, bank := range []struct{ name, input string }{
				{"short", part.short},
				{"nondigit", "12345678901x\n"},
				{"blank", "\n"},
			} {
				t.Run(bank.name, func(t *testing.T) {
					input := strings.NewReader("987654321111111\n" + bank.input)
					if result, err := part.solve(input); err == nil {
						t.Fatalf("got %q without error for malformed bank", result)
					}
				})
			}
		})
	}
}
