package aoc2025

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day08p01_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`),
			Want: "40",
		},
	}
	aoc.AOCTest(t, day08p01(10), tests)
}

func Test_day08p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2025, 8),
			Want:  "96672",
		},
	}
	aoc.AOCTest(t, day08p01(1000), tests)
}

func Test_day08p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`),
			Want: "25272",
		},
		{
			Input: aoc.FileInput(t, 2025, 8),
			Want:  "22517595",
		},
	}
	aoc.AOCTest(t, day08p02, tests)
}
