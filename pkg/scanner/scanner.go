package scanner

import (
	"bufio"
	"bytes"
	"io"
	"iter"
)

type FromBytes[T any] func([]byte) (T, error)

type Scanner[T any] struct {
	scanner *bufio.Scanner
	convert func([]byte) (T, error)
	err     error
}

func NewScanner[T any](r io.Reader, converter FromBytes[T]) *Scanner[T] {
	scanner := bufio.NewScanner(r)
	return &Scanner[T]{
		scanner: scanner,
		convert: converter,
	}
}

func NewScannerWithSplit[T any](r io.Reader, splitFunc bufio.SplitFunc, convert FromBytes[T]) *Scanner[T] {
	s := NewScanner(bufio.NewReader(r), convert)
	s.scanner.Split(splitFunc)
	return s
}

func (s *Scanner[T]) Values() iter.Seq[T] {
	return func(yield func(v T) bool) {
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

func (s *Scanner[T]) Err() error {
	if s.err != nil {
		return s.err
	}
	return s.scanner.Err()
}

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
