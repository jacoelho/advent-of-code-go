package aoc2023

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	workflowAccept = "A"
	workflowReject = "R"
	workflowStart  = "in"
)

// Category represents a rating category
type Category byte

const (
	CategoryNone Category = 0
	CategoryX    Category = 'x'
	CategoryM    Category = 'm'
	CategoryA    Category = 'a'
	CategoryS    Category = 's'
)

// ParseCategory parses a category from a string
func ParseCategory(s string) (Category, error) {
	if len(s) == 1 {
		switch Category(s[0]) {
		case CategoryX, CategoryM, CategoryA, CategoryS:
			return Category(s[0]), nil
		}
	}
	return CategoryNone, fmt.Errorf("invalid category: %s", s)
}

// Operator represents a comparison operator
type Operator byte

const (
	OperatorNone        Operator = 0
	OperatorGreaterThan Operator = '>'
	OperatorLessThan    Operator = '<'
)

// Evaluate performs the comparison operation
func (o Operator) Evaluate(a, b int) bool {
	switch o {
	case OperatorGreaterThan:
		return a > b
	case OperatorLessThan:
		return a < b
	default:
		panic(fmt.Sprintf("unknown operator: %c", o))
	}
}

// Part represents a part with four ratings
type Part struct {
	X, M, A, S int
}

// Get returns the rating for the specified category
func (p Part) Get(c Category) int {
	switch c {
	case CategoryX:
		return p.X
	case CategoryM:
		return p.M
	case CategoryA:
		return p.A
	case CategoryS:
		return p.S
	default:
		panic(fmt.Sprintf("unknown category: %c", c))
	}
}

// Sum returns the sum of all ratings
func (p Part) Sum() int {
	return p.X + p.M + p.A + p.S
}

// Set sets the rating for the specified category
func (p *Part) Set(c Category, value int) {
	switch c {
	case CategoryX:
		p.X = value
	case CategoryM:
		p.M = value
	case CategoryA:
		p.A = value
	case CategoryS:
		p.S = value
	}
}

// Rule represents a workflow rule
type Rule struct {
	Category            Category
	Operator            Operator
	Value               int
	DestinationWorkflow string
}

// Matches returns true if the rule matches the given part
func (r Rule) Matches(p Part) bool {
	if r.Category == CategoryNone {
		return true
	}
	return r.Operator.Evaluate(p.Get(r.Category), r.Value)
}

// Workflow represents a workflow with a list of rules
type Workflow struct {
	Rules []Rule
}

// Workflows is a map of workflow name to workflow
type Workflows map[string]Workflow

// InclusiveRange represents an inclusive integer range [Min, Max]
type InclusiveRange struct {
	Min, Max int
}

// Size returns the number of integers in the range
func (r InclusiveRange) Size() int {
	if !r.IsValid() {
		return 0
	}
	return r.Max - r.Min + 1
}

// IsValid returns true if the range is non-empty
func (r InclusiveRange) IsValid() bool {
	return r.Min <= r.Max
}

// Split splits range into (matching, non-matching) for the given operator
func (r InclusiveRange) Split(op Operator, value int) (matching, nonMatching InclusiveRange) {
	switch op {
	case OperatorGreaterThan:
		return InclusiveRange{Min: value + 1, Max: r.Max}, InclusiveRange{Min: r.Min, Max: value}
	case OperatorLessThan:
		return InclusiveRange{Min: r.Min, Max: value - 1}, InclusiveRange{Min: value, Max: r.Max}
	default:
		return r, InclusiveRange{}
	}
}

// RatingRange represents ranges for all four rating categories
type RatingRange struct {
	X, M, A, S InclusiveRange
}

// NewRatingRange creates a new RatingRange with all categories set to the same range
func NewRatingRange(min, max int) RatingRange {
	r := InclusiveRange{Min: min, Max: max}
	return RatingRange{X: r, M: r, A: r, S: r}
}

// Size returns the number of distinct combinations in this range
func (r RatingRange) Size() int {
	return r.X.Size() * r.M.Size() * r.A.Size() * r.S.Size()
}

// IsValid returns true if all rating ranges are valid
func (r RatingRange) IsValid() bool {
	return r.X.IsValid() && r.M.IsValid() && r.A.IsValid() && r.S.IsValid()
}

// Get returns the range for the specified category
func (r RatingRange) Get(c Category) InclusiveRange {
	switch c {
	case CategoryX:
		return r.X
	case CategoryM:
		return r.M
	case CategoryA:
		return r.A
	case CategoryS:
		return r.S
	default:
		panic("invalid category")
	}
}

// Set returns a new RatingRange with the specified category set to rng
func (r RatingRange) Set(c Category, rng InclusiveRange) RatingRange {
	result := r
	switch c {
	case CategoryX:
		result.X = rng
	case CategoryM:
		result.M = rng
	case CategoryA:
		result.A = rng
	case CategoryS:
		result.S = rng
	}
	return result
}

// Split returns (matching, non-matching) ranges for a rule
func (r RatingRange) Split(category Category, operator Operator, value int) (matching, nonMatching RatingRange) {
	categoryRange := r.Get(category)
	matchingRange, nonMatchingRange := categoryRange.Split(operator, value)

	matching = r.Set(category, matchingRange)
	nonMatching = r.Set(category, nonMatchingRange)
	return
}

// cutAny tries to cut the string using any of the provided delimiters.
// Returns the parts, the delimiter that was found, and whether a cut was made.
func cutAny(s string, delimiters ...string) (before, after, delimiter string, found bool) {
	for _, delim := range delimiters {
		if b, a, ok := strings.Cut(s, delim); ok {
			return b, a, delim, true
		}
	}
	return "", "", "", false
}

func day19p01(r io.Reader) (string, error) {
	workflows, parts, err := parseInput(r)
	if err != nil {
		return "", err
	}

	total := 0
	for _, part := range parts {
		if processPart(part, workflows) {
			total += part.Sum()
		}
	}

	return strconv.Itoa(total), nil
}

func day19p02(r io.Reader) (string, error) {
	workflows, _, err := parseInput(r)
	if err != nil {
		return "", err
	}

	count := countAccepted(workflows)
	return strconv.Itoa(count), nil
}

// countAccepted calculates the total number of accepted rating combinations
func countAccepted(workflows Workflows) int {
	initial := NewRatingRange(1, 4000)
	return processRange(initial, workflowStart, workflows)
}

// processRange recursively processes a rating range through workflows
func processRange(rng RatingRange, workflowName string, workflows Workflows) int {
	if workflowName == workflowAccept {
		return rng.Size()
	}
	if workflowName == workflowReject {
		return 0
	}

	workflow := workflows[workflowName]
	total := 0
	current := rng

	for _, rule := range workflow.Rules {
		if rule.Category == CategoryNone {
			// default rule: send all remaining to destination
			total += processRange(current, rule.DestinationWorkflow, workflows)
			break
		}

		// split range
		matching, nonMatching := current.Split(rule.Category, rule.Operator, rule.Value)

		// process matching range
		if matching.IsValid() {
			total += processRange(matching, rule.DestinationWorkflow, workflows)
		}

		// continue with non-matching range
		current = nonMatching
		if !current.IsValid() {
			break
		}
	}

	return total
}

// parseInput parses the entire input into workflows and parts
func parseInput(r io.Reader) (Workflows, []Part, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}

	sections := strings.Split(strings.TrimSpace(string(content)), "\n\n")
	if len(sections) != 2 {
		return nil, nil, fmt.Errorf("expected 2 sections separated by blank line, got %d", len(sections))
	}

	workflows, err := parseWorkflows(sections[0])
	if err != nil {
		return nil, nil, err
	}

	parts, err := parseParts(sections[1])
	if err != nil {
		return nil, nil, err
	}

	return workflows, parts, nil
}

// parseWorkflows parses workflow definitions
func parseWorkflows(input string) (Workflows, error) {
	workflows := make(Workflows)
	lines := strings.SplitSeq(strings.TrimSpace(input), "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// find the opening brace
		name, rest, found := strings.Cut(line, "{")
		if !found {
			return nil, fmt.Errorf("invalid workflow format: %s", line)
		}

		rulesStr := strings.TrimSuffix(rest, "}")

		rules, err := parseRules(rulesStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rules for workflow %s: %w", name, err)
		}

		workflows[name] = Workflow{
			Rules: rules,
		}
	}

	return workflows, nil
}

// parseRules parses a comma-separated list of rules
func parseRules(rulesStr string) ([]Rule, error) {
	ruleStrs := strings.Split(rulesStr, ",")
	rules := make([]Rule, 0, len(ruleStrs))

	for _, ruleStr := range ruleStrs {
		ruleStr = strings.TrimSpace(ruleStr)
		rule, err := parseRule(ruleStr)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// parseRule parses a single rule
func parseRule(ruleStr string) (Rule, error) {
	// check if this is a conditional rule (contains :)
	condition, destination, hasCondition := strings.Cut(ruleStr, ":")
	if !hasCondition {
		// default rule (no condition)
		return Rule{
			Category:            CategoryNone,
			Operator:            OperatorNone,
			Value:               0,
			DestinationWorkflow: ruleStr,
		}, nil
	}

	// parse conditional rule: category>value or category<value
	catStr, valStr, delim, found := cutAny(condition, ">", "<")
	if !found {
		return Rule{}, fmt.Errorf("invalid condition format: %s", condition)
	}

	// determine operator from delimiter
	var operator Operator
	switch delim {
	case ">":
		operator = OperatorGreaterThan
	case "<":
		operator = OperatorLessThan
	}

	// common parsing logic
	category, err := ParseCategory(catStr)
	if err != nil {
		return Rule{}, fmt.Errorf("invalid category in condition: %s", condition)
	}

	value, err := strconv.Atoi(valStr)
	if err != nil {
		return Rule{}, fmt.Errorf("invalid value in condition: %s", condition)
	}

	return Rule{
		Category:            category,
		Operator:            operator,
		Value:               value,
		DestinationWorkflow: destination,
	}, nil
}

// parseParts parses part definitions
func parseParts(input string) ([]Part, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	parts := make([]Part, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		part, err := parsePart(line)
		if err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	return parts, nil
}

// parsePart parses a single part from format {x=val,m=val,a=val,s=val}
func parsePart(line string) (Part, error) {
	content, found := strings.CutPrefix(line, "{")
	if !found {
		return Part{}, fmt.Errorf("invalid part format: %s", line)
	}
	content, found = strings.CutSuffix(content, "}")
	if !found {
		return Part{}, fmt.Errorf("invalid part format: %s", line)
	}
	ratings := strings.Split(content, ",")

	if len(ratings) != 4 {
		return Part{}, fmt.Errorf("expected 4 ratings, got %d: %s", len(ratings), line)
	}

	var part Part
	for _, rating := range ratings {
		rating = strings.TrimSpace(rating)
		category, valueStr, found := strings.Cut(rating, "=")
		if !found {
			return Part{}, fmt.Errorf("invalid rating format: %s", rating)
		}

		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return Part{}, fmt.Errorf("invalid rating value: %s", rating)
		}

		cat, err := ParseCategory(category)
		if err != nil {
			return Part{}, err
		}
		part.Set(cat, value)
	}

	return part, nil
}

// processPart processes a part through the workflow system
func processPart(part Part, workflows Workflows) bool {
	currentWorkflow := workflowStart

workflowLoop:
	for {
		workflow, exists := workflows[currentWorkflow]
		if !exists {
			panic(fmt.Sprintf("workflow %s not found", currentWorkflow))
		}

		for _, rule := range workflow.Rules {
			if rule.Matches(part) {
				switch rule.DestinationWorkflow {
				case workflowAccept:
					return true
				case workflowReject:
					return false
				default:
					currentWorkflow = rule.DestinationWorkflow
					continue workflowLoop
				}
			}
		}
		panic(fmt.Sprintf("no rule matched in workflow %s", currentWorkflow))
	}
}
