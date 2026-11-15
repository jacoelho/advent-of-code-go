package aoc2023

import (
	"io"
	"strings"
	"testing"
)

func requireRejected(t *testing.T, inputs []string, solvers ...func(io.Reader) (string, error)) {
	t.Helper()
	for _, input := range inputs {
		for part, solve := range solvers {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("part %d accepted %q with answer %q", part+1, input, got)
			}
		}
	}
}
