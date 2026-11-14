package collections

import (
	"iter"
)

// Deque represents a double-ended queue implemented as a circular buffer.
type Deque[T any] struct {
	data       []T
	head, tail int
	size       int
}

// NewDeque returns an empty deque with initial capacity max(1, capacity).
func NewDeque[T any](capacity int) *Deque[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &Deque[T]{data: make([]T, capacity)}
}

func (d *Deque[T]) grow() {
	data := make([]T, 2*len(d.data))
	n := copy(data, d.data[d.head:])
	copy(data[n:], d.data[:d.head])
	d.data = data
	d.head = 0
	d.tail = d.size
}

// PushFront adds an element to the front of the deque.
func (d *Deque[T]) PushFront(value T) {
	if d.size == len(d.data) {
		d.grow()
	}
	d.head = (d.head - 1 + len(d.data)) % len(d.data)
	d.data[d.head] = value
	d.size++
}

// PushBack adds an element to the back of the deque.
func (d *Deque[T]) PushBack(value T) {
	if d.size == len(d.data) {
		d.grow()
	}
	d.data[d.tail] = value
	d.tail = (d.tail + 1) % len(d.data)
	d.size++
}

// PopFront removes and returns the front value, or (zero, false) when empty. It
// clears the removed slot so unused storage does not retain references.
func (d *Deque[T]) PopFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	value := d.data[d.head]
	var zero T
	d.data[d.head] = zero
	d.head = (d.head + 1) % len(d.data)
	d.size--

	return value, true
}

// PopBack removes and returns the back value, or (zero, false) when empty. It
// clears the removed slot so unused storage does not retain references.
func (d *Deque[T]) PopBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	d.tail = (d.tail - 1 + len(d.data)) % len(d.data)
	value := d.data[d.tail]
	var zero T
	d.data[d.tail] = zero
	d.size--

	return value, true
}

// PeekFront returns the front value without removing it, or (zero, false) when
// empty.
func (d *Deque[T]) PeekFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.data[d.head], true
}

// PeekBack returns the back value without removing it, or (zero, false) when
// empty.
func (d *Deque[T]) PeekBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.data[(d.tail-1+len(d.data))%len(d.data)], true
}

// Size returns the number of elements in the deque.
func (d *Deque[T]) Size() int {
	return d.size
}

// IterFront returns an iterator that yields elements from front to back without modifying the deque.
func (d *Deque[T]) IterFront() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.size; i++ {
			idx := (d.head + i) % len(d.data)
			if !yield(d.data[idx]) {
				return
			}
		}
	}
}

// IterBack returns an iterator that yields elements from back to front without modifying the deque.
func (d *Deque[T]) IterBack() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < d.size; i++ {
			idx := (d.tail - 1 - i + len(d.data)) % len(d.data)
			if !yield(d.data[idx]) {
				return
			}
		}
	}
}
