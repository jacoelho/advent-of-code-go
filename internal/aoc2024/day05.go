package aoc2024

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func parseSleighLaunchSafetyRules(r io.Reader) (map[[2]int]bool, [][]int, error) {
	rules := make(map[[2]int]bool)
	var updates [][]int
	readingUpdates := false
	lines, work := 0, 0
	s := bufio.NewScanner(r)
	for s.Scan() {
		lines++
		if lines > 100_000 {
			return nil, nil, fmt.Errorf("page ordering input exceeds 100000 lines")
		}
		text := s.Text()
		if text == "" {
			readingUpdates = true
			continue
		}
		if !readingUpdates {
			left, right, ok := strings.Cut(text, "|")
			if !ok || strings.Contains(right, "|") {
				return nil, nil, fmt.Errorf("invalid page rule %q", text)
			}
			before, err := convert.ScanNumber[int]([]byte(left))
			if err != nil {
				return nil, nil, fmt.Errorf("page rule: %w", err)
			}
			after, err := convert.ScanNumber[int]([]byte(right))
			if err != nil {
				return nil, nil, fmt.Errorf("page rule: %w", err)
			}
			if before < 0 || after < 0 {
				return nil, nil, fmt.Errorf("page numbers must be nonnegative")
			}
			rules[[2]int{before, after}] = true
			continue
		}

		var update []int
		seen := make(map[int]bool)
		for field := range strings.SplitSeq(text, ",") {
			page, err := convert.ScanNumber[int]([]byte(field))
			if err != nil {
				return nil, nil, fmt.Errorf("page update: %w", err)
			}
			if page < 0 || seen[page] {
				return nil, nil, fmt.Errorf("update pages must be nonnegative and unique")
			}
			if len(update) == 1023 {
				return nil, nil, fmt.Errorf("update exceeds 1023 pages")
			}
			seen[page] = true
			update = append(update, page)
		}
		if len(update)%2 == 0 {
			return nil, nil, fmt.Errorf("update needs an odd number of pages")
		}
		work += len(update) * len(update)
		if work > 10_000_000 {
			return nil, nil, fmt.Errorf("page ordering exceeds ten million page pairs")
		}
		updates = append(updates, update)
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if len(updates) == 0 {
		return nil, nil, fmt.Errorf("missing page updates")
	}
	return rules, updates, nil
}

func isOrdered(rules map[[2]int]bool, update []int) bool {
	for i, page := range update {
		if rules[[2]int{page, page}] {
			return false
		}
		for _, later := range update[i+1:] {
			if rules[[2]int{later, page}] {
				return false
			}
		}
	}
	return true
}

func orderPages(rules map[[2]int]bool, update []int) ([]int, error) {
	incoming := make([]int, len(update))
	for i, page := range update {
		for _, before := range update {
			if rules[[2]int{before, page}] {
				incoming[i]++
			}
		}
	}
	ordered := make([]int, 0, len(update))
	for len(ordered) < len(update) {
		first := -1
		// Original update order breaks ties among currently available pages.
		for i, count := range incoming {
			if count == 0 {
				first = i
				break
			}
		}
		if first == -1 {
			return nil, fmt.Errorf("page update contains a rule cycle")
		}
		incoming[first] = -1
		page := update[first]
		ordered = append(ordered, page)
		for i, after := range update {
			if incoming[i] > 0 && rules[[2]int{page, after}] {
				incoming[i]--
			}
		}
	}
	return ordered, nil
}

func day05p01(r io.Reader) (string, error) {
	rules, updates, err := parseSleighLaunchSafetyRules(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, update := range updates {
		if isOrdered(rules, update) {
			middle := int64(update[len(update)/2])
			if middle > math.MaxInt64-total {
				return "", fmt.Errorf("middle page total overflows int64")
			}
			total += middle
		}
	}
	return strconv.FormatInt(total, 10), nil
}

func day05p02(r io.Reader) (string, error) {
	rules, updates, err := parseSleighLaunchSafetyRules(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, update := range updates {
		if isOrdered(rules, update) {
			continue
		}
		ordered, err := orderPages(rules, update)
		if err != nil {
			return "", err
		}
		middle := int64(ordered[len(ordered)/2])
		if middle > math.MaxInt64-total {
			return "", fmt.Errorf("middle page total overflows int64")
		}
		total += middle
	}
	return strconv.FormatInt(total, 10), nil
}
