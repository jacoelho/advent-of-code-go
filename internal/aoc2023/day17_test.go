package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day17p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`),
			Want: "102",
		},
		{
			Input: aoc.FileInput(t, 2023, 17),
			Want:  "755",
		},
	}
	aoc.AOCTest(t, day17p01, tests)
}

func Test_day17p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`2413432311323
3215453535623
3255245654254
3446585845452
4546657867536
1438598798454
4457876987766
3637877979653
4654967986887
4564679986453
1224686865563
2546548887735
4322674655533`),
			Want: "94",
		},
		{
			Input: strings.NewReader(`111111111111
999999999991
999999999991
999999999991
999999999991`),
			Want: "71",
		},
		{
			Input: aoc.FileInput(t, 2023, 17),
			Want:  "881",
		},
	}
	aoc.AOCTest(t, day17p02, tests)
}
