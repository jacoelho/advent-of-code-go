// Package scanner reads typed values from tokenized input.
package scanner

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

// FromBytes converts one token into a value.
type FromBytes[T any] func([]byte) (T, error)

// Scanner converts the tokens of a reader into values of type T.
type Scanner[T any] struct {
	scanner *bufio.Scanner
	convert func([]byte) (T, error)
	err     error
}

// NewScanner returns a Scanner that splits r into lines and converts each line
// with converter. Lines longer than bufio.MaxScanTokenSize fail.
func NewScanner[T any](r io.Reader, converter FromBytes[T]) *Scanner[T] {
	scanner := bufio.NewScanner(r)
	return &Scanner[T]{
		scanner: scanner,
		convert: converter,
	}
}

// NewScannerWithSplit returns a Scanner that splits r with splitFunc and
// converts each token with convert.
func NewScannerWithSplit[T any](r io.Reader, splitFunc bufio.SplitFunc, convert FromBytes[T]) *Scanner[T] {
	s := NewScanner(r, convert)
	s.scanner.Split(splitFunc)
	return s
}

// Values returns a single-use traversal. Stopping early permits a later traversal
// to resume; conversion failures are terminal. The converter's input bytes are
// borrowed until the next scan and must be copied if retained in its result.
// Err reports only failures encountered while consuming input, not unread input.
func (s *Scanner[T]) Values() iter.Seq[T] {
	return func(yield func(v T) bool) {
		if s.err != nil {
			return
		}
		for s.scanner.Scan() {
			var v T
			v, s.err = s.convert(s.scanner.Bytes())
			if s.err != nil {
				return
			}
			if !yield(v) {
				return
			}
		}
	}
}

// Err returns the first conversion error, or else the error from reading input.
func (s *Scanner[T]) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}

// SplitBySeparator returns a bufio.SplitFunc that yields the tokens between
// occurrences of separator. Adjacent separators yield an empty token; a trailing
// separator does not. separator must be non-empty.
func SplitBySeparator(separator []byte) func(data []byte, atEOF bool) (int, []byte, error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if i := bytes.Index(data, separator); i != -1 {
			return i + len(separator), data[:i], nil
		}
		if atEOF && len(data) > 0 {
			return len(data), data, nil
		}
		return 0, nil, nil
	}
}
