package aoc2021

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func polymerInput(r io.Reader) (string, map[[2]byte]byte, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		if err := s.Err(); err != nil {
			return "", nil, err
		}
		return "", nil, fmt.Errorf("empty polymer input")
	}
	template := s.Text()
	if len(template) == 0 {
		return "", nil, fmt.Errorf("empty polymer template")
	}
	for _, element := range template {
		if element < 'A' || element > 'Z' {
			return "", nil, fmt.Errorf("invalid polymer element %q", element)
		}
	}
	rules := make(map[[2]byte]byte)
	for s.Scan() {
		if s.Text() == "" {
			continue
		}
		pair, inserted, ok := strings.Cut(s.Text(), " -> ")
		if !ok || len(pair) != 2 || len(inserted) != 1 || pair[0] < 'A' || pair[0] > 'Z' ||
			pair[1] < 'A' || pair[1] > 'Z' || inserted[0] < 'A' || inserted[0] > 'Z' {
			return "", nil, fmt.Errorf("invalid polymer rule %q", s.Text())
		}
		key := [2]byte{pair[0], pair[1]}
		if previous, exists := rules[key]; exists && previous != inserted[0] {
			return "", nil, fmt.Errorf("conflicting polymer rule for %q", pair)
		}
		rules[key] = inserted[0]
	}
	return template, rules, s.Err()
}

func addPolymerPair(pairs map[[2]byte]uint64, pair [2]byte, count uint64) error {
	if pairs[pair] > math.MaxUint64-count {
		return fmt.Errorf("polymer pair count overflows uint64")
	}
	pairs[pair] += count
	return nil
}

func polymerDifference(r io.Reader, steps int) (uint64, error) {
	template, rules, err := polymerInput(r)
	if err != nil {
		return 0, err
	}
	pairs := make(map[[2]byte]uint64)
	for i := 0; i+1 < len(template); i++ {
		pairs[[2]byte{template[i], template[i+1]}]++
	}
	for range steps {
		next := make(map[[2]byte]uint64)
		for pair, count := range pairs {
			if element, found := rules[pair]; found {
				if err := addPolymerPair(next, [2]byte{pair[0], element}, count); err != nil {
					return 0, err
				}
				if err := addPolymerPair(next, [2]byte{element, pair[1]}, count); err != nil {
					return 0, err
				}
			} else if err := addPolymerPair(next, pair, count); err != nil {
				return 0, err
			}
		}
		pairs = next
	}
	var counts [26]uint64
	counts[template[len(template)-1]-'A'] = 1
	for pair, count := range pairs {
		index := pair[0] - 'A'
		if counts[index] > math.MaxUint64-count {
			return 0, fmt.Errorf("polymer element count overflows uint64")
		}
		counts[index] += count
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
