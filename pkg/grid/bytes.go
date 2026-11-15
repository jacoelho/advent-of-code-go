package grid

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"iter"
)

// Bytes is a non-empty rectangular grid of byte cells stored row by row. The
// zero value is an empty grid with no cells.
type Bytes struct {
	width int
	cells []byte
}

// ReadBytes reads one grid row per line of r. Rows must be non-empty and of
// equal width, every cell must be one of the bytes in alphabet, and there must
// be at least one row.
func ReadBytes(r io.Reader, alphabet string) (Bytes, error) {
	var allowed [256]bool
	for i := range len(alphabet) {
		allowed[alphabet[i]] = true
	}
	var g Bytes
	s := bufio.NewScanner(r)
	for y := 0; s.Scan(); y++ {
		line := s.Bytes()
		if len(line) == 0 || y > 0 && len(line) != g.width {
			return Bytes{}, fmt.Errorf("grid row %d has width %d, want a non-empty rectangular grid", y+1, len(line))
		}
		for x, c := range line {
			if !allowed[c] {
				return Bytes{}, fmt.Errorf("grid cell %d,%d is %q, want one of %q", x, y, c, alphabet)
			}
		}
		g.width = len(line)
		g.cells = append(g.cells, line...)
	}
	if err := s.Err(); err != nil {
		return Bytes{}, err
	}
	if len(g.cells) == 0 {
		return Bytes{}, fmt.Errorf("empty grid")
	}
	return g, nil
}

// Width returns the number of columns.
func (g Bytes) Width() int { return g.width }

// Height returns the number of rows.
func (g Bytes) Height() int {
	if g.width == 0 {
		return 0
	}
	return len(g.cells) / g.width
}

// Cells returns every cell in row-major order. The slice shares storage with
// g, so writes to it change g.
func (g Bytes) Cells() []byte { return g.cells }

// Row returns row y. The slice shares storage with g.
func (g Bytes) Row(y int) []byte { return g.cells[y*g.width : (y+1)*g.width : (y+1)*g.width] }

// Rows returns every row. The rows share storage with g.
func (g Bytes) Rows() [][]byte {
	rows := make([][]byte, g.Height())
	for y := range rows {
		rows[y] = g.Row(y)
	}
	return rows
}

// Contains reports whether p lies inside g.
func (g Bytes) Contains(p Position2D[int]) bool {
	return 0 <= p.X && p.X < g.width && 0 <= p.Y && p.Y*g.width < len(g.cells)
}

// Index returns the offset of p in Cells. p must lie inside g.
func (g Bytes) Index(p Position2D[int]) int { return p.Y*g.width + p.X }

// Position returns the position of the cell at offset i in Cells.
func (g Bytes) Position(i int) Position2D[int] {
	return Position2D[int]{X: i % g.width, Y: i / g.width}
}

// At returns the cell at p. It panics if p lies outside g.
func (g Bytes) At(p Position2D[int]) byte {
	if !g.Contains(p) {
		panic(fmt.Sprintf("grid: position %d,%d outside %dx%d grid", p.X, p.Y, g.width, g.Height()))
	}
	return g.cells[g.Index(p)]
}

// Find returns the position of the only cell holding c, or an error if no
// cell or more than one cell holds it.
func (g Bytes) Find(c byte) (Position2D[int], error) {
	first := bytes.IndexByte(g.cells, c)
	if first < 0 {
		return Position2D[int]{}, fmt.Errorf("grid has no %q", c)
	}
	if second := bytes.IndexByte(g.cells[first+1:], c); second >= 0 {
		p, q := g.Position(first), g.Position(first+1+second)
		return Position2D[int]{}, fmt.Errorf("grid has %q at %d,%d and %d,%d", c, p.X, p.Y, q.X, q.Y)
	}
	return g.Position(first), nil
}

// Neighbours4 yields the orthogonal neighbours of p that lie inside g.
func (g Bytes) Neighbours4(p Position2D[int]) iter.Seq[Position2D[int]] {
	return func(yield func(Position2D[int]) bool) {
		for _, offset := range offsets4 {
			if q := p.Add(offset); g.Contains(q) && !yield(q) {
				return
			}
		}
	}
}

// Neighbours8 yields the orthogonal and diagonal neighbours of p that lie
// inside g.
func (g Bytes) Neighbours8(p Position2D[int]) iter.Seq[Position2D[int]] {
	return func(yield func(Position2D[int]) bool) {
		for _, offset := range offsets8 {
			if q := p.Add(offset); g.Contains(q) && !yield(q) {
				return
			}
		}
	}
}

var (
	offsets4 = [...]Position2D[int]{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	offsets8 = [...]Position2D[int]{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}}
)
