package aoc2025

import (
	"io"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseRollsPaper(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})
	g := grid.NewGrid2D[int](slices.Collect(s.Values()))
	maps.DeleteFunc(g, func(_ grid.Position2D[int], v rune) bool {
		return v != '@'
	})
	return g, s.Err()
}

func accessibleRolls(g grid.Grid2D[int, rune]) []grid.Position2D[int] {
	var removable []grid.Position2D[int]
	for pos := range g {
		neighborCount := xiter.Len(g.ValidNeighbours8(pos))
		if neighborCount < 4 {
			removable = append(removable, pos)
		}
	}
	return removable
}

func day04p01(r io.Reader) (string, error) {
	g, err := parseRollsPaper(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(len(accessibleRolls(g))), nil
}

func day04p02(r io.Reader) (string, error) {
	g, err := parseRollsPaper(r)
	if err != nil {
		return "", err
	}

	total := 0
	for {
		removable := accessibleRolls(g)
		if len(removable) == 0 {
			break
		}
		for _, pos := range removable {
			delete(g, pos)
		}
		total += len(removable)
	}
	return strconv.Itoa(total), nil
}
