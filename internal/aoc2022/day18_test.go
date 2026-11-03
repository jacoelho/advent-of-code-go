package aoc2022

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(lavaExample), Want: "64"},
		{Input: strings.NewReader("0,0,0\n1,0,0"), Want: "10"},
		{Input: aoc.FileInput(t, 2022, 18), Want: "4340"},
	}

	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(lavaExample), Want: "58"},
		{Input: strings.NewReader("0,0,0\n1,0,0"), Want: "10"},
		{Input: aoc.FileInput(t, 2022, 18), Want: "2468"},
	}

	aoc.AOCTest(t, day18p02, tests)
}

func TestLavaEnclosedAir(t *testing.T) {
	for _, offset := range []int{0, -10} {
		var input strings.Builder
		for x := range 3 {
			for y := range 3 {
				for z := range 3 {
					if x != 1 || y != 1 || z != 1 {
						fmt.Fprintf(&input, "%d,%d,%d\n", x+offset, y+offset, z+offset)
					}
				}
			}
		}
		aoc.AOCTest(t, day18p01, []aoc.TestInput{{Input: strings.NewReader(input.String()), Want: "60"}})
		aoc.AOCTest(t, day18p02, []aoc.TestInput{{Input: strings.NewReader(input.String()), Want: "54"}})
	}
}

func TestLavaErrors(t *testing.T) {
	for _, input := range []string{"", "1,2", "1,2,x", "-2147483648,0,0\n2147483647,0,0"} {
		if _, err := day18p02(strings.NewReader(input)); err == nil {
			t.Errorf("expected lava input or bounds error for %q", input)
		}
	}
}

const lavaExample = `2,2,2
1,2,2
3,2,2
2,1,2
2,3,2
2,2,1
2,2,3
2,2,4
2,2,6
1,2,5
3,2,5
2,1,5
2,3,5`

func TestDay18ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day18p01},
		{"part2", day18p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
