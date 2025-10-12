package xiter

import (
	"slices"
	"testing"
)

func TestPermutations(t *testing.T) {
	input := []int{1, 2, 3}
	expected := [][]int{
		{1, 2, 3},
		{2, 1, 3},
		{3, 1, 2},
		{1, 3, 2},
		{2, 3, 1},
		{3, 2, 1},
	}

	result := slices.Collect(Permutations(input))

	if len(result) != len(expected) {
		t.Fatalf("expected %d permutations, got %d", len(expected), len(result))
	}

	// Check that all expected permutations are present
	for _, exp := range expected {
		found := false
		for _, res := range result {
			if slices.Equal(exp, res) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected permutation %v not found in results", exp)
		}
	}
}

func TestPermutationsEmpty(t *testing.T) {
	input := []int{}
	result := slices.Collect(Permutations(input))

	if len(result) != 0 {
		t.Errorf("expected 0 permutations for empty slice, got %d", len(result))
	}
}

func TestPermutationsSingle(t *testing.T) {
	input := []int{42}
	result := slices.Collect(Permutations(input))

	if len(result) != 1 {
		t.Fatalf("expected 1 permutation for single element, got %d", len(result))
	}

	if !slices.Equal(result[0], input) {
		t.Errorf("expected %v, got %v", input, result[0])
	}
}
