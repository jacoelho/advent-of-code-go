package collections

import (
	"math/rand/v2"
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

func TestDequeMatchesSliceModelAcrossGrowthAndWraparound(t *testing.T) {
	random := rand.New(rand.NewPCG(3, 4))
	d := NewDeque[int](1)
	var model []int
	for step := range 5000 {
		switch op := random.IntN(4); {
		case op == 0:
			d.PushBack(step)
			model = append(model, step)
		case op == 1:
			d.PushFront(step)
			model = append([]int{step}, model...)
		case op == 2 && len(model) > 0:
			got, ok := d.PopFront()
			if !ok || got != model[0] {
				t.Fatalf("step %d: PopFront = %d, %v; want %d", step, got, ok, model[0])
			}
			model = model[1:]
		case op == 3 && len(model) > 0:
			got, ok := d.PopBack()
			if !ok || got != model[len(model)-1] {
				t.Fatalf("step %d: PopBack = %d, %v; want %d", step, got, ok, model[len(model)-1])
			}
			model = model[:len(model)-1]
		}
		if d.Size() != len(model) {
			t.Fatalf("step %d: Size = %d; want %d", step, d.Size(), len(model))
		}
	}
	if got := slices.Collect(d.IterFront()); !slices.Equal(got, model) {
		t.Fatalf("IterFront = %v; want %v", got, model)
	}
	slices.Reverse(model)
	if got := slices.Collect(d.IterBack()); !slices.Equal(got, model) {
		t.Fatalf("IterBack = %v; want %v", got, model)
	}
}
