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
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type card struct {
	winning collections.Set[int]
	numbers []int
}

func parseCard(line []byte) (card, error) {
	parts := strings.Split(string(line), ": ")
	numberParts := strings.Split(parts[1], " | ")

	winning := collections.NewSetFromIter(xiter.Map(aoc.MustAtoi, strings.FieldsSeq(numberParts[0])))
	numbers := slices.Collect(xiter.Map(aoc.MustAtoi, strings.FieldsSeq(numberParts[1])))

	return card{winning: winning, numbers: numbers}, nil
}

func (c card) matches() int {
	return xslices.CountFunc(func(n int) bool {
		return c.winning.Contains(n)
	}, c.numbers)
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

	total := xiter.Sum(xiter.Map(func(c card) int {
		return c.score()
	}, s.Values()))

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

	return strconv.Itoa(xslices.Sum(copies)), nil
}
