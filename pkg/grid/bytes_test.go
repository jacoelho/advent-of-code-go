package grid

import (
	"errors"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestReadBytesLaysOutRowsInOrder(t *testing.T) {
	g, err := ReadBytes(strings.NewReader("#.S\n..#\n"), "#.S")
	if err != nil {
		t.Fatal(err)
	}
	if g.Width() != 3 || g.Height() != 2 || string(g.Cells()) != "#.S..#" {
		t.Fatalf("got %dx%d %q", g.Width(), g.Height(), g.Cells())
	}
	if string(g.Row(1)) != "..#" || g.At(NewPosition2D(2, 0)) != 'S' {
		t.Fatalf("row 1 %q, (2,0) %q", g.Row(1), g.At(NewPosition2D(2, 0)))
	}
	for i := range g.Cells() {
		if g.Index(g.Position(i)) != i {
			t.Fatalf("index %d does not round-trip through %v", i, g.Position(i))
		}
	}
}

func TestReadBytesRejectsMalformedGrids(t *testing.T) {
	for _, input := range []string{"", "#.\n#\n", "#.\n\n#.\n", "#x\n", "\n"} {
		if _, err := ReadBytes(strings.NewReader(input), "#."); err == nil {
			t.Fatalf("input %q accepted", input)
		}
	}
	failure := errors.New("read failed")
	if _, err := ReadBytes(io.MultiReader(strings.NewReader("#.\n"), iotest.ErrReader(failure)), "#."); !errors.Is(err, failure) {
		t.Fatalf("got %v; want reader failure", err)
	}
}

func TestBytesRowWritesDoNotSpillIntoNextRow(t *testing.T) {
	g, err := ReadBytes(strings.NewReader("ab\ncd\n"), "abcdx")
	if err != nil {
		t.Fatal(err)
	}
	row := append(g.Row(0), 'x')
	g.Rows()[1][0] = 'x'
	if string(g.Cells()) != "abxd" || string(row) != "abx" {
		t.Fatalf("cells %q, appended row %q", g.Cells(), row)
	}
}

func TestBytesFindRequiresExactlyOneCell(t *testing.T) {
	g, err := ReadBytes(strings.NewReader("..E\nSS.\n"), ".SE")
	if err != nil {
		t.Fatal(err)
	}
	if p, err := g.Find('E'); err != nil || p != NewPosition2D(2, 0) {
		t.Fatalf("E: got %v, %v", p, err)
	}
	if _, err := g.Find('S'); err == nil {
		t.Fatal("two S cells accepted")
	}
	if _, err := g.Find('#'); err == nil {
		t.Fatal("missing cell accepted")
	}
}

func TestBytesNeighboursStayInsideGrid(t *testing.T) {
	g, err := ReadBytes(strings.NewReader("...\n...\n"), ".")
	if err != nil {
		t.Fatal(err)
	}
	corner := slices.Collect(g.Neighbours4(NewPosition2D(0, 0)))
	if len(corner) != 2 {
		t.Fatalf("corner 4-neighbours %v", corner)
	}
	edge := slices.Collect(g.Neighbours8(NewPosition2D(1, 1)))
	if len(edge) != 5 {
		t.Fatalf("bottom-edge 8-neighbours %v", edge)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("At outside the grid did not panic")
		}
	}()
	g.At(NewPosition2D(3, 0))
}
