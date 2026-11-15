package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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
	if len(ratings) == 0 {
		return nil, fmt.Errorf("no adapters")
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
	ways := make([]int64, ratings[len(ratings)-1]+1)
	ways[0] = 1
	for _, rating := range ratings {
		for gap := 1; gap <= min(3, rating); gap++ {
			sum, ok := xmath.CheckedAdd(ways[rating], ways[rating-gap])
			if !ok {
				return "", fmt.Errorf("adapter arrangements overflow int64 at rating %d", rating)
			}
			ways[rating] = sum
		}
	}
	return strconv.FormatInt(ways[len(ways)-1], 10), nil
}
