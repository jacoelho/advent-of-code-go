package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
)

func seatingHappiness(r io.Reader, extraGuests int) (string, error) {
	pattern := regexp.MustCompile(`^(\w+) would (gain|lose) (\d+) happiness units by sitting next to (\w+)\.$`)
	guests := make(map[string]int)
	pairs := make(map[[2]int]int64)
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := pattern.FindStringSubmatch(s.Text())
		if match == nil || match[1] == match[4] {
			return "", fmt.Errorf("invalid happiness statement %q", s.Text())
		}
		value, err := strconv.ParseInt(match[3], 10, 32)
		if err != nil {
			return "", fmt.Errorf("invalid happiness value: %w", err)
		}
		if match[2] == "lose" {
			value = -value
		}
		for _, name := range []string{match[1], match[4]} {
			if _, ok := guests[name]; !ok {
				guests[name] = len(guests)
			}
		}
		pair := [2]int{guests[match[1]], guests[match[4]]}
		if _, ok := pairs[pair]; ok {
			return "", fmt.Errorf("duplicate happiness statement for %s and %s", match[1], match[4])
		}
		pairs[pair] = value
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(guests) < 2 || len(pairs) != len(guests)*(len(guests)-1) {
		return "", fmt.Errorf("incomplete happiness statements")
	}
	n := len(guests) + extraGuests
	if n > 20 {
		return "", fmt.Errorf("seating search supports at most 20 guests")
	}
	weights := make([][]int64, n)
	for i := range n {
		weights[i] = make([]int64, n)
		for j := range n {
			weights[i][j] = pairs[[2]int{i, j}] + pairs[[2]int{j, i}]
		}
	}
	// Fix guest zero to remove rotations; each state ends at one seated guest.
	remaining := n - 1
	states := 1 << remaining
	best := make([]int64, states*remaining)
	for i := range best {
		best[i] = math.MinInt64
	}
	for last := range remaining {
		best[(1<<last)*remaining+last] = weights[0][last+1]
	}
	for mask := 1; mask < states; mask++ {
		for last := range remaining {
			value := best[mask*remaining+last]
			if value == math.MinInt64 {
				continue
			}
			for next := range remaining {
				if mask&(1<<next) == 0 {
					index := (mask|1<<next)*remaining + next
					best[index] = max(best[index], value+weights[last+1][next+1])
				}
			}
		}
	}
	answer := int64(math.MinInt64)
	for last := range remaining {
		answer = max(answer, best[(states-1)*remaining+last]+weights[last+1][0])
	}
	return strconv.FormatInt(answer, 10), nil
}

func day13p01(r io.Reader) (string, error) {
	return seatingHappiness(r, 0)
}

func day13p02(r io.Reader) (string, error) {
	return seatingHappiness(r, 1)
}
