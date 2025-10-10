package aoc2020

import (
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
	"github.com/jacoelho/advent-of-code-go/internal/xslices"
)

var foodRegex = regexp.MustCompile(`^(.+) \(contains (.+)\)$`)

type food struct {
	ingredients []string
	allergens   collections.Set[string]
}

func parseFood(r io.Reader) ([]food, error) {
	s := scanner.NewScanner(r, func(line []byte) (food, error) {
		matches := foodRegex.FindSubmatch(line)
		if matches == nil {
			return food{}, fmt.Errorf("invalid food line: %s", string(line))
		}

		ingredients := strings.Fields(string(matches[1]))
		allergensList := strings.Split(string(matches[2]), ", ")

		return food{
			ingredients: ingredients,
			allergens:   collections.NewSet(allergensList...),
		}, nil
	})

	return slices.Collect(s.Values()), s.Err()
}

func allergenCandidates(foods []food) map[string]collections.Set[string] {
	allergenCandidates := make(map[string]collections.Set[string])

	for _, f := range foods {
		ingredientSet := collections.NewSet(f.ingredients...)

		for allergen := range f.allergens.Iter() {
			if existing, ok := allergenCandidates[allergen]; ok {
				allergenCandidates[allergen] = existing.Intersect(ingredientSet)
			} else {
				allergenCandidates[allergen] = ingredientSet.Clone()
			}
		}
	}

	return allergenCandidates
}

func day21p01(r io.Reader) (string, error) {
	foods, err := parseFood(r)
	if err != nil {
		return "", err
	}

	possibleAllergenIngredients := collections.NewSet[string]()
	for _, candidates := range allergenCandidates(foods) {
		possibleAllergenIngredients = possibleAllergenIngredients.Union(candidates)
	}

	count := xslices.Sum(xslices.Map(func(f food) int {
		return xslices.CountFunc(func(ingredient string) bool {
			return !possibleAllergenIngredients.Contains(ingredient)
		}, f.ingredients)
	}, foods))

	return strconv.Itoa(count), nil
}

func day21p02(r io.Reader) (string, error) {
	foods, err := parseFood(r)
	if err != nil {
		return "", err
	}

	allergenCandidates := allergenCandidates(foods)

	// resolve using constraint propagation
	allergenToIngredient := make(map[string]string)

	for len(allergenCandidates) > 0 {
		pair, found := xmaps.Find(allergenCandidates, func(_ string, candidates collections.Set[string]) bool {
			return candidates.Len() == 1
		})
		if !found {
			break
		}

		ingredient, _ := pair.V.Next()
		allergenToIngredient[pair.K] = ingredient
		delete(allergenCandidates, pair.K)

		for a := range allergenCandidates {
			allergenCandidates[a].Remove(ingredient)
		}
	}

	allergens := slices.Collect(maps.Keys(allergenToIngredient))
	slices.Sort(allergens)
	ingredients := xslices.Map(func(allergen string) string { return allergenToIngredient[allergen] }, allergens)

	return strings.Join(ingredients, ","), nil
}
