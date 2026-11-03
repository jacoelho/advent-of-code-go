package scanner

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func TestNumberOverflowStopsScanner(t *testing.T) {
	s := NewScanner(strings.NewReader("+12\n-128\n128\n7\n"), convert.ScanNumber[int8])
	got := slices.Collect(s.Values())
	if !slices.Equal(got, []int8{12, -128}) || !errors.Is(s.Err(), strconv.ErrRange) {
		t.Fatalf("values = %v, error = %v; want [12 -128], ErrRange", got, s.Err())
	}
}
