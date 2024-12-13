package aoc2024

import (
	"bufio"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

type machine struct {
	A     grid.Position2D[int]
	B     grid.Position2D[int]
	Prize grid.Position2D[int]
}

func parseMachine(text string) machine {
	digits := convert.ExtractDigits[int](text)
	return machine{
		A:     grid.Position2D[int]{X: digits[0], Y: digits[1]},
		B:     grid.Position2D[int]{X: digits[2], Y: digits[3]},
		Prize: grid.Position2D[int]{X: digits[4], Y: digits[5]},
	}
}

func parseArcadeLayout(r io.Reader) ([]machine, error) {
	var result []machine
	s := bufio.NewScanner(r)
	s.Split(scanner.SplitBySeparator([]byte{'\n', '\n'}))

	for s.Scan() {
		result = append(result, parseMachine(s.Text()))
	}
	return result, s.Err()
}

func calculatePresses(m machine) int {
	determinantB := m.B.Y*m.Prize.X - m.B.X*m.Prize.Y
	determinantA := -m.A.Y*m.Prize.X + m.A.X*m.Prize.Y
	determinantAB := (m.A.X * m.B.Y) - (m.A.Y * m.B.X)

	// determinantAB == 0  are linearly dependent
	// determinantB%determinantAB != 0 we want integer solutions
	// determinantA%determinantAB != 0 we want integer solutions
	if determinantAB == 0 || determinantB%determinantAB != 0 || determinantA%determinantAB != 0 {
		return 0
	}
	return 3*determinantB/determinantAB + determinantA/determinantAB
}

func day13p01(r io.Reader) (string, error) {
	layout := aoc.Must(parseArcadeLayout(r))

	total := xslices.Reduce(func(sum int, m machine) int {
		return sum + calculatePresses(m)
	}, 0, layout)

	return strconv.Itoa(total), nil
}

func day13p02(r io.Reader) (string, error) {
	layout := xslices.Map(func(m machine) machine {
		return machine{
			A: m.A,
			B: m.B,
			Prize: m.Prize.Add(grid.Position2D[int]{
				X: 10000000000000,
				Y: 10000000000000,
			}),
		}
	}, aoc.Must(parseArcadeLayout(r)))

	total := xslices.Reduce(func(sum int, m machine) int {
		return sum + calculatePresses(m)
	}, 0, layout)

	return strconv.Itoa(total), nil
}
