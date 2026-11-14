package collections

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestHeapPopsInLessOrderAcrossInterleavedOperations(t *testing.T) {
	random := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		h := NewHeap(func(a, b int) bool { return a < b })
		var reference []int
		for range random.IntN(64) {
			if len(reference) > 0 && random.IntN(3) == 0 {
				slices.Sort(reference)
				want := reference[0]
				reference = reference[1:]
				if got, ok := h.Pop(); !ok || got != want {
					t.Fatalf("trial %d: Pop = %d, %v; want %d", trial, got, ok, want)
				}
				continue
			}
			value := random.IntN(16)
			h.Push(value)
			reference = append(reference, value)
		}
		slices.Sort(reference)
		if h.Len() != len(reference) {
			t.Fatalf("trial %d: Len = %d; want %d", trial, h.Len(), len(reference))
		}
		if len(reference) > 0 {
			if top, ok := h.Peek(); !ok || top != reference[0] {
				t.Fatalf("trial %d: Peek = %d, %v; want %d", trial, top, ok, reference[0])
			}
		}
		if got := slices.Collect(h.PopSeq()); !slices.Equal(got, reference) {
			t.Fatalf("trial %d: drained %v; want %v", trial, got, reference)
		}
	}
}

func TestHeapEmpty(t *testing.T) {
	h := NewHeap(func(a, b int) bool { return a < b })
	if v, ok := h.Pop(); ok || v != 0 {
		t.Fatalf("Pop on empty heap = %d, %v", v, ok)
	}
	if v, ok := h.Peek(); ok || v != 0 {
		t.Fatalf("Peek on empty heap = %d, %v", v, ok)
	}
}

func TestHeapPopSeqStopsWithoutDrainingRemainder(t *testing.T) {
	h := NewHeap(func(a, b int) bool { return a < b })
	for _, v := range []int{5, 1, 4, 2, 3} {
		h.Push(v)
	}
	for v := range h.PopSeq() {
		if v == 2 {
			break
		}
	}
	if got := slices.Collect(h.PopSeq()); !slices.Equal(got, []int{3, 4, 5}) {
		t.Fatalf("remaining = %v; want [3 4 5]", got)
	}
}

func TestHeapPopReleasesReferences(t *testing.T) {
	h := NewHeap(func(a, b *int) bool { return *a < *b })
	values := []*int{new(3), new(1), new(2)}
	for _, v := range values {
		h.Push(v)
	}
	for range values {
		got, _ := h.Pop()
		for _, held := range h.items[:cap(h.items)] {
			if held == got {
				t.Fatalf("heap still retains popped value %d", *got)
			}
		}
	}
}
