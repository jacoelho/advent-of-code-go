package aoc2021

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type paperFold struct {
	axis byte
	line int
}

func paperCoordinate(text string) (int, error) {
	value, err := convert.ScanNumber[int]([]byte(text))
	if err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, fmt.Errorf("negative paper coordinate %d", value)
	}
	return value, nil
}

func transparentPaper(r io.Reader) (collections.Set[grid.Position2D[int]], []paperFold, error) {
	points := collections.NewSet[grid.Position2D[int]]()
	var folds []paperFold
	sections, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(sections) > 2 || len(sections) == 2 && len(sections[1]) == 0 {
		return nil, nil, fmt.Errorf("blank line among paper folds")
	}
	var dots, foldLines []string
	if len(sections) > 0 {
		dots = sections[0]
	}
	if len(sections) > 1 {
		foldLines = sections[1]
	}
	for _, line := range dots {
		point, err := grid.ParsePosition2D[int](line)
		if err != nil {
			return nil, nil, err
		}
		if point.X < 0 || point.Y < 0 {
			return nil, nil, fmt.Errorf("negative paper point %q", line)
		}
		points.Add(point)
	}
	for _, line := range foldLines {
		text, ok := strings.CutPrefix(line, "fold along ")
		if !ok {
			return nil, nil, fmt.Errorf("invalid paper fold %q", line)
		}
		axis, value, ok := strings.Cut(text, "=")
		if !ok || axis != "x" && axis != "y" {
			return nil, nil, fmt.Errorf("invalid paper fold %q", line)
		}
		position, err := paperCoordinate(value)
		if err != nil {
			return nil, nil, err
		}
		folds = append(folds, paperFold{axis: axis[0], line: position})
	}
	if points.IsEmpty() {
		return nil, nil, fmt.Errorf("paper has no dots")
	}
	return points, folds, nil
}

func foldPaper(points collections.Set[grid.Position2D[int]], fold paperFold) (collections.Set[grid.Position2D[int]], error) {
	result := collections.NewSet[grid.Position2D[int]]()
	for point := range points {
		coordinate := &point.X
		if fold.axis == 'y' {
			coordinate = &point.Y
		}
		if *coordinate > fold.line {
			*coordinate = fold.line - (*coordinate - fold.line)
		}
		if *coordinate < 0 {
			return nil, fmt.Errorf("dot folds past the paper edge along %c=%d", fold.axis, fold.line)
		}
		result.Add(point)
	}
	return result, nil
}

func renderPaper(points collections.Set[grid.Position2D[int]]) ([]string, error) {
	var width, height int
	for point := range points {
		if point.X == math.MaxInt {
			return nil, fmt.Errorf("paper width overflows int")
		}
		width, height = max(width, point.X+1), max(height, point.Y+1)
	}
	rows := make([][]byte, height)
	for y := range rows {
		rows[y] = []byte(strings.Repeat(".", width))
	}
	for point := range points {
		rows[point.Y][point.X] = '#'
	}
	display := make([]string, height)
	for y, row := range rows {
		display[y] = string(row)
	}
	return display, nil
}

func day13p01(r io.Reader) (string, error) {
	points, folds, err := transparentPaper(r)
	if err != nil {
		return "", err
	}
	if len(folds) == 0 {
		return "", fmt.Errorf("paper has no fold instructions")
	}
	points, err = foldPaper(points, folds[0])
	if err != nil {
		return "", err
	}
	return strconv.Itoa(points.Len()), nil
}

func day13p02(r io.Reader) (string, error) {
	points, folds, err := transparentPaper(r)
	if err != nil {
		return "", err
	}
	for _, fold := range folds {
		if points, err = foldPaper(points, fold); err != nil {
			return "", err
		}
	}
	for point := range points {
		if point.Y >= 6 {
			return "", fmt.Errorf("folded message requires a six-row display, dot at row %d", point.Y)
		}
	}
	rows, err := renderPaper(points)
	if err != nil {
		return "", err
	}
	if len(rows) != 6 {
		return "", fmt.Errorf("folded message requires a six-row display, got %d rows", len(rows))
	}
	padding := (5 - len(rows[0])%5) % 5
	for i := range rows {
		rows[i] += strings.Repeat(".", padding)
	}
	return aoc.Decode6RowDisplay(rows, 5)
}
