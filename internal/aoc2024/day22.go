package aoc2024

import (
	"fmt"
	"io"
	"iter"
	"maps"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/conc"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func parseSecretNumbers(r io.Reader) ([]int, error) {
	s := scanner.NewScanner[int](r, convert.ScanNumber)
	return slices.Collect(s.Values()), s.Err()
}

func secretNumberIter(n int) iter.Seq[int] {
	return xiter.Apply(n, func(n int) int {
		n = (n ^ (n << 6)) & 16777215
		n = (n ^ (n >> 5)) & 16777215
		n = (n ^ (n << 11)) & 16777215
		return n
	})
}

func secretNumberIterationN(n int) int {
	v, ok := xiter.From(secretNumberIter(n)).Skip(2000).First()
	if !ok {
		panic("next not found")
	}
	return v
}

func day22p01(r io.Reader) (string, error) {
	secretNumbers, err := parseSecretNumbers(r)
	if err != nil {
		return "", err
	}

	var total int
	for _, secretNumber := range secretNumbers {
		total += secretNumberIterationN(secretNumber)
	}
	return strconv.Itoa(total), nil
}

func sellSequences(n int) map[[4]int]int {
	prices := xiter.From(secretNumberIter(n)).Take(2001).Map(func(v int) int { return v % 10 })
	total := make(map[[4]int]int)
	var changes [4]int
	previous := n % 10
	for i, price := range xiter.Enumerate(prices.All()) {
		if i == 0 {
			continue
		}
		changes = [4]int{changes[1], changes[2], changes[3], price - previous}
		previous = price
		if i < len(changes) {
			continue
		}
		if _, seen := total[changes]; !seen {
			total[changes] = price
		}
	}
	return total
}

func day22p02(r io.Reader) (string, error) {
	secretNumbers, err := parseSecretNumbers(r)
	if err != nil {
		return "", err
	}

	total := make(map[[4]int]int)
	for _, m := range conc.Map(sellSequences, secretNumbers) {
		for k, v := range m {
			total[k] += v
		}
	}
	maxBananas, ok := xiter.Max(maps.Values(total))
	if !ok {
		return "", fmt.Errorf("no secret numbers")
	}
	return strconv.Itoa(maxBananas), nil
}
