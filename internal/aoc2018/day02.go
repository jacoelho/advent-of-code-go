package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func boxIDs(r io.Reader) ([]string, error) {
	s := bufio.NewScanner(r)
	var ids []string
	for s.Scan() {
		ids = append(ids, s.Text())
	}
	return ids, s.Err()
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
			if len(left) != len(right) {
				continue
			}
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
