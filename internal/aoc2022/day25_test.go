package aoc2022

import (
	"math/big"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader("1=-0-2\n12111\n2=0=\n21\n2=01\n111\n20012\n112\n1=-1=\n1-12\n12\n1=\n122"), Want: "2=-1=0"},
		{Input: strings.NewReader("0"), Want: "0"},
		{Input: strings.NewReader("="), Want: "="},
		{Input: strings.NewReader("=\n-"), Want: "-2"},
		{Input: strings.NewReader("2\n="), Want: "0"},
		{Input: strings.NewReader("1" + strings.Repeat("0", 100) + "\n="), Want: "1" + strings.Repeat("0", 99) + "="},
		{Input: aoc.FileInput(t, 2022, 25), Want: "2011-=2=-1020-1===-1"},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func TestSNAFUEncoding(t *testing.T) {
	tests := []struct {
		value int64
		want  string
	}{
		{value: -3, want: "-2"},
		{value: -2, want: "="},
		{value: -1, want: "-"},
		{value: 0, want: "0"},
		{value: 3, want: "1="},
		{value: 4, want: "1-"},
		{value: 5, want: "10"},
		{value: 2022, want: "1=11-2"},
		{value: 12345, want: "1-0---0"},
		{value: 314159265, want: "1121-1110-1=0"},
	}
	for _, tt := range tests {
		if got := encodeSNAFU(big.NewInt(tt.value)); got != tt.want {
			t.Fatalf("encode %d: got %q; want %q", tt.value, got, tt.want)
		}
	}
}

func TestSNAFUInvalidInput(t *testing.T) {
	for _, input := range []string{"", "3", "1\n\n2", strings.Repeat("1", 1025)} {
		if _, err := day25p01(strings.NewReader(input)); err == nil {
			t.Fatalf("expected SNAFU error for %q", input)
		}
	}
}
