package aoc

import (
	"fmt"
	"strings"
)

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
		glyph := displayGlyph(rows, start, start+cellWidth, 5, 0)
		letter, ok := sixRowLetter(glyph)
		if !ok {
			return "", fmt.Errorf("unknown display glyph in cell %d:%s", start/cellWidth, displayGlyph(rows, start, start+cellWidth, cellWidth, 0))
		}
		text = append(text, letter)
	}
	return string(text), nil
}

func validateDisplay(rows []string) (int, error) {
	if len(rows) == 0 {
		return 0, fmt.Errorf("empty display")
	}

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

// displayGlyph normalizes dark pixels and pads each row to the font's width.
// Surrounding newlines align the rows in the lookup tables' raw string literals.
func displayGlyph(rows []string, start, end, width, leftPadding int) string {
	var glyph strings.Builder
	glyph.Grow(len(rows)*(width+1) + 1)
	glyph.WriteByte('\n')
	for _, row := range rows {
		for x := range width {
			pixel := byte('.')
			column := start + x - leftPadding
			if column >= start && column < end && row[column] == '#' {
				pixel = '#'
			}
			glyph.WriteByte(pixel)
		}
		glyph.WriteByte('\n')
	}

	return glyph.String()
}

// DecodeSeparatedDisplay reads the six-, eight- or ten-row AoC font.
// Glyphs have no internal blank columns and are separated by dark columns.
// Outer dark columns are ignored.
func DecodeSeparatedDisplay(rows []string) (string, error) {
	if len(rows) != 6 && len(rows) != 8 && len(rows) != 10 {
		return "", fmt.Errorf("separated display must have six, eight or ten rows")
	}
	width, err := validateDisplay(rows)
	if err != nil {
		return "", err
	}

	glyphWidth := 6
	if len(rows) != 10 {
		glyphWidth = 5
	}
	columnLit := func(x int) bool {
		for _, row := range rows {
			if row[x] == '#' {
				return true
			}
		}
		return false
	}

	var text []byte
	for x := 0; x < width; {
		if !columnLit(x) {
			x++
			continue
		}
		start := x
		for x < width && columnLit(x) {
			x++
		}
		if x-start > glyphWidth {
			return "", fmt.Errorf("display glyph at column %d is too wide", start)
		}
		glyph := displayGlyph(rows, start, x, glyphWidth, 0)
		var letter byte
		var ok bool
		if len(rows) == 6 {
			// Cropping can remove leading blank columns, as in I and T.
			letter, ok = sixRowLetter(glyph)
			for padding := 1; !ok && padding <= glyphWidth-(x-start); padding++ {
				letter, ok = sixRowLetter(displayGlyph(rows, start, x, glyphWidth, padding))
			}
		} else {
			letter, ok = tallDisplayLetter(len(rows), glyph)
		}
		if !ok {
			return "", fmt.Errorf("unknown display glyph at column %d:%s", start, displayGlyph(rows, start, x, x-start, 0))
		}
		text = append(text, letter)
	}

	return string(text), nil
}

// Five columns preserve full-width letters such as Y; most letters leave the last column dark.
func sixRowLetter(glyph string) (byte, bool) {
	switch glyph {
	case `
.##..
#..#.
#..#.
####.
#..#.
#..#.
`:
		return 'A', true
	case `
###..
#..#.
###..
#..#.
#..#.
###..
`:
		return 'B', true
	case `
.##..
#..#.
#....
#....
#..#.
.##..
`:
		return 'C', true
	case `
###..
#..#.
#..#.
#..#.
#..#.
###..
`:
		return 'D', true
	case `
####.
#....
###..
#....
#....
####.
`:
		return 'E', true
	case `
####.
#....
###..
#....
#....
#....
`:
		return 'F', true
	case `
.##..
#..#.
#....
#.##.
#..#.
.###.
`:
		return 'G', true
	case `
#..#.
#..#.
####.
#..#.
#..#.
#..#.
`:
		return 'H', true
	case `
.###.
..#..
..#..
..#..
..#..
.###.
`:
		return 'I', true
	case `
..##.
...#.
...#.
...#.
#..#.
.##..
`:
		return 'J', true
	case `
#..#.
#.#..
##...
#.#..
#.#..
#..#.
`:
		return 'K', true
	case `
#....
#....
#....
#....
#....
####.
`:
		return 'L', true
	case `
#...#
##.##
#.#.#
#...#
#...#
#...#
`:
		return 'M', true
	case `
#..#.
#..#.
##.#.
#.##.
#..#.
#..#.
`:
		return 'N', true
	case `
.##..
#..#.
#..#.
#..#.
#..#.
.##..
`:
		return 'O', true
	case `
###..
#..#.
#..#.
###..
#....
#....
`:
		return 'P', true
	case `
.##..
#..#.
#..#.
#.##.
#..#.
.##.#
`:
		return 'Q', true
	case `
###..
#..#.
#..#.
###..
#.#..
#..#.
`:
		return 'R', true
	case `
.###.
#....
#....
.##..
...#.
###..
`:
		return 'S', true
	case `
.###.
..#..
..#..
..#..
..#..
..#..
`:
		return 'T', true
	case `
#..#.
#..#.
#..#.
#..#.
#..#.
.##..
`:
		return 'U', true
	case `
#...#
#...#
#...#
.#.#.
.#.#.
..#..
`:
		return 'V', true
	case `
#...#
#...#
#.#.#
#.#.#
#.#.#
.#.#.
`:
		return 'W', true
	case `
#...#
.#.#.
..#..
.#.#.
.#.#.
#...#
`:
		return 'X', true
	case `
#...#
#...#
.#.#.
..#..
..#..
..#..
`:
		return 'Y', true
	case `
####.
...#.
..#..
.#...
#....
####.
`:
		return 'Z', true
	case `
.....
.....
.....
.....
.....
.....
`:
		return ' ', true
	default:
		return 0, false
	}
}

func tallDisplayLetter(height int, glyph string) (byte, bool) {
	if height == 8 {
		switch glyph {
		case `
#...#
#...#
#...#
#####
#...#
#...#
#...#
#...#
`:
			return 'H', true
		case `
###..
.#...
.#...
.#...
.#...
.#...
.#...
###..
`:
			return 'I', true
		}
		return 0, false
	}

	switch glyph {
	case `
..##..
.#..#.
#....#
#....#
#....#
######
#....#
#....#
#....#
#....#
`:
		return 'A', true
	case `
#####.
#....#
#....#
#....#
#####.
#....#
#....#
#....#
#....#
#####.
`:
		return 'B', true
	case `
.####.
#....#
#.....
#.....
#.....
#.....
#.....
#.....
#....#
.####.
`:
		return 'C', true
	case `
######
#.....
#.....
#.....
#####.
#.....
#.....
#.....
#.....
######
`:
		return 'E', true
	case `
######
#.....
#.....
#.....
#####.
#.....
#.....
#.....
#.....
#.....
`:
		return 'F', true
	case `
.####.
#....#
#.....
#.....
#.....
#..###
#....#
#....#
#...##
.###.#
`:
		return 'G', true
	case `
#....#
#....#
#....#
#....#
######
#....#
#....#
#....#
#....#
#....#
`:
		return 'H', true
	case `
...###
....#.
....#.
....#.
....#.
....#.
....#.
#...#.
#...#.
.###..
`:
		return 'J', true
	case `
#....#
#...#.
#..#..
#.#...
##....
##....
#.#...
#..#..
#...#.
#....#
`:
		return 'K', true
	case `
#.....
#.....
#.....
#.....
#.....
#.....
#.....
#.....
#.....
######
`:
		return 'L', true
	case `
#....#
##...#
##...#
#.#..#
#.#..#
#..#.#
#..#.#
#...##
#...##
#....#
`:
		return 'N', true
	case `
#####.
#....#
#....#
#....#
#####.
#.....
#.....
#.....
#.....
#.....
`:
		return 'P', true
	case `
#####.
#....#
#....#
#....#
#####.
#..#..
#...#.
#...#.
#....#
#....#
`:
		return 'R', true
	case `
#....#
#....#
.#..#.
.#..#.
..##..
..##..
.#..#.
.#..#.
#....#
#....#
`:
		return 'X', true
	case `
######
.....#
.....#
....#.
...#..
..#...
.#....
#.....
#.....
######
`:
		return 'Z', true
	default:
		return 0, false
	}
}
