package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

func countPassphrases(r io.Reader, canonical func(string) string) (string, error) {
	s := bufio.NewScanner(r)
	count, lines := 0, 0
	for s.Scan() {
		lines++
		if strings.TrimSpace(s.Text()) == "" {
			return "", fmt.Errorf("empty passphrase on line %d", lines)
		}
		if i := strings.IndexFunc(s.Text(), func(c rune) bool { return c != ' ' && (c < 'a' || c > 'z') }); i >= 0 {
			return "", fmt.Errorf("invalid passphrase character %q on line %d", s.Text()[i], lines)
		}
		seen := collections.NewSet[string]()
		valid := true
		for word := range strings.FieldsSeq(s.Text()) {
			key := canonical(word)
			if seen.Contains(key) {
				valid = false
				break
			}
			seen.Add(key)
		}
		if valid {
			count++
		}
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if lines == 0 {
		return "", fmt.Errorf("no passphrases")
	}

	return strconv.Itoa(count), nil
}

func day04p01(r io.Reader) (string, error) {
	return countPassphrases(r, func(word string) string { return word })
}

func day04p02(r io.Reader) (string, error) {
	return countPassphrases(r, func(word string) string {
		letters := []byte(word)
		slices.Sort(letters)
		return string(letters)
	})
}
