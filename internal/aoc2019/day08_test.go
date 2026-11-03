package aoc2019

import (
	"errors"
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

func Test_day08p02Errors(t *testing.T) {
	t.Run("unknown glyph", func(t *testing.T) {
		input := strings.Repeat("11111"+strings.Repeat("0", 20), 6)
		if _, err := day08p02(strings.NewReader(input)); err == nil {
			t.Fatal("expected unknown glyph error")
		}
	})
	t.Run("reader failure", func(t *testing.T) {
		failure := errors.New("input failed")
		if _, err := day08p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
			t.Fatalf("got error %v, want input failure", err)
		}
	})
}
