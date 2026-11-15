package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func rockSegments(r io.Reader) ([][2]grid.Position2D[int], error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	var segments [][2]grid.Position2D[int]
	for line := range strings.SplitSeq(input, "\n") {
		var previous grid.Position2D[int]
		for i, text := range strings.Split(line, " -> ") {
			rock, err := grid.ParsePosition2D[int32](text)
			if err != nil {
				return nil, err
			}
			if rock.X < 0 || rock.Y < 0 {
				return nil, fmt.Errorf("negative rock coordinate %q", text)
			}
			point := grid.Position2D[int]{X: int(rock.X), Y: int(rock.Y)}
			if i == 0 {
				previous = point
			}
			if previous.X != point.X && previous.Y != point.Y {
				return nil, fmt.Errorf("diagonal rock segment in %q", line)
			}
			segments = append(segments, [2]grid.Position2D[int]{previous, point})
			previous = point
		}
	}
	return segments, nil
}

func sandUnits(r io.Reader, floor bool) (string, error) {
	segments, err := rockSegments(r)
	if err != nil {
		return "", err
	}
	bottom := 0
	for _, segment := range segments {
		bottom = max(bottom, segment[0].Y, segment[1].Y)
	}
	// Sand moves at most one column per row, so nothing outside this triangle
	// above the floor at bottom+2 can be reached from the source.
	left, width, height := 500-bottom-2, 2*bottom+5, bottom+2
	area, ok := xmath.CheckedMul(width, height)
	if !ok {
		return "", fmt.Errorf("sand area for depth %d overflows int", bottom)
	}
	occupied := make([]bool, area)
	for _, segment := range segments {
		low, high := segment[0], segment[1]
		low.X, high.X = min(low.X, high.X), max(low.X, high.X)
		low.Y, high.Y = min(low.Y, high.Y), max(low.Y, high.Y)
		for y := low.Y; y <= high.Y; y++ {
			for x := max(low.X, left); x <= min(high.X, left+width-1); x++ {
				occupied[y*width+x-left] = true
			}
		}
	}
	source := grid.Position2D[int]{X: 500}
	if occupied[source.X-left] {
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
		if at.Y+1 < height {
			for _, dx := range [3]int{0, -1, 1} {
				next := grid.Position2D[int]{X: at.X + dx, Y: at.Y + 1}
				if !occupied[next.Y*width+next.X-left] {
					path = append(path, next)
					moved = true
					break
				}
			}
		}
		if moved {
			continue
		}
		occupied[at.Y*width+at.X-left] = true
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
