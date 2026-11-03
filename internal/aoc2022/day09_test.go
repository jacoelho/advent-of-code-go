package aoc2022

import (
	"strings"
	"testing"

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
