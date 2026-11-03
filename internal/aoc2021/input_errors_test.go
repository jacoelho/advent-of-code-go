package aoc2021

import (
	"errors"
	"io"
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
