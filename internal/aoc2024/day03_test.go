package aoc2024

import (
	"io"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day03p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))`),
			Want:  "161",
		},
		{
			Input: aoc.FileInput(t, 2024, 3),
			Want:  "173517243",
		},
	}

	aoc.AOCTest(t, day03p01, tests)
}

func Test_day03p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))`),
			Want:  "48",
		},
		{
			Input: aoc.FileInput(t, 2024, 3),
			Want:  "100450138",
		},
	}

	aoc.AOCTest(t, day03p02, tests)
}

func TestCorruptedMemoryGrammar(t *testing.T) {
	tests := []struct {
		input string
		want  [2]string
	}{
		{"mul(1000,2)mul(3,4)", [2]string{"12", "12"}},
		{"mul(1,2)mul(999,999)mul(2,1000)", [2]string{"998003", "998003"}},
		{"don't()mul(1000,2)mul(2,3)do()mul(3,4)mul(1,1000)", [2]string{"18", "12"}},
	}
	for _, tt := range tests {
		for part, solve := range []func(io.Reader) (string, error){day03p01, day03p02} {
			got, err := solve(strings.NewReader(tt.input))
			if err != nil || got != tt.want[part] {
				t.Errorf("part %d input %q: got %q, %v; want %q", part+1, tt.input, got, err, tt.want[part])
			}
		}
	}
}
