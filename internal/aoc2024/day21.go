package aoc2024

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/search"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
	"github.com/jacoelho/advent-of-code-go/internal/xstrings"
)

func parseCodes(r io.Reader) ([]string, error) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) {
		return string(b), nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func numericKeyPad() grid.Grid2D[int, rune] {
	/*
		+---+---+---+
		| 7 | 8 | 9 |
		+---+---+---+
		| 4 | 5 | 6 |
		+---+---+---+
		| 1 | 2 | 3 |
		+---+---+---+
		|   | 0 | A |
		+---+---+---+
	*/
	return grid.Grid2D[int, rune]{
		{0, 0}: '7', {1, 0}: '8', {2, 0}: '9',
		{0, 1}: '4', {1, 1}: '5', {2, 1}: '6',
		{0, 2}: '1', {1, 2}: '2', {2, 2}: '3',
		/*     empty    */ {1, 3}: '0', {2, 3}: 'A',
	}
}

func directionalKeyPad() grid.Grid2D[int, rune] {
	/*
		+---+---+---+
		|	| ^ | A |
		+---+---+---+
		| < | v | > |
		+---+---+---+
	*/
	return grid.Grid2D[int, rune]{
		/*     empty    */ {1, 0}: '^', {2, 0}: 'A',
		{0, 1}: '<', {1, 1}: 'v', {2, 1}: '>',
	}
}

func offsetToMovement(offset grid.Position2D[int]) rune {
	switch offset {
	case grid.Position2D[int]{0, -1}:
		return '^'
	case grid.Position2D[int]{0, 1}:
		return 'v'
	case grid.Position2D[int]{-1, 0}:
		return '<'
	case grid.Position2D[int]{1, 0}:
		return '>'
	default:
		panic("invalid offset")
	}
}

func shortestPathsBetweenKeys(g grid.Grid2D[int, rune]) map[string][]string {
	neighbours := func(position grid.Position2D[int]) []grid.Position2D[int] {
		return slices.Collect(xiter.Filter(g.Contains, grid.Neighbours4(position)))
	}
	targetValue := func(target grid.Position2D[int]) func(position grid.Position2D[int]) int {
		return func(position grid.Position2D[int]) int {
			if position == target {
				return 0
			}
			return 1
		}
	}
	result := make(map[string][]string)
	for pair := range xslices.PairwiseSelf(slices.Collect(maps.Keys(g))) {
		_, paths, found := search.AStarBag(pair.V1, neighbours, targetValue(pair.V2), search.ConstantStepCost)
		if !found {
			panic("no paths found")
		}
		pathsRunes := make([]string, len(paths))
		for i, path := range paths {
			var runes []rune
			for i := 0; i < len(path)-1; i++ {
				offset := path[i+1].Sub(path[i])
				runes = append(runes, offsetToMovement(offset))
			}
			pathsRunes[i] = string(runes)
		}
		result[string([]rune{g[pair.V1], g[pair.V2]})] = pathsRunes
	}
	return result
}

func findMoves(code string, depth int) int {
	numericKeyPadPaths := shortestPathsBetweenKeys(numericKeyPad())
	directionalKeyPadPaths := shortestPathsBetweenKeys(directionalKeyPad())

	type pair struct {
		V1 string
		V2 int
	}

	cache := make(map[pair]int)

	var findMovesRec func(transitions map[string][]string, code string, depth int) int
	findMovesRec = func(transitions map[string][]string, code string, depth int) int {
		key := pair{code, depth}
		if v, found := cache[key]; found {
			return v
		}

		var total int
		for _, move := range xstrings.Pairs(code) {
			paths, found := transitions[move]
			if !found {
				panic("no paths found" + move)
			}
			if depth == 0 {
				v := slices.Min(xslices.Map(func(path string) int { return len(path) }, paths))
				total += v
				continue
			}
			vv := xslices.Map(func(path string) int {
				return findMovesRec(directionalKeyPadPaths, path, depth-1)
			}, xslices.Filter(func(e string) bool {
				return e != ""
			}, paths))
			if len(vv) > 0 {
				v := slices.Min(vv)
				total += v
			}

		}

		cache[key] = total
		return total
	}

	return findMovesRec(numericKeyPadPaths, "A"+code, depth)
}

func day21p01(r io.Reader) (string, error) {
	codes := aoc.Must(parseCodes(r))

	var total int
	for _, code := range codes {
		codeInt := convert.ExtractDigits[int](code)[0]
		fmt.Println(findMoves(code, 2))
		total += findMoves(code, 2) * codeInt
	}
	return strconv.Itoa(total), nil
}

func day21p02(r io.Reader) (string, error) {
	return "", nil
}
