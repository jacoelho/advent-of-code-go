// ABOUTME: Binary search for finding boundaries in sorted/monotonic ranges.
// ABOUTME: Finds the smallest value where a predicate transitions from false to true.
package search

import "golang.org/x/exp/constraints"

// BinarySearch finds the smallest value in [min, max] where predicate returns true.
// the predicate must be monotonic: false for all values below some threshold,
// and true for all values at or above that threshold.
func BinarySearch[T constraints.Integer](min, max T, predicate func(T) bool) T {
	if predicate(min) {
		return min
	}
	if !predicate(max) {
		return max + 1
	}

	for min < max {
		mid := min + (max-min)/2
		if predicate(mid) {
			max = mid
		} else {
			min = mid + 1
		}
	}
	return min
}
