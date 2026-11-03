package aoc2015

import (
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func niceString(word string) bool {
	vowels, doubled := 0, false
	for i := range len(word) {
		if strings.ContainsRune("aeiou", rune(word[i])) {
			vowels++
		}
		if i > 0 {
			pair := word[i-1 : i+1]
			switch pair {
			case "ab", "cd", "pq", "xy":
				return false
			}
			doubled = doubled || word[i] == word[i-1]
		}
	}
	return vowels >= 3 && doubled
}

func nicerString(word string) bool {
	firstPair := make(map[[2]byte]int)
	repeatedPair, repeatedLetter := false, false
	for i := range len(word) {
		if i >= 2 && word[i] == word[i-2] {
			repeatedLetter = true
		}
		if i+1 < len(word) {
			pair := [2]byte{word[i], word[i+1]}
			if previous, ok := firstPair[pair]; ok {
				repeatedPair = repeatedPair || i-previous >= 2
			} else {
				firstPair[pair] = i
			}
		}
	}
	return repeatedPair && repeatedLetter
}

func countNiceStrings(r io.Reader, nice func(string) bool) (string, error) {
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	count := xiter.From(s.Values()).Filter(nice).Count()
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(count), nil
}

func day05p01(r io.Reader) (string, error) {
	return countNiceStrings(r, niceString)
}

func day05p02(r io.Reader) (string, error) {
	return countNiceStrings(r, nicerString)
}
