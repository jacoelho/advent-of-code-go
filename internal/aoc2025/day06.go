package aoc2025

import (
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type problem struct {
	op      rune
	numbers []int
}

func readProblems(r io.Reader) ([][]rune, error) {
	s := scanner.NewScanner(r, func(b []byte) ([]rune, error) {
		return []rune(string(b)), nil
	})
	return slices.Collect(s.Values()), s.Err()
}

func isSeparatorColumn(grid [][]rune, col int) bool {
	return xslices.Every(func(row []rune) bool { return row[col] == ' ' }, grid)
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

func extractProblems(grid [][]rune, extract numberExtractor) []problem {
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
		problems = append(problems, newProblem(startCol, col))
	}

	return problems
}

func parseMathHomework(r io.Reader, extract numberExtractor) ([]problem, error) {
	grid, err := readProblems(r)
	if err != nil {
		return nil, err
	}
	return extractProblems(grid, extract), nil
}

func totalProblems(problems []problem) int {
	total := 0
	for _, p := range problems {
		switch p.op {
		case '+':
			total += xslices.Sum(p.numbers)
		case '*':
			total += xslices.Product(p.numbers)
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
