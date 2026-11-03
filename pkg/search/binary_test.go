package search

import (
	"errors"
	"math"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	for _, tt := range []struct {
		name     string
		min, max int64
		want     int64
	}{
		{"interior", 0, 10, 5},
		{"minimum", 0, 10, 0},
		{"maximum", 0, 10, 10},
		{"singleton", 5, 5, 5},
		{"full signed range", math.MinInt64, math.MaxInt64, 0},
		{"largest signed value", math.MinInt64, math.MaxInt64, math.MaxInt64},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BinarySearch(tt.min, tt.max, func(value int64) (bool, error) {
				return value >= tt.want, nil
			})
			if err != nil || got != tt.want {
				t.Fatalf("BinarySearch() = %v, %v; want %v", got, err, tt.want)
			}
		})
	}
	got, err := BinarySearch(uint64(0), uint64(math.MaxUint64), func(value uint64) (bool, error) {
		return value == math.MaxUint64, nil
	})
	if err != nil || got != math.MaxUint64 {
		t.Fatalf("unsigned search = %v, %v; want MaxUint64", got, err)
	}
}

func TestBinarySearchNoMatch(t *testing.T) {
	for _, interval := range [][2]int64{{0, 10}, {5, 5}, {0, math.MaxInt64}} {
		_, err := BinarySearch(interval[0], interval[1], func(int64) (bool, error) {
			return false, nil
		})
		if err == nil {
			t.Errorf("expected no-match error for %v", interval)
		}
	}
	_, err := BinarySearch(uint64(0), uint64(math.MaxUint64), func(uint64) (bool, error) {
		return false, nil
	})
	if err == nil {
		t.Fatal("expected unsigned no-match error")
	}
}

func TestBinarySearchInvalidInterval(t *testing.T) {
	_, err := BinarySearch(2, 1, func(int) (bool, error) {
		t.Fatal("predicate called for invalid interval")
		return false, nil
	})
	if err == nil {
		t.Fatal("expected invalid interval error")
	}
}

func TestBinarySearchPredicateErrors(t *testing.T) {
	want := errors.New("query failed")
	for _, threshold := range []int{0, 5, 10} {
		_, err := BinarySearch(0, 10, func(value int) (bool, error) {
			if value == threshold {
				return false, want
			}
			return value > threshold, nil
		})
		if !errors.Is(err, want) {
			t.Errorf("threshold %d: error = %v; want %v", threshold, err, want)
		}
	}
}
