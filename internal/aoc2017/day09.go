package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func groupStreamTotals(r io.Reader) (score, garbage int, err error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}

	depth, inGarbage, canceled, started := 0, false, false, false
	for _, symbol := range bytes.TrimSpace(input) {
		if canceled {
			canceled = false
			continue
		}
		if inGarbage {
			switch symbol {
			case '!':
				canceled = true
			case '>':
				inGarbage = false
			default:
				garbage++
			}
			continue
		}
		if depth == 0 {
			if started || symbol != '{' && symbol != '<' {
				return 0, 0, fmt.Errorf("group stream must be one group or garbage, got %q", symbol)
			}
			started = true
		}
		switch symbol {
		case '{':
			depth++
			score += depth
		case '}':
			depth--
		case '<':
			inGarbage = true
		case ',':
		default:
			return 0, 0, fmt.Errorf("invalid group-stream symbol %q", symbol)
		}
	}
	if !started || depth != 0 || inGarbage || canceled {
		return 0, 0, fmt.Errorf("unfinished group stream")
	}

	return score, garbage, nil
}

func day09p01(r io.Reader) (string, error) {
	score, _, err := groupStreamTotals(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(score), nil
}

func day09p02(r io.Reader) (string, error) {
	_, garbage, err := groupStreamTotals(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(garbage), nil
}
