package aoc2021

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const polymerElements = 26

// polymerRules holds the inserted element plus one for each pair, or zero for no rule.
type polymerRules [polymerElements][polymerElements]byte

func polymerElement(c byte) (byte, bool) {
	return c - 'A', c >= 'A' && c <= 'Z'
}

func polymerInput(r io.Reader) ([]byte, *polymerRules, error) {
	sections, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(sections) == 0 {
		return nil, nil, fmt.Errorf("empty polymer input")
	}
	if len(sections[0]) != 1 {
		return nil, nil, fmt.Errorf("expected one polymer template line followed by a blank line")
	}
	if len(sections) > 2 || len(sections) == 2 && len(sections[1]) == 0 {
		return nil, nil, fmt.Errorf("blank line among polymer rules")
	}
	template := []byte(sections[0][0])
	for i, c := range template {
		element, ok := polymerElement(c)
		if !ok {
			return nil, nil, fmt.Errorf("invalid polymer element %q", c)
		}
		template[i] = element
	}
	var ruleLines []string
	if len(sections) == 2 {
		ruleLines = sections[1]
	}
	rules := new(polymerRules)
	for _, line := range ruleLines {
		pair, inserted, ok := strings.Cut(line, " -> ")
		if !ok || len(pair) != 2 || len(inserted) != 1 {
			return nil, nil, fmt.Errorf("invalid polymer rule %q", line)
		}
		left, leftOK := polymerElement(pair[0])
		right, rightOK := polymerElement(pair[1])
		element, elementOK := polymerElement(inserted[0])
		if !leftOK || !rightOK || !elementOK {
			return nil, nil, fmt.Errorf("invalid polymer rule %q", line)
		}
		if previous := rules[left][right]; previous != 0 && previous != element+1 {
			return nil, nil, fmt.Errorf("conflicting polymer rule for %q", pair)
		}
		rules[left][right] = element + 1
	}
	return template, rules, nil
}

func polymerDifference(r io.Reader, steps int) (uint64, error) {
	template, rules, err := polymerInput(r)
	if err != nil {
		return 0, err
	}
	// The element total, (len(template)-1)<<steps + 1, bounds every count.
	if steps >= 64 || uint64(len(template)-1) > (math.MaxUint64-1)>>steps {
		return 0, fmt.Errorf("polymer of %d elements overflows after %d steps", len(template), steps)
	}
	var pairs [polymerElements][polymerElements]uint64
	for i := 0; i+1 < len(template); i++ {
		pairs[template[i]][template[i+1]]++
	}
	for range steps {
		var next [polymerElements][polymerElements]uint64
		for left := range polymerElements {
			for right, count := range pairs[left] {
				if count == 0 {
					continue
				}
				if inserted := rules[left][right]; inserted != 0 {
					element := inserted - 1
					next[left][element] += count
					next[element][right] += count
				} else {
					next[left][right] += count
				}
			}
		}
		pairs = next
	}
	var counts [polymerElements]uint64
	counts[template[len(template)-1]] = 1
	for left := range polymerElements {
		for _, count := range pairs[left] {
			counts[left] += count
		}
	}
	least, most := uint64(math.MaxUint64), uint64(0)
	for _, count := range counts {
		if count != 0 {
			least, most = min(least, count), max(most, count)
		}
	}
	return most - least, nil
}

func day14p01(r io.Reader) (string, error) {
	difference, err := polymerDifference(r, 10)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(difference, 10), nil
}

func day14p02(r io.Reader) (string, error) {
	difference, err := polymerDifference(r, 40)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(difference, 10), nil
}
