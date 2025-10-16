package aoc2023

import (
	"io"
	"slices"
	"strconv"
	"unicode"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type number struct {
	value     int
	positions []grid.Position2D[int]
	neighbors collections.Set[grid.Position2D[int]]
}

func buildSymbolSet(schematic grid.Grid2D[int, rune]) collections.Set[grid.Position2D[int]] {
	symbolSet := collections.NewSet[grid.Position2D[int]]()
	for pos, ch := range schematic {
		if !unicode.IsDigit(ch) && ch != '.' {
			symbolSet.Add(pos)
		}
	}
	return symbolSet
}

func parseSchematic(r io.Reader) grid.Grid2D[int, rune] {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		return []rune(string(b)), nil
	})
	return grid.NewGrid2D[int](slices.Collect(s.Values()))
}

func findNumbers(schematic grid.Grid2D[int, rune]) func(yield func(number) bool) {
	return func(yield func(number) bool) {
		for row := 0; ; row++ {
			rowPos := grid.Position2D[int]{Y: row, X: 0}
			if _, exists := schematic[rowPos]; !exists {
				break
			}

			for col := 0; ; col++ {
				pos := grid.Position2D[int]{Y: row, X: col}
				ch, exists := schematic[pos]
				if !exists {
					break
				}
				if !unicode.IsDigit(ch) {
					continue
				}

				num := number{value: 0}

				for {
					pos := grid.Position2D[int]{Y: row, X: col}
					ch, exists := schematic[pos]
					if !exists || !unicode.IsDigit(ch) {
						break
					}
					num.value = num.value*10 + int(ch-'0')
					num.positions = append(num.positions, pos)
					col++
				}
				col--

				num.neighbors = collections.NewSet[grid.Position2D[int]]()
				for _, pos := range num.positions {
					num.neighbors.Add(slices.Collect(grid.Neighbours8(pos))...)
				}

				if !yield(num) {
					return
				}
			}
		}
	}
}

func day03p01(r io.Reader) (string, error) {
	schematic := parseSchematic(r)
	symbolSet := buildSymbolSet(schematic)

	sum := xiter.Sum(xiter.Map(func(num number) int {
		if !num.neighbors.Intersect(symbolSet).IsEmpty() {
			return num.value
		}
		return 0
	}, findNumbers(schematic)))

	return strconv.Itoa(sum), nil
}

func day03p02(r io.Reader) (string, error) {
	schematic := parseSchematic(r)
	numbers := slices.Collect(findNumbers(schematic))

	sum := 0
	for pos, ch := range schematic {
		if ch == '*' {
			adjacentNumbers := xslices.Filter(func(num number) bool {
				return num.neighbors.Contains(pos)
			}, numbers)

			if len(adjacentNumbers) == 2 {
				sum += xslices.Product(xslices.Map(func(n number) int { return n.value }, adjacentNumbers))
			}
		}
	}

	return strconv.Itoa(sum), nil
}
