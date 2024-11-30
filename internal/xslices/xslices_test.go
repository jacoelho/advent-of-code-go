package xslices

import (
	"reflect"
	"testing"
)

func TestWindow(t *testing.T) {
	cases := []struct {
		name   string
		s      []int
		n      int
		chunks [][]int
	}{
		{
			name:   "nil",
			s:      nil,
			n:      1,
			chunks: nil,
		},
		{
			name:   "empty",
			s:      []int{},
			n:      1,
			chunks: nil,
		},
		{
			name:   "short",
			s:      []int{1, 2},
			n:      3,
			chunks: [][]int{{1, 2}},
		},
		{
			name:   "one",
			s:      []int{1, 2},
			n:      2,
			chunks: [][]int{{1, 2}},
		},
		{
			name:   "even",
			s:      []int{1, 2, 3, 4},
			n:      2,
			chunks: [][]int{{1, 2}, {2, 3}, {3, 4}},
		},
		{
			name:   "odd",
			s:      []int{1, 2, 3, 4, 5},
			n:      2,
			chunks: [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var chunks [][]int
			for c := range Window(tc.s, tc.n) {
				chunks = append(chunks, c)
			}

			if !reflect.DeepEqual(chunks, tc.chunks) {
				t.Errorf("got %v, want %v", chunks, tc.chunks)
			}
		})
	}
}
