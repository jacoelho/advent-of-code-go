package aoc2022

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func rearrangeCrates(r io.Reader, preserve bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	drawing, procedure, ok := strings.Cut(input, "\n\n")
	if !ok {
		return "", fmt.Errorf("missing crate procedure")
	}
	rows := strings.Split(drawing, "\n")
	if len(rows) > 10_000 {
		return "", fmt.Errorf("crate drawing exceeds ten thousand rows")
	}
	labels := rows[len(rows)-1]
	var columns []int
	indices := make(map[int]int)
	for i := 0; i < len(labels); {
		if labels[i] == ' ' {
			i++
			continue
		}
		start := i
		for i < len(labels) && labels[i] >= '0' && labels[i] <= '9' {
			i++
		}
		if start == i {
			return "", fmt.Errorf("invalid stack label")
		}
		label, err := convert.ScanNumber[int]([]byte(labels[start:i]))
		if err != nil {
			return "", err
		}
		if _, exists := indices[label]; exists || label <= 0 {
			return "", fmt.Errorf("invalid or duplicate stack label")
		}
		indices[label] = len(columns)
		columns = append(columns, i-1)
		if len(columns) > 128 {
			return "", fmt.Errorf("too many crate stacks")
		}
	}
	if len(columns) == 0 {
		return "", fmt.Errorf("no stacks")
	}
	stacks := make([][]byte, len(columns))
	for row := len(rows) - 2; row >= 0; row-- {
		line := rows[row]
		used := make([]bool, len(line))
		for stack, column := range columns {
			if column >= len(line) || line[column] == ' ' {
				continue
			}
			if column == 0 || column+1 >= len(line) || line[column-1] != '[' || line[column+1] != ']' || line[column] < 'A' || line[column] > 'Z' {
				return "", fmt.Errorf("invalid crate drawing")
			}
			for j := column - 1; j <= column+1; j++ {
				if used[j] {
					return "", fmt.Errorf("overlapping crate columns")
				}
				used[j] = true
			}
			stacks[stack] = append(stacks[stack], line[column])
		}
		for i, c := range []byte(line) {
			if c != ' ' && !used[i] {
				return "", fmt.Errorf("unlabelled crate")
			}
		}
	}
	moved := 0
	for line := range strings.SplitSeq(procedure, "\n") {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 6 || fields[0] != "move" || fields[2] != "from" || fields[4] != "to" {
			return "", fmt.Errorf("invalid move %q", line)
		}
		var values [3]int
		for i := range values {
			value, err := convert.ScanNumber[int]([]byte(fields[2*i+1]))
			if err != nil {
				return "", err
			}
			values[i] = value
		}
		count := values[0]
		from, fromOK := indices[values[1]]
		to, toOK := indices[values[2]]
		if !fromOK || !toOK || count < 0 || count > len(stacks[from]) {
			return "", fmt.Errorf("invalid crate move")
		}
		if from == to || count == 0 {
			continue
		}
		if count > 10_000_000-moved {
			return "", fmt.Errorf("crate move work limit exceeded")
		}
		moved += count
		cut := len(stacks[from]) - count
		block := stacks[from][cut:]
		if !preserve {
			slices.Reverse(block)
		}
		stacks[to] = append(stacks[to], block...)
		stacks[from] = stacks[from][:cut]
	}
	result := make([]byte, len(stacks))
	for i, stack := range stacks {
		if len(stack) == 0 {
			return "", fmt.Errorf("stack has no top crate")
		}
		result[i] = stack[len(stack)-1]
	}
	return string(result), nil
}

func day05p01(r io.Reader) (string, error) {
	return rearrangeCrates(r, false)
}

func day05p02(r io.Reader) (string, error) {
	return rearrangeCrates(r, true)
}
