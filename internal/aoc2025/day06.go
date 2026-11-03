package aoc2025

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type problem struct {
	op      rune
	numbers []int
}

func readProblems(r io.Reader) ([][]rune, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		return []rune(string(b)), nil
	})
	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("homework needs numbers and an operator row")
	}
	width := 0
	for _, row := range rows {
		width = max(width, len(row))
	}
	if width == 0 {
		return nil, fmt.Errorf("empty homework")
	}
	for i, row := range rows {
		for _, cell := range row {
			valid := cell == ' ' || cell >= '0' && cell <= '9'
			if i == len(rows)-1 {
				valid = cell == ' ' || cell == '+' || cell == '*'
			}
			if !valid {
				return nil, fmt.Errorf("invalid homework character %q on row %d", cell, i+1)
			}
		}
		for len(row) < width {
			row = append(row, ' ')
		}
		rows[i] = row
	}
	return rows, nil
}

func isSeparatorColumn(grid [][]rune, col int) bool {
	return xiter.FromSlice(grid).Every(func(row []rune) bool { return row[col] == ' ' })
}

func digitsToNumber(chars []rune) int {
	num := 0
	for _, ch := range chars {
		if ch >= '0' && ch <= '9' {
			num = num*10 + int(ch-'0')
		}
	}
	return num
}

type numberExtractor func(grid [][]rune, numRows, startCol, endCol int) []int

func extractRowNumbers(grid [][]rune, numRows, startCol, endCol int) []int {
	nums := make([]int, 0, numRows)
	for row := range numRows {
		nums = append(nums, digitsToNumber(grid[row][startCol:endCol]))
	}
	return nums
}

func extractColumnNumbers(grid [][]rune, numRows, startCol, endCol int) []int {
	nums := make([]int, 0, endCol-startCol)
	for col := startCol; col < endCol; col++ {
		colChars := make([]rune, numRows)
		for row := range numRows {
			colChars[row] = grid[row][col]
		}
		nums = append(nums, digitsToNumber(colChars))
	}
	return nums
}

func extractProblems(grid [][]rune, extract numberExtractor) ([]problem, error) {
	operatorRow := len(grid) - 1
	numRows := operatorRow
	maxColumn := len(grid[0])

	newProblem := func(startCol, endCol int) problem {
		return problem{
			op:      grid[operatorRow][startCol],
			numbers: extract(grid, numRows, startCol, endCol),
		}
	}

	var problems []problem
	col := 0
	for {
		for col < maxColumn && isSeparatorColumn(grid, col) {
			col++
		}
		if col >= maxColumn {
			break
		}
		startCol := col
		for col < maxColumn && !isSeparatorColumn(grid, col) {
			col++
		}
		if op := grid[operatorRow][startCol]; op != '+' && op != '*' {
			return nil, fmt.Errorf("homework problem at column %d needs an operator", startCol+1)
		}
		for _, cell := range grid[operatorRow][startCol+1 : col] {
			if cell != ' ' {
				return nil, fmt.Errorf("homework problem at column %d has multiple operators", startCol+1)
			}
		}
		problems = append(problems, newProblem(startCol, col))
	}

	if len(problems) == 0 {
		return nil, fmt.Errorf("empty homework")
	}
	return problems, nil
}

func parseMathHomework(r io.Reader, extract numberExtractor) ([]problem, error) {
	grid, err := readProblems(r)
	if err != nil {
		return nil, err
	}
	return extractProblems(grid, extract)
}

func totalProblems(problems []problem) int {
	total := 0
	for _, p := range problems {
		switch p.op {
		case '+':
			total += xiter.Sum(slices.Values(p.numbers))
		case '*':
			total += xiter.Product(slices.Values(p.numbers))

		default:
			panic("unknown operator")
		}
	}
	return total
}

func day06p01(r io.Reader) (string, error) {
	problems, err := parseMathHomework(r, extractRowNumbers)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(totalProblems(problems)), nil
}

func day06p02(r io.Reader) (string, error) {
	problems, err := parseMathHomework(r, extractColumnNumbers)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(totalProblems(problems)), nil
}
