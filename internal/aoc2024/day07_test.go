package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day07p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`),
			Want: "3749",
		},
		{
			Input: aoc.FileInput(t, 2024, 7),
			Want:  "945512582195",
		},
	}

	aoc.AOCTest(t, day07p01, tests)
}

func Test_day07p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`190: 10 19
3267: 81 40 27
83: 17 5
156: 15 6
7290: 6 8 6 15
161011: 16 10 13
192: 17 8 14
21037: 9 7 18 13
292: 11 6 16 20`),
			Want: "11387",
		},
		{
			Input: aoc.FileInput(t, 2024, 7),
			Want:  "271691107779347",
		},
	}

	aoc.AOCTest(t, day07p02, tests)
}

func Benchmark_day07(b *testing.B) {
	in := aoc.FileInput(b, 2024, 7)

	b.ResetTimer()
	for range b.N {
		v, err := day07p02(in)
		if err != nil {
			panic(err)
		}
		_ = v
	}
}
