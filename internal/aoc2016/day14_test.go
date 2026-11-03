package aoc2016

import (
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day14p01(t *testing.T) {
	aoc.AOCTest(t, day14p01, []aoc.TestInput{
		{Input: strings.NewReader("abc"), Want: "22728"},
		{Input: aoc.FileInput(t, 2016, 14), Want: "15035"},
	})
}

func Test_day14p02(t *testing.T) {
	if testing.Short() {
		t.Skip("full stretched-MD5 fixture")
	}
	aoc.AOCTest(t, day14p02, []aoc.TestInput{{Input: aoc.FileInput(t, 2016, 14), Want: "19968"}})
}

func TestPadHashFirstTripleAndEveryQuintuple(t *testing.T) {
	features := inspectPadHash([]byte("aaabbbbbcccccdddddeee"))
	if features.triple != 10 || features.quintuples != 1<<11|1<<12|1<<13 {
		t.Errorf("features=%+v; want first triple a and quintuple digits b,c,d", features)
	}
}

func TestPadLookaheadBoundaries(t *testing.T) {
	for _, tt := range []struct {
		name     string
		features map[int]padHashFeatures
		want     int
	}{
		{"include 1000", map[int]padHashFeatures{0: {triple: 1}, 1000: {triple: -1, quintuples: 1 << 1}}, 0},
		{"exclude 1001", map[int]padHashFeatures{0: {triple: 1}, 1: {triple: 1}, 1001: {triple: -1, quintuples: 1 << 1}}, 1},
		{"exclude self", map[int]padHashFeatures{0: {triple: 1, quintuples: 1 << 1}, 1: {triple: 2}, 2: {triple: -1, quintuples: 1 << 2}}, 1},
		{"match later quintuple", map[int]padHashFeatures{0: {triple: 12}, 1: inspectPadHash([]byte("bbbbbaaaaaccccc"))}, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findPadKey(func(index int) padHashFeatures {
				if features, ok := tt.features[index]; ok {
					return features
				}
				return padHashFeatures{triple: -1}
			}, 1)
			if err != nil || got != tt.want {
				t.Errorf("got %d, %v; want %d", got, err, tt.want)
			}
		})
	}
}

func TestPadFirstStretchedKey(t *testing.T) {
	got, err := findPadKey(func(index int) padHashFeatures { return padHash([]byte("abc"), index, 2017) }, 1)
	if err != nil || got != 10 {
		t.Errorf("got %d, %v; want 10", got, err)
	}
}

func TestPadErrors(t *testing.T) {
	for _, input := range []string{"", "a b"} {
		if _, err := day14p01(strings.NewReader(input)); err == nil {
			t.Errorf("expected error for %q", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := day14p01(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
	if _, err := day14p02(iotest.ErrReader(failure)); !errors.Is(err, failure) {
		t.Errorf("got %v, want read error", err)
	}
}
