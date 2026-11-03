package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day05Example = `0,9 -> 5,9
8,0 -> 0,8
9,4 -> 3,4
2,2 -> 2,1
7,0 -> 7,4
6,4 -> 2,0
0,9 -> 2,9
3,4 -> 1,4
0,0 -> 8,8
5,5 -> 8,2`

func Test_day05p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day05Example), Want: "5"},
		{Input: strings.NewReader(`0,0 -> 0,0
0,0 -> 0,0`), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 5), Want: "6841"},
	}

	aoc.AOCTest(t, day05p01, tests)
}

func Test_day05p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day05Example), Want: "12"},
		{Input: strings.NewReader(`0,0 -> 0,0
0,0 -> 0,0`), Want: "1"},
		{Input: aoc.FileInput(t, 2021, 5), Want: "19258"},
	}

	aoc.AOCTest(t, day05p02, tests)
}

func TestDay05UnsupportedSlope(t *testing.T) {
	if _, err := day05p02(strings.NewReader(`0,0 -> 2,1`)); err == nil {
		t.Fatal("expected unsupported slope error")
	}
}

func TestDay05MalformedEndpoint(t *testing.T) {
	if _, err := day05p02(strings.NewReader("0,0 -> 1,1junk")); err == nil {
		t.Fatal("expected malformed endpoint error")
	}
}

func TestDay05RasterLimit(t *testing.T) {
	if _, err := day05p02(strings.NewReader("-9223372036854775808,0 -> 9223372036854775807,0")); err == nil {
		t.Fatal("expected raster limit error")
	}
}
