package aoc2024

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
	"github.com/jacoelho/advent-of-code-go/pkg/xstrings"
)

var numericKeyPad = map[rune]grid.Position2D[int]{
	'7': {X: 0, Y: 0}, '8': {X: 1, Y: 0}, '9': {X: 2, Y: 0},
	'4': {X: 0, Y: 1}, '5': {X: 1, Y: 1}, '6': {X: 2, Y: 1},
	'1': {X: 0, Y: 2}, '2': {X: 1, Y: 2}, '3': {X: 2, Y: 2},
	' ': {X: 0, Y: 3}, '0': {X: 1, Y: 3}, 'A': {X: 2, Y: 3},
}

var directionalKeyPad = map[rune]grid.Position2D[int]{
	' ': {X: 0, Y: 0}, '^': {X: 1, Y: 0}, 'A': {X: 2, Y: 0},
	'<': {X: 0, Y: 1}, 'v': {X: 1, Y: 1}, '>': {X: 2, Y: 1},
}

type keypadCode struct {
	keys   string
	number int
}

func parseCodes(r io.Reader) ([]keypadCode, error) {
	s := scanner.NewScanner(r, func(b []byte) (keypadCode, error) {
		if len(b) < 2 || b[len(b)-1] != 'A' {
			return keypadCode{}, fmt.Errorf("keypad code needs digits followed by A")
		}
		for _, digit := range b[:len(b)-1] {
			if digit < '0' || digit > '9' {
				return keypadCode{}, fmt.Errorf("invalid keypad code %q", b)
			}
		}
		number, err := convert.ScanNumber[int](b[:len(b)-1])
		if err != nil {
			return keypadCode{}, fmt.Errorf("keypad code: %w", err)
		}
		return keypadCode{keys: string(b), number: number}, nil
	})
	codes := slices.Collect(s.Values())
	return codes, s.Err()
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
	for range firstDistance {
		sb.WriteString(firstPress)
	}

	for range secondDistance {
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
		codes, err := parseCodes(r)
		if err != nil {
			return "", err
		}
		var total int
		for _, code := range codes {
			presses := minKeyPresses(robots, code.keys)
			total += code.number * presses
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
