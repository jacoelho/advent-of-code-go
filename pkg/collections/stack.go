package collections

// Stack is a last-in, first-out stack. The zero value is an empty stack.
type Stack[T any] []T

// NewStack returns a stack holding el, with the last element on top.
func NewStack[T any](el ...T) *Stack[T] {
	s := new(Stack[T])
	for _, t := range el {
		s.Push(t)
	}
	return s
}

// Push places value on top of the stack.
func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}

// Pop removes and returns the top value, or (zero, false) when empty. It clears
// the removed slot so unused storage does not retain references.
func (s *Stack[T]) Pop() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	el := (*s)[len(*s)-1]
	var zero T
	(*s)[len(*s)-1] = zero
	*s = (*s)[:len(*s)-1]
	return el, true
}

// Peek returns the top value without removing it, or (zero, false) when
// empty.
func (s *Stack[T]) Peek() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	return (*s)[len(*s)-1], true
}

// Len returns the number of values in the stack.
func (s *Stack[T]) Len() int {
	return len(*s)
}

// IsEmpty reports whether the stack has no values.
func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
