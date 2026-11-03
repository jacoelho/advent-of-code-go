package aoc2022

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day10Example = `addx 15
addx -11
addx 6
addx -3
addx 5
addx -1
addx -8
addx 13
addx 4
noop
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx 5
addx -1
addx -35
addx 1
addx 24
addx -19
addx 1
addx 16
addx -11
noop
noop
addx 21
addx -15
noop
noop
addx -3
addx 9
addx 1
addx -3
addx 8
addx 1
addx 5
noop
noop
noop
noop
noop
addx -36
noop
addx 1
addx 7
noop
noop
noop
addx 2
addx 6
noop
noop
noop
noop
noop
addx 1
noop
noop
addx 7
addx 1
noop
addx -13
addx 13
addx 7
noop
addx 1
addx -33
noop
noop
noop
addx 2
noop
noop
noop
addx 8
noop
addx -1
addx 2
addx 1
noop
addx 17
addx -9
addx 1
addx 1
addx -3
addx 11
noop
noop
addx 1
noop
addx 1
noop
noop
addx -13
addx -19
addx 1
addx 3
addx 26
addx -30
addx 12
addx -1
addx 3
addx 1
noop
noop
noop
addx -9
addx 18
addx 1
addx 2
noop
noop
addx 9
noop
noop
noop
addx -1
addx 2
addx -37
addx 1
addx 3
noop
addx 15
addx -21
addx 22
addx -6
addx 1
noop
addx 2
addx 1
noop
addx -10
noop
noop
addx 20
addx 1
addx 2
addx 2
addx -6
addx -11
noop
noop
noop`

func Test_day10p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day10Example), Want: "13140"},
		{Input: aoc.FileInput(t, 2022, 10), Want: "14240"},
	}
	aoc.AOCTest(t, day10p01, tests)
}

func Test_day10p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: aoc.FileInput(t, 2022, 10), Want: "PLULKBZH"},
	}
	aoc.AOCTest(t, day10p02, tests)
}

func TestRegisterCycles(t *testing.T) {
	values, cycles, err := registerCycles(strings.NewReader("noop\naddx 3\naddx -5"))
	if err != nil {
		t.Fatal(err)
	}
	want := []int64{1, 1, 1, 4, 4}
	if cycles != len(want) {
		t.Fatalf("got %d cycles, want %d", cycles, len(want))
	}
	for i, value := range want {
		if values[i] != value {
			t.Errorf("cycle %d got %d, want %d", i+1, values[i], value)
		}
	}
}

func TestCRTDisplay(t *testing.T) {
	rows, err := crtDisplay(strings.NewReader(day10Example))
	if err != nil {
		t.Fatal(err)
	}
	const want = `##..##..##..##..##..##..##..##..##..##..
###...###...###...###...###...###...###.
####....####....####....####....####....
#####.....#####.....#####.....#####.....
######......######......######......####
#######.......#######.......#######.....`
	if got := strings.Join(rows, "\n"); got != want {
		t.Errorf("display:\n%s\nwant:\n%s", got, want)
	}
	if _, err := day10p01(strings.NewReader("noop")); err == nil {
		t.Fatal("expected insufficient signal cycles error")
	}
	if _, err := day10p02(strings.NewReader("noop")); err == nil {
		t.Fatal("expected insufficient CRT cycles error")
	}
}
