package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

type infectionState byte

const (
	clean infectionState = iota
	weakened
	infected
	flagged
)

type virusMode int

const (
	simpleVirus virusMode = iota
	evolvedVirus
)

func parseInfections(r io.Reader) (map[grid.Position2D[int]]infectionState, error) {
	var rows [][]byte
	s := bufio.NewScanner(r)
	width := 0
	for s.Scan() {
		row := append([]byte(nil), s.Bytes()...)
		if len(rows) == 0 {
			width = len(row)
		}
		if len(row) != width || width%2 != 1 {
			return nil, fmt.Errorf("infection grid must have equal odd-width rows")
		}
		rows = append(rows, row)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(rows)%2 != 1 {
		return nil, fmt.Errorf("infection grid must have an odd number of rows")
	}

	nodes := make(map[grid.Position2D[int]]infectionState)
	for y, row := range rows {
		for x, cell := range row {
			switch cell {
			case '#':
				nodes[grid.Position2D[int]{X: x - width/2, Y: y - len(rows)/2}] = infected
			case '.':
			default:
				return nil, fmt.Errorf("invalid infection cell %q", cell)
			}
		}
	}

	return nodes, nil
}

func virusBursts(r io.Reader, mode virusMode, bursts int) (int, error) {
	if bursts < 0 || bursts > 10_000_000 {
		return 0, fmt.Errorf("virus burst count must be between 0 and 10000000")
	}
	nodes, err := parseInfections(r)
	if err != nil {
		return 0, err
	}

	pos, dir := grid.Position2D[int]{}, grid.Position2D[int]{Y: -1}
	infections := 0
	for range bursts {
		state := nodes[pos]
		switch state {
		case clean:
			dir = dir.TurnLeft()
		case infected:
			dir = dir.TurnRight()
		case flagged:
			dir = grid.Position2D[int]{X: -dir.X, Y: -dir.Y}
		}
		if mode == evolvedVirus {
			state = (state + 1) % 4
		} else if state == clean {
			state = infected
		} else {
			state = clean
		}
		if state == infected {
			infections++
		}
		if state == clean {
			delete(nodes, pos)
		} else {
			nodes[pos] = state
		}
		pos = pos.Add(dir)
	}

	return infections, nil
}

func day22p01(r io.Reader) (string, error) {
	count, err := virusBursts(r, simpleVirus, 10_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}

func day22p02(r io.Reader) (string, error) {
	count, err := virusBursts(r, evolvedVirus, 10_000_000)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(count), nil
}
