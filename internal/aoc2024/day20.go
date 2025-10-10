package aoc2024

import (
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
)

func parseRaceTrack(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		result := make([]rune, 0, len(b))
		for _, v := range b {
			result = append(result, rune(v))
		}
		return result, nil
	})
	return grid.NewGrid2D[int](slices.Collect(s.Values())), s.Err()
}

func raceStartPosition(g grid.Grid2D[int, rune]) grid.Position2D[int] {
	element, ok := xmaps.Find(g, func(p grid.Position2D[int], v rune) bool { return v == 'S' })
	if !ok {
		panic("start not found")
	}
	return element.K
}

func day20Heuristic(raceTrack grid.Grid2D[int, rune]) func(position grid.Position2D[int]) int {
	return func(position grid.Position2D[int]) int {
		if raceTrack[position] == 'E' {
			return 0
		}
		return 1
	}
}

func day20Neighbours(raceTrack grid.Grid2D[int, rune]) func(position grid.Position2D[int]) []grid.Position2D[int] {
	return func(position grid.Position2D[int]) []grid.Position2D[int] {
		var result []grid.Position2D[int]
		for p := range grid.Neighbours4(position) {
			if raceTrack[p] == '.' || raceTrack[p] == 'E' {
				result = append(result, p)
			}
		}
		return result
	}
}

func cheatDistanceOffsets(cheatDistance int) iter.Seq[grid.Position2D[int]] {
	return func(yield func(grid.Position2D[int]) bool) {
		for dy := -cheatDistance; dy <= cheatDistance; dy++ {
			for dx := -cheatDistance + xmath.Abs(dy); dx <= cheatDistance-xmath.Abs(dy); dx++ {
				if !yield(grid.Position2D[int]{X: dx, Y: dy}) {
					return
				}
			}
		}
	}
}

func day20(cheatDistance, saveGoal int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		raceTrack := aoc.Must(parseRaceTrack(r))
		startPosition := raceStartPosition(raceTrack)
		neighbours := day20Neighbours(raceTrack)
		heuristic := day20Heuristic(raceTrack)

		_, path, _ := search.AStar(startPosition, neighbours, heuristic, search.ConstantStepCost)

		distances := make(grid.Grid2D[int, int], len(path))
		for distance, p := range path {
			distances[p] = distance
		}

		offsets := slices.Collect(cheatDistanceOffsets(cheatDistance))

		var count int
		for start, startDistance := range distances {
			for _, offset := range offsets {
				end := start.Add(offset)
				if endDistance, found := distances[end]; found {
					if endDistance-startDistance-start.Distance(end) >= saveGoal {
						count++
					}
				}
			}
		}
		return strconv.Itoa(count), nil
	}
}

func day20p01(saveGoal int) func(r io.Reader) (string, error) {
	return day20(2, saveGoal)
}

func day20p02(saveGoal int) func(r io.Reader) (string, error) {
	return day20(20, saveGoal)
}
