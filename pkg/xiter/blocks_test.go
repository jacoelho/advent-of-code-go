package xiter

import (
	"slices"
	"testing"
)

func TestBlocksSplitsOnBlankLines(t *testing.T) {
	for _, tt := range []struct {
		lines []string
		want  [][]string
	}{
		{nil, nil},
		{[]string{"a"}, [][]string{{"a"}}},
		{[]string{"a", "b", "", "c"}, [][]string{{"a", "b"}, {"c"}}},
		{[]string{"a", "", "", "b"}, [][]string{{"a"}, nil, {"b"}}},
		{[]string{"", "a"}, [][]string{nil, {"a"}}},
		{[]string{"a", ""}, [][]string{{"a"}}},
		{[]string{"a", "", ""}, [][]string{{"a"}, nil}},
		{[]string{""}, [][]string{nil}},
	} {
		got := slices.Collect(Blocks(slices.Values(tt.lines)))
		if !slices.EqualFunc(got, tt.want, slices.Equal) {
			t.Errorf("%q: got %q; want %q", tt.lines, got, tt.want)
		}
	}
}

func TestBlocksStopsWhenConsumerStops(t *testing.T) {
	pulled := 0
	lines := func(yield func(string) bool) {
		for _, line := range []string{"a", "", "b", "", "c"} {
			pulled++
			if !yield(line) {
				return
			}
		}
	}
	for range Blocks(lines) {
		break
	}
	if pulled != 2 {
		t.Fatalf("pulled %d lines; want 2", pulled)
	}
}

func TestBlocksDoesNotShareStorageBetweenBlocks(t *testing.T) {
	blocks := slices.Collect(Blocks(slices.Values([]string{"a", "", "b"})))
	blocks[0] = append(blocks[0], "x")
	if !slices.Equal(blocks[1], []string{"b"}) {
		t.Fatalf("second block changed to %q", blocks[1])
	}
}
