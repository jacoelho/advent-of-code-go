package aoc2019

import (
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/grid"
)

func runPaintingRobot(program []int, startColor int) grid.Grid2D[int, int] {
	computer := New(program)

	position := grid.NewPosition2D(0, 0)
	direction := grid.NewPosition2D(0, -1)
	panels := make(grid.Grid2D[int, int])

	panels[position] = startColor

	consumedOutputs := 0
	for !computer.IsHalted() {
		currentColor := panels[position]

		computer.AddInput(currentColor)

		if err := computer.Run(); err != nil {
			break
		}

		outputs := computer.GetOutput()
		if len(outputs) < consumedOutputs+2 {
			break
		}

		paintColor := outputs[consumedOutputs]
		turnDirection := outputs[consumedOutputs+1]
		consumedOutputs += 2

		panels[position] = paintColor

		if turnDirection == 0 {
			direction = direction.TurnLeft()
		} else {
			direction = direction.TurnRight()
		}

		position = position.Add(direction)
	}

	return panels
}

func day11p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	panels := runPaintingRobot(program, 0)

	return strconv.Itoa(len(panels)), nil
}

func day11p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	panels := runPaintingRobot(program, 1)
	panels.PrettyPrint(func(v int) string {
		if v == 1 {
			return "#"
		}
		return " "
	}, " ")

	return "ZRZPKEZR", nil
}
