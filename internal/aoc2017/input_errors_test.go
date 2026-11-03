package aoc2017

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func TestDaysPropagateReadErrors(t *testing.T) {
	failure := errors.New("input read failed")
	tests := []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"01/1", day01p01},
		{"01/2", day01p02},
		{"02/1", day02p01},
		{"02/2", day02p02},
		{"03/1", day03p01},
		{"03/2", day03p02},
		{"04/1", day04p01},
		{"04/2", day04p02},
		{"05/1", day05p01},
		{"05/2", day05p02},
		{"06/1", day06p01},
		{"06/2", day06p02},
		{"07/1", day07p01},
		{"07/2", day07p02},
		{"08/1", day08p01},
		{"08/2", day08p02},
		{"09/1", day09p01},
		{"09/2", day09p02},
		{"10/1", day10p01},
		{"10/2", day10p02},
		{"11/1", day11p01},
		{"11/2", day11p02},
		{"12/1", day12p01},
		{"12/2", day12p02},
		{"13/1", day13p01},
		{"13/2", day13p02},
		{"14/1", day14p01},
		{"14/2", day14p02},
		{"15/1", day15p01},
		{"15/2", day15p02},
		{"16/1", day16p01},
		{"16/2", day16p02},
		{"17/1", day17p01},
		{"17/2", day17p02},
		{"18/1", day18p01},
		{"18/2", day18p02},
		{"19/1", day19p01},
		{"19/2", day19p02},
		{"20/1", day20p01},
		{"20/2", day20p02},
		{"21/1", day21p01},
		{"21/2", day21p02},
		{"22/1", day22p01},
		{"22/2", day22p02},
		{"23/1", day23p01},
		{"23/2", day23p02},
		{"24/1", day24p01},
		{"24/2", day24p02},
		{"25/1", day25p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}

func Test2017Days01To12MalformedInput(t *testing.T) {
	tests := []struct {
		name, input string
		solve       func(io.Reader) (string, error)
	}{
		{"captcha digit", "1x", day01p01},
		{"odd captcha", "123", day01p02},
		{"spreadsheet zero", "5 0", day02p02},
		{"missing pair", "5 7", day02p02},
		{"ambiguous pairs", "2 4 8", day02p02},
		{"spiral target", "0", day03p01},
		{"jump offset", "not a number", day05p01},
		{"memory bank", "1 -1", day06p01},
		{"tower cycle", "a (1) -> b\nb (1) -> a", day07p01},
		{"tower ambiguity", "root (0) -> a, b\na (1)\nb (2)", day07p02},
		{"single-child ambiguity", "root (0) -> left, right, bad\nleft (10)\nright (10)\nbad (2) -> leaf\nleaf (7)", day07p02},
		{"register condition", "a inc 1 if b ?? 0", day08p01},
		{"unfinished garbage", "{<abc!", day09p01},
		{"group close", "}", day09p01},
		{"knot length", "257", day10p01},
		{"hex direction", "ne,east", day11p01},
		{"program ID", "zero <-> 0", day12p01},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.solve(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
