package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func boxIDs(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	var ids []string
	for s.Scan() {
		id := s.Text()
		if len(id) == 0 || strings.IndexFunc(id, func(r rune) bool { return r < 'a' || r > 'z' }) >= 0 {
			return nil, fmt.Errorf("invalid box ID %q", id)
		}
		if len(ids) > 0 && len(id) != len(ids[0]) {
			return nil, fmt.Errorf("box ID %q differs in length from %q", id, ids[0])
		}
		ids = append(ids, id)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("no box IDs")
	}
	return ids, nil
}

func day02p01(r io.Reader) (string, error) {
	ids, err := boxIDs(r)
	if err != nil {
		return "", err
	}

	var twice, thrice int
	for _, id := range ids {
		var counts [256]int
		for i := range len(id) {
			counts[id[i]]++
		}
		var hasTwo, hasThree bool
		for _, count := range counts {
			hasTwo = hasTwo || count == 2
			hasThree = hasThree || count == 3
		}
		if hasTwo {
			twice++
		}
		if hasThree {
			thrice++
		}
	}
	return strconv.Itoa(twice * thrice), nil
}

func day02p02(r io.Reader) (string, error) {
	ids, err := boxIDs(r)
	if err != nil {
		return "", err
	}

	for i, left := range ids {
		for _, right := range ids[i+1:] {
			difference := -1
			for j := range len(left) {
				if left[j] == right[j] {
					continue
				}
				if difference != -1 {
					difference = -1
					break
				}
				difference = j
			}
			if difference != -1 {
				return left[:difference] + left[difference+1:], nil
			}
		}
	}
	return "", fmt.Errorf("no box IDs differ at exactly one position")
}
