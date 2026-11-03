package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
)

func diagnosticReport(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	var rows []string
	for s.Scan() {
		row := s.Text()
		if len(row) == 0 || len(row) > 64 || len(rows) > 0 && len(row) != len(rows[0]) {
			return nil, fmt.Errorf("invalid diagnostic row width")
		}
		for _, bit := range row {
			if bit != '0' && bit != '1' {
				return nil, fmt.Errorf("invalid diagnostic bit %q", bit)
			}
		}
		rows = append(rows, row)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty diagnostic report")
	}
	return rows, nil
}

func diagnosticProduct(a, b uint64) (string, error) {
	if b != 0 && a > math.MaxUint64/b {
		return "", fmt.Errorf("diagnostic product overflows uint64")
	}
	return strconv.FormatUint(a*b, 10), nil
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

func gasRating(report []string, gas diagnosticGas) uint64 {
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
	value, _ := strconv.ParseUint(rows[0], 2, 64)
	return value
}

func day03p02(r io.Reader) (string, error) {
	rows, err := diagnosticReport(r)
	if err != nil {
		return "", err
	}
	return diagnosticProduct(gasRating(rows, oxygen), gasRating(rows, carbonDioxide))
}
