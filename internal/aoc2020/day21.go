package aoc2020

import (
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
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

	count := xiter.FromSlice(foods).SumBy(func(f food) int {
		return xiter.FromSlice(f.ingredients).Filter(func(ingredient string) bool {
			return !possibleAllergenIngredients.Contains(ingredient)
		}).Count()
	})

	return strconv.Itoa(count), nil
}

func day21p02(r io.Reader) (string, error) {
	foods, err := parseFood(r)
	if err != nil {
		return "", err
	}

	allergenCandidates := allergenCandidates(foods)

	allergenToIngredient := make(map[string]string)

	for len(allergenCandidates) > 0 {
		for allergen, candidates := range allergenCandidates {
			if candidates.Len() == 0 {
				return "", fmt.Errorf("allergen %s has no possible ingredient", allergen)
			}
		}
		allergen, found := xiter.From(maps.Keys(allergenCandidates)).Filter(func(name string) bool {
			return allergenCandidates[name].Len() == 1
		}).First()
		if !found {
			return "", fmt.Errorf("cannot resolve allergen assignment from non-singleton candidates")
		}

		ingredient, _ := allergenCandidates[allergen].Next()
		allergenToIngredient[allergen] = ingredient
		delete(allergenCandidates, allergen)

		for a := range allergenCandidates {
			allergenCandidates[a].Remove(ingredient)
		}
	}

	allergens := slices.Sorted(maps.Keys(allergenToIngredient))
	ingredients := make([]string, len(allergens))
	for i, allergen := range allergens {
		ingredients[i] = allergenToIngredient[allergen]
	}

	return strings.Join(ingredients, ","), nil
}
