package aoc2024

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const (
	day25Width  = 5
	day25Height = 7
	day25Row    = 1<<day25Width - 1
)

// parseSchematics returns each 5×7 schematic as a bitmask with bit y*5+x set
// for every '#'.
func parseSchematics(r io.Reader) (locks, keys []uint64, err error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(blocks) == 0 {
		return nil, nil, fmt.Errorf("no schematics")
	}
	for n, rows := range blocks {
		if len(rows) != day25Height {
			return nil, nil, fmt.Errorf("schematic %d has %d rows, want %d", n+1, len(rows), day25Height)
		}
		var mask uint64
		for y, row := range rows {
			if len(row) != day25Width {
				return nil, nil, fmt.Errorf("schematic %d row %d has %d cells, want %d", n+1, y+1, len(row), day25Width)
			}
			for x, c := range row {
				switch c {
				case '#':
					mask |= 1 << (y*day25Width + x)
				case '.':
				default:
					return nil, nil, fmt.Errorf("schematic %d row %d has invalid character %q", n+1, y+1, c)
				}
			}
		}
		top := mask & day25Row
		bottom := mask >> ((day25Height - 1) * day25Width) & day25Row
		switch {
		case top == day25Row && bottom == 0:
			locks = append(locks, mask)
		case top == 0 && bottom == day25Row:
			keys = append(keys, mask)
		default:
			return nil, nil, fmt.Errorf("schematic %d is neither a lock nor a key", n+1)
		}
	}
	return locks, keys, nil
}

func day25p01(r io.Reader) (string, error) {
	locks, keys, err := parseSchematics(r)
	if err != nil {
		return "", err
	}
	var total int
	for _, lock := range locks {
		for _, key := range keys {
			if lock&key == 0 {
				total++
			}
		}
	}
	return strconv.Itoa(total), nil
}
