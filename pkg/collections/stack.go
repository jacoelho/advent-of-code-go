package collections

type Stack[T any] []T

func NewStack[T any](el ...T) *Stack[T] {
	s := new(Stack[T])
	for _, t := range el {
		s.Push(t)
	}
	return s
}

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

func (s *Stack[T]) Peek() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	return (*s)[len(*s)-1], true
}

func (s *Stack[T]) Len() int {
	return len(*s)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}
