package aoc2020

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type messageRule struct {
	literal      byte
	alternatives [][]int
}

func parseMessageGrammar(r io.Reader) (map[int]messageRule, []string, error) {
	rules := make(map[int]messageRule)
	var messages []string
	readingRules := true
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		if line == "" {
			readingRules = false
			continue
		}
		if !readingRules {
			messages = append(messages, line)
			continue
		}
		label, body, ok := strings.Cut(line, ":")
		id, err := strconv.Atoi(label)
		if !ok || err != nil || id < 0 {
			return nil, nil, fmt.Errorf("invalid rule %q", line)
		}
		if _, exists := rules[id]; exists {
			return nil, nil, fmt.Errorf("duplicate rule %d", id)
		}
		body = strings.TrimSpace(body)
		rule := messageRule{}
		if strings.Contains(body, `"`) {
			if len(body) != 3 || body[0] != '"' || body[2] != '"' || body[1] < 'a' || body[1] > 'z' {
				return nil, nil, fmt.Errorf("invalid literal in rule %d", id)
			}
			rule.literal = body[1]
		} else {
			for alternative := range strings.SplitSeq(body, "|") {
				var sequence []int
				for token := range strings.FieldsSeq(alternative) {
					reference, err := strconv.Atoi(token)
					if err != nil || reference < 0 {
						return nil, nil, fmt.Errorf("invalid reference %q in rule %d", token, id)
					}
					sequence = append(sequence, reference)
				}
				if len(sequence) == 0 {
					return nil, nil, fmt.Errorf("empty alternative in rule %d", id)
				}
				rule.alternatives = append(rule.alternatives, sequence)
			}
		}
		rules[id] = rule
	}
	if err := s.Err(); err != nil {
		return nil, nil, err
	}
	if _, ok := rules[0]; !ok {
		return nil, nil, fmt.Errorf("missing rule 0")
	}
	for id, rule := range rules {
		for _, alternative := range rule.alternatives {
			for _, reference := range alternative {
				if _, ok := rules[reference]; !ok {
					return nil, nil, fmt.Errorf("rule %d references missing rule %d", id, reference)
				}
			}
		}
	}
	return rules, messages, nil
}

type messageState struct {
	rule, alternative, next, start int
}

type messageChart struct {
	states []messageState
	seen   map[messageState]bool
}

func (c *messageChart) add(state messageState) {
	if c.seen == nil {
		c.seen = make(map[messageState]bool)
	}
	if !c.seen[state] {
		c.seen[state] = true
		c.states = append(c.states, state)
	}
}

func matchesMessage(rules map[int]messageRule, message string) bool {
	sequence := func(state messageState) []int {
		if state.rule == -1 {
			return []int{0}
		}
		return rules[state.rule].alternatives[state.alternative]
	}
	// An Earley chart deduplicates partial matches at each input offset, so recursive
	// and ambiguous rules terminate without imposing an arbitrary recursion limit.
	// Empty alternatives are rejected, so completed rules always consume input.
	chart := make([]messageChart, len(message)+1)
	chart[0].add(messageState{rule: -1})
	for position := range chart {
		for i := 0; i < len(chart[position].states); i++ {
			state := chart[position].states[i]
			symbols := sequence(state)
			if state.next == len(symbols) {
				if state.rule == -1 {
					if position == len(message) {
						return true
					}
					continue
				}
				for _, previous := range chart[state.start].states {
					expected := sequence(previous)
					if previous.next < len(expected) && expected[previous.next] == state.rule {
						previous.next++
						chart[position].add(previous)
					}
				}
				continue
			}
			id := symbols[state.next]
			rule := rules[id]
			if rule.literal != 0 {
				if position < len(message) && message[position] == rule.literal {
					state.next++
					chart[position+1].add(state)
				}
			} else {
				for alternative := range rule.alternatives {
					chart[position].add(messageState{rule: id, alternative: alternative, start: position})
				}
			}
		}
	}
	return false
}

func countGrammarMatches(rules map[int]messageRule, messages []string) string {
	count := xiter.FromSlice(messages).Filter(func(message string) bool {
		return matchesMessage(rules, message)
	}).Count()
	return strconv.Itoa(count)
}

func day19p01(r io.Reader) (string, error) {
	rules, messages, err := parseMessageGrammar(r)
	if err != nil {
		return "", err
	}
	return countGrammarMatches(rules, messages), nil
}

func day19p02(r io.Reader) (string, error) {
	rules, messages, err := parseMessageGrammar(r)
	if err != nil {
		return "", err
	}
	for _, id := range []int{42, 31} {
		if _, ok := rules[id]; !ok {
			return "", fmt.Errorf("recursive rules require rule %d", id)
		}
	}
	rules[8] = messageRule{alternatives: [][]int{{42}, {42, 8}}}
	rules[11] = messageRule{alternatives: [][]int{{42, 31}, {42, 11, 31}}}
	return countGrammarMatches(rules, messages), nil
}
