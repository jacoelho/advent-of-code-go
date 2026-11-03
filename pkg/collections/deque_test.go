package collections

import (
	"slices"
	"testing"
)

func TestDequePopsReleaseReferences(t *testing.T) {
	a, b, c, d, e, f := new(1), new(2), new(3), new(4), new(5), new(6)
	queue := NewDeque[*int](2)
	for _, value := range []*int{a, b, c, d} {
		queue.PushBack(value)
	}
	var removed []*int
	checkPop := func(pop func() (*int, bool), want *int, remaining []*int) {
		t.Helper()
		got, ok := pop()
		if !ok || got != want {
			t.Fatalf("pop = %v, %v; want %v", got, ok, want)
		}
		removed = append(removed, got)
		if live := slices.Collect(queue.IterFront()); !slices.Equal(live, remaining) {
			t.Fatalf("surviving values = %v, want %v", live, remaining)
		}
		for _, held := range queue.data {
			if slices.Contains(removed, held) {
				t.Fatal("deque still retains a removed value")
			}
		}
	}
	checkPop(queue.PopFront, a, []*int{b, c, d})
	checkPop(queue.PopFront, b, []*int{c, d})
	queue.PushBack(e)
	queue.PushBack(f)
	checkPop(queue.PopBack, f, []*int{c, d, e})
	checkPop(queue.PopFront, c, []*int{d, e})
	checkPop(queue.PopBack, e, []*int{d})
	checkPop(queue.PopFront, d, nil)
	for _, pop := range []func() (*int, bool){queue.PopFront, queue.PopBack} {
		if got, ok := pop(); ok || got != nil {
			t.Fatalf("empty pop = %v, %v", got, ok)
		}
	}
}
