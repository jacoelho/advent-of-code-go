package aoc2023

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type card struct {
	winning collections.Set[int]
	numbers []int
}

func parseCard(line []byte) (card, error) {
	parts := strings.Split(string(line), ": ")
	numberParts := strings.Split(parts[1], " | ")

	winning := collections.NewSetFromIter(xiter.From(strings.FieldsSeq(numberParts[0])).Map(aoc.MustAtoi).All())
	numbers := xiter.From(strings.FieldsSeq(numberParts[1])).Map(aoc.MustAtoi).Collect()

	return card{winning: winning, numbers: numbers}, nil
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
