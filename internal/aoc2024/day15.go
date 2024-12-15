package aoc2024

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

func convertDirection(r rune) grid.Position2D[int] {
	switch r {
	case '<':
		return grid.Position2D[int]{X: -1, Y: 0}
	case '>':
		return grid.Position2D[int]{X: 1, Y: 0}
	case 'v':
		return grid.Position2D[int]{X: 0, Y: 1}
	case '^':
		return grid.Position2D[int]{X: 0, Y: -1}
	default:
		panic("unreachable")
	}
}

func warehouseIdentityTile(r rune) []rune { return []rune{r} }

func warehouseScaleUpTile(r rune) []rune {
	switch r {
	case '#', '.':
		return []rune{r, r}
	case 'O':
		return []rune{'[', ']'}
	case '@':
		return []rune{'@', '.'}
	}
	panic("unreachable")
}

func parseWarehouse(
	r io.Reader,
	transform func(r rune) []rune,
) (grid.Grid2D[int, rune], []grid.Position2D[int]) {
	s := bufio.NewScanner(r)
	s.Split(scanner.SplitBySeparator([]byte{'\n', '\n'}))

	_, mapString := s.Scan(), s.Text()
	_, movementString := s.Scan(), s.Text()
	if s.Err() != nil {
		panic(s.Err())
	}

	mapLines := xslices.Map(func(in string) []rune {
		var result []rune
		for _, r := range in {
			result = append(result, transform(r)...)
		}
		return result
	}, strings.Split(mapString, "\n"))

	var movementLines []grid.Position2D[int]
	for _, line := range strings.Split(movementString, "\n") {
		for _, ch := range line {
			movementLines = append(movementLines, convertDirection(ch))
		}
	}

	return grid.NewGrid2D[int, rune](mapLines), movementLines
}

func robotStartPosition(g grid.Grid2D[int, rune]) grid.Position2D[int] {
	element, ok := xmaps.Find(g, func(p grid.Position2D[int], v rune) bool { return v == '@' })
	if !ok {
		panic("robot not found")
	}
	return element.K
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
	warehouse, movements := parseWarehouse(r, warehouseIdentityTile)
	robotPosition := robotStartPosition(warehouse)

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
	warehouse, movements := parseWarehouse(r, warehouseScaleUpTile)
	robotPosition := robotStartPosition(warehouse)

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
