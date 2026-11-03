package aoc2021

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day11Example = `5483143223
2745854711
5264556173
6141336146
6357385478
4167524645
2176841721
6882881134
4846848554
5283751526`

func Test_day11p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "1656"},
		{Input: strings.NewReader(`0`), Want: "10"},
		{Input: aoc.FileInput(t, 2021, 11), Want: "1739"},
	}

	aoc.AOCTest(t, day11p01, tests)
}

func Test_day11p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day11Example), Want: "195"},
		{Input: strings.NewReader(`0`), Want: "10"},
		{Input: aoc.FileInput(t, 2021, 11), Want: "324"},
	}

	aoc.AOCTest(t, day11p02, tests)
}

func TestDay11NoSynchronization(t *testing.T) {
	if _, err := day11p02(strings.NewReader(`05`)); err == nil {
		t.Fatal("expected no synchronization error")
	}
}
