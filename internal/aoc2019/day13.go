package aoc2019

import (
	"cmp"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type arcadeCabinet struct {
	tiles          grid.Grid2D[int, int]
	score          int
	ballX, paddleX int
}

func (cabinet *arcadeCabinet) update(output []int) error {
	if len(output)%3 != 0 {
		return fmt.Errorf("arcade output contains a partial triple")
	}
	for i := 0; i < len(output); i += 3 {
		x, y, tile := output[i], output[i+1], output[i+2]
		if x == -1 && y == 0 {
			cabinet.score = tile
			continue
		}
		if x < 0 || y < 0 || x > 1_000_000 || y > 1_000_000 || tile < 0 || tile > 4 {
			return fmt.Errorf("invalid arcade tile (%d,%d,%d)", x, y, tile)
		}
		position := grid.Position2D[int]{X: x, Y: y}
		if _, exists := cabinet.tiles[position]; !exists && len(cabinet.tiles) == 100_000 {
			return fmt.Errorf("arcade exceeds 100000 tile positions")
		}
		cabinet.tiles[position] = tile
		switch tile {
		case 3:
			cabinet.paddleX = x
		case 4:
			cabinet.ballX = x
		}
	}
	return nil
}

func playArcade(program []int, freePlay bool) (arcadeCabinet, error) {
	computer := New(program)
	cabinet := arcadeCabinet{tiles: make(grid.Grid2D[int, int])}
	if freePlay {
		if err := computer.SetMemory(0, 2); err != nil {
			return arcadeCabinet{}, err
		}
	}
	updates := 0
	for range 100_000 {
		if err := computer.Run(); err != nil {
			return arcadeCabinet{}, err
		}
		output := computer.GetOutput()
		updates += len(output) / 3
		if updates > 1_000_000 {
			return arcadeCabinet{}, fmt.Errorf("arcade exceeds one million updates")
		}
		if err := cabinet.update(output); err != nil {
			return arcadeCabinet{}, err
		}
		if computer.IsHalted() {
			return cabinet, nil
		}
		if !freePlay {
			return arcadeCabinet{}, fmt.Errorf("arcade scan requested joystick input")
		}
		computer.AddInput(cmp.Compare(cabinet.ballX, cabinet.paddleX))
	}
	return arcadeCabinet{}, fmt.Errorf("arcade exceeds 100000 input requests")
}

func day13p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}
	cabinet, err := playArcade(program, false)
	if err != nil {
		return "", err
	}
	blocks := 0
	for _, tile := range cabinet.tiles {
		if tile == 2 {
			blocks++
		}
	}
	return strconv.Itoa(blocks), nil
}

func day13p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}
	cabinet, err := playArcade(program, true)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(cabinet.score), nil
}
