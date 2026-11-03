package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"
)

type syntaxLine struct {
	corrupt   uint64
	unmatched []int
}

func syntaxLines(r io.Reader) ([]syntaxLine, error) {
	s := bufio.NewScanner(r)
	var lines []syntaxLine
	for s.Scan() {
		var line syntaxLine
		for _, delimiter := range s.Bytes() {
			if opening := strings.IndexByte("([{<", delimiter); opening >= 0 {
				if line.corrupt == 0 {
					line.unmatched = append(line.unmatched, opening)
				}
				continue
			}
			closing := strings.IndexByte(")]}>", delimiter)
			if closing < 0 {
				return nil, fmt.Errorf("invalid syntax delimiter %q", delimiter)
			}
			if line.corrupt != 0 {
				continue
			}
			if len(line.unmatched) == 0 || line.unmatched[len(line.unmatched)-1] != closing {
				line.corrupt = [4]uint64{3, 57, 1197, 25137}[closing]
			} else {
				line.unmatched = line.unmatched[:len(line.unmatched)-1]
			}
		}
		lines = append(lines, line)
	}
	return lines, s.Err()
}

func day10p01(r io.Reader) (string, error) {
	lines, err := syntaxLines(r)
	if err != nil {
		return "", err
	}
	var total uint64
	for _, line := range lines {
		total += line.corrupt
	}
	return strconv.FormatUint(total, 10), nil
}

func day10p02(r io.Reader) (string, error) {
	lines, err := syntaxLines(r)
	if err != nil {
		return "", err
	}
	var scores []uint64
	for _, line := range lines {
		if line.corrupt != 0 || len(line.unmatched) == 0 {
			continue
		}
		var score uint64
		for i := len(line.unmatched) - 1; i >= 0; i-- {
			value := uint64(line.unmatched[i] + 1)
			if score > (math.MaxUint64-value)/5 {
				return "", fmt.Errorf("completion score overflows uint64")
			}
			score = score*5 + value
		}
		scores = append(scores, score)
	}
	if len(scores) == 0 || len(scores)%2 == 0 {
		return "", fmt.Errorf("completion median requires an odd number of incomplete lines")
	}
	slices.Sort(scores)
	return strconv.FormatUint(scores[len(scores)/2], 10), nil
}
