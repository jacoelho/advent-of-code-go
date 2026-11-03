package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func spreadsheetRows(r io.Reader) *scanner.Scanner[[]int] {
	return scanner.NewScanner(r, func(line []byte) ([]int, error) {
		var row []int
		for field := range bytes.FieldsSeq(line) {
			value, err := convert.ScanNumber[int](field)
			if err != nil {
				return nil, err
			}
			if value <= 0 {
				return nil, fmt.Errorf("spreadsheet value must be positive, got %d", value)
			}
			row = append(row, value)
		}
		if len(row) == 0 {
			return nil, fmt.Errorf("empty spreadsheet row")
		}
		return row, nil
	})
}

func day02p01(r io.Reader) (string, error) {
	s := spreadsheetRows(r)
	total := xiter.From(s.Values()).SumBy(func(row []int) int {
		return slices.Max(row) - slices.Min(row)
	})
	if err := s.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(total), nil
}

func day02p02(r io.Reader) (string, error) {
	s := spreadsheetRows(r)
	total := 0
	for row := range s.Values() {
		quotient := 0
		for i, a := range row {
			for _, b := range row[i+1:] {
				larger, smaller := max(a, b), min(a, b)
				if larger%smaller == 0 {
					if quotient != 0 {
						return "", fmt.Errorf("multiple divisible pairs in spreadsheet row")
					}
					quotient = larger / smaller
				}
			}
		}
		if quotient == 0 {
			return "", fmt.Errorf("no divisible pair in spreadsheet row")
		}
		total += quotient
	}
	if err := s.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(total), nil
}
