package aoc2021

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func diagnosticReport(r io.Reader) ([][]byte, error) {
	report, err := grid.ReadBytes(r, "01")
	if err != nil {
		return nil, err
	}
	if report.Width() > 64 {
		return nil, fmt.Errorf("diagnostic rows have %d bits, want at most 64", report.Width())
	}
	return report.Rows(), nil
}

func diagnosticProduct(a, b uint64) (string, error) {
	product, ok := xmath.CheckedMul(a, b)
	if !ok {
		return "", fmt.Errorf("diagnostic product overflows uint64")
	}
	return strconv.FormatUint(product, 10), nil
}

func day03p01(r io.Reader) (string, error) {
	rows, err := diagnosticReport(r)
	if err != nil {
		return "", err
	}

	var gamma uint64
	for bit := range len(rows[0]) {
		ones := 0
		for _, row := range rows {
			if row[bit] == '1' {
				ones++
			}
		}
		gamma <<= 1
		if ones >= len(rows)-ones {
			gamma |= 1
		}
	}
	mask := uint64(1)<<len(rows[0]) - 1
	return diagnosticProduct(gamma, gamma^mask)
}

type diagnosticGas byte

const (
	oxygen diagnosticGas = iota
	carbonDioxide
)

func gasRating(report [][]byte, gas diagnosticGas) uint64 {
	rows := slices.Clone(report)
	for bit := 0; bit < len(report[0]) && len(rows) > 1; bit++ {
		ones := 0
		for _, row := range rows {
			if row[bit] == '1' {
				ones++
			}
		}
		zeros := len(rows) - ones
		wanted := byte('0')
		if zeros == 0 || ones != 0 && (gas == oxygen && ones >= zeros || gas == carbonDioxide && ones < zeros) {
			wanted = '1'
		}
		kept := rows[:0]
		for _, row := range rows {
			if row[bit] == wanted {
				kept = append(kept, row)
			}
		}
		rows = kept
	}
	value, _ := strconv.ParseUint(string(rows[0]), 2, 64)
	return value
}

func day03p02(r io.Reader) (string, error) {
	rows, err := diagnosticReport(r)
	if err != nil {
		return "", err
	}
	return diagnosticProduct(gasRating(rows, oxygen), gasRating(rows, carbonDioxide))
}
