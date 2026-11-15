package aoc2020

import (
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func parseCombat(r io.Reader) (deck1, deck2 []int, err error) {
	sections, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(sections) != 2 {
		return nil, nil, fmt.Errorf("expected 2 decks, got %d", len(sections))
	}
	seen := make(map[int]bool)
	var decks [2][]int
	for i, lines := range sections {
		if len(lines) == 0 {
			return nil, nil, fmt.Errorf("deck %d is missing", i+1)
		}
		if want := fmt.Sprintf("Player %d:", i+1); lines[0] != want {
			return nil, nil, fmt.Errorf("deck %d: header %q, want %q", i+1, lines[0], want)
		}
		if len(lines) == 1 {
			return nil, nil, fmt.Errorf("deck %d is empty", i+1)
		}
		for _, line := range lines[1:] {
			card, err := strconv.Atoi(line)
			if err != nil {
				return nil, nil, fmt.Errorf("deck %d: %w", i+1, err)
			}
			if card <= 0 {
				return nil, nil, fmt.Errorf("deck %d: card %d is not positive", i+1, card)
			}
			if seen[card] {
				return nil, nil, fmt.Errorf("duplicate card %d", card)
			}
			seen[card] = true
			decks[i] = append(decks[i], card)
		}
	}
	return decks[0], decks[1], nil
}

func scoreDeck(deck []int) int {
	total := 0
	for i, card := range deck {
		total += card * (len(deck) - i)
	}
	return total
}

func day22p01(r io.Reader) (string, error) {
	deck1, deck2, err := parseCombat(r)
	if err != nil {
		return "", err
	}
	for len(deck1) > 0 && len(deck2) > 0 {
		card1, card2 := deck1[0], deck2[0]
		deck1, deck2 = deck1[1:], deck2[1:]
		if card1 > card2 {
			deck1 = append(deck1, card1, card2)
		} else {
			deck2 = append(deck2, card2, card1)
		}
	}
	return strconv.Itoa(scoreDeck(slices.Concat(deck1, deck2))), nil
}

// recursiveCombat plays a game of Recursive Combat and reports whether player 1
// wins, together with the winner's final deck.
func recursiveCombat(deck1, deck2 []int, subgame bool) (bool, []int) {
	// Cards are distinct and positive, so the highest card exceeds the number of
	// other cards in play and can never trigger a subgame. Whoever holds it in a
	// subgame never loses it, so player 1 holding it wins by emptying player 2's
	// deck or by the repetition rule.
	if subgame && slices.Max(deck1) > slices.Max(deck2) {
		return true, nil
	}
	seen := make(map[string]struct{})
	var key []byte
	for len(deck1) > 0 && len(deck2) > 0 {
		key = key[:0]
		for _, card := range deck1 {
			key = binary.AppendUvarint(key, uint64(card))
		}
		key = append(key, 0)
		for _, card := range deck2 {
			key = binary.AppendUvarint(key, uint64(card))
		}
		if _, ok := seen[string(key)]; ok {
			return true, deck1
		}
		seen[string(key)] = struct{}{}

		card1, card2 := deck1[0], deck2[0]
		deck1, deck2 = deck1[1:], deck2[1:]
		player1Wins := card1 > card2
		if len(deck1) >= card1 && len(deck2) >= card2 {
			player1Wins, _ = recursiveCombat(slices.Clone(deck1[:card1]), slices.Clone(deck2[:card2]), true)
		}
		if player1Wins {
			deck1 = append(deck1, card1, card2)
		} else {
			deck2 = append(deck2, card2, card1)
		}
	}
	if len(deck1) > 0 {
		return true, deck1
	}
	return false, deck2
}

func day22p02(r io.Reader) (string, error) {
	deck1, deck2, err := parseCombat(r)
	if err != nil {
		return "", err
	}
	_, deck := recursiveCombat(deck1, deck2, false)
	return strconv.Itoa(scoreDeck(deck)), nil
}
