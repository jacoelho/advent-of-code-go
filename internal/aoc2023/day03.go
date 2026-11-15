package aoc2023

import (
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const schematicAlphabet = "0123456789.!\"#$%&'()*+,-/:;<=>?@[\\]^_`{|}~ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

type schematicNumber struct {
	value           int
	row, start, end int
}

type schematicCell struct {
	x, y int
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func schematicNumbers(rows [][]byte) iter.Seq[schematicNumber] {
	return func(yield func(schematicNumber) bool) {
		for y, line := range rows {
			for x := 0; x < len(line); x++ {
				if !isDigit(line[x]) {
					continue
				}
				number := schematicNumber{row: y, start: x}
				for ; x < len(line) && isDigit(line[x]); x++ {
					number.value = number.value*10 + int(line[x]-'0')
				}
				number.end = x
				if !yield(number) {
					return
				}
			}
		}
	}
}

// adjacentSymbols yields the symbol cells that touch number, diagonals included.
func adjacentSymbols(rows [][]byte, number schematicNumber) iter.Seq[schematicCell] {
	return func(yield func(schematicCell) bool) {
		for y := max(0, number.row-1); y <= min(len(rows)-1, number.row+1); y++ {
			for x := max(0, number.start-1); x <= min(len(rows[y])-1, number.end); x++ {
				if c := rows[y][x]; c != '.' && !isDigit(c) && !yield(schematicCell{x, y}) {
					return
				}
			}
		}
	}
}

func day03p01(r io.Reader) (string, error) {
	g, err := grid.ReadBytes(r, schematicAlphabet)
	if err != nil {
		return "", err
	}
	rows := g.Rows()
	sum := 0
	for number := range schematicNumbers(rows) {
		for range adjacentSymbols(rows, number) {
			sum += number.value
			break
		}
	}
	return strconv.Itoa(sum), nil
}

func day03p02(r io.Reader) (string, error) {
	g, err := grid.ReadBytes(r, schematicAlphabet)
	if err != nil {
		return "", err
	}
	rows := g.Rows()
	gears := make(map[schematicCell][]int)
	for number := range schematicNumbers(rows) {
		for cell := range adjacentSymbols(rows, number) {
			if rows[cell.y][cell.x] == '*' {
				gears[cell] = append(gears[cell], number.value)
			}
		}
	}
	sum := 0
	for _, values := range gears {
		if len(values) == 2 {
			sum += values[0] * values[1]
		}
	}
	return strconv.Itoa(sum), nil
}
