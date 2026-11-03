package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
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
	if value < 0 || value > 1_000_000 {
		return 0, fmt.Errorf("paper coordinate %d is outside 0..1000000", value)
	}
	return value, nil
}

func transparentPaper(r io.Reader) (collections.Set[grid.Position2D[int]], []paperFold, error) {
	points := collections.NewSet[grid.Position2D[int]]()
	var folds []paperFold
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		if text, ok := strings.CutPrefix(line, "fold along "); ok {
			axis, value, ok := strings.Cut(text, "=")
			if !ok || axis != "x" && axis != "y" {
				return nil, nil, fmt.Errorf("invalid paper fold %q", line)
			}
			position, err := paperCoordinate(value)
			if err != nil {
				return nil, nil, err
			}
			folds = append(folds, paperFold{axis: axis[0], line: position})
			continue
		}
		xText, yText, ok := strings.Cut(line, ",")
		if !ok {
			return nil, nil, fmt.Errorf("invalid paper point %q", line)
		}
		x, err := paperCoordinate(xText)
		if err != nil {
			return nil, nil, err
		}
		y, err := paperCoordinate(yText)
		if err != nil {
			return nil, nil, err
		}
		points.Add(grid.Position2D[int]{X: x, Y: y})
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if points.IsEmpty() {
		return nil, nil, fmt.Errorf("paper has no dots")
	}
	return points, folds, nil
}

func foldPaper(points collections.Set[grid.Position2D[int]], fold paperFold) collections.Set[grid.Position2D[int]] {
	result := collections.NewSet[grid.Position2D[int]]()
	for point := range points {
		if fold.axis == 'x' && point.X > fold.line {
			point.X = fold.line - (point.X - fold.line)
		}
		if fold.axis == 'y' && point.Y > fold.line {
			point.Y = fold.line - (point.Y - fold.line)
		}
		result.Add(point)
	}
	return result
}

func renderPaper(points collections.Set[grid.Position2D[int]]) ([]string, error) {
	var left, top, right, bottom int
	for point := range points {
		left, top = min(left, point.X), min(top, point.Y)
		right, bottom = max(right, point.X), max(bottom, point.Y)
	}
	width, height := right-left+1, bottom-top+1
	if width <= 0 || height <= 0 || width > 1_000_000/height {
		return nil, fmt.Errorf("folded paper exceeds display budget")
	}
	rows := make([]string, height)
	for y := range height {
		row := []byte(strings.Repeat(".", width))
		for x := range width {
			if points.Contains(grid.Position2D[int]{X: left + x, Y: top + y}) {
				row[x] = '#'
			}
		}
		rows[y] = string(row)
	}
	return rows, nil
}

func day13p01(r io.Reader) (string, error) {
	points, folds, err := transparentPaper(r)
	if err != nil {
		return "", err
	}
	if len(folds) == 0 {
		return "", fmt.Errorf("paper has no fold instructions")
	}
	return strconv.Itoa(foldPaper(points, folds[0]).Len()), nil
}

func day13p02(r io.Reader) (string, error) {
	points, folds, err := transparentPaper(r)
	if err != nil {
		return "", err
	}
	for _, fold := range folds {
		points = foldPaper(points, fold)
	}
	rows, err := renderPaper(points)
	if err != nil {
		return "", err
	}
	if len(rows) != 6 {
		return "", fmt.Errorf("folded message requires a supported six-row display")
	}
	padding := (5 - len(rows[0])%5) % 5
	for i := range rows {
		rows[i] += strings.Repeat(".", padding)
	}
	return aoc.Decode6RowDisplay(rows, 5)
}
