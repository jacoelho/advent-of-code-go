package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func transpose(rows [][]byte) [][]byte {
	columns := make([][]byte, len(rows[0]))
	for x := range columns {
		columns[x] = make([]byte, len(rows))
		for y, row := range rows {
			columns[x][y] = row[x]
		}
	}
	return columns
}

// reflectionLine returns how many rows precede the line that mirrors rows with
// exactly smudges differing cells, or 0 when no such line exists.
func reflectionLine(rows [][]byte, smudges int) int {
	for line := 1; line < len(rows); line++ {
		differences := 0
		for top, bottom := line-1, line; top >= 0 && bottom < len(rows) && differences <= smudges; top, bottom = top-1, bottom+1 {
			for x := range rows[top] {
				if rows[top][x] != rows[bottom][x] {
					differences++
				}
			}
		}
		if differences == smudges {
			return line
		}
	}
	return 0
}

func day13(r io.Reader, smudges int) (string, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return "", err
	}
	if len(blocks) == 0 {
		return "", fmt.Errorf("no patterns")
	}
	total := 0
	for i, block := range blocks {
		g, err := grid.ReadBytes(strings.NewReader(strings.Join(block, "\n")), ".#")
		if err != nil {
			return "", fmt.Errorf("pattern %d: %w", i+1, err)
		}
		rows := g.Rows()
		if column := reflectionLine(transpose(rows), smudges); column > 0 {
			total += column
		} else if row := reflectionLine(rows, smudges); row > 0 {
			total += 100 * row
		} else {
			return "", fmt.Errorf("pattern %d has no reflection with %d smudges", i+1, smudges)
		}
	}
	return strconv.Itoa(total), nil
}

func day13p01(r io.Reader) (string, error) {
	return day13(r, 0)
}

func day13p02(r io.Reader) (string, error) {
	return day13(r, 1)
}
