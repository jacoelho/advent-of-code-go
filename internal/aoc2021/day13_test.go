package aoc2021

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day13Example = `6,10
0,14
9,10
0,3
10,4
4,11
6,0
6,12
4,1
0,13
10,12
3,4
3,0
8,4
1,10
2,14
8,10
9,0

fold along y=7
fold along x=5`

func Test_day13p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day13Example), Want: "17"},
		{Input: aoc.FileInput(t, 2021, 13), Want: "695"},
	}

	aoc.AOCTest(t, day13p01, tests)
}

func Test_day13p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2021, 13), Want: "GJZGLUPJ"},
	}

	aoc.AOCTest(t, day13p02, tests)
}

func TestDay13TrailingCoordinateText(t *testing.T) {
	if _, err := day13p01(strings.NewReader(`1,2x

fold along x=3`)); err == nil {
		t.Fatal("expected trailing coordinate text error")
	}
}

func TestDay13CoordinateBudget(t *testing.T) {
	if _, err := day13p01(strings.NewReader(`9223372036854775807,0

fold along x=1`)); err == nil {
		t.Fatal("expected coordinate budget error")
	}
}

func TestFoldedPaperExample(t *testing.T) {
	points, folds, err := transparentPaper(strings.NewReader(day13Example))
	if err != nil {
		t.Fatal(err)
	}
	for _, fold := range folds {
		points = foldPaper(points, fold)
	}
	rows, err := renderPaper(points)
	if err != nil {
		t.Fatal(err)
	}
	want := "#####\n#...#\n#...#\n#...#\n#####"
	if got := strings.Join(rows, "\n"); got != want {
		t.Errorf("got display:\n%s\nwant:\n%s", got, want)
	}
}

func TestDay13FoldedLetter(t *testing.T) {
	rows := []string{".##.", "#..#", "#..#", "####", "#..#", "#..#"}
	var input strings.Builder
	for y, row := range rows {
		for x, pixel := range row {
			if pixel == '#' {
				fmt.Fprintf(&input, "%d,%d\n", x, y)
			}
		}
	}
	input.WriteString("\nfold along x=10\n")
	got, err := day13p02(strings.NewReader(input.String()))
	if err != nil {
		t.Fatal(err)
	}
	if got != "A" {
		t.Errorf("got %q, want A", got)
	}
}
