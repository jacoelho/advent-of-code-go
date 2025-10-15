package aoc2019

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type Chemical struct {
	name     string
	quantity int
}

type Reaction struct {
	output Chemical
	inputs []Chemical
}

func parseChemical(s string) (Chemical, error) {
	s = strings.TrimSpace(s)
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return Chemical{}, fmt.Errorf("invalid chemical format: %s", s)
	}
	quantity, err := strconv.Atoi(parts[0])
	if err != nil {
		return Chemical{}, err
	}
	return Chemical{name: parts[1], quantity: quantity}, nil
}

func parseReaction(line string) (Reaction, error) {
	parts := strings.Split(line, " => ")
	if len(parts) != 2 {
		return Reaction{}, fmt.Errorf("invalid reaction format: %s", line)
	}

	output, err := parseChemical(parts[1])
	if err != nil {
		return Reaction{}, err
	}

	inputParts := strings.Split(parts[0], ", ")
	inputs := make([]Chemical, 0, len(inputParts))
	for _, inputStr := range inputParts {
		input, err := parseChemical(inputStr)
		if err != nil {
			return Reaction{}, err
		}
		inputs = append(inputs, input)
	}

	return Reaction{output: output, inputs: inputs}, nil
}

func parseReactions(r io.Reader) (map[string]Reaction, error) {
	s := scanner.NewScanner(r, func(b []byte) (string, error) {
		return string(b), nil
	})

	reactions := make(map[string]Reaction)
	for line := range s.Values() {
		reaction, err := parseReaction(line)
		if err != nil {
			return nil, err
		}
		reactions[reaction.output.name] = reaction
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return reactions, nil
}

func calculateOreForFuel(reactions map[string]Reaction, fuelAmount int) int {
	needed := map[string]int{"FUEL": fuelAmount}
	leftover := make(map[string]int)

	for len(needed) > 0 {
		var current string
		for chem, amount := range needed {
			if chem != "ORE" && amount > 0 {
				current = chem
				break
			}
		}

		if current == "" {
			break
		}

		need := needed[current]
		delete(needed, current)

		if leftover[current] >= need {
			leftover[current] -= need
			continue
		}

		need -= leftover[current]
		leftover[current] = 0

		reaction := reactions[current]
		timesNeeded := (need + reaction.output.quantity - 1) / reaction.output.quantity
		produced := timesNeeded * reaction.output.quantity

		if produced > need {
			leftover[current] = produced - need
		}

		for _, input := range reaction.inputs {
			needed[input.name] += timesNeeded * input.quantity
		}
	}

	return needed["ORE"]
}

func day14p01(r io.Reader) (string, error) {
	reactions, err := parseReactions(r)
	if err != nil {
		return "", err
	}

	ore := calculateOreForFuel(reactions, 1)
	return strconv.Itoa(ore), nil
}

func day14p02(r io.Reader) (string, error) {
	reactions, err := parseReactions(r)
	if err != nil {
		return "", err
	}

	const totalOre = 1_000_000_000_000

	low, high := 0, totalOre
	for low < high {
		mid := (low + high + 1) / 2
		ore := calculateOreForFuel(reactions, mid)
		if ore <= totalOre {
			low = mid
		} else {
			high = mid - 1
		}
	}

	return strconv.Itoa(low), nil
}
