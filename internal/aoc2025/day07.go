package aoc2025

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/decimal"
)

type tachyonManifold struct {
	rows  [][]byte
	start grid.Position2D[int]
}

func parseTachyonManifoldDiagram(r io.Reader) (tachyonManifold, error) {
	g, err := grid.ReadBytes(r, ".^S")
	if err != nil {
		return tachyonManifold{}, fmt.Errorf("manifold: %w", err)
	}
	start, err := g.Find('S')
	if err != nil {
		return tachyonManifold{}, fmt.Errorf("manifold: %w", err)
	}
	return tachyonManifold{rows: g.Rows(), start: start}, nil
}

// trace sends the beam down from the start one row at a time. A splitter
// replaces each arriving beam with one in each neighbouring column, which
// continue from the next row; a beam leaving the sides or the bottom ends its
// timeline. It returns how many splitters are hit and how many timelines end,
// or an error if the timeline count exceeds the Decimal range.
func (m tachyonManifold) trace() (int, decimal.Decimal, error) {
	width := len(m.rows[0])
	beams, next := make([]decimal.Decimal, width), make([]decimal.Decimal, width)
	beams[m.start.X] = decimal.FromInt64(1)
	splits := 0
	var timelines decimal.Decimal
	add := func(total *decimal.Decimal, count decimal.Decimal) error {
		sum, err := total.Add(count)
		if err != nil {
			return fmt.Errorf("tachyon timelines: %w", err)
		}
		*total = sum
		return nil
	}
	for _, row := range m.rows[m.start.Y+1:] {
		clear(next)
		for x, count := range beams {
			if count.IsZero() {
				continue
			}
			if row[x] != '^' {
				if err := add(&next[x], count); err != nil {
					return 0, decimal.Decimal{}, err
				}
				continue
			}
			splits++
			for _, side := range [2]int{x - 1, x + 1} {
				total := &timelines
				if 0 <= side && side < width {
					total = &next[side]
				}
				if err := add(total, count); err != nil {
					return 0, decimal.Decimal{}, err
				}
			}
		}
		beams, next = next, beams
	}
	for _, count := range beams {
		if err := add(&timelines, count); err != nil {
			return 0, decimal.Decimal{}, err
		}
	}
	return splits, timelines, nil
}

func day07p01(r io.Reader) (string, error) {
	manifold, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}
	splits, _, err := manifold.trace()
	if err != nil {
		return "", err
	}
	return strconv.Itoa(splits), nil
}

func day07p02(r io.Reader) (string, error) {
	manifold, err := parseTachyonManifoldDiagram(r)
	if err != nil {
		return "", err
	}
	_, timelines, err := manifold.trace()
	if err != nil {
		return "", err
	}
	return timelines.String(), nil
}
