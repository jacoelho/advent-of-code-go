package aoc2024

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

func parseCalibrationEquations(r io.Reader) (map[int][]int, error) {
	equations := make(map[int][]int)
	s := bufio.NewScanner(r)
	for s.Scan() {
		text := s.Text()

		values := strings.Split(text, ":")
		testValue := aoc.MustAtoi(values[0])

		equations[testValue] = xslices.Map(func(in string) int {
			return aoc.MustAtoi(in)
		}, strings.Fields(values[1]))
	}
	return equations, s.Err()
}

func day07Heuristic(equation []int) func(in [2]int) int {
	return func(in [2]int) int {
		if in[0] == 0 && equation[in[0]] == in[1] {
			return 0
		}
		return 1
	}
}

func day07Neighbours(equation []int) func(state [2]int) [][2]int {
	return func(state [2]int) [][2]int {
		idx, t := state[0], state[1]
		if idx < 0 {
			return nil
		}

		v := equation[idx]
		var n [][2]int

		if t%v == 0 {
			n = append(n, [2]int{idx - 1, t / v})
		}
		if t > v {
			n = append(n, [2]int{idx - 1, t - v})
		}

		return n
	}
}

func day07p01(r io.Reader) (string, error) {
	equations := aoc.Must(parseCalibrationEquations(r))

	var total int
	for target, equation := range equations {
		start := [2]int{len(equation) - 1, target}

		_, _, found := search.AStar(
			start,
			day07Neighbours(equation),
			day07Heuristic(equation),
			search.ConstantStepCost,
		)
		if found {
			total += target
		}
	}
	return strconv.Itoa(total), nil
}

func day07p02(r io.Reader) (string, error) {
	equations := aoc.Must(parseCalibrationEquations(r))

	var total int
	for target, equation := range equations {
		initialNeighbours := day07Neighbours(equation)

		neighbours := func(state [2]int) [][2]int {
			n := initialNeighbours(state)
			if n == nil {
				return nil
			}

			v, t := equation[state[0]], state[1]

			if newT, found := strings.CutSuffix(strconv.Itoa(t), strconv.Itoa(v)); found {
				if newT, err := strconv.Atoi(newT); err == nil {
					n = append(n, [2]int{state[0] - 1, newT})
				}
			}

			return n
		}

		start := [2]int{len(equation) - 1, target}

		_, _, found := search.AStar(
			start,
			neighbours,
			day07Heuristic(equation),
			search.ConstantStepCost,
		)
		if found {
			total += target
		}
	}
	return strconv.Itoa(total), nil
}