package aoc2022

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const day09Example = `R 4
U 4
L 3
D 1
R 4
D 1
L 5
R 2`

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day09Example), Want: "13"},
		{Input: aoc.FileInput(t, 2022, 9), Want: "6269"},
	}
	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(day09Example), Want: "1"},
		{Input: aoc.FileInput(t, 2022, 9), Want: "2557"},
	}
	aoc.AOCTest(t, day09p02, tests)
}

func TestRopeOrigin(t *testing.T) {
	aoc.AOCTest(t, day09p01, []aoc.TestInput{{Input: strings.NewReader("R 0"), Want: "1"}})
	aoc.AOCTest(t, day09p02, []aoc.TestInput{{Input: strings.NewReader("R 0"), Want: "1"}, {Input: strings.NewReader("R 5\nU 8\nL 8\nD 3\nR 17\nD 10\nL 25\nU 20"), Want: "36"}})
}

func TestDay09ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name  string
		solve func(io.Reader) (string, error)
	}{
		{"part1", day09p01},
		{"part2", day09p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.solve(iotest.ErrReader(failure))
			if !errors.Is(err, failure) {
				t.Fatalf("got %v; want input read failure", err)
			}
		})
	}
}
