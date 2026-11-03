package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func sandUnits(r io.Reader, floor bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	occupied := collections.NewSet[grid.Position2D[int]]()
	bottom := 0
	work := 0
	for line := range strings.SplitSeq(input, "\n") {
		var previous grid.Position2D[int]
		for i, text := range strings.Split(line, " -> ") {
			coordinates := strings.Split(text, ",")
			if len(coordinates) != 2 {
				return "", fmt.Errorf("invalid rock point")
			}
			var values [2]int
			for j, coordinate := range coordinates {
				value, err := convert.ScanNumber[int]([]byte(coordinate))
				if err != nil {
					return "", err
				}
				if value < -1_000_000 || value > 1_000_000 {
					return "", fmt.Errorf("rock coordinate exceeds one million")
				}
				values[j] = value
			}
			point := grid.Position2D[int]{X: values[0], Y: values[1]}
			if point.Y < 0 {
				return "", fmt.Errorf("rock above sand source")
			}
			bottom = max(bottom, point.Y)
			if i == 0 {
				previous = point
			}
			if previous.X != point.X && previous.Y != point.Y {
				return "", fmt.Errorf("diagonal rock segment")
			}
			dx, dy := max(-1, min(1, point.X-previous.X)), max(-1, min(1, point.Y-previous.Y))
			for {
				work++
				if work > 2_000_000 {
					return "", fmt.Errorf("rock raster limit exceeded")
				}
				occupied.Add(previous)
				if previous == point {
					break
				}
				previous.X += dx
				previous.Y += dy
			}
		}
	}
	source := grid.Position2D[int]{X: 500}
	if occupied.Contains(source) {
		return "0", nil
	}
	path := []grid.Position2D[int]{source}
	settled := 0
	for len(path) > 0 {
		at := path[len(path)-1]
		if !floor && at.Y > bottom {
			return strconv.Itoa(settled), nil
		}
		moved := false
		for _, dx := range []int{0, -1, 1} {
			next := grid.Position2D[int]{X: at.X + dx, Y: at.Y + 1}
			if floor && next.Y == bottom+2 || occupied.Contains(next) {
				continue
			}
			path = append(path, next)
			moved = true
			break
		}
		if moved {
			if len(path) > 1_000_000 {
				return "", fmt.Errorf("sand depth limit exceeded")
			}
			continue
		}
		if settled >= 1_000_000 {
			return "", fmt.Errorf("sand exceeds one million grains")
		}
		occupied.Add(at)
		settled++
		path = path[:len(path)-1]
	}
	return strconv.Itoa(settled), nil
}

func day14p01(r io.Reader) (string, error) {
	return sandUnits(r, false)
}

func day14p02(r io.Reader) (string, error) {
	return sandUnits(r, true)
}
