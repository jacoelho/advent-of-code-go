package search

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// BinarySearch finds the smallest value in [min, max] where predicate returns true.
// The predicate must be false below a threshold and true at or above it.
// An invalid interval is rejected before calling predicate. No match and predicate
// failures return errors; predicate errors remain available through errors.Is/As.
func BinarySearch[T constraints.Integer](min, max T, predicate func(T) (bool, error)) (T, error) {
	if min > max {
		return 0, fmt.Errorf("invalid search interval [%v, %v]", min, max)
	}
	match, err := predicate(max)
	if err != nil {
		return 0, fmt.Errorf("search predicate at %v: %w", max, err)
	}
	if !match {
		return 0, fmt.Errorf("no match in search interval [%v, %v]", min, max)
	}

	for min < max {
		// This floor average cannot overflow, even across the full signed range.
		mid := (min & max) + ((min ^ max) >> 1)
		match, err := predicate(mid)
		if err != nil {
			return 0, fmt.Errorf("search predicate at %v: %w", mid, err)
		}
		if match {
			max = mid
		} else {
			min = mid + 1
		}
	}
	return min, nil
}
