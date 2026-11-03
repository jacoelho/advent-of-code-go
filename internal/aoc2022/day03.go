package aoc2022

import (
	"fmt"
	"io"
	"math/bits"
	"strconv"
	"strings"
)

func itemTypes(items string) (uint64, error) {
	var mask uint64
	for i := range len(items) {
		c := items[i]
		var priority byte
		switch {
		case c >= 'a' && c <= 'z':
			priority = c - 'a'
		case c >= 'A' && c <= 'Z':
			priority = c - 'A' + 26
		default:
			return 0, fmt.Errorf("invalid item %q", c)
		}
		mask |= 1 << priority
	}
	return mask, nil
}

func rucksackPriority(r io.Reader, badges bool) (string, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return "", err
	}
	lines := strings.Split(input, "\n")
	total := 0
	if badges && len(lines)%3 != 0 {
		return "", fmt.Errorf("incomplete badge group")
	}
	for i := 0; i < len(lines); {
		var common uint64
		if badges {
			common = (1 << 52) - 1
			for _, line := range lines[i : i+3] {
				if len(line)%2 != 0 {
					return "", fmt.Errorf("odd compartment length")
				}
				mask, err := itemTypes(line)
				if err != nil {
					return "", err
				}
				common &= mask
			}
			i += 3
		} else {
			line := lines[i]
			if len(line)%2 != 0 {
				return "", fmt.Errorf("odd compartment length")
			}
			left, err := itemTypes(line[:len(line)/2])
			if err != nil {
				return "", err
			}
			right, err := itemTypes(line[len(line)/2:])
			if err != nil {
				return "", err
			}
			common = left & right
			i++
		}
		if bits.OnesCount64(common) != 1 {
			return "", fmt.Errorf("expected one shared item")
		}
		total += bits.TrailingZeros64(common) + 1
	}
	return strconv.Itoa(total), nil
}

func day03p01(r io.Reader) (string, error) {
	return rucksackPriority(r, false)
}

func day03p02(r io.Reader) (string, error) {
	return rucksackPriority(r, true)
}
