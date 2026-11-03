package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day09p01(t *testing.T) {
	aoc.AOCTest(t, day09p01, []aoc.TestInput{
		{Input: strings.NewReader("ADVENT"), Want: "6"},
		{Input: strings.NewReader("A(1x5)BC"), Want: "7"},
		{Input: strings.NewReader("(3x3)XYZ"), Want: "9"},
		{Input: strings.NewReader("A(2x2)BCD(2x2)EFG"), Want: "11"},
		{Input: strings.NewReader("(6x1)(1x3)A"), Want: "6"},
		{Input: strings.NewReader("X(8x2)(3x3)ABCY"), Want: "18"},
		{Input: aoc.FileInput(t, 2016, 9), Want: "150914"},
	})
}

func Test_day09p02(t *testing.T) {
	aoc.AOCTest(t, day09p02, []aoc.TestInput{
		{Input: strings.NewReader("X(8x2)(3x3)ABCY"), Want: "20"},
		{Input: strings.NewReader("(27x12)(20x12)(13x14)(7x10)(1x12)A"), Want: "241920"},
		{Input: strings.NewReader("(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN"), Want: "445"},
		{Input: strings.NewReader(" A \n (1x5) B\t C "), Want: "7"},
		{Input: aoc.FileInput(t, 2016, 9), Want: "11052855125"},
	})
}

func TestDecompressionErrors(t *testing.T) {
	for _, input := range []string{"(2x3)A", "(2x3", "(x3)AB", "(2x-1)AB", "(6x2)(4x2)A"} {
		if _, err := day09p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day09p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day09p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
