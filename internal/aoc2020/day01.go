package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseExpenses(r io.Reader) ([]int, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	expenses := slices.Collect(s.Values())
	slices.Sort(expenses)
	return expenses, s.Err()
}

func expensePair(expenses []int, target int) (int, bool) {
	for left, right := 0, len(expenses)-1; left < right; {
		sum := expenses[left] + expenses[right]
		switch {
		case sum == target:
			return expenses[left] * expenses[right], true
		case sum < target:
			left++
		default:
			right--
		}
	}
	return 0, false
}

func day01p01(r io.Reader) (string, error) {
	expenses, err := parseExpenses(r)
	if err != nil {
		return "", err
	}
	if product, ok := expensePair(expenses, 2020); ok {
		return strconv.Itoa(product), nil
	}
	return "", fmt.Errorf("no two expense entries sum to 2020")
}

func day01p02(r io.Reader) (string, error) {
	expenses, err := parseExpenses(r)
	if err != nil {
		return "", err
	}
	for i, expense := range expenses {
		if product, ok := expensePair(expenses[i+1:], 2020-expense); ok {
			return strconv.Itoa(expense * product), nil
		}
	}
	return "", fmt.Errorf("no three expense entries sum to 2020")
}
