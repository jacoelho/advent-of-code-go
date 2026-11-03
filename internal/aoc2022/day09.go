package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func ropePositions(r io.Reader, knots int) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	rope := make([]grid.Position2D[int], knots)
	visited := collections.NewSet(rope[knots-1])
	steps := 0
	for line := range strings.SplitSeq(input, "\n") {
		direction, countText, ok := strings.Cut(line, " ")
		if !ok {
			return "", fmt.Errorf("invalid rope motion")
		}
		count, err := convert.ScanNumber[int]([]byte(countText))
		if err != nil {
			return "", err
		}
		if count < 0 || count > 1_000_000-steps {
			return "", fmt.Errorf("rope exceeds one million steps")
		}
		steps += count
		var delta grid.Position2D[int]
		switch direction {
		case "R":
			delta.X = 1
		case "L":
			delta.X = -1
		case "U":
			delta.Y = 1
		case "D":
			delta.Y = -1
		default:
			return "", fmt.Errorf("invalid rope direction")
		}
		for range count {
			rope[0] = rope[0].Add(delta)
			for i := 1; i < len(rope); i++ {
				difference := rope[i-1].Sub(rope[i])
				if max(xmath.Abs(difference.X), xmath.Abs(difference.Y)) <= 1 {
					continue
				}
				rope[i].X += max(-1, min(1, difference.X))
				rope[i].Y += max(-1, min(1, difference.Y))
			}
			visited.Add(rope[knots-1])
		}
	}
	return strconv.Itoa(visited.Len()), nil
}

func day09p01(r io.Reader) (string, error) {
	return ropePositions(r, 2)
}

func day09p02(r io.Reader) (string, error) {
	return ropePositions(r, 10)
}
