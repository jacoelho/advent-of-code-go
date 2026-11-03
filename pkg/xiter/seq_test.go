package xiter_test

import (
	"iter"
	"math"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

func TestChainIsLazyAndStopsAtRequiredPrefix(t *testing.T) {
	started, cleaned := false, false
	visited, mapped := 0, 0
	source := func(yield func(int) bool) {
		started = true
		defer func() { cleaned = true }()
		for _, v := range []int{1, 2, 3, 4, 5} {
			visited++
			if !yield(v) {
				return
			}
		}
	}
	seq := xiter.From(source).
		Map(func(v int) string { mapped++; return strconv.Itoa(v * 2) }).
		Filter(func(v string) bool { return v != "4" }).
		Take(2)
	if started || mapped != 0 {
		t.Fatal("building a chain consumed its source")
	}
	if got := seq.Collect(); !slices.Equal(got, []string{"2", "6"}) {
		t.Fatalf("Collect() = %v, want [2 6]", got)
	}
	if visited != 3 || mapped != 3 || !cleaned {
		t.Fatalf("visited=%d, mapped=%d, cleaned=%v; want 3, 3, true", visited, mapped, cleaned)
	}
}

func TestTakeAndSkipBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name string
		n    int
		take []int
		skip []int
	}{
		{"negative", -1, nil, []int{1, 2, 3}},
		{"zero", 0, nil, []int{1, 2, 3}},
		{"middle", 2, []int{1, 2}, []int{3}},
		{"beyond", 5, []int{1, 2, 3}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			seq := xiter.FromSlice([]int{1, 2, 3})
			if got := seq.Take(tc.n).Collect(); !slices.Equal(got, tc.take) {
				t.Errorf("Take(%d) = %v, want %v", tc.n, got, tc.take)
			}
			if got := seq.Skip(tc.n).Collect(); !slices.Equal(got, tc.skip) {
				t.Errorf("Skip(%d) = %v, want %v", tc.n, got, tc.skip)
			}
		})
	}
	for _, n := range []int{-1, 0} {
		seq := xiter.From(func(func(int) bool) { t.Fatal("Take invoked the source") })
		if got := seq.Take(n).Collect(); got != nil {
			t.Errorf("Take(%d) = %v, want nil", n, got)
		}
	}
}

func TestReplayAndConsumptiveSources(t *testing.T) {
	seq := xiter.From(xiter.Apply(1, func(v int) int { return v * 2 })).Skip(1).Take(3)
	for range 2 {
		if got := seq.Collect(); !slices.Equal(got, []int{2, 4, 8}) {
			t.Fatalf("replayed chain = %v, want [2 4 8]", got)
		}
	}

	values, cursor := []int{1, 2, 3, 4}, 0
	singleUse := xiter.From(func(yield func(int) bool) {
		for cursor < len(values) {
			v := values[cursor]
			cursor++
			if !yield(v) {
				return
			}
		}
	}).Take(2)
	for _, want := range [][]int{{1, 2}, {3, 4}, nil} {
		if got := singleUse.Collect(); !slices.Equal(got, want) {
			t.Fatalf("consumptive chain = %v, want %v", got, want)
		}
	}
}

func TestEmptyTerminals(t *testing.T) {
	for _, seq := range []xiter.Seq[int]{{}, xiter.From[int](nil), xiter.FromSlice([]int{})} {
		if seq.Collect() != nil || seq.Count() != 0 {
			t.Fatal("empty collection/count")
		}
		if got := seq.Reduce("seed", func(string, int) string { t.Fatal("empty reduce called reducer"); return "" }); got != "seed" {
			t.Fatalf("empty Reduce = %q", got)
		}
		predicate := func(int) bool { t.Fatal("empty terminal called predicate"); return false }
		if seq.Any(predicate) || !seq.Every(predicate) {
			t.Fatal("empty Any/Every identities")
		}
		if v, ok := seq.First(); v != 0 || ok {
			t.Fatalf("empty First = %v, %v", v, ok)
		}
		if got := seq.SumBy(func(v int) int { return v }); got != 0 {
			t.Fatalf("empty SumBy = %d", got)
		}
	}
}

func TestTerminalResultsAndShortCircuit(t *testing.T) {
	seq := xiter.FromSlice([]int{2, 3, 4})
	if got := seq.Reduce("values:", func(acc string, v int) string { return acc + strconv.Itoa(v) }); got != "values:234" {
		t.Fatalf("Reduce = %q", got)
	}
	if got := seq.SumBy(func(v int) int64 { return int64(v * v) }); got != 29 {
		t.Fatalf("SumBy = %d, want 29", got)
	}
	if got := seq.Filter(func(v int) bool { return v%2 == 0 }).Count(); got != 2 {
		t.Fatalf("Count = %d, want 2", got)
	}
	for _, tc := range []struct {
		name    string
		run     func(xiter.Seq[int]) bool
		visited int
	}{
		{"Any", func(s xiter.Seq[int]) bool { return s.Any(func(v int) bool { return v == 3 }) }, 2},
		{"Every", func(s xiter.Seq[int]) bool { return !s.Every(func(v int) bool { return v < 3 }) }, 2},
		{"First", func(s xiter.Seq[int]) bool { v, ok := s.First(); return ok && v == 2 }, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			visited := 0
			if !tc.run(seq.Map(func(v int) int { visited++; return v })) || visited != tc.visited {
				t.Fatalf("wrong terminal result or consumed %d values; want %d", visited, tc.visited)
			}
		})
	}
}

func TestCollectOwnsSliceAndFromSliceBorrows(t *testing.T) {
	input := []int{1, 2, 3}
	seq := xiter.FromSlice(input)
	collected := seq.Collect()
	collected[0] = 9
	if input[0] != 1 {
		t.Fatal("Collect aliased the input slice")
	}
	input[1] = 8
	if got := seq.Collect(); !slices.Equal(got, []int{1, 8, 3}) {
		t.Fatalf("FromSlice did not observe input mutation: %v", got)
	}
}

func singleTraversal[T any](t *testing.T, values []T) iter.Seq[T] {
	t.Helper()
	started := false
	return func(yield func(T) bool) {
		if started {
			t.Fatal("terminal started its source twice")
		}
		started = true
		for _, value := range values {
			if !yield(value) {
				return
			}
		}
	}
}

func TestExtremaConsumeOnce(t *testing.T) {
	for _, tc := range []struct {
		values []int
		want   int
		ok     bool
	}{{nil, 0, false}, {[]int{-5}, -5, true}, {[]int{-5, -2, -9}, -2, true}} {
		if got, ok := xiter.Max(singleTraversal(t, tc.values)); got != tc.want || ok != tc.ok {
			t.Errorf("Max(%v) = %d, %v; want %d, %v", tc.values, got, ok, tc.want, tc.ok)
		}
		if got, ok := xiter.MaxBy(func(a, b int) bool { return a < b }, singleTraversal(t, tc.values)); got != tc.want || ok != tc.ok {
			t.Errorf("MaxBy(%v) = %d, %v; want %d, %v", tc.values, got, ok, tc.want, tc.ok)
		}
	}
	if got, ok := xiter.MaxBy(func(a, b string) bool { return len(a) < len(b) }, singleTraversal(t, []string{"x", "abc", "def"})); got != "abc" || !ok {
		t.Fatalf("MaxBy tie = %q, %v; want first greatest abc", got, ok)
	}
	if got, ok := xiter.Max(slices.Values([]float64{1, math.NaN(), 2})); !ok || !math.IsNaN(got) {
		t.Fatalf("Max should propagate NaN: %v, %v", got, ok)
	}
}

func TestZipKeepsUnmatchedValues(t *testing.T) {
	for _, tc := range []struct {
		x, y []int
		want []xiter.Zipped[int, int]
	}{
		{[]int{2}, []int{5, 7}, []xiter.Zipped[int, int]{{V1: 2, Ok1: true, V2: 5, Ok2: true}, {V2: 7, Ok2: true}}},
		{[]int{2, 3}, []int{5}, []xiter.Zipped[int, int]{{V1: 2, Ok1: true, V2: 5, Ok2: true}, {V1: 3, Ok1: true}}},
	} {
		if got := slices.Collect(xiter.Zip(slices.Values(tc.x), slices.Values(tc.y))); !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("Zip = %v, want %v", got, tc.want)
		}
	}
}

func TestDotProductStopsAtShorterSource(t *testing.T) {
	for _, tc := range []struct {
		x, y         []int
		want         int
		readX, readY int
	}{
		{[]int{2, 3}, []int{5, 7, 11, 13}, 31, 2, 2},
		{[]int{2, 3, 4, 5}, []int{5, 7}, 31, 3, 2},
		{nil, []int{5, 7}, 0, 0, 0},
		{[]int{2, 3}, nil, 0, 1, 0},
	} {
		readsX, readsY := 0, 0
		startedX, startedY, closedX, closedY := false, false, false, false
		source := func(values []int, reads *int, started, closed *bool) iter.Seq[int] {
			return func(yield func(int) bool) {
				*started = true
				defer func() { *closed = true }()
				for _, value := range values {
					*reads++
					if !yield(value) {
						return
					}
				}
			}
		}
		got := xiter.DotProduct(source(tc.x, &readsX, &startedX, &closedX), source(tc.y, &readsY, &startedY, &closedY))
		if got != tc.want || readsX != tc.readX || readsY != tc.readY {
			t.Fatalf("DotProduct(%v, %v) = %d, reads %d/%d; want %d, reads %d/%d", tc.x, tc.y, got, readsX, readsY, tc.want, tc.readX, tc.readY)
		}
		if startedX != closedX || startedY != closedY {
			t.Fatal("source cleanup did not run")
		}
	}
}

func TestProductIdentityAndFullConsumption(t *testing.T) {
	type factor int64
	var product factor = xiter.Product(slices.Values([]factor{-2, 3, 4}))
	if product != -24 {
		t.Fatalf("Product = %d, want -24", product)
	}
	if got := xiter.Product(slices.Values([]factor{})); got != 1 {
		t.Fatalf("empty Product = %d, want 1", got)
	}
	visited := 0
	source := xiter.FromSlice([]int{0, 2, 3}).Map(func(value int) int {
		visited++
		return value
	})
	if got := xiter.Product(source.All()); got != 0 || visited != 3 {
		t.Fatalf("Product = %d, visited %d; want 0 and full consumption", got, visited)
	}
}

func TestNthBoundariesAndConsumption(t *testing.T) {
	for _, tc := range []struct {
		index, want, visited int
		ok                   bool
	}{
		{-1, 0, 0, false},
		{0, 4, 1, true},
		{1, 8, 2, true},
		{3, 0, 3, false},
	} {
		t.Run(strconv.Itoa(tc.index), func(t *testing.T) {
			visited := 0
			started, cleaned := false, false
			source := func(yield func(int) bool) {
				started = true
				defer func() { cleaned = true }()
				for _, value := range []int{4, 8, 12} {
					visited++
					if !yield(value) {
						return
					}
				}
			}
			got, ok := xiter.Nth(source, tc.index)
			if got != tc.want || ok != tc.ok || visited != tc.visited {
				t.Fatalf("Nth = %d, %v; visited %d; want %d, %v; visited %d", got, ok, visited, tc.want, tc.ok, tc.visited)
			}
			if started != (tc.index >= 0) || cleaned != started {
				t.Fatalf("started=%v, cleaned=%v for index %d", started, cleaned, tc.index)
			}
		})
	}
}
