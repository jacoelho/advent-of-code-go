package aoc2023

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type card struct {
	winning collections.Set[int]
	numbers []int
}

func parseCard(line []byte) (card, error) {
	header, body, ok := strings.Cut(string(line), ":")
	header, isCard := strings.CutPrefix(header, "Card ")
	winningText, numbersText, hasBar := strings.Cut(body, " | ")
	if !ok || !isCard || !hasBar {
		return card{}, fmt.Errorf("invalid card %q", line)
	}
	if id, err := strconv.Atoi(strings.TrimSpace(header)); err != nil || id < 0 {
		return card{}, fmt.Errorf("invalid card number in %q", line)
	}
	winning, err := parseNaturals(winningText)
	if err != nil {
		return card{}, fmt.Errorf("card %q: %w", line, err)
	}
	numbers, err := parseNaturals(numbersText)
	if err != nil {
		return card{}, fmt.Errorf("card %q: %w", line, err)
	}
	return card{winning: collections.NewSet(winning...), numbers: numbers}, nil
}

func (c card) matches() int {
	return xiter.FromSlice(c.numbers).Filter(func(n int) bool {
		return c.winning.Contains(n)
	}).Count()
}

func (c card) score() int {
	matches := c.matches()
	if matches == 0 {
		return 0
	}
	return 1 << (matches - 1)
}

func day04p01(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseCard)

	total := xiter.From(s.Values()).SumBy(func(c card) int {
		return c.score()
	})

	if err := s.Err(); err != nil {
		return "", err
	}

	return strconv.Itoa(total), nil
}

func day04p02(r io.Reader) (string, error) {
	s := scanner.NewScanner(r, parseCard)
	cards := slices.Collect(s.Values())

	if err := s.Err(); err != nil {
		return "", err
	}

	copies := make([]int, len(cards))
	for i := range copies {
		copies[i] = 1
	}

	for i, c := range cards {
		matches := c.matches()
		for j := 1; j <= matches && i+j < len(cards); j++ {
			copies[i+j] += copies[i]
		}
	}

	return strconv.Itoa(xiter.Sum(slices.Values(copies))), nil
}
