package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type turingTransition struct {
	write bool
	move  int
	next  byte
}

type turingBlueprint struct {
	initial byte
	steps   int
	rules   [26][2]turingTransition
}

func blueprintField(line, prefix, suffix string) (string, error) {
	value, ok := strings.CutPrefix(line, prefix)
	if !ok {
		return "", fmt.Errorf("invalid blueprint line %q", line)
	}
	value, ok = strings.CutSuffix(value, suffix)
	if !ok || value == "" {
		return "", fmt.Errorf("invalid blueprint line %q", line)
	}

	return value, nil
}

func blueprintState(line, prefix, suffix string) (byte, error) {
	value, err := blueprintField(line, prefix, suffix)
	if err != nil {
		return 0, err
	}
	if len(value) != 1 || value[0] < 'A' || value[0] > 'Z' {
		return 0, fmt.Errorf("invalid blueprint state %q", value)
	}

	return value[0] - 'A', nil
}

func parseBlueprint(r io.Reader) (turingBlueprint, error) {
	var lines []string
	s := bufio.NewScanner(r)
	for s.Scan() {
		if line := strings.TrimSpace(s.Text()); line != "" {
			lines = append(lines, line)
		}
	}
	if err := s.Err(); err != nil {
		return turingBlueprint{}, err
	}
	if len(lines) < 11 || (len(lines)-2)%9 != 0 {
		return turingBlueprint{}, fmt.Errorf("incomplete Turing blueprint")
	}

	var blueprint turingBlueprint
	initial, err := blueprintState(lines[0], "Begin in state ", ".")
	if err != nil {
		return turingBlueprint{}, err
	}
	count, err := blueprintField(lines[1], "Perform a diagnostic checksum after ", " steps.")
	if err != nil {
		return turingBlueprint{}, err
	}
	steps, err := convert.ScanNumber[int]([]byte(count))
	if err != nil {
		return turingBlueprint{}, err
	}
	if steps < 0 || steps > 100_000_000 {
		return turingBlueprint{}, fmt.Errorf("turing step count must be between 0 and 100000000")
	}
	blueprint.initial, blueprint.steps = initial, steps

	var defined [26]bool
	for i := 2; i < len(lines); i += 9 {
		state, err := blueprintState(lines[i], "In state ", ":")
		if err != nil {
			return turingBlueprint{}, err
		}
		if defined[state] {
			return turingBlueprint{}, fmt.Errorf("duplicate Turing state %c", state+'A')
		}
		defined[state] = true
		var cases [2]bool
		for j := range 2 {
			start := i + 1 + 4*j
			value, err := blueprintField(lines[start], "If the current value is ", ":")
			if err != nil || value != "0" && value != "1" {
				return turingBlueprint{}, fmt.Errorf("invalid Turing condition %q", lines[start])
			}
			bit := value[0] - '0'
			if cases[bit] {
				return turingBlueprint{}, fmt.Errorf("duplicate Turing condition %q", lines[start])
			}
			cases[bit] = true
			write, err := blueprintField(lines[start+1], "- Write the value ", ".")
			if err != nil || write != "0" && write != "1" {
				return turingBlueprint{}, fmt.Errorf("invalid Turing write %q", lines[start+1])
			}
			direction, err := blueprintField(lines[start+2], "- Move one slot to the ", ".")
			if err != nil || direction != "left" && direction != "right" {
				return turingBlueprint{}, fmt.Errorf("invalid Turing move %q", lines[start+2])
			}
			move := 1
			if direction == "left" {
				move = -1
			}
			next, err := blueprintState(lines[start+3], "- Continue with state ", ".")
			if err != nil {
				return turingBlueprint{}, err
			}
			blueprint.rules[state][bit] = turingTransition{write: write == "1", move: move, next: next}
		}
	}
	if !defined[blueprint.initial] {
		return turingBlueprint{}, fmt.Errorf("undefined initial Turing state %c", blueprint.initial+'A')
	}
	for state, exists := range defined {
		if !exists {
			continue
		}
		for _, transition := range blueprint.rules[state] {
			if !defined[transition.next] {
				return turingBlueprint{}, fmt.Errorf("undefined next Turing state %c", transition.next+'A')
			}
		}
	}

	return blueprint, nil
}

func day25p01(r io.Reader) (string, error) {
	blueprint, err := parseBlueprint(r)
	if err != nil {
		return "", err
	}

	ones := collections.NewSet[int]()
	pos, state := 0, blueprint.initial
	for range blueprint.steps {
		bit := 0
		if ones.Contains(pos) {
			bit = 1
		}
		transition := blueprint.rules[state][bit]
		if transition.write {
			ones.Add(pos)
		} else {
			ones.Remove(pos)
		}
		pos += transition.move
		state = transition.next
	}

	return strconv.Itoa(ones.Len()), nil
}
