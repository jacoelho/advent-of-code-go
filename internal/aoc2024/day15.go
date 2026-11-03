package aoc2024

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func convertDirection(r rune) (grid.Position2D[int], error) {
	switch r {
	case '<':
		return grid.Position2D[int]{X: -1, Y: 0}, nil
	case '>':
		return grid.Position2D[int]{X: 1, Y: 0}, nil
	case 'v':
		return grid.Position2D[int]{X: 0, Y: 1}, nil
	case '^':
		return grid.Position2D[int]{X: 0, Y: -1}, nil
	default:
		return grid.Position2D[int]{}, fmt.Errorf("invalid warehouse movement %q", r)
	}
}

func warehouseIdentityTile(r rune) []rune {
	return []rune{r}
}

func warehouseScaleUpTile(r rune) []rune {
	switch r {
	case 'O':
		return []rune{'[', ']'}
	case '@':
		return []rune{'@', '.'}
	}
	return []rune{r, r}
}

func parseWarehouse(
	r io.Reader,
	transform func(r rune) []rune,
) (grid.Grid2D[int, rune], []grid.Position2D[int], error) {
	const inputLimit = 128 << 10
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, nil, err
	}
	if len(data) > inputLimit {
		return nil, nil, fmt.Errorf("warehouse input exceeds 128 KiB")
	}
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	mapString, movementString, ok := strings.Cut(text, "\n\n")
	if !ok {
		return nil, nil, fmt.Errorf("warehouse needs a map and movement section")
	}

	rows := strings.Split(mapString, "\n")
	width := len(rows[0])
	if width < 3 || len(rows) < 3 {
		return nil, nil, fmt.Errorf("warehouse needs an enclosed map")
	}
	mapLines := make([][]rune, len(rows))
	for y, row := range rows {
		if len(row) != width {
			return nil, nil, fmt.Errorf("warehouse row %d has inconsistent width", y+1)
		}
		for x, tile := range row {
			if tile != '#' && tile != '.' && tile != 'O' && tile != '@' {
				return nil, nil, fmt.Errorf("invalid warehouse tile %q", tile)
			}
			if (x == 0 || x == width-1 || y == 0 || y == len(rows)-1) && tile != '#' {
				return nil, nil, fmt.Errorf("warehouse boundary must be walls")
			}
			mapLines[y] = append(mapLines[y], transform(tile)...)
		}
	}

	var movements []grid.Position2D[int]
	for _, ch := range movementString {
		if ch == '\n' {
			continue
		}
		direction, err := convertDirection(ch)
		if err != nil {
			return nil, nil, err
		}
		movements = append(movements, direction)
	}
	return grid.NewGrid2D[int](mapLines), movements, nil
}

func robotStartPosition(g grid.Grid2D[int, rune]) (grid.Position2D[int], error) {
	var position grid.Position2D[int]
	found := false
	for p, tile := range g {
		if tile != '@' {
			continue
		}
		if found {
			return position, fmt.Errorf("warehouse has multiple robots")
		}
		position, found = p, true
	}
	if !found {
		return position, fmt.Errorf("warehouse robot not found")
	}
	return position, nil
}

func move(
	g grid.Grid2D[int, rune],
	robot grid.Position2D[int],
	direction grid.Position2D[int],
) grid.Position2D[int] {
	var boxesToMove []grid.Position2D[int]

LOOP:
	for currentPosition := robot.Add(direction); ; currentPosition = currentPosition.Add(direction) {
		switch g[currentPosition] {
		case '.':
			break LOOP
		case '#':
			return robot
		case 'O':
			boxesToMove = append(boxesToMove, currentPosition)
		}
	}

	g[robot] = '.'
	g[robot.Add(direction)] = '@'

	for i := len(boxesToMove) - 1; i >= 0; i-- {
		box := boxesToMove[i]
		g[box.Add(direction)] = '.'
		g[box.Add(direction)] = 'O'
	}

	return robot.Add(direction)
}

func day15p01(r io.Reader) (string, error) {
	warehouse, movements, err := parseWarehouse(r, warehouseIdentityTile)
	if err != nil {
		return "", err
	}
	robotPosition, err := robotStartPosition(warehouse)
	if err != nil {
		return "", err
	}

	for _, direction := range movements {
		robotPosition = move(warehouse, robotPosition, direction)
	}

	var total int
	for p, v := range warehouse {
		if v == 'O' {
			total += 100*p.Y + p.X
		}
	}
	return strconv.Itoa(total), nil
}

func moveMultiple(
	g grid.Grid2D[int, rune],
	robot grid.Position2D[int],
	direction grid.Position2D[int],
) grid.Position2D[int] {
	visited := collections.NewSet[grid.Position2D[int]]()
	frontier := collections.NewDeque[grid.Position2D[int]](10)
	frontier.PushBack(robot)

	type pair struct {
		position grid.Position2D[int]
		symbol   rune
	}

	var boxesToMove []pair

	movePair := func(pos grid.Position2D[int], sym rune) {
		oppositeBracket := func(bracket rune) rune {
			if bracket == '[' {
				return ']'
			}
			return '['
		}

		offsetBracket := func(bracket rune) grid.Position2D[int] {
			if bracket == '[' {
				return grid.Position2D[int]{X: 1, Y: 0}
			}
			return grid.Position2D[int]{X: -1, Y: 0}
		}

		p1 := pair{position: pos, symbol: sym}
		p2 := pair{position: pos.Add(offsetBracket(sym)), symbol: oppositeBracket(sym)}

		boxesToMove = append(boxesToMove, p1, p2)

		for _, p := range []grid.Position2D[int]{p1.position, p2.position} {
			if !visited.Contains(p) {
				visited.Add(p)
				frontier.PushBack(p)
			}
		}
	}

	for frontier.Size() > 0 {
		node, _ := frontier.PopFront()
		node = node.Add(direction)

		switch g[node] {
		case '#':
			return robot
		case '[':
			movePair(node, '[')
		case ']':
			movePair(node, ']')
		}
	}

	g[robot] = '.'
	g[robot.Add(direction)] = '@'

	for i := len(boxesToMove) - 1; i >= 0; i-- {
		box := boxesToMove[i]
		g[box.position] = '.'
		g[box.position.Add(direction)] = box.symbol
	}

	return robot.Add(direction)
}

func day15p02(r io.Reader) (string, error) {
	warehouse, movements, err := parseWarehouse(r, warehouseScaleUpTile)
	if err != nil {
		return "", err
	}
	robotPosition, err := robotStartPosition(warehouse)
	if err != nil {
		return "", err
	}

	for _, direction := range movements {
		robotPosition = moveMultiple(warehouse, robotPosition, direction)
	}

	var total int
	for p, v := range warehouse {
		if v == '[' {
			total += 100*p.Y + p.X
		}
	}
	return strconv.Itoa(total), nil
}
