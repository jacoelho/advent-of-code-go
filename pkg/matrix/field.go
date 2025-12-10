package matrix

// Field is a constraint requiring field operations.
// When using Matrix[T], T must implement Field[T] (e.g., rat implements Field[rat]).
type Field[T any] interface {
	Add(other T) T
	Sub(other T) T
	Mul(other T) T
	Div(other T) T
	Neg() T
	IsZero() bool
	Zero() T // returns zero element of the field
	One() T  // returns multiplicative identity
}
