package aoc2023

import (
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
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

type hand struct {
	cards string
	bid   int
	typ   handType
}

var cardValues = map[byte]int{
	'A': 14, 'K': 13, 'Q': 12, 'J': 11, 'T': 10,
	'9': 9, '8': 8, '7': 7, '6': 6, '5': 5, '4': 4, '3': 3, '2': 2,
}

func cardValue(card byte, jokerRule bool) int {
	if jokerRule && card == 'J' {
		return 1
	}
	return cardValues[card]
}

func classifyHand(cards string, jokerRule bool) handType {
	jokers := 0
	filtered := cards

	if jokerRule {
		jokers = strings.Count(cards, "J")
		filtered = strings.ReplaceAll(cards, "J", "")
	}

	frequencies := xslices.Frequencies([]rune(filtered))
	counts := slices.Sorted(maps.Values(frequencies))
	slices.Reverse(counts)

	if len(counts) == 0 {
		counts = []int{0}
	}
	counts[0] += jokers

	switch {
	case counts[0] == 5:
		return fiveOfKind
	case counts[0] == 4:
		return fourOfKind
	case counts[0] == 3 && len(counts) > 1 && counts[1] == 2:
		return fullHouse
	case counts[0] == 3:
		return threeOfKind
	case counts[0] == 2 && len(counts) > 1 && counts[1] == 2:
		return twoPair
	case counts[0] == 2:
		return onePair
	default:
		return highCard
	}
}

func compareHands(a, b hand, jokerRule bool) int {
	if a.typ != b.typ {
		return int(a.typ) - int(b.typ)
	}

	for i := 0; i < len(a.cards); i++ {
		aVal := cardValue(a.cards[i], jokerRule)
		bVal := cardValue(b.cards[i], jokerRule)
		if aVal != bVal {
			return aVal - bVal
		}
	}
	return 0
}

func parseLine(jokerRule bool) func([]byte) (hand, error) {
	return func(line []byte) (hand, error) {
		parts := strings.Fields(string(line))
		if len(parts) != 2 {
			return hand{}, nil
		}

		bid := aoc.MustAtoi(parts[1])
		return hand{
			cards: parts[0],
			bid:   bid,
			typ:   classifyHand(parts[0], jokerRule),
		}, nil
	}
}

func parseHands(r io.Reader, jokerRule bool) ([]hand, error) {
	s := scanner.NewScanner(r, parseLine(jokerRule))
	hands := slices.Collect(s.Values())
	return hands, s.Err()
}

func calculateWinnings(hands []hand, jokerRule bool) int {
	slices.SortFunc(hands, func(a, b hand) int {
		return compareHands(a, b, jokerRule)
	})

	total := 0
	for rank, h := range hands {
		total += h.bid * (rank + 1)
	}
	return total
}

func day07p01(r io.Reader) (string, error) {
	hands, err := parseHands(r, false)
	if err != nil {
		return "", err
	}

	result := calculateWinnings(hands, false)
	return strconv.Itoa(result), nil
}

func day07p02(r io.Reader) (string, error) {
	hands, err := parseHands(r, true)
	if err != nil {
		return "", err
	}

	result := calculateWinnings(hands, true)
	return strconv.Itoa(result), nil
}
