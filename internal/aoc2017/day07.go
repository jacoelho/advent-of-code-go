package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type towerProgram struct {
	weight   int
	children []string
}

type programTower struct {
	root     string
	programs map[string]towerProgram
	totals   map[string]int
	balanced map[string]bool
}

func parseTower(r io.Reader) (programTower, error) {
	tower := programTower{
		programs: make(map[string]towerProgram),
		totals:   make(map[string]int),
		balanced: make(map[string]bool),
	}
	s := bufio.NewScanner(r)
	for s.Scan() {
		body, children, hasChildren := strings.Cut(s.Text(), " -> ")
		fields := strings.Fields(body)
		if len(fields) != 2 || !strings.HasPrefix(fields[1], "(") || !strings.HasSuffix(fields[1], ")") {
			return tower, fmt.Errorf("invalid tower program %q", s.Text())
		}
		weight, err := convert.ScanNumber[int]([]byte(fields[1][1 : len(fields[1])-1]))
		if err != nil {
			return tower, err
		}
		if weight < 0 {
			return tower, fmt.Errorf("negative program weight")
		}
		if _, exists := tower.programs[fields[0]]; exists {
			return tower, fmt.Errorf("duplicate program %q", fields[0])
		}
		program := towerProgram{weight: weight}
		if hasChildren {
			program.children = strings.Split(children, ", ")
		}
		tower.programs[fields[0]] = program
	}
	if err := s.Err(); err != nil {
		return tower, err
	}

	parents := make(map[string]string)
	for name, program := range tower.programs {
		for _, child := range program.children {
			if _, exists := tower.programs[child]; !exists {
				return tower, fmt.Errorf("undefined child %q", child)
			}
			if _, exists := parents[child]; exists {
				return tower, fmt.Errorf("program %q has multiple parents", child)
			}
			parents[child] = name
		}
	}

	for name := range tower.programs {
		if _, hasParent := parents[name]; !hasParent {
			if tower.root != "" {
				return tower, fmt.Errorf("tower has multiple roots")
			}
			tower.root = name
		}
	}
	if tower.root == "" {
		return tower, fmt.Errorf("tower has no root")
	}

	seen := make(map[string]bool)
	var total func(string) (int, error)
	total = func(name string) (int, error) {
		if seen[name] {
			return 0, fmt.Errorf("cycle at program %q", name)
		}
		seen[name] = true
		program := tower.programs[name]
		weight := program.weight
		balanced := true
		for i, child := range program.children {
			childWeight, err := total(child)
			if err != nil {
				return 0, err
			}
			if !tower.balanced[child] || i > 0 && childWeight != tower.totals[program.children[0]] {
				balanced = false
			}
			weight += childWeight
		}
		tower.totals[name] = weight
		tower.balanced[name] = balanced
		return weight, nil
	}
	if _, err := total(tower.root); err != nil {
		return tower, err
	}
	if len(seen) != len(tower.programs) {
		return tower, fmt.Errorf("tower contains disconnected programs")
	}

	return tower, nil
}

func (tower programTower) balancedCorrections(current string, delta int) (int, int) {
	var weight, count int
	for {
		program := tower.programs[current]
		if corrected := program.weight + delta; corrected >= 0 {
			weight = corrected
			count++
		}
		// Changing a descendant preserves balance only through a single-child chain.
		if count == 2 || len(program.children) != 1 {
			return weight, count
		}
		current = program.children[0]
	}
}

func (tower programTower) correctedWeight() (int, error) {
	current := tower.root
	var expected *int
	for {
		program := tower.programs[current]
		if tower.balanced[current] {
			if expected == nil {
				return 0, fmt.Errorf("tower is already balanced")
			}
			weight, count := tower.balancedCorrections(current, *expected-tower.totals[current])
			if count != 1 {
				return 0, fmt.Errorf("no unique nonnegative correction below %q", current)
			}
			return weight, nil
		}
		if len(program.children) == 1 {
			if expected != nil {
				expected = new(*expected - program.weight)
			}
			current = program.children[0]
			continue
		}

		groups := make(map[int][]string)
		for _, child := range program.children {
			groups[tower.totals[child]] = append(groups[tower.totals[child]], child)
		}
		if len(groups) != 2 {
			return 0, fmt.Errorf("ambiguous child weights at %q", current)
		}
		var wanted int
		known := false
		for weight, children := range groups {
			if len(children) > 1 {
				if known {
					return 0, fmt.Errorf("multiple unbalanced children at %q", current)
				}
				wanted, known = weight, true
			}
		}
		if !known {
			if expected != nil {
				// The parent's expected total can distinguish two unequal child weights.
				if (*expected-program.weight)%len(program.children) != 0 {
					return 0, fmt.Errorf("ambiguous correction at %q", current)
				}
				wanted = (*expected - program.weight) / len(program.children)
			} else {
				left, right := program.children[0], program.children[1]
				if tower.balanced[left] && tower.balanced[right] {
					delta := tower.totals[right] - tower.totals[left]
					leftWeight, leftCount := tower.balancedCorrections(left, delta)
					rightWeight, rightCount := tower.balancedCorrections(right, -delta)
					if leftCount+rightCount != 1 {
						return 0, fmt.Errorf("no unique nonnegative correction below %q", current)
					}
					if leftCount == 1 {
						return leftWeight, nil
					}
					return rightWeight, nil
				}
				if !tower.balanced[left] && !tower.balanced[right] {
					return 0, fmt.Errorf("multiple unbalanced children at %q", current)
				}
				wanted = tower.totals[left]
				if tower.balanced[right] {
					wanted = tower.totals[right]
				}
			}
		}
		if expected != nil && program.weight+wanted*len(program.children) != *expected {
			return 0, fmt.Errorf("inconsistent correction at %q", current)
		}
		wrong := ""
		for _, child := range program.children {
			if tower.totals[child] != wanted {
				if wrong != "" {
					return 0, fmt.Errorf("multiple unbalanced children at %q", current)
				}
				wrong = child
			} else if !tower.balanced[child] {
				return 0, fmt.Errorf("multiple unbalanced children at %q", current)
			}
		}
		if wrong == "" {
			return 0, fmt.Errorf("no unbalanced child at %q", current)
		}
		current, expected = wrong, new(wanted)
	}
}

func day07p01(r io.Reader) (string, error) {
	tower, err := parseTower(r)
	if err != nil {
		return "", err
	}

	return tower.root, nil
}

func day07p02(r io.Reader) (string, error) {
	tower, err := parseTower(r)
	if err != nil {
		return "", err
	}

	weight, err := tower.correctedWeight()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(weight), nil
}
