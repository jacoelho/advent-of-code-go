package aoc2024

import (
	"io"
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
	return grid.NewGrid2D[int, rune](slices.Collect(s.Values())), s.Err()
}

func raceStartPosition(g grid.Grid2D[int, rune]) grid.Position2D[int] {
	element, ok := xmaps.Find(g, func(p grid.Position2D[int], v rune) bool { return v == 'S' })
	if !ok {
		panic("start not found")
	}
	return element.K
}

func day20(cheatDistance, save int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		raceTrack := aoc.Must(parseRaceTrack(r))
		startPosition := raceStartPosition(raceTrack)

		neighbours := func(position grid.Position2D[int]) []grid.Position2D[int] {
			var result []grid.Position2D[int]
			for p := range grid.Neighbours4(position) {
				if raceTrack[p] == '.' || raceTrack[p] == 'E' {
					result = append(result, p)
				}
			}
			return result
		}

		heuristic := func(position grid.Position2D[int]) int {
			if raceTrack[position] == 'E' {
				return 0
			}
			return 1
		}

		_, path, _ := search.AStar(startPosition, neighbours, heuristic, search.ConstantStepCost)
		distances := make(grid.Grid2D[int, int], len(path))
		for distance, p := range path {
			distances[p] = distance
		}

		savings := []struct {
			Distance      int
			CheatDistance int
		}{}
		for start := range distances {
			for rowOffset := -cheatDistance; rowOffset <= cheatDistance; rowOffset++ {
				for colOffset := xmath.Abs(rowOffset) - cheatDistance; colOffset <= cheatDistance-xmath.Abs(rowOffset); colOffset++ {
					end := start.Add(grid.Position2D[int]{X: colOffset, Y: rowOffset})
					if _, found := distances[end]; !found {
						continue
					}
					savings = append(savings, struct {
						Distance      int
						CheatDistance int
					}{
						Distance:      distances[end] - distances[start],
						CheatDistance: start.Distance(end),
					})
				}
			}
		}
		count := 0
		for _, saving := range savings {
			if saving.CheatDistance <= cheatDistance && saving.Distance-saving.CheatDistance >= save {
				count++
			}
		}
		return strconv.Itoa(count), nil
	}
}

func day20p01(cheatDistance, save int) func(r io.Reader) (string, error) {
	return day20(cheatDistance, save)
}

func day20p02(cheatDistance, save int) func(r io.Reader) (string, error) {
	return day20(cheatDistance, save)
}
