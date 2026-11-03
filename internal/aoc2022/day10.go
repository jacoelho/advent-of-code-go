package aoc2022

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

func registerCycles(r io.Reader) ([240]int64, int, error) {
	var values [240]int64
	input, err := puzzleInput(r)
	if err != nil {
		return values, 0, err
	}
	x, cycle := int64(1), 0
	for line := range strings.SplitSeq(input, "\n") {
		duration, delta := 1, int64(0)
		if line != "noop" {
			operand, ok := strings.CutPrefix(line, "addx ")
			if !ok {
				return values, 0, fmt.Errorf("invalid CPU instruction")
			}
			delta, err = convert.ScanNumber[int64]([]byte(operand))
			if err != nil {
				return values, 0, err
			}
			if delta < -1_000_000 || delta > 1_000_000 {
				return values, 0, fmt.Errorf("CPU operand exceeds one million")
			}
			duration = 2
		}
		if cycle > 1_000_000-duration {
			return values, 0, fmt.Errorf("CPU exceeds one million cycles")
		}
		for range duration {
			if cycle < len(values) {
				values[cycle] = x
			}
			cycle++
		}
		x += delta
	}
	return values, cycle, nil
}

func crtDisplay(r io.Reader) ([]string, error) {
	values, cycles, err := registerCycles(r)
	if err != nil {
		return nil, err
	}
	if cycles < len(values) {
		return nil, fmt.Errorf("CRT requires 240 cycles")
	}
	rows := make([]string, 6)
	for y := range rows {
		row := make([]byte, 40)
		for x := range row {
			row[x] = '.'
			if xmath.Abs(int64(x)-values[y*40+x]) <= 1 {
				row[x] = '#'
			}
		}
		rows[y] = string(row)
	}
	return rows, nil
}

func day10p01(r io.Reader) (string, error) {
	values, cycles, err := registerCycles(r)
	if err != nil {
		return "", err
	}
	if cycles < 220 {
		return "", fmt.Errorf("signal strength requires 220 cycles")
	}
	var total int64
	for cycle := 20; cycle <= 220; cycle += 40 {
		total += int64(cycle) * values[cycle-1]
	}
	return strconv.FormatInt(total, 10), nil
}

func day10p02(r io.Reader) (string, error) {
	rows, err := crtDisplay(r)
	if err != nil {
		return "", err
	}
	return aoc.Decode6RowDisplay(rows, 5)
}
