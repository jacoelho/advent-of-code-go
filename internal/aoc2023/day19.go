package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type ratingPart [4]int64

type workflowRule struct {
	category    int
	operator    byte
	value       int64
	destination string
}

type workflowSet map[string][]workflowRule

type ratingInterval struct {
	min, max int64
}

type ratingRanges [4]ratingInterval

func (r ratingInterval) split(operator byte, value int64) (ratingInterval, ratingInterval) {
	empty := ratingInterval{min: 1, max: 0}
	if r.min > r.max {
		return empty, empty
	}
	// Classify all/none before adjusting the threshold, including int64 extremes.
	if operator == '<' {
		if value <= r.min {
			return empty, r
		}
		if value > r.max {
			return r, empty
		}
		return ratingInterval{min: r.min, max: value - 1}, ratingInterval{min: value, max: r.max}
	}
	if value < r.min {
		return r, empty
	}
	if value >= r.max {
		return empty, r
	}
	return ratingInterval{min: value + 1, max: r.max}, ratingInterval{min: r.min, max: value}
}

func (r ratingRanges) size() int64 {
	size := int64(1)
	for _, interval := range r {
		if interval.min > interval.max {
			return 0
		}
		size *= interval.max - interval.min + 1
	}
	return size
}

func ratingCategory(name string) (int, error) {
	if len(name) == 1 {
		if category := strings.IndexByte("xmas", name[0]); category >= 0 {
			return category, nil
		}
	}
	return 0, fmt.Errorf("invalid rating category %q", name)
}

func workflowName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for _, c := range name {
		if c < 'a' || c > 'z' {
			return false
		}
	}
	return true
}

func parseWorkflows(lines []string) (workflowSet, error) {
	workflows := make(workflowSet)
	for _, line := range lines {
		name, body, ok := strings.Cut(strings.TrimSpace(line), "{")
		if !ok || !workflowName(name) {
			return nil, fmt.Errorf("invalid workflow declaration %q", line)
		}
		body, ok = strings.CutSuffix(body, "}")
		if !ok {
			return nil, fmt.Errorf("unclosed workflow %q", name)
		}
		if _, exists := workflows[name]; exists {
			return nil, fmt.Errorf("duplicate workflow %q", name)
		}
		texts := strings.Split(body, ",")
		var rules []workflowRule
		for i, text := range texts {
			condition, destination, conditional := strings.Cut(strings.TrimSpace(text), ":")
			rule := workflowRule{category: -1, destination: condition}
			if conditional {
				if i == len(texts)-1 || len(condition) < 3 || condition[1] != '<' && condition[1] != '>' {
					return nil, fmt.Errorf("invalid workflow condition %q", condition)
				}
				category, err := ratingCategory(condition[:1])
				if err != nil {
					return nil, err
				}
				value, err := convert.ScanNumber[int64]([]byte(condition[2:]))
				if err != nil {
					return nil, err
				}
				rule = workflowRule{category: category, operator: condition[1], value: value, destination: destination}
			} else if i != len(texts)-1 {
				return nil, fmt.Errorf("workflow default must be its final rule")
			}
			if rule.destination != "A" && rule.destination != "R" && !workflowName(rule.destination) {
				return nil, fmt.Errorf("invalid workflow destination %q", rule.destination)
			}
			rules = append(rules, rule)
		}
		workflows[name] = rules
	}
	if _, exists := workflows["in"]; !exists {
		return nil, fmt.Errorf("missing entry workflow")
	}
	for _, rules := range workflows {
		for _, rule := range rules {
			if rule.destination == "A" || rule.destination == "R" {
				continue
			}
			if _, exists := workflows[rule.destination]; !exists {
				return nil, fmt.Errorf("missing workflow %q", rule.destination)
			}
		}
	}
	return workflows, nil
}

func parseRatingPart(line string) (ratingPart, error) {
	body, ok := strings.CutPrefix(strings.TrimSpace(line), "{")
	if !ok {
		return ratingPart{}, fmt.Errorf("invalid rating part")
	}
	body, ok = strings.CutSuffix(body, "}")
	if !ok {
		return ratingPart{}, fmt.Errorf("unclosed rating part")
	}
	var part ratingPart
	var seen [4]bool
	for field := range strings.SplitSeq(body, ",") {
		name, text, ok := strings.Cut(strings.TrimSpace(field), "=")
		if !ok {
			return ratingPart{}, fmt.Errorf("invalid rating field")
		}
		category, err := ratingCategory(name)
		if err != nil {
			return ratingPart{}, err
		}
		if seen[category] {
			return ratingPart{}, fmt.Errorf("duplicate rating category %q", name)
		}
		value, err := convert.ScanNumber[int64]([]byte(text))
		if err != nil {
			return ratingPart{}, err
		}
		if value < 1 || value > 4000 {
			return ratingPart{}, fmt.Errorf("rating must be between 1 and 4000")
		}
		part[category], seen[category] = value, true
	}
	for _, present := range seen {
		if !present {
			return ratingPart{}, fmt.Errorf("missing rating category")
		}
	}
	return part, nil
}

func parseWorkflowInput(r io.Reader) (workflowSet, []ratingPart, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(blocks) != 2 || len(blocks[1]) == 0 {
		return nil, nil, fmt.Errorf("expected workflow and part sections")
	}
	workflows, err := parseWorkflows(blocks[0])
	if err != nil {
		return nil, nil, err
	}
	var parts []ratingPart
	for _, line := range blocks[1] {
		part, err := parseRatingPart(line)
		if err != nil {
			return nil, nil, err
		}
		parts = append(parts, part)
	}
	return workflows, parts, nil
}

func processPart(part ratingPart, workflows workflowSet) (bool, error) {
	name := "in"
	for range len(workflows) + 1 {
		if name == "A" || name == "R" {
			return name == "A", nil
		}
		for _, rule := range workflows[name] {
			if rule.category == -1 || rule.operator == '<' && part[rule.category] < rule.value || rule.operator == '>' && part[rule.category] > rule.value {
				name = rule.destination
				break
			}
		}
	}
	return false, fmt.Errorf("part enters a workflow cycle")
}

func countAccepted(workflows workflowSet, limit int64) (int64, error) {
	if limit < 1 || limit > 4000 {
		return 0, fmt.Errorf("rating limit must be between 1 and 4000")
	}
	var initial ratingRanges
	for i := range initial {
		initial[i] = ratingInterval{min: 1, max: limit}
	}
	active := make(map[string]bool)
	var visit func(ratingRanges, string) (int64, error)
	visit = func(current ratingRanges, name string) (int64, error) {
		if name == "A" {
			return current.size(), nil
		}
		if name == "R" {
			return 0, nil
		}
		// Ratings never change, so returning to a workflow traps this subset.
		if active[name] {
			return 0, fmt.Errorf("ratings enter a workflow cycle at %s", name)
		}
		active[name] = true
		defer delete(active, name)
		var total int64
		for _, rule := range workflows[name] {
			if rule.category == -1 {
				count, err := visit(current, rule.destination)
				return total + count, err
			}
			matching, rest := current[rule.category].split(rule.operator, rule.value)
			if matching.min <= matching.max {
				branch := current
				branch[rule.category] = matching
				count, err := visit(branch, rule.destination)
				if err != nil {
					return 0, err
				}
				total += count
			}
			if rest.min > rest.max {
				return total, nil
			}
			current[rule.category] = rest
		}
		return 0, fmt.Errorf("workflow has no default rule")
	}
	return visit(initial, "in")
}

func day19p01(r io.Reader) (string, error) {
	workflows, parts, err := parseWorkflowInput(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, part := range parts {
		accepted, err := processPart(part, workflows)
		if err != nil {
			return "", err
		}
		if accepted {
			for _, value := range part {
				total += value
			}
		}
	}
	return strconv.FormatInt(total, 10), nil
}

func day19p02(r io.Reader) (string, error) {
	workflows, _, err := parseWorkflowInput(r)
	if err != nil {
		return "", err
	}
	count, err := countAccepted(workflows, 4000)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(count, 10), nil
}
