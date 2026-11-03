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

func parseSchematic(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		return []rune(string(b)), nil
	})
	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	return grid.NewGrid2D[int](rows), nil
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
					for neighbor := range grid.Neighbours8(pos) {
						num.neighbors.Add(neighbor)
					}
				}

				if !yield(num) {
					return
				}
			}
		}
	}
}

func day03p01(r io.Reader) (string, error) {
	schematic, err := parseSchematic(r)
	if err != nil {
		return "", err
	}
	symbolSet := buildSymbolSet(schematic)

	sum := xiter.From(findNumbers(schematic)).SumBy(func(num number) int {
		if xiter.From(num.neighbors.Iter()).Any(symbolSet.Contains) {
			return num.value
		}
		return 0
	})

	return strconv.Itoa(sum), nil
}

func day03p02(r io.Reader) (string, error) {
	schematic, err := parseSchematic(r)
	if err != nil {
		return "", err
	}
	numbers := slices.Collect(findNumbers(schematic))

	sum := 0
	for pos, ch := range schematic {
		if ch == '*' {
			type gearRatio struct {
				count, product int
			}
			adjacent := xiter.FromSlice(numbers).Filter(func(num number) bool {
				return num.neighbors.Contains(pos)
			}).Take(3).Reduce(gearRatio{product: 1}, func(ratio gearRatio, num number) gearRatio {
				return gearRatio{count: ratio.count + 1, product: ratio.product * num.value}
			})
			if adjacent.count == 2 {
				sum += adjacent.product
			}
		}
	}

	return strconv.Itoa(sum), nil
}
