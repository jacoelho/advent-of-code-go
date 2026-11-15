package scanner

import (
	"io"
	"slices"

	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

// ReadBlocks reads the lines of r and groups them into blank-line-separated
// blocks as xiter.Blocks does.
func ReadBlocks(r io.Reader) ([][]string, error) {
	s := NewScanner(r, func(b []byte) (string, error) { return string(b), nil })
	blocks := slices.Collect(xiter.Blocks(s.Values()))
	if err := s.Err(); err != nil {
		return nil, err
	}
	return blocks, nil
}
