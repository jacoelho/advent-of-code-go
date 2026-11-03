package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"
)

func Test2024ReaderErrors(t *testing.T) {
	tests := []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"01 part 1", day01p01, "1 2\n"},
		{"01 part 2", day01p02, "1 2\n"},
		{"02 part 1", day02p01, "1 2 3\n"},
		{"02 part 2", day02p02, "1 2 3\n"},
		{"03 part 1", day03p01, "mul(2,3)"},
		{"03 part 2", day03p02, "mul(2,3)"},
		{"04 part 1", day04p01, "XMAS\n"},
		{"04 part 2", day04p02, "XMAS\n"},
		{"05 part 1", day05p01, "1|2\n\n1,2,3\n"},
		{"05 part 2", day05p02, "1|2\n\n1,2,3\n"},
		{"06 part 1", day06p01, "^\n"},
		{"06 part 2", day06p02, "^\n"},
		{"07 part 1", day07p01, "10: 5 5\n"},
		{"07 part 2", day07p02, "10: 5 5\n"},
		{"08 part 1", day08p01, "A.A\n"},
		{"08 part 2", day08p02, "A.A\n"},
		{"09 part 1", day09p01, "12345\n"},
		{"09 part 2", day09p02, "12345\n"},
		{"10 part 1", day10p01, "0123456789\n"},
		{"10 part 2", day10p02, "0123456789\n"},
		{"11 part 1", day11p01, "125 17\n"},
		{"11 part 2", day11p02, "125 17\n"},
		{"12 part 1", day12p01, "A\n"},
		{"12 part 2", day12p02, "A\n"},
		{"13 part 1", day13p01, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=3\n\n"},
		{"13 part 2", day13p02, "Button A: X+1, Y+0\nButton B: X+0, Y+1\nPrize: X=3, Y=3\n\n"},
		{"14 part 1", day14p01(101, 103), "p=0,0 v=0,0\n"},
		{"14 part 2", day14p02, "p=0,0 v=0,0\n"},
		{"15 part 1", day15p01, "###\n#@#\n###\n\n^\n"},
		{"15 part 2", day15p02, "###\n#@#\n###\n\n^\n"},
		{"16 part 1", day16p01, "SE\n"},
		{"16 part 2", day16p02, "SE\n"},
		{"17 part 1", day17p01, "Register A: 0\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n"},
		{"17 part 2", day17p02, "Register A: 0\nRegister B: 0\nRegister C: 0\n\nProgram: 0,3,5,4,3,0\n"},
		{"18 part 1", day18p01(1, 0), "0,0\n"},
		{"18 part 2", day18p02(1, 0), "0,0\n"},
		{"19 part 1", day19p01, "r, b\n\nr\n"},
		{"19 part 2", day19p02, "r, b\n\nr\n"},
		{"20 part 1", day20p01(1), "SE\n"},
		{"20 part 2", day20p02(1), "SE\n"},
		{"21 part 1", day21p01, "029A\n"},
		{"21 part 2", day21p02, "029A\n"},
		{"22 part 1", day22p01, "1\n"},
		{"22 part 2", day22p02, "1\n"},
		{"23 part 1", day23p01, "ta-tb\n"},
		{"23 part 2", day23p02, "ta-tb\n"},
		{"24 part 1", day24p01, "x00: 0\ny00: 0\n\nx00 XOR y00 -> z00\nx00 AND y00 -> z01\n"},
		{"24 part 2", day24p02, "x00: 0\ny00: 0\n\nx00 XOR y00 -> z00\nx00 AND y00 -> z01\n"},
		{"25 part 1", day25p01, "#####\n.####\n.####\n.####\n.#.#.\n.#...\n.....\n"},
	}
	want := errors.New("input read failure")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
