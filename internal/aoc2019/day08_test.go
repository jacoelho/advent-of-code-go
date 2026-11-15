package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 8),
			Want:  "2684",
		},
	}
	aoc.AOCTest(t, day08p01, tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2019, 8),
			Want:  "YGRYZ",
		},
	}
	aoc.AOCTest(t, day08p02, tests)
}

func Test_day08p02Layers(t *testing.T) {
	rows := []string{
		"#...#" + ".##.." + "####." + ".##.." + "####.",
		"#...#" + "#..#." + "...#." + "#..#." + "...#.",
		".#.#." + "#..#." + "..#.." + "#..#." + "..#..",
		"..#.." + "####." + ".#..." + "####." + ".#...",
		"..#.." + "#..#." + "#...." + "#..#." + "#....",
		"..#.." + "#..#." + "####." + "#..#." + "####.",
	}
	var top, bottom strings.Builder
	// Reveal Y through transparency; opaque pixels must hide the inverted lower layer.
	for _, row := range rows {
		for x, pixel := range row {
			value := byte('0')
			if pixel == '#' {
				value = '1'
			}
			if x < 5 {
				top.WriteByte('2')
				bottom.WriteByte(value)
			} else {
				top.WriteByte(value)
				bottom.WriteByte('0' + '1' - value)
			}
		}
	}
	aoc.AOCTest(t, day08p02, []aoc.TestInput{
		{Input: strings.NewReader(top.String() + bottom.String()), Want: "YAZAZ"},
	})
}

func TestImageLayerErrors(t *testing.T) {
	for _, input := range []string{"", strings.Repeat("0", 151), strings.Repeat("0", 149) + "3"} {
		for _, solve := range []func(io.Reader) (string, error){day08p01, day08p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("expected image error for %d pixels", len(input))
			}
		}
	}
}

func Test_day08p02Errors(t *testing.T) {
	t.Run("unknown glyph", func(t *testing.T) {
		input := strings.Repeat("11111"+strings.Repeat("0", 20), 6)
		if _, err := day08p02(strings.NewReader(input)); err == nil {
			t.Fatal("expected unknown glyph error")
		}
	})
	t.Run("fully transparent pixel", func(t *testing.T) {
		if _, err := day08p02(strings.NewReader(strings.Repeat("2", 150))); err == nil {
			t.Fatal("expected transparent pixel error")
		}
	})
	t.Run("reader failure", func(t *testing.T) {
		failure := errors.New("input failed")
		if _, err := day08p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Fatalf("got error %v, want input failure", err)
		}
	})
}
