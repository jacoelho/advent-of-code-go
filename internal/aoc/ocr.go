package aoc

import (
	"fmt"
	"strings"
)

// Each font lists its letters side by side, separated by a dark column, in the
// order given by its letters string.
var (
	// The six-row letters are those listed by
	// https://github.com/bsoyka/advent-of-code-ocr; Y is five columns wide, as
	// it appears in 2019 day 8.
	sixRowFont = newFont("ABCEFGHIJKLOPRSUYZ", `
.##..###...##..####.####..##..#..#.###...##.#..#.#.....##..###..###...###.#..#.#...#.####
#..#.#..#.#..#.#....#....#..#.#..#..#.....#.#.#..#....#..#.#..#.#..#.#....#..#.#...#....#
#..#.###..#....###..###..#....####..#.....#.##...#....#..#.#..#.#..#.#....#..#..#.#....#.
####.#..#.#....#....#....#.##.#..#..#.....#.#.#..#....#..#.###..###...##..#..#...#....#..
#..#.#..#.#..#.#....#....#..#.#..#..#..#..#.#.#..#....#..#.#....#.#.....#.#..#...#...#...
#..#.###...##..####.#.....###.#..#.###..##..#..#.####..##..#....#..#.###...##....#...####`)

	// The eight-row letters are those of the 2018 day 10 example.
	eightRowFont = newFont("HI", `
#...#.###
#...#..#.
#...#..#.
#####..#.
#...#..#.
#...#..#.
#...#..#.
#...#.###`)

	tenRowFont = newFont("ABCEFGHJKLNPRXZ", `
..##...#####...####..######.######..####..#....#....###.#....#.#......#....#.#####..#####..#....#.######
.#..#..#....#.#....#.#......#......#....#.#....#.....#..#...#..#......##...#.#....#.#....#.#....#......#
#....#.#....#.#......#......#......#......#....#.....#..#..#...#......##...#.#....#.#....#..#..#.......#
#....#.#....#.#......#......#......#......#....#.....#..#.#....#......#.#..#.#....#.#....#..#..#......#.
#....#.#####..#......#####..#####..#......######.....#..##.....#......#.#..#.#####..#####....##......#..
######.#....#.#......#......#......#..###.#....#.....#..##.....#......#..#.#.#......#..#.....##.....#...
#....#.#....#.#......#......#......#....#.#....#.....#..#.#....#......#..#.#.#......#...#...#..#...#....
#....#.#....#.#......#......#......#....#.#....#.#...#..#..#...#......#...##.#......#...#...#..#..#.....
#....#.#....#.#....#.#......#......#...##.#....#.#...#..#...#..#......#...##.#......#....#.#....#.#.....
#....#.#####...####..######.#.......###.#.#....#..###...#....#.######.#....#.#......#....#.#....#.######`)
)

// maxGlyphWidth keeps a glyph's pixels and width within one uint64 key.
const maxGlyphWidth = 6

// font maps each letter's glyph, cropped to its lit columns, to the letter.
type font map[uint64]byte

func newFont(letters, drawing string) font {
	rows := strings.Split(strings.TrimPrefix(drawing, "\n"), "\n")
	f := make(font, len(letters))
	runs := litRuns(rows, 0, len(rows[0]))
	if len(runs) != len(letters) {
		panic(fmt.Sprintf("font drawing has %d glyphs for letters %q", len(runs), letters))
	}
	for i, run := range runs {
		f[glyphKey(rows, run[0], run[1])] = letters[i]
	}
	if len(f) != len(letters) {
		panic(fmt.Sprintf("font drawing repeats a glyph among letters %q", letters))
	}
	return f
}

// glyphKey packs the pixels of columns [start, end) row by row, above the
// glyph width.
func glyphKey(rows []string, start, end int) uint64 {
	key := uint64(end - start)
	for _, row := range rows {
		for x := start; x < end; x++ {
			key <<= 1
			if row[x] == '#' {
				key |= 1
			}
		}
	}
	return key
}

// litRuns returns the maximal runs [start, end) of columns in [from, to)
// holding at least one lit pixel.
func litRuns(rows []string, from, to int) [][2]int {
	lit := func(x int) bool {
		for _, row := range rows {
			if row[x] == '#' {
				return true
			}
		}
		return false
	}
	var runs [][2]int
	for x := from; x < to; {
		if !lit(x) {
			x++
			continue
		}
		start := x
		for x < to && lit(x) {
			x++
		}
		runs = append(runs, [2]int{start, x})
	}
	return runs
}

func (f font) letter(rows []string, start, end int) (byte, error) {
	if end-start > maxGlyphWidth {
		return 0, fmt.Errorf("display glyph at column %d is too wide", start)
	}
	letter, ok := f[glyphKey(rows, start, end)]
	if !ok {
		return 0, fmt.Errorf("unknown display glyph at column %d:%s", start, renderGlyph(rows, start, end))
	}
	return letter, nil
}

func renderGlyph(rows []string, start, end int) string {
	var glyph strings.Builder
	glyph.WriteByte('\n')
	for _, row := range rows {
		for x := start; x < end; x++ {
			pixel := byte('.')
			if row[x] == '#' {
				pixel = '#'
			}
			glyph.WriteByte(pixel)
		}
		glyph.WriteByte('\n')
	}
	return glyph.String()
}

// Decode6RowDisplay reads fixed four- or five-column cells in the AoC six-row
// font. '#' is lit; '.' and space are dark. Blank cells decode as spaces.
func Decode6RowDisplay(rows []string, cellWidth int) (string, error) {
	if cellWidth != 4 && cellWidth != 5 {
		return "", fmt.Errorf("display cell width must be 4 or 5, got %d", cellWidth)
	}
	if len(rows) != 6 {
		return "", fmt.Errorf("display must have six rows, got %d", len(rows))
	}
	width, err := validateDisplay(rows)
	if err != nil {
		return "", err
	}
	if width%cellWidth != 0 {
		return "", fmt.Errorf("display width %d is not a multiple of cell width %d", width, cellWidth)
	}
	text := make([]byte, 0, width/cellWidth)
	for start := 0; start < width; start += cellWidth {
		runs := litRuns(rows, start, start+cellWidth)
		switch len(runs) {
		case 0:
			text = append(text, ' ')
			continue
		case 1:
			if letter, err := sixRowFont.letter(rows, runs[0][0], runs[0][1]); err == nil {
				text = append(text, letter)
				continue
			}
		}
		return "", fmt.Errorf("unknown display glyph in cell %d:%s", start/cellWidth, renderGlyph(rows, start, start+cellWidth))
	}
	return string(text), nil
}

// DecodeSeparatedDisplay reads the six-, eight- or ten-row AoC font.
// Glyphs have no internal blank columns and are separated by dark columns.
// Outer dark columns are ignored.
func DecodeSeparatedDisplay(rows []string) (string, error) {
	var f font
	switch len(rows) {
	case 6:
		f = sixRowFont
	case 8:
		f = eightRowFont
	case 10:
		f = tenRowFont
	default:
		return "", fmt.Errorf("separated display must have six, eight or ten rows")
	}
	width, err := validateDisplay(rows)
	if err != nil {
		return "", err
	}
	var text []byte
	for _, run := range litRuns(rows, 0, width) {
		letter, err := f.letter(rows, run[0], run[1])
		if err != nil {
			return "", err
		}
		text = append(text, letter)
	}
	return string(text), nil
}

func validateDisplay(rows []string) (int, error) {
	width := len(rows[0])
	for y, row := range rows {
		if len(row) != width {
			return 0, fmt.Errorf("display row %d has width %d, want %d", y, len(row), width)
		}
		for _, pixel := range row {
			if pixel != '#' && pixel != '.' && pixel != ' ' {
				return 0, fmt.Errorf("invalid display pixel %q in row %d", pixel, y)
			}
		}
	}
	return width, nil
}
