package aoc2019

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day04p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("111111-111111"), Want: "1"},
		{Input: strings.NewReader("223450-223450"), Want: "0"},
		{Input: strings.NewReader("123789-123789"), Want: "0"},
		{
			Input: aoc.FileInput(t, 2019, 4),
			Want:  "1660",
		},
	}
	aoc.AOCTest(t, day04p01, tests)
}

func Test_day04p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("112233-112233"), Want: "1"},
		{Input: strings.NewReader("123444-123444"), Want: "0"},
		{Input: strings.NewReader("111122-111122"), Want: "1"},
		{
			Input: aoc.FileInput(t, 2019, 4),
			Want:  "1135",
		},
	}
	aoc.AOCTest(t, day04p02, tests)
}

func TestPasswordRangeErrors(t *testing.T) {
	for _, input := range []string{"", "123456", "x-5", "5-x", "-5-10", "20-10", "1-2-3"} {
		for _, solve := range []func(io.Reader) (string, error){day04p01, day04p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected password range error for %q", input)
			}
		}
	}
}
