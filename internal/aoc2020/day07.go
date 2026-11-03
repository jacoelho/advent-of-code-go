package aoc2020

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type bagContents struct {
	color string
	count int64
}

func parseBagRules(r io.Reader) (map[string][]bagContents, error) {
	rules := make(map[string][]bagContents)
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		color, contents, ok := strings.Cut(line, " bags contain ")
		contents, period := strings.CutSuffix(contents, ".")
		if !ok || !period || color == "" {
			return nil, fmt.Errorf("invalid bag rule %q", line)
		}
		if _, exists := rules[color]; exists {
			return nil, fmt.Errorf("duplicate bag rule for %q", color)
		}
		rules[color] = nil
		if contents == "no other bags" {
			continue
		}
		for entry := range strings.SplitSeq(contents, ", ") {
			quantity, child, ok := strings.Cut(entry, " ")
			count, err := strconv.ParseInt(quantity, 10, 64)
			name, suffix := strings.CutSuffix(child, " bags")
			if !suffix {
				name, suffix = strings.CutSuffix(child, " bag")
			}
			if !ok || err != nil || count <= 0 || !suffix || name == "" {
				return nil, fmt.Errorf("invalid contained bag %q", entry)
			}
			rules[color] = append(rules[color], bagContents{name, count})
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	for color, contents := range rules {
		for _, child := range contents {
			if _, exists := rules[child.color]; !exists {
				return nil, fmt.Errorf("bag %q refers to undefined color %q", color, child.color)
			}
		}
	}
	return rules, nil
}

func day07p01(r io.Reader) (string, error) {
	rules, err := parseBagRules(r)
	if err != nil {
		return "", err
	}
	parents := make(map[string][]string)
	for color, contents := range rules {
		for _, child := range contents {
			parents[child.color] = append(parents[child.color], color)
		}
	}
	count := -1 // The root is not one of its containing bags.
	for range search.BFS("shiny gold", func(color string) iter.Seq[string] {
		return slices.Values(parents[color])
	}) {
		count++
	}
	return strconv.Itoa(count), nil
}

func day07p02(r io.Reader) (string, error) {
	rules, err := parseBagRules(r)
	if err != nil {
		return "", err
	}
	if _, exists := rules["shiny gold"]; !exists {
		return "", fmt.Errorf("no shiny gold bag rule")
	}
	totals := make(map[string]int64)
	visiting := make(map[string]bool)
	var count func(string) (int64, error)
	count = func(color string) (int64, error) {
		if total, ok := totals[color]; ok {
			return total, nil
		}
		if visiting[color] {
			return 0, fmt.Errorf("cycle in bag rules at %q", color)
		}
		visiting[color] = true
		var total int64
		for _, child := range rules[color] {
			nested, err := count(child.color)
			if err != nil {
				return 0, err
			}
			total += child.count * (1 + nested)
		}
		delete(visiting, color)
		totals[color] = total
		return total, nil
	}
	total, err := count("shiny gold")
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(total, 10), nil
}
