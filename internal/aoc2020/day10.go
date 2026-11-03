package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseAdapters(r io.Reader) ([]int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int, error) {
		rating, err := strconv.Atoi(string(line))
		if err != nil || rating <= 0 {
			return 0, fmt.Errorf("invalid adapter rating %q", line)
		}
		return rating, nil
	})
	ratings := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	slices.Sort(ratings)
	previous := 0
	for _, rating := range ratings {
		if rating <= previous || rating-previous > 3 {
			return nil, fmt.Errorf("cannot connect adapter %d after %d", rating, previous)
		}
		previous = rating
	}
	return ratings, nil
}

func day10p01(r io.Reader) (string, error) {
	ratings, err := parseAdapters(r)
	if err != nil {
		return "", err
	}
	ones, threes, previous := 0, 1, 0
	for _, rating := range ratings {
		switch rating - previous {
		case 1:
			ones++
		case 3:
			threes++
		}
		previous = rating
	}
	return strconv.Itoa(ones * threes), nil
}

func day10p02(r io.Reader) (string, error) {
	ratings, err := parseAdapters(r)
	if err != nil {
		return "", err
	}
	ways := map[int]int64{0: 1}
	last := 0
	for _, rating := range ratings {
		ways[rating] = ways[rating-1] + ways[rating-2] + ways[rating-3]
		last = rating
	}
	return strconv.FormatInt(ways[last], 10), nil
}
