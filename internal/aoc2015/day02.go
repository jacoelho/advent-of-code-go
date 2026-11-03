package aoc2015

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func wrappingSupplies(r io.Reader) (paper, ribbon int, err error) {
	s := scanner.NewScanner(r, func(line []byte) ([3]int, error) {
		var dimensions [3]int
		parts := bytes.Split(line, []byte{'x'})
		if len(parts) != len(dimensions) {
			return dimensions, fmt.Errorf("invalid box dimensions %q", line)
		}
		for i, part := range parts {
			value, err := convert.ScanNumber[int](part)
			if err != nil {
				return dimensions, fmt.Errorf("invalid box dimension %q: %w", part, err)
			}
			if value <= 0 {
				return dimensions, fmt.Errorf("invalid box dimension %q", part)
			}
			dimensions[i] = value
		}
		return dimensions, nil
	})
	for dimensions := range s.Values() {
		l, w, h := dimensions[0], dimensions[1], dimensions[2]
		lw, wh, hl := l*w, w*h, h*l
		paper += 2*(lw+wh+hl) + min(lw, wh, hl)
		ribbon += 2*min(l+w, w+h, h+l) + l*w*h
	}
	return paper, ribbon, s.Err()
}

func day02p01(r io.Reader) (string, error) {
	paper, _, err := wrappingSupplies(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(paper), nil
}

func day02p02(r io.Reader) (string, error) {
	_, ribbon, err := wrappingSupplies(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(ribbon), nil
}
