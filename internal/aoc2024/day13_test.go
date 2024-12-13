package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`),
			Want: "480",
		},
		{
			Input: aoc.FileInput(t, 2024, 13),
			Want:  "28138", // too high
		},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400

Button A: X+26, Y+66
Button B: X+67, Y+21
Prize: X=12748, Y=12176

Button A: X+17, Y+86
Button B: X+84, Y+37
Prize: X=7870, Y=6450

Button A: X+69, Y+23
Button B: X+27, Y+71
Prize: X=18641, Y=10279`),
			Want: "875318608908",
		},
		{
			Input: aoc.FileInput(t, 2024, 13),
			Want:  "108394825772874",
		},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func Benchmark_day13p02(b *testing.B) {
	in := aoc.FileInput(b, 2024, 13)

	b.ResetTimer()
	for range b.N {
		v, err := day13p02(in)
		if err != nil {
			panic(err)
		}
		_ = v
	}
}
