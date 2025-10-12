package search

import "testing"

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name      string
		min       int
		max       int
		predicate func(int) bool
		want      int
	}{
		{
			name: "find threshold at 5",
			min:  0,
			max:  10,
			predicate: func(x int) bool {
				return x >= 5
			},
			want: 5,
		},
		{
			name: "predicate true at min",
			min:  0,
			max:  10,
			predicate: func(x int) bool {
				return true
			},
			want: 0,
		},
		{
			name: "predicate never true",
			min:  0,
			max:  10,
			predicate: func(x int) bool {
				return false
			},
			want: 11,
		},
		{
			name: "single element range - true",
			min:  5,
			max:  5,
			predicate: func(x int) bool {
				return x >= 5
			},
			want: 5,
		},
		{
			name: "single element range - false",
			min:  5,
			max:  5,
			predicate: func(x int) bool {
				return x >= 6
			},
			want: 6,
		},
		{
			name: "find threshold at max",
			min:  0,
			max:  10,
			predicate: func(x int) bool {
				return x >= 10
			},
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BinarySearch(tt.min, tt.max, tt.predicate)
			if got != tt.want {
				t.Errorf("BinarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
