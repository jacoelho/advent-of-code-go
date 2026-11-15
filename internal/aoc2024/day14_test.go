package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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

func TestBathroomTreeIsFrameWhereRobotsConverge(t *testing.T) {
	const width, height, tree = 101, 103, 7777
	var input strings.Builder
	for v := 1; v <= 20; v++ {
		x := xmath.Modulo(50-v*tree, width)
		y := xmath.Modulo(40-2*v*tree, height)
		fmt.Fprintf(&input, "p=%d,%d v=%d,%d\n", x, y, v, 2*v)
	}
	got, err := day14p02(strings.NewReader(input.String()))
	if err != nil || got != strconv.Itoa(tree) {
		t.Fatalf("got %q, %v; want %d", got, err, tree)
	}
}

func TestRobotInputErrors(t *testing.T) {
	for _, input := range []string{
		"",
		"p=0,0",
		"p=0,0 v=1",
		"p=0,x v=0,0",
		"p=0,0 v=9223372036854775808,0",
		"p=-1,0 v=0,0",
		"p=101,0 v=0,0",
		"p=0,103 v=0,0",
	} {
		for _, solve := range []func(io.Reader) (string, error){day14p01(101, 103), day14p02} {
			if _, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: expected error", input)
			}
		}
	}
	for _, size := range [][2]int{{0, 103}, {101, 0}} {
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
