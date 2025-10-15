package aoc2020

import (
	"bytes"
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type combat struct {
	player1 *collections.Deque[int]
	player2 *collections.Deque[int]
}

type winner int8

const (
	winnerPlayer1 winner = iota
	winnerPlayer2
)

func parseDeck(b []byte) ([]int, error) {
	lines := bytes.Split(b, []byte("\n"))

	var deck []int
	for _, line := range lines[1:] { // skip header
		if len(line) == 0 {
			continue
		}
		n, err := convert.ScanNumber[int](line)
		if err != nil {
			return nil, err
		}
		deck = append(deck, n)
	}
	return deck, nil
}

func parseCombat(r io.Reader) (*combat, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte{'\n', '\n'}), parseDeck)
	decks := slices.Collect(s.Values())
	if len(decks) < 2 {
		return nil, errors.New("not enough decks")
	}
	p1 := collections.NewDeque[int](len(decks[0]))
	for _, v := range decks[0] {
		p1.PushBack(v)
	}
	p2 := collections.NewDeque[int](len(decks[1]))
	for _, v := range decks[1] {
		p2.PushBack(v)
	}
	return &combat{
		player1: p1,
		player2: p2,
	}, s.Err()
}

func day22p01(r io.Reader) (string, error) {
	c, err := parseCombat(r)
	if err != nil {
		return "", err
	}
	_, winner := playCombat(c)
	score := scoreDeck(winner)
	return strconv.Itoa(score), nil
}

func day22p02(r io.Reader) (string, error) {
	c, err := parseCombat(r)
	if err != nil {
		return "", err
	}
	_, deck := playRecursiveCombat(c)
	return strconv.Itoa(scoreDeck(deck)), nil
}

// awardRound appends the winner's card followed by the loser's card to the winner's deck.
func awardRound(winner *collections.Deque[int], winnerCard, loserCard int) {
	winner.PushBack(winnerCard)
	winner.PushBack(loserCard)
}

// playCombat runs the standard combat game and returns the winner and their deck
func playCombat(c *combat) (winner, *collections.Deque[int]) {
	for c.player1.Size() > 0 && c.player2.Size() > 0 {
		c1, _ := c.player1.PopFront()
		c2, _ := c.player2.PopFront()
		if c1 > c2 {
			awardRound(c.player1, c1, c2)
		} else {
			awardRound(c.player2, c2, c1)
		}
	}
	if c.player1.Size() > 0 {
		return winnerPlayer1, c.player1
	}
	return winnerPlayer2, c.player2
}

// scoreDeck computes the score by summing value*position from bottom=1 to top=N
func scoreDeck(d *collections.Deque[int]) int {
	total := 0
	multiplier := 1
	for v := range d.IterBack() {
		total += v * multiplier
		multiplier++
	}
	return total
}

func subDeckN(d *collections.Deque[int], n int) *collections.Deque[int] {
	nd := collections.NewDeque[int](n)
	i := 0
	for v := range d.IterFront() {
		if i >= n {
			break
		}
		nd.PushBack(v)
		i++
	}
	return nd
}

func playRecursiveCombat(c *combat) (winner, *collections.Deque[int]) {
	seen := map[string]struct{}{}

	for c.player1.Size() > 0 && c.player2.Size() > 0 {
		key := deckState(c.player1, c.player2)
		if _, ok := seen[key]; ok {
			return winnerPlayer1, c.player1
		}
		seen[key] = struct{}{}

		c1, _ := c.player1.PopFront()
		c2, _ := c.player2.PopFront()

		var roundWinner winner
		if hasAtLeastCards(c.player1, c1) && hasAtLeastCards(c.player2, c2) {
			subWinner, _ := playRecursiveCombat(&combat{
				player1: subDeckN(c.player1, c1),
				player2: subDeckN(c.player2, c2),
			})
			roundWinner = subWinner
		} else if c1 > c2 {
			roundWinner = winnerPlayer1
		} else {
			roundWinner = winnerPlayer2
		}

		if roundWinner == winnerPlayer1 {
			awardRound(c.player1, c1, c2)
		} else {
			awardRound(c.player2, c2, c1)
		}
	}
	if c.player1.Size() > 0 {
		return winnerPlayer1, c.player1
	}
	return winnerPlayer2, c.player2
}

// deckState returns a string representation of the deck state.
func deckState(p1, p2 *collections.Deque[int]) string {
	var b strings.Builder
	writeDeck(&b, p1)
	b.WriteByte('|')
	writeDeck(&b, p2)
	return b.String()
}

// writeDeck appends a comma-separated front-to-back representation of the deck into the builder.
func writeDeck(b *strings.Builder, d *collections.Deque[int]) {
	for v := range d.IterFront() {
		b.WriteString(strconv.Itoa(v))
		b.WriteByte(',')
	}
}

// hasAtLeastCards returns true if the deque has at least n elements.
func hasAtLeastCards(d *collections.Deque[int], n int) bool { return d.Size() >= n }
