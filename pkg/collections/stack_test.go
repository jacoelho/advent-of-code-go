package collections

import "testing"

func TestStackPopReleasesReferences(t *testing.T) {
	a, b, c := new(1), new(2), new(3)
	stack := NewStack(a, b, c)
	for i, want := range []*int{c, b, a} {
		got, ok := stack.Pop()
		if !ok || got != want {
			t.Fatalf("pop %d: got %v, %v; want %v", i, got, ok, want)
		}
		for _, held := range (*stack)[:cap(*stack)] {
			if held == got {
				t.Fatalf("pop %d: stack still retains removed value", i)
			}
		}
		if i < 2 {
			if top, ok := stack.Peek(); !ok || top != []*int{b, a}[i] {
				t.Fatalf("pop %d: surviving top = %v, %v", i, top, ok)
			}
		}
	}
	if got, ok := stack.Pop(); ok || got != nil {
		t.Fatalf("empty pop = %v, %v", got, ok)
	}
}
