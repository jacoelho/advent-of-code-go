package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type auntComparison uint8

const (
	exactAuntCounts auntComparison = iota
	rangedAuntCounts
)

func identifyAunt(r io.Reader, comparison auntComparison) (string, error) {
	expected := map[string]int{
		"children": 3, "cats": 7, "samoyeds": 2, "pomeranians": 3, "akitas": 0,
		"vizslas": 0, "goldfish": 5, "trees": 3, "cars": 2, "perfumes": 1,
	}
	seen := make(map[int]bool)
	answer := 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		header, details, ok := strings.Cut(s.Text(), ": ")
		idText, prefixOK := strings.CutPrefix(header, "Sue ")
		id, err := strconv.Atoi(idText)
		if !ok || !prefixOK || err != nil || id <= 0 || seen[id] {
			return "", fmt.Errorf("invalid or duplicate aunt %q", s.Text())
		}
		seen[id] = true
		properties := make(map[string]bool)
		matches := true
		for detail := range strings.SplitSeq(details, ", ") {
			name, countText, ok := strings.Cut(detail, ": ")
			count, err := strconv.Atoi(countText)
			want, known := expected[name]
			if !ok || err != nil || count < 0 || !known || properties[name] {
				return "", fmt.Errorf("invalid property for Sue %d: %q", id, detail)
			}
			properties[name] = true
			switch {
			case comparison == rangedAuntCounts && (name == "cats" || name == "trees"):
				matches = matches && count > want
			case comparison == rangedAuntCounts && (name == "pomeranians" || name == "goldfish"):
				matches = matches && count < want
			default:
				matches = matches && count == want
			}
		}
		if matches {
			if answer != 0 {
				return "", fmt.Errorf("multiple matching aunts: %d and %d", answer, id)
			}
			answer = id
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if answer == 0 {
		return "", fmt.Errorf("no matching aunt")
	}
	return strconv.Itoa(answer), nil
}

func day16p01(r io.Reader) (string, error) {
	return identifyAunt(r, exactAuntCounts)
}

func day16p02(r io.Reader) (string, error) {
	return identifyAunt(r, rangedAuntCounts)
}
