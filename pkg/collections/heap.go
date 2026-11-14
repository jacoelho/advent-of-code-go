package collections

import "iter"

// A Heap is a min-heap backed by a slice.
type Heap[E any] struct {
	items []E
	less  func(E, E) bool
}

// NewHeap constructs a new Heap with a comparison function.
func NewHeap[E any](less func(E, E) bool) *Heap[E] {
	return &Heap[E]{items: make([]E, 0, 10), less: less}
}

// Push pushes an element onto the heap. The complexity is O(log n)
// where n = h.Len().
func (h *Heap[E]) Push(elem E) {
	h.items = append(h.items, elem)
	h.up(len(h.items) - 1)
}

// Pop removes and returns the minimum value according to less, or (zero, false)
// when empty. Its time complexity is O(log n), where n = h.Len().
func (h *Heap[E]) Pop() (E, bool) {
	var zero E
	if len(h.items) == 0 {
		return zero, false
	}
	top := h.items[0]
	last := len(h.items) - 1
	h.items[0] = h.items[last]
	// Clear the vacated slot so the backing array does not retain the value.
	h.items[last] = zero
	h.items = h.items[:last]
	h.down(0)
	return top, true
}

// PopSeq returns an iterator that pops and yields values in order of less
// until the heap is empty. It consumes the heap: values yielded are removed,
// and stopping early leaves the rest in place.
func (h *Heap[E]) PopSeq() iter.Seq[E] {
	return func(yield func(E) bool) {
		for h.Len() > 0 {
			v, _ := h.Pop()
			if !yield(v) {
				return
			}
		}
	}
}

// Peek returns the minimum value according to less without removing it, or
// (zero, false) when empty. Its time complexity is O(1).
func (h *Heap[E]) Peek() (E, bool) {
	if len(h.items) == 0 {
		var zero E
		return zero, false
	}
	return h.items[0], true
}

// Len returns the number of elements in the heap.
func (h *Heap[E]) Len() int {
	return len(h.items)
}

func (h *Heap[E]) up(i int) {
	items := h.items
	for i > 0 {
		parent := (i - 1) / 2
		if !h.less(items[i], items[parent]) {
			return
		}
		items[i], items[parent] = items[parent], items[i]
		i = parent
	}
}

func (h *Heap[E]) down(i int) {
	items := h.items
	for {
		child := 2*i + 1
		if child >= len(items) {
			return
		}
		if right := child + 1; right < len(items) && h.less(items[right], items[child]) {
			child = right
		}
		if !h.less(items[child], items[i]) {
			return
		}
		items[i], items[child] = items[child], items[i]
		i = child
	}
}
