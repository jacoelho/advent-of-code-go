package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day22p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1
10
100
2024`),
			Want: "37327623",
		},
		{
			Input: aoc.FileInput(t, 2024, 22),
			Want:  "12759339434",
		},
	}
	aoc.AOCTest(t, day22p01, tests)
}

func Test_day22p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`1
2
3
2024`),
			Want: "23",
		},
		{
			Input: aoc.FileInput(t, 2024, 22),
			Want:  "1405",
		},
	}
	aoc.AOCTest(t, day22p02, tests)
}
