package scanner

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadBlocksGroupsLinesAndIgnoresOneTrailingBlankLine(t *testing.T) {
	got, err := ReadBlocks(strings.NewReader("a\nb\n\nc\n\n"))
	want := [][]string{{"a", "b"}, {"c"}}
	if err != nil || !slices.EqualFunc(got, want, slices.Equal) {
		t.Fatalf("got %q, %v; want %q", got, err, want)
	}
}

func TestReadBlocksReportsReadFailure(t *testing.T) {
	failure := errors.New("read failed")
	if _, err := ReadBlocks(io.MultiReader(strings.NewReader("a\n\nb\n"), iotest.ErrReader(failure))); !errors.Is(err, failure) {
		t.Fatalf("got %v; want read failure", err)
	}
}
