package aoc2015

import (
	"errors"
	"io"
	"strings"
	"testing"
)

type lateDaysReadFailure struct{ err error }

func (r lateDaysReadFailure) Read([]byte) (int, error) { return 0, r.err }

func TestLateDaysReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){
		day19p01, day19p02,
		day20p01, day20p02, day21p01, day21p02, day22p01, day22p02,
		day23p01, day23p02, day24p01, day24p02, day25p01,
	} {
		if _, err := solve(lateDaysReadFailure{failure}); !errors.Is(err, failure) {
			t.Errorf("read error = %v; want wrapped input failure", err)
		}
	}
}

func TestInvalidPresentTargets(t *testing.T) {
	for _, input := range []string{"", "0", "-1", "-", "abc", "9223372036854775808"} {
		if _, err := day20p01(strings.NewReader(input)); err == nil {
			t.Errorf("accepted target %q", input)
		}
	}
}
