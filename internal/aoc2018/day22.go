package aoc2018

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type rescueCave struct {
	depth   int
	target  grid.Position2D[int]
	erosion [][]uint16
}

func parseRescueCave(r io.Reader) (*rescueCave, error) {
	s := bufio.NewScanner(r)
	cave := &rescueCave{}
	line := 0
	for s.Scan() {
		switch line {
		case 0:
			text, ok := bytes.CutPrefix(s.Bytes(), []byte("depth: "))
			if !ok {
				return nil, fmt.Errorf("missing cave depth")
			}
			depth, err := convert.ScanNumber[int64](text)
			if err != nil {
				return nil, err
			}
			if depth < 0 {
				return nil, fmt.Errorf("negative cave depth")
			}
			cave.depth = int(depth % 20183)
		case 1:
			text, ok := bytes.CutPrefix(s.Bytes(), []byte("target: "))
			if !ok {
				return nil, fmt.Errorf("missing cave target")
			}
			fields := bytes.Split(text, []byte(","))
			if len(fields) != 2 {
				return nil, fmt.Errorf("expected two cave target coordinates")
			}
			for i, field := range fields {
				value, err := convert.ScanNumber[int32](field)
				if err != nil {
					return nil, err
				}
				if value < 0 || value > 100_000 {
					return nil, fmt.Errorf("cave coordinates must be between 0 and 100000")
				}
				if i == 0 {
					cave.target.X = int(value)
				} else {
					cave.target.Y = int(value)
				}
			}
		default:
			return nil, fmt.Errorf("unexpected cave input %q", s.Text())
		}
		line++
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if line != 2 {
		return nil, fmt.Errorf("expected cave depth and target")
	}

	return cave, nil
}

func (cave *rescueCave) region(pos grid.Position2D[int]) (int, error) {
	if pos.Y < len(cave.erosion) && pos.X < len(cave.erosion[pos.Y]) {
		return int(cave.erosion[pos.Y][pos.X] % 3), nil
	}
	width, height := pos.X+1, max(pos.Y+1, len(cave.erosion))
	if len(cave.erosion) != 0 {
		width = max(width, len(cave.erosion[0]))
	}
	if width > 4_000_000/height {
		return 0, fmt.Errorf("cave exploration exceeded 4000000 erosion cells")
	}
	for len(cave.erosion) < height {
		cave.erosion = append(cave.erosion, nil)
	}
	for y := range height {
		for x := len(cave.erosion[y]); x < width; x++ {
			var index int64
			switch {
			case x == 0 && y == 0 || x == cave.target.X && y == cave.target.Y:
			case y == 0:
				index = int64(x) * 16807
			case x == 0:
				index = int64(y) * 48271
			default:
				index = int64(cave.erosion[y][x-1]) * int64(cave.erosion[y-1][x])
			}
			value := uint16((index + int64(cave.depth)) % 20183)
			cave.erosion[y] = append(cave.erosion[y], value)
		}
	}

	return int(cave.erosion[pos.Y][pos.X] % 3), nil
}

func day22p01(r io.Reader) (string, error) {
	cave, err := parseRescueCave(r)
	if err != nil {
		return "", err
	}
	if _, err := cave.region(cave.target); err != nil {
		return "", err
	}

	risk := 0
	for _, row := range cave.erosion {
		for _, erosion := range row {
			risk += int(erosion % 3)
		}
	}

	return strconv.Itoa(risk), nil
}

func day22p02(r io.Reader) (string, error) {
	cave, err := parseRescueCave(r)
	if err != nil {
		return "", err
	}
	const torch = 1
	if cave.depth%3 == torch {
		return "", fmt.Errorf("torch is not permitted at the cave entrance and target")
	}

	type state struct {
		pos  grid.Position2D[int]
		tool int
	}
	var searchErr error
	expanded := 0
	neighbours := func(current state) []state {
		if searchErr != nil {
			return nil
		}
		expanded++
		if expanded > 1_000_000 {
			searchErr = fmt.Errorf("cave search exceeded 1000000 expansions")
			return nil
		}
		region, err := cave.region(current.pos)
		if err != nil {
			searchErr = err
			return nil
		}
		var next []state
		for tool := range 3 {
			if tool != current.tool && tool != region {
				next = append(next, state{pos: current.pos, tool: tool})
			}
		}
		for pos := range grid.Neighbours4(current.pos) {
			region, err := cave.region(pos)
			if err != nil {
				searchErr = err
				return nil
			}
			if current.tool != region {
				next = append(next, state{pos: pos, tool: current.tool})
			}
		}
		return next
	}
	heuristic := func(current state) int {
		distance := current.pos.Distance(cave.target)
		if current.tool != torch {
			distance += 7
		}
		return distance
	}
	stepCost := func(from, to state) int {
		if from.pos == to.pos {
			return 7
		}
		return 1
	}
	minutes, _, found := search.AStar(state{tool: torch}, neighbours, heuristic, stepCost)
	if searchErr != nil {
		return "", searchErr
	}
	if !found {
		return "", fmt.Errorf("cave target is unreachable with the torch")
	}

	return strconv.Itoa(minutes), nil
}
