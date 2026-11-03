package aoc2020

import (
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

type ticketRule struct {
	name   string
	ranges [2][2]int
}

func (rule ticketRule) accepts(value int) bool {
	for _, span := range rule.ranges {
		if value >= span[0] && value <= span[1] {
			return true
		}
	}
	return false
}

type ticketNotes struct {
	rules  []ticketRule
	yours  []int
	nearby [][]int
}

func parseTicketNotes(r io.Reader) (ticketNotes, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return ticketNotes{}, err
	}
	sections := strings.Split(strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n")), "\n\n")
	if len(sections) != 3 {
		return ticketNotes{}, fmt.Errorf("expected rules, your ticket, and nearby tickets")
	}

	pattern := regexp.MustCompile(`^([a-z ]+): ([0-9]+)-([0-9]+) or ([0-9]+)-([0-9]+)$`)
	var notes ticketNotes
	names := make(map[string]bool)
	for line := range strings.SplitSeq(sections[0], "\n") {
		matches := pattern.FindStringSubmatch(line)
		if matches == nil {
			return ticketNotes{}, fmt.Errorf("invalid ticket rule %q", line)
		}
		if names[matches[1]] {
			return ticketNotes{}, fmt.Errorf("duplicate ticket rule %q", matches[1])
		}
		names[matches[1]] = true
		rule := ticketRule{name: matches[1]}
		for i := range 4 {
			value, err := strconv.Atoi(matches[i+2])
			if err != nil {
				return ticketNotes{}, fmt.Errorf("ticket rule %q: %w", rule.name, err)
			}
			rule.ranges[i/2][i%2] = value
		}
		for _, span := range rule.ranges {
			if span[0] > span[1] {
				return ticketNotes{}, fmt.Errorf("reversed range in ticket rule %q", rule.name)
			}
		}
		notes.rules = append(notes.rules, rule)
	}

	parseTicket := func(line string) ([]int, error) {
		var values []int
		for item := range strings.SplitSeq(line, ",") {
			value, err := strconv.Atoi(item)
			if err != nil || value < 0 {
				return nil, fmt.Errorf("invalid ticket value %q", item)
			}
			values = append(values, value)
		}
		if len(values) != len(notes.rules) {
			return nil, fmt.Errorf("ticket has %d fields, want %d", len(values), len(notes.rules))
		}
		return values, nil
	}

	yours := strings.Split(sections[1], "\n")
	if len(yours) != 2 || yours[0] != "your ticket:" {
		return ticketNotes{}, fmt.Errorf("invalid your-ticket section")
	}
	notes.yours, err = parseTicket(yours[1])
	if err != nil {
		return ticketNotes{}, err
	}

	nearby := strings.Split(sections[2], "\n")
	if nearby[0] != "nearby tickets:" {
		return ticketNotes{}, fmt.Errorf("invalid nearby-tickets section")
	}
	for _, line := range nearby[1:] {
		ticket, err := parseTicket(line)
		if err != nil {
			return ticketNotes{}, err
		}
		notes.nearby = append(notes.nearby, ticket)
	}
	return notes, nil
}

func validTicketValue(rules []ticketRule, value int) bool {
	for _, rule := range rules {
		if rule.accepts(value) {
			return true
		}
	}
	return false
}

func day16p01(r io.Reader) (string, error) {
	notes, err := parseTicketNotes(r)
	if err != nil {
		return "", err
	}
	sum := new(big.Int)
	for _, ticket := range notes.nearby {
		for _, value := range ticket {
			if !validTicketValue(notes.rules, value) {
				sum.Add(sum, big.NewInt(int64(value)))
			}
		}
	}
	return sum.String(), nil
}

// Match every rule to one column, optionally excluding one edge. Rechecking
// each selected edge establishes uniqueness without enumerating permutations.
func matchTicketFields(candidates [][]bool, excludedRule, excludedColumn int) ([]int, bool) {
	owners := make([]int, len(candidates))
	for i := range owners {
		owners[i] = -1
	}
	var assign func(int, []bool) bool
	assign = func(rule int, visited []bool) bool {
		for column, possible := range candidates[rule] {
			if !possible || visited[column] || rule == excludedRule && column == excludedColumn {
				continue
			}
			visited[column] = true
			if owners[column] < 0 || assign(owners[column], visited) {
				owners[column] = rule
				return true
			}
		}
		return false
	}
	for rule := range candidates {
		if !assign(rule, make([]bool, len(candidates))) {
			return nil, false
		}
	}
	columns := make([]int, len(candidates))
	for column, rule := range owners {
		columns[rule] = column
	}
	return columns, true
}

func day16p02(r io.Reader) (string, error) {
	notes, err := parseTicketNotes(r)
	if err != nil {
		return "", err
	}
	tickets := [][]int{notes.yours}
	for _, value := range notes.yours {
		if !validTicketValue(notes.rules, value) {
			return "", fmt.Errorf("your ticket contains invalid value %d", value)
		}
	}
	for _, ticket := range notes.nearby {
		valid := true
		for _, value := range ticket {
			if !validTicketValue(notes.rules, value) {
				valid = false
				break
			}
		}
		if valid {
			tickets = append(tickets, ticket)
		}
	}

	candidates := make([][]bool, len(notes.rules))
	for i, rule := range notes.rules {
		candidates[i] = make([]bool, len(notes.rules))
		for column := range notes.rules {
			candidates[i][column] = true
			for _, ticket := range tickets {
				if !rule.accepts(ticket[column]) {
					candidates[i][column] = false
					break
				}
			}
		}
	}

	columns, ok := matchTicketFields(candidates, -1, -1)
	if !ok {
		return "", fmt.Errorf("ticket fields have no complete assignment")
	}
	for rule, column := range columns {
		if _, alternative := matchTicketFields(candidates, rule, column); alternative {
			return "", fmt.Errorf("ticket field assignment is ambiguous")
		}
	}

	product := big.NewInt(1)
	for rule, column := range columns {
		if strings.HasPrefix(notes.rules[rule].name, "departure") {
			product.Mul(product, big.NewInt(int64(notes.yours[column])))
		}
	}
	return product.String(), nil
}
