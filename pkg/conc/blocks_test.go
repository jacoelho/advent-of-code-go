package conc

import (
	"math"
	"slices"
	"testing"
)

func TestBlocksYieldsEveryIndexInOrder(t *testing.T) {
	identity := func(start, end int) []int {
		var indices []int
		for i := start; i < end; i++ {
			indices = append(indices, i)
		}
		return indices
	}
	var got []int
	for i := range Blocks(5, 3, identity) {
		got = append(got, i)
		if len(got) == 1000 {
			break
		}
	}
	for i, v := range got {
		if v != 5+i {
			t.Fatalf("value %d is %d; want %d", i, v, 5+i)
		}
	}
}

func TestBlocksEndsAtMaxInt(t *testing.T) {
	var blocks [][2]int
	for b := range Blocks(math.MaxInt-10, 4, func(start, end int) [][2]int { return [][2]int{{start, end}} }) {
		blocks = append(blocks, b)
	}
	want := [][2]int{{math.MaxInt - 10, math.MaxInt - 6}, {math.MaxInt - 6, math.MaxInt - 2}, {math.MaxInt - 2, math.MaxInt}}
	if !slices.Equal(blocks, want) {
		t.Fatalf("got %v; want %v", blocks, want)
	}
}

func TestBlocksFromNegativeStart(t *testing.T) {
	var blocks [][2]int
	for b := range Blocks(-5, 4, func(start, end int) [][2]int { return [][2]int{{start, end}} }) {
		blocks = append(blocks, b)
		if len(blocks) == 3 {
			break
		}
	}
	want := [][2]int{{-5, -1}, {-1, 3}, {3, 7}}
	if !slices.Equal(blocks, want) {
		t.Fatalf("got %v; want %v", blocks, want)
	}
}
