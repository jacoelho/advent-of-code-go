package aoc2024

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`),
			Want: "22",
		},
	}
	aoc.AOCTest(t, day18p01(6, 12), tests)
}

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 18),
			Want:  "280",
		},
	}
	aoc.AOCTest(t, day18p01(70, 1024), tests)
}

func Test_day18p02_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`5,4
4,2
4,5
3,0
2,1
6,3
2,4
1,5
0,6
3,3
2,6
5,1
1,2
5,5
2,5
6,5
1,4
0,4
6,4
1,1
6,1
1,0
0,5
1,6
2,0`),
			Want: "6,1",
		},
	}
	aoc.AOCTest(t, day18p02(6, 12), tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 18),
			Want:  "28,56",
		},
	}
	aoc.AOCTest(t, day18p02(70, 1024), tests)
}
