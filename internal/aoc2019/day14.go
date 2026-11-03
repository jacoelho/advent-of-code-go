package aoc2019

import (
	"fmt"
	"io"
	"iter"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/graph"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

type chemical struct {
	name     string
	quantity int64
}

type reaction struct {
	output chemical
	inputs []chemical
}

type reactionSystem struct {
	recipes map[string]reaction
	order   []string
}

func parseChemical(s string) (chemical, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return chemical{}, fmt.Errorf("invalid chemical %q", s)
	}
	quantity, err := convert.ScanNumber[int64]([]byte(parts[0]))
	if err != nil {
		return chemical{}, err
	}
	if quantity <= 0 || len(parts[1]) > 128 {
		return chemical{}, fmt.Errorf("chemical needs a positive quantity and bounded name")
	}
	for _, c := range parts[1] {
		if c < 'A' || c > 'Z' {
			return chemical{}, fmt.Errorf("invalid chemical name %q", parts[1])
		}
	}
	return chemical{name: parts[1], quantity: quantity}, nil
}

func parseReaction(line string) (reaction, error) {
	inputs, output, ok := strings.Cut(line, " => ")
	if !ok {
		return reaction{}, fmt.Errorf("invalid reaction %q", line)
	}
	product, err := parseChemical(output)
	if err != nil {
		return reaction{}, err
	}
	result := reaction{output: product}
	for input := range strings.SplitSeq(inputs, ",") {
		ingredient, err := parseChemical(input)
		if err != nil {
			return reaction{}, err
		}
		result.inputs = append(result.inputs, ingredient)
	}
	return result, nil
}

func parseReactions(r io.Reader) (reactionSystem, error) {
	const inputLimit = 4 << 20
	limited := &io.LimitedReader{R: r, N: inputLimit + 1}
	s := scanner.NewScanner(limited, func(line []byte) (reaction, error) {
		return parseReaction(string(line))
	})
	system := reactionSystem{recipes: make(map[string]reaction)}
	var names []string
	ingredients := 0
	for recipe := range s.Values() {
		name := recipe.output.name
		if _, exists := system.recipes[name]; exists || name == "ORE" {
			return reactionSystem{}, fmt.Errorf("duplicate or invalid recipe for %s", name)
		}
		ingredients += len(recipe.inputs)
		if len(names) == 4096 || ingredients > 65_536 {
			return reactionSystem{}, fmt.Errorf("reaction system exceeds recipe or ingredient limit")
		}
		system.recipes[name] = recipe
		names = append(names, name)
	}
	if err := s.Err(); err != nil {
		return reactionSystem{}, err
	}
	if limited.N == 0 {
		return reactionSystem{}, fmt.Errorf("reactions exceed 4 MiB")
	}
	if _, exists := system.recipes["FUEL"]; !exists {
		return reactionSystem{}, fmt.Errorf("missing FUEL recipe")
	}

	for _, name := range names {
		for _, ingredient := range system.recipes[name].inputs {
			if _, exists := system.recipes[ingredient.name]; !exists && ingredient.name != "ORE" {
				return reactionSystem{}, fmt.Errorf("missing recipe for %s", ingredient.name)
			}
		}
	}
	// Demand flows from consumers to their ingredients, opposite evaluation order.
	order, err := graph.TopologicalOrder(names, func(name string) iter.Seq[string] {
		return func(yield func(string) bool) {
			for _, ingredient := range system.recipes[name].inputs {
				if ingredient.name != "ORE" && !yield(ingredient.name) {
					return
				}
			}
		}
	})
	if err != nil {
		return reactionSystem{}, fmt.Errorf("reaction dependencies: %w", err)
	}
	system.order = order
	return system, nil
}

func (system reactionSystem) oreForFuel(fuel int64) (int64, error) {
	needed := map[string]int64{"FUEL": fuel}
	// Combine all consumers' demand before rounding an ingredient to whole batches.
	for _, name := range system.order {
		need := needed[name]
		if need == 0 {
			continue
		}
		recipe := system.recipes[name]
		batches := need / recipe.output.quantity
		if need%recipe.output.quantity != 0 {
			batches++
		}
		for _, ingredient := range recipe.inputs {
			if batches > math.MaxInt64/ingredient.quantity {
				return 0, fmt.Errorf("reaction demand for %s overflows int64", ingredient.name)
			}
			quantity := batches * ingredient.quantity
			if needed[ingredient.name] > math.MaxInt64-quantity {
				return 0, fmt.Errorf("reaction demand for %s overflows int64", ingredient.name)
			}
			needed[ingredient.name] += quantity
		}
	}
	return needed["ORE"], nil
}

func day14p01(r io.Reader) (string, error) {
	system, err := parseReactions(r)
	if err != nil {
		return "", err
	}
	ore, err := system.oreForFuel(1)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(ore, 10), nil
}

func day14p02(r io.Reader) (string, error) {
	system, err := parseReactions(r)
	if err != nil {
		return "", err
	}

	const totalOre = 1_000_000_000_000
	low, high := int64(0), int64(1)
	for {
		ore, err := system.oreForFuel(high)
		if err != nil {
			return "", err
		}
		if ore > totalOre {
			break
		}
		low = high
		if high == math.MaxInt64 {
			return "", fmt.Errorf("fuel search has no infeasible int64 upper bound")
		}
		if high > math.MaxInt64/2 {
			high = math.MaxInt64
		} else {
			high *= 2
		}
	}
	firstInfeasible, err := search.BinarySearch(low+1, high, func(fuel int64) (bool, error) {
		ore, err := system.oreForFuel(fuel)
		return ore > totalOre, err
	})
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(firstInfeasible-1, 10), nil
}
