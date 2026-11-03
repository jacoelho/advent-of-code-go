package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func Test_day14p01_example(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`p=0,4 v=3,-3
p=6,3 v=-1,-3
p=10,3 v=-1,2
p=2,0 v=2,-1
p=0,0 v=1,3
p=3,0 v=-2,-2
p=7,6 v=-1,-3
p=3,0 v=-1,-2
p=9,3 v=2,3
p=7,3 v=-1,2
p=2,4 v=2,-3
p=9,5 v=-3,-3`),
			Want: "12",
		},
	}
	aoc.AOCTest(t, day14p01(11, 7), tests)
}

func Test_day14p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 14),
			Want:  "229868730",
		},
	}
	aoc.AOCTest(t, day14p01(101, 103), tests)
}

func Test_day14p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: aoc.FileInput(t, 2024, 14),
			Want:  "7861",
		},
	}
	aoc.AOCTest(t, day14p02, tests)
}

func TestRobotFrames(t *testing.T) {
	var robots []robot
	for y := range 3 {
		for x := range 2 {
			if x != 0 || y != 0 {
				robots = append(robots, robot{position: grid.NewPosition2D(x, y)})
			}
		}
	}
	robots = append(robots, robot{velocity: grid.NewPosition2D(1, 1)})
	got, err := uniqueRobotFrame(robots, 2, 3)
	if err != nil || got != 0 {
		t.Fatalf("unique frame = %d, %v; want 0", got, err)
	}
	for _, input := range [][]robot{{{}, {}}, {{}}, append(robots, robot{})} {
		if _, err := uniqueRobotFrame(input, 2, 3); err == nil {
			t.Fatal("expected non-unique or absent frame error")
		}
	}
}

func TestRobotInputErrors(t *testing.T) {
	for _, input := range []string{"p=0,0", "p=0,0 v=1", "p=0,x v=0,0", "p=0,0 v=9223372036854775808,0"} {
		for _, solve := range []func(io.Reader) (string, error){day14p01(101, 103), day14p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
	for _, size := range [][2]int{{0, 103}, {101, 0}, {1_000_001, 103}} {
		if _, err := day14p01(size[0], size[1])(strings.NewReader("p=0,0 v=0,0")); err == nil {
			t.Fatalf("dimensions %v: expected error", size)
		}
	}
}

func TestDay14ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day14p01(101, 103), "p=0,0 v=0,0\n"},
		{"part2", day14p02, "p=0,0 v=0,0\n"},
	} {
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
