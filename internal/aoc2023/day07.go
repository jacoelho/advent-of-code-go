package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type handType int

const (
	highCard handType = iota
	onePair
	twoPair
	threeOfKind
	fullHouse
	fourOfKind
	fiveOfKind
)

const (
	cardOrder      = "23456789TJQKA"
	jokerCardOrder = "J23456789TQKA"
)

type hand struct {
	strengths [5]int
	bid       int
	typ       handType
}

func classifyHand(strengths [5]int, jokers int) handType {
	var counts [len(cardOrder)]int
	for _, strength := range strengths {
		counts[strength]++
	}
	// Jokers have strength 0 under the joker rule; they join the largest group.
	counts[0] -= jokers
	slices.Sort(counts[:])
	largest, second := counts[len(counts)-1]+jokers, counts[len(counts)-2]

	switch {
	case largest == 5:
		return fiveOfKind
	case largest == 4:
		return fourOfKind
	case largest == 3 && second == 2:
		return fullHouse
	case largest == 3:
		return threeOfKind
	case largest == 2 && second == 2:
		return twoPair
	case largest == 2:
		return onePair
	default:
		return highCard
	}
}

func parseHand(order string, jokerRule bool) func([]byte) (hand, error) {
	return func(line []byte) (hand, error) {
		cards, bidText, ok := strings.Cut(string(line), " ")
		if !ok || len(cards) != 5 {
			return hand{}, fmt.Errorf("invalid hand %q", line)
		}
		bid, err := strconv.Atoi(bidText)
		if err != nil || bid < 0 {
			return hand{}, fmt.Errorf("invalid bid in %q", line)
		}
		var result hand
		jokers := 0
		for i := range cards {
			strength := strings.IndexByte(order, cards[i])
			if strength < 0 {
				return hand{}, fmt.Errorf("invalid card %q in %q", cards[i], line)
			}
			result.strengths[i] = strength
			if jokerRule && cards[i] == 'J' {
				jokers++
			}
		}
		result.bid = bid
		result.typ = classifyHand(result.strengths, jokers)
		return result, nil
	}
}

func totalWinnings(r io.Reader, order string, jokerRule bool) (string, error) {
	s := scanner.NewScanner(r, parseHand(order, jokerRule))
	hands := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}

	slices.SortFunc(hands, func(a, b hand) int {
		return cmp.Or(cmp.Compare(a.typ, b.typ), slices.Compare(a.strengths[:], b.strengths[:]))
	})
	total := 0
	for rank, h := range hands {
		total += h.bid * (rank + 1)
	}
	return strconv.Itoa(total), nil
}

func day07p01(r io.Reader) (string, error) {
	return totalWinnings(r, cardOrder, false)
}

func day07p02(r io.Reader) (string, error) {
	return totalWinnings(r, jokerCardOrder, true)
}
