package aoc2024

import (
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
	"github.com/jacoelho/advent-of-code-go/internal/xstrings"
)

var numericKeyPad = map[rune]grid.Position2D[int]{
	'7': {0, 0}, '8': {1, 0}, '9': {2, 0},
	'4': {0, 1}, '5': {1, 1}, '6': {2, 1},
	'1': {0, 2}, '2': {1, 2}, '3': {2, 2},
	' ': {0, 3}, '0': {1, 3}, 'A': {2, 3},
}

var directionalKeyPad = map[rune]grid.Position2D[int]{
	' ': {0, 0}, '^': {1, 0}, 'A': {2, 0},
	'<': {0, 1}, 'v': {1, 1}, '>': {2, 1},
}

func parseCodes(r io.Reader) ([]string, error) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) {
		return string(b), nil
	})
	return slices.Collect(s.Values()), s.Err()
}

type robotCacheItem struct {
	robot int
	from  rune
	to    rune
}

func sequencePressCost(cache map[robotCacheItem]int, robot int, sequence string) int {
	return xslices.Reduce(func(total int, pair []rune) int {
		if robot == 0 {
			return total + 1
		}
		k := robotCacheItem{robot: robot, from: pair[0], to: pair[1]}
		return total + cache[k]
	}, 0, xstrings.Pairs("A"+sequence))
}

func fillRobotBestPath(
	cache map[robotCacheItem]int,
	robot int,
	keyPad map[rune]grid.Position2D[int],
) {
	emptyPosition := keyPad[' ']
	for start, startPos := range keyPad {
		for end, endPos := range keyPad {
			distanceX := xmath.Abs(endPos.X - startPos.X)
			distanceY := xmath.Abs(endPos.Y - startPos.Y)

			horizontalKey := ">"
			if endPos.X < startPos.X {
				horizontalKey = "<"
			}
			verticalKey := "v"
			if endPos.Y < startPos.Y {
				verticalKey = "^"
			}

			minHorizontalCost := math.MaxInt
			if emptyPosition != (grid.Position2D[int]{X: endPos.X, Y: startPos.Y}) {
				minHorizontalCost = sequencePressCost(
					cache,
					robot-1,
					buildPressSequence(horizontalKey, verticalKey, distanceX, distanceY),
				)
			}

			minVerticalCost := math.MaxInt
			if emptyPosition != (grid.Position2D[int]{X: startPos.X, Y: endPos.Y}) {
				minVerticalCost = sequencePressCost(
					cache,
					robot-1,
					buildPressSequence(verticalKey, horizontalKey, distanceY, distanceX),
				)
			}

			cache[robotCacheItem{robot: robot, from: start, to: end}] = min(minHorizontalCost, minVerticalCost)
		}
	}
}

func buildPressSequence(firstPress, secondPress string, firstDistance, secondDistance int) string {
	var sb strings.Builder

	size := firstDistance + secondDistance + 1

	sb.Grow(size)
	for i := 0; i < firstDistance; i++ {
		sb.WriteString(firstPress)
	}

	for i := 0; i < secondDistance; i++ {
		sb.WriteString(secondPress)
	}

	sb.WriteString("A")

	return sb.String()
}

func robotMovementCache(count int) map[robotCacheItem]int {
	cache := make(map[robotCacheItem]int)
	for i := 1; i <= count; i++ {
		fillRobotBestPath(cache, i, directionalKeyPad)
	}
	fillRobotBestPath(cache, count+1, numericKeyPad)

	return cache
}

func minKeyPresses(count int, code string) int {
	cache := robotMovementCache(count)
	return sequencePressCost(cache, count+1, code)
}

func day21(robots int) func(io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		codes := aoc.Must(parseCodes(r))
		var total int
		for _, code := range codes {
			presses := minKeyPresses(robots, code)
			numericCode := aoc.Must(strconv.Atoi(code[:len(code)-1]))
			total += numericCode * presses
		}
		return strconv.Itoa(total), nil
	}
}

func day21p01(r io.Reader) (string, error) {
	return day21(2)(r)
}

func day21p02(r io.Reader) (string, error) {
	return day21(25)(r)
}
