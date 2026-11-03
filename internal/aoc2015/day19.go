package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type medicineRule struct{ from, to string }

func readMedicine(r io.Reader) ([]medicineRule, string, error) {
	var rules []medicineRule
	var molecule string
	readingRules := true
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			readingRules = false
			continue
		}
		if readingRules {
			from, to, ok := strings.Cut(line, " => ")
			if !ok || from == "" || to == "" || strings.ContainsAny(from+to, " \t") {
				return nil, "", fmt.Errorf("invalid replacement %q", line)
			}
			rules = append(rules, medicineRule{from, to})
		} else {
			if molecule != "" || strings.ContainsAny(line, " \t") {
				return nil, "", fmt.Errorf("invalid medicine molecule %q", line)
			}
			molecule = line
		}
	}
	if err := s.Err(); err != nil {
		return nil, "", err
	}
	if len(rules) == 0 || molecule == "" {
		return nil, "", fmt.Errorf("missing replacements or medicine molecule")
	}
	return rules, molecule, nil
}

func medicineProducts(molecule string, rules []medicineRule) iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, rule := range rules {
			for start := 0; start <= len(molecule)-len(rule.from); {
				position := strings.Index(molecule[start:], rule.from)
				if position < 0 {
					break
				}
				position += start
				if !yield(molecule[:position] + rule.to + molecule[position+len(rule.from):]) {
					return
				}
				start = position + 1
			}
		}
	}
}

func day19p01(r io.Reader) (string, error) {
	rules, molecule, err := readMedicine(r)
	if err != nil {
		return "", err
	}
	products := collections.NewSetFromIter(medicineProducts(molecule, rules))
	return strconv.Itoa(products.Len()), nil
}

func medicineAtoms(molecule string) ([]string, bool) {
	if molecule == "e" {
		return []string{"e"}, true
	}
	var atoms []string
	for i := 0; i < len(molecule); {
		start := i
		if molecule[i] < 'A' || molecule[i] > 'Z' {
			return nil, false
		}
		i++
		if i < len(molecule) && molecule[i] >= 'a' && molecule[i] <= 'z' {
			i++
		}
		atoms = append(atoms, molecule[start:i])
	}
	return atoms, len(atoms) > 0
}

func medicineGrammar(rules []medicineRule, molecule string) (map[string]int, [][]string, []string, bool) {
	target, ok := medicineAtoms(molecule)
	if !ok {
		return nil, nil, nil, false
	}
	ids := map[string]int{"e": 0}
	add := func(atom string) {
		if _, exists := ids[atom]; !exists {
			ids[atom] = len(ids)
		}
	}
	for _, atom := range target {
		add(atom)
	}
	right := make([][]string, len(rules))
	for i, rule := range rules {
		left, valid := medicineAtoms(rule.from)
		if !valid || len(left) != 1 {
			return nil, nil, nil, false
		}
		atoms, valid := medicineAtoms(rule.to)
		if !valid {
			return nil, nil, nil, false
		}
		add(rule.from)
		for _, atom := range atoms {
			add(atom)
		}
		right[i] = atoms
	}
	// Only use atom parsing when every literal rule match must coincide with an atom.
	for _, rule := range rules {
		for atom := range ids {
			if atom != rule.from && strings.Contains(atom, rule.from) {
				return nil, nil, nil, false
			}
		}
	}
	return ids, right, target, true
}

func medicineParseCost(rules []medicineRule, ids map[string]int, right [][]string, target []string) (int, error) {
	type binaryRule struct{ parent, left, right, cost int }
	type unitRule struct{ parent, child int }
	var binary []binaryRule
	var units []unitRule
	symbols := len(ids)
	for i, rule := range rules {
		parent := ids[rule.from]
		rhs := right[i]
		if len(rhs) == 1 {
			units = append(units, unitRule{parent, ids[rhs[0]]})
			continue
		}
		cost := 1
		for j := 0; j < len(rhs)-2; j++ {
			binary = append(binary, binaryRule{parent, ids[rhs[j]], symbols, cost})
			parent = symbols
			symbols++
			cost = 0
		}
		binary = append(binary, binaryRule{parent, ids[rhs[len(rhs)-2]], ids[rhs[len(rhs)-1]], cost})
	}

	const infinity = math.MaxInt / 4
	unit := make([][]int, symbols)
	for parent := range unit {
		unit[parent] = make([]int, symbols)
		for child := range symbols {
			unit[parent][child] = infinity
		}
		unit[parent][parent] = 0
	}
	for _, rule := range units {
		unit[rule.parent][rule.child] = min(unit[rule.parent][rule.child], 1)
	}
	for through := range symbols {
		for parent := range symbols {
			for child := range symbols {
				unit[parent][child] = min(unit[parent][child], unit[parent][through]+unit[through][child])
			}
		}
	}

	byLeft := make([][]binaryRule, symbols)
	for _, rule := range binary {
		byLeft[rule.left] = append(byLeft[rule.left], rule)
	}

	width := len(target) + 1
	chart := make([]map[int]int, width*width)
	for length := 1; length <= len(target); length++ {
		for start := 0; start+length <= len(target); start++ {
			end := start + length
			costs := make(map[int]int)
			if length == 1 {
				costs[ids[target[start]]] = 0
			}
			for split := start + 1; split < end; split++ {
				for left, leftCost := range chart[start*width+split] {
					for _, rule := range byLeft[left] {
						rightCost, exists := chart[split*width+end][rule.right]
						if !exists {
							continue
						}
						cost := leftCost + rightCost + rule.cost
						if previous, exists := costs[rule.parent]; !exists || cost < previous {
							costs[rule.parent] = cost
						}
					}
				}
			}
			for child, cost := range costs {
				for parent := range symbols {
					if unit[parent][child] == infinity {
						continue
					}
					total := cost + unit[parent][child]
					if previous, exists := costs[parent]; !exists || total < previous {
						costs[parent] = total
					}
				}
			}
			chart[start*width+end] = costs
		}
	}

	if cost, ok := chart[len(target)][ids["e"]]; ok {
		return cost, nil
	}
	return 0, fmt.Errorf("medicine molecule is unreachable from e")
}

func medicineLiteralCost(rules []medicineRule, target string) (int, error) {
	reversed := make([]medicineRule, len(rules))
	for i, rule := range rules {
		if len(rule.from) > len(rule.to) {
			return 0, fmt.Errorf("literal search requires noncontracting rules: %s => %s", rule.from, rule.to)
		}
		reversed[i] = medicineRule{rule.to, rule.from}
	}
	type step struct {
		molecule string
		cost     int
	}
	queue := collections.NewDeque[step](16)
	queue.PushBack(step{target, 0})
	seen := collections.NewSet(target)
	for queue.Size() > 0 {
		current, _ := queue.PopFront()
		if current.molecule == "e" {
			return current.cost, nil
		}
		for next := range medicineProducts(current.molecule, reversed) {
			if seen.Contains(next) {
				continue
			}
			seen.Add(next)
			queue.PushBack(step{next, current.cost + 1})
		}
	}
	return 0, fmt.Errorf("medicine molecule is unreachable from e")
}

func day19p02(r io.Reader) (string, error) {
	rules, molecule, err := readMedicine(r)
	if err != nil {
		return "", err
	}
	if molecule == "e" {
		return "0", nil
	}
	var cost int
	if ids, right, target, ok := medicineGrammar(rules, molecule); ok {
		cost, err = medicineParseCost(rules, ids, right, target)
	} else {
		cost, err = medicineLiteralCost(rules, molecule)
	}
	if err != nil {
		return "", err
	}
	return strconv.Itoa(cost), nil
}
