package aoc2020

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

func parseCups(r io.Reader) ([]int, error) {
	sc := bufio.NewScanner(r)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	line := sc.Bytes()
	cups := make([]int, 0, len(line))
	for _, b := range line {
		if b >= '0' && b <= '9' {
			cups = append(cups, int(b-'0'))
		}
	}
	return cups, nil
}

// cupRing initializes the cups ring from the given cups, extending sequentially up to 'total' if needed.
func cupRing(cups []int, total int) ([]int, int) {
	next := make([]int, total+1)
	for i := 0; i < len(cups)-1; i++ {
		next[cups[i]] = cups[i+1]
	}
	last := cups[len(cups)-1]
	maxCup := xslices.Max(cups)
	if total > len(cups) {
		next[last] = maxCup + 1
		for x := maxCup + 1; x < total; x++ {
			next[x] = x + 1
		}
		next[total] = cups[0]
	} else {
		next[last] = cups[0]
	}
	return next, cups[0]
}

// simulate performs 'moves' crab-cups moves on the cups ring.
func simulate(next []int, current, moves, maxCup int) {
	for i := 0; i < moves; i++ {
		a := next[current]
		b := next[a]
		c := next[b]

		next[current] = next[c]

		dest := current - 1
		if dest == 0 {
			dest = maxCup
		}
		for dest == a || dest == b || dest == c {
			dest--
			if dest == 0 {
				dest = maxCup
			}
		}

		next[c] = next[dest]
		next[dest] = a
		current = next[current]
	}
}

func day23p01(r io.Reader) (string, error) {
	cups, err := parseCups(r)
	if err != nil {
		return "", err
	}
	if len(cups) == 0 {
		return "", nil
	}

	total := len(cups)
	next, current := cupRing(cups, total)
	simulate(next, current, 100, total)

	var b strings.Builder
	for x := next[1]; x != 1; x = next[x] {
		b.WriteString(strconv.Itoa(x))
	}
	return b.String(), nil
}

func day23p02(r io.Reader) (string, error) {
	cups, err := parseCups(r)
	if err != nil {
		return "", err
	}
	if len(cups) == 0 {
		return "", nil
	}

	const (
		total = 1_000_000
		moves = 10_000_000
	)

	next, current := cupRing(cups, total)
	simulate(next, current, moves, total)

	x := next[1]
	y := next[x]
	prod := uint64(x) * uint64(y)
	return strconv.FormatUint(prod, 10), nil
}
