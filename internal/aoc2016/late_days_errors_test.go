package aoc2016

import (
	"errors"
	"io"
	"testing"
)

type lateReadFailure2016 struct{ err error }

func (r lateReadFailure2016) Read([]byte) (int, error) { return 0, r.err }

func Test2016LateDayReadErrors(t *testing.T) {
	failure := errors.New("input failed")
	for _, solve := range []func(io.Reader) (string, error){day17p01, day17p02, day18p01, day18p02, day19p01, day19p02, day20p01, day20p02, day21p01, day21p02, day22p01, day22p02, day23p01, day23p02, day24p01, day24p02, day25p01} {
		if _, err := solve(lateReadFailure2016{failure}); !errors.Is(err, failure) {
			t.Errorf("read error%v lost input cause", err)
		}
	}
}
