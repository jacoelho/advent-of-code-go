package aoc2019

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type hullPanel struct {
	color   int
	painted bool
}

func runPaintingRobot(program []int, startColor int) (grid.Grid2D[int, hullPanel], error) {
	computer := New(program)
	position := grid.NewPosition2D(0, 0)
	direction := grid.NewPosition2D(0, -1)
	panels := make(grid.Grid2D[int, hullPanel])
	panels[position] = hullPanel{color: startColor}

	for range 100_000 {
		computer.AddInput(panels[position].color)
		if err := computer.Run(); err != nil {
			return nil, err
		}
		output := computer.GetOutput()
		if len(output) == 0 && computer.IsHalted() {
			return panels, nil
		}
		if len(output) != 2 || output[0] < 0 || output[0] > 1 || output[1] < 0 || output[1] > 1 {
			return nil, fmt.Errorf("painting robot requires one paint/turn pair of zero or one")
		}

		panels[position] = hullPanel{color: output[0], painted: true}
		if output[1] == 0 {
			direction = direction.TurnLeft()
		} else {
			direction = direction.TurnRight()
		}
		position = position.Add(direction)
		if computer.IsHalted() {
			return panels, nil
		}
	}
	return nil, fmt.Errorf("painting robot exceeds 100000 commands")
}

func paintingDisplay(panels grid.Grid2D[int, hullPanel]) ([]string, error) {
	var minX, maxX, minY, maxY int
	found := false
	for position, panel := range panels {
		if panel.color != 1 {
			continue
		}
		if !found {
			minX, maxX, minY, maxY = position.X, position.X, position.Y, position.Y
			found = true
		}
		minX, maxX = min(minX, position.X), max(maxX, position.X)
		minY, maxY = min(minY, position.Y), max(maxY, position.Y)
	}
	if !found {
		return nil, fmt.Errorf("painting has no white panels")
	}
	width, height := maxX-minX+1, maxY-minY+1
	if width > 1_000_000 || height > 1_000_000/width {
		return nil, fmt.Errorf("painting display exceeds one million pixels")
	}
	rows := make([]string, height)
	for y := range height {
		row := make([]byte, width)
		for x := range width {
			row[x] = '.'
			if panels[grid.Position2D[int]{X: minX + x, Y: minY + y}].color == 1 {
				row[x] = '#'
			}
		}
		rows[y] = string(row)
	}
	return rows, nil
}

func day11p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}
	panels, err := runPaintingRobot(program, 0)
	if err != nil {
		return "", err
	}
	painted := 0
	for _, panel := range panels {
		if panel.painted {
			painted++
		}
	}
	return strconv.Itoa(painted), nil
}

func day11p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}
	panels, err := runPaintingRobot(program, 1)
	if err != nil {
		return "", err
	}
	rows, err := paintingDisplay(panels)
	if err != nil {
		return "", err
	}
	return aoc.DecodeSeparatedDisplay(rows)
}
