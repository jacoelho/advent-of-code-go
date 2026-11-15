package aoc2020

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// parseCups returns cup labels, which form a permutation of 1..n with n >= 5 so
// that each move can pick up three cups and still find a destination.
func parseCups(r io.Reader) ([]int32, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	line := strings.TrimSpace(string(data))
	if len(line) < 5 {
		return nil, fmt.Errorf("need at least 5 cups, got %d", len(line))
	}
	seen := make([]bool, len(line)+1)
	cups := make([]int32, len(line))
	for i, c := range []byte(line) {
		label := int(c - '0')
		if c < '1' || c > '9' || label > len(line) {
			return nil, fmt.Errorf("invalid cup label %q: labels must be 1..%d", c, len(line))
		}
		if seen[label] {
			return nil, fmt.Errorf("duplicate cup label %d", label)
		}
		seen[label] = true
		cups[i] = int32(label)
	}
	return cups, nil
}

// playCups plays moves on a ring of total cups, the given cups followed by
// labels len(cups)+1..total, and returns the successor of each cup label.
func playCups(cups []int32, total int32, moves int) []int32 {
	next := make([]int32, total+1)
	last := cups[len(cups)-1]
	for i, cup := range cups[:len(cups)-1] {
		next[cup] = cups[i+1]
	}
	for cup := int32(len(cups)) + 1; cup <= total; cup++ {
		next[last] = cup
		last = cup
	}
	next[last] = cups[0]

	current := cups[0]
	for range moves {
		a := next[current]
		b := next[a]
		c := next[b]
		next[current] = next[c]
		destination := current
		for {
			destination--
			if destination == 0 {
				destination = total
			}
			if destination != a && destination != b && destination != c {
				break
			}
		}
		next[c] = next[destination]
		next[destination] = a
		current = next[current]
	}
	return next
}

func day23p01(r io.Reader) (string, error) {
	cups, err := parseCups(r)
	if err != nil {
		return "", err
	}
	next := playCups(cups, int32(len(cups)), 100)
	var b strings.Builder
	for cup := next[1]; cup != 1; cup = next[cup] {
		b.WriteByte(byte('0' + cup))
	}
	return b.String(), nil
}

func day23p02(r io.Reader) (string, error) {
	cups, err := parseCups(r)
	if err != nil {
		return "", err
	}
	next := playCups(cups, 1_000_000, 10_000_000)
	first := next[1]
	return strconv.FormatInt(int64(first)*int64(next[first]), 10), nil
}
