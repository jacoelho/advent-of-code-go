package collections

// Deque represents a double-ended queue implemented as a circular buffer.
type Deque[T any] struct {
	data       []T
	head, tail int
	size, cap  int
}

// NewDeque creates a new deque with a specified initial capacity.
func NewDeque[T any](capacity int) *Deque[T] {
	if capacity < 1 {
		capacity = 1
	}
	return &Deque[T]{
		data: make([]T, capacity),
		cap:  capacity,
	}
}

// resize adjusts the capacity of the deque.
func (d *Deque[T]) resize(newCap int) {
	newData := make([]T, newCap)
	for i := 0; i < d.size; i++ {
		newData[i] = d.data[(d.head+i)%d.cap]
	}
	d.data = newData
	d.head = 0
	d.tail = d.size
	d.cap = newCap
}

// PushFront adds an element to the front of the deque.
func (d *Deque[T]) PushFront(value T) {
	if d.size == d.cap {
		d.resize(d.cap * 2)
	}
	d.head = (d.head - 1 + d.cap) % d.cap
	d.data[d.head] = value
	d.size++
}

// PushBack adds an element to the back of the deque.
func (d *Deque[T]) PushBack(value T) {
	if d.size == d.cap {
		d.resize(d.cap * 2)
	}
	d.data[d.tail] = value
	d.tail = (d.tail + 1) % d.cap
	d.size++
}

// PopFront removes and returns the front element of the deque.
func (d *Deque[T]) PopFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	value := d.data[d.head]
	d.head = (d.head + 1) % d.cap
	d.size--

	// Shrink capacity if necessary
	if d.size > 0 && d.size <= d.cap/4 {
		d.resize(d.cap / 2)
	}
	return value, true
}

// PopBack removes and returns the back element of the deque.
func (d *Deque[T]) PopBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	d.tail = (d.tail - 1 + d.cap) % d.cap
	value := d.data[d.tail]
	d.size--

	// Shrink capacity if necessary
	if d.size > 0 && d.size <= d.cap/4 {
		d.resize(d.cap / 2)
	}
	return value, true
}

// PeekFront returns the front element without removing it.
func (d *Deque[T]) PeekFront() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.data[d.head], true
}

// PeekBack returns the back element without removing it.
func (d *Deque[T]) PeekBack() (T, bool) {
	if d.size == 0 {
		var zero T
		return zero, false
	}
	return d.data[(d.tail-1+d.cap)%d.cap], true
}

// Size returns the number of elements in the deque.
func (d *Deque[T]) Size() int {
	return d.size
}
