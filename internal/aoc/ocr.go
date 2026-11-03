package aoc

import "fmt"

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
		mask := uint32(displayMask(rows, start, start+cellWidth, 5))
		letter, ok := sixRowLetter(mask)
		if !ok {
			return "", fmt.Errorf("unknown display glyph in cell %d (mask %#x)", start/cellWidth, mask)
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

func displayMask(rows []string, start, end, maskWidth int) uint64 {
	var mask uint64
	for _, row := range rows {
		for x := range maskWidth {
			mask <<= 1
			if start+x < end && row[start+x] == '#' {
				mask |= 1
			}
		}
	}

	return mask
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

	maskWidth := 6
	if len(rows) != 10 {
		maskWidth = 5
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
		if x-start > maskWidth {
			return "", fmt.Errorf("display glyph at column %d is too wide", start)
		}
		mask := displayMask(rows, start, x, maskWidth)
		var letter byte
		var ok bool
		if len(rows) == 6 {
			// Cropping can remove leading blank columns, as in I and T.
			for padding := 0; padding <= maskWidth-(x-start); padding++ {
				letter, ok = sixRowLetter(uint32(mask >> padding))
				if ok {
					break
				}
			}
		} else {
			letter, ok = tallDisplayLetter(len(rows), mask)
		}
		if !ok {
			return "", fmt.Errorf("unknown display glyph at column %d (mask %#x)", start, mask)
		}
		text = append(text, letter)
	}

	return string(text), nil
}

// Five bits per row retain full-width glyphs such as Y; most letters leave the fifth column dark.
func sixRowLetter(mask uint32) (byte, bool) {
	switch mask {
	case 0x19297A52:
		return 'A', true
	case 0x392E4A5C:
		return 'B', true
	case 0x1928424C:
		return 'C', true
	case 0x39294A5C:
		return 'D', true
	case 0x3D0E421E:
		return 'E', true
	case 0x3D0E4210:
		return 'F', true
	case 0x19285A4E:
		return 'G', true
	case 0x252F4A52:
		return 'H', true
	case 0x1C42108E:
		return 'I', true
	case 0x0C210A4C:
		return 'J', true
	case 0x254C5292:
		return 'K', true
	case 0x2108421E:
		return 'L', true
	case 0x23BAC631:
		return 'M', true
	case 0x252D5A52:
		return 'N', true
	case 0x19294A4C:
		return 'O', true
	case 0x39297210:
		return 'P', true
	case 0x19295A4D:
		return 'Q', true
	case 0x39297292:
		return 'R', true
	case 0x1D08305C:
		return 'S', true
	case 0x1C421084:
		return 'T', true
	case 0x25294A4C:
		return 'U', true
	case 0x2318A944:
		return 'V', true
	case 0x231AD6AA:
		return 'W', true
	case 0x22A22951:
		return 'X', true
	case 0x23151084:
		return 'Y', true
	case 0x3C22221E:
		return 'Z', true
	case 0:
		return ' ', true
	default:
		return 0, false
	}
}

// Ten-row glyphs: mstksg/advent-of-code-ocr, LetterMap.hs, rawLetterforms2.
func tallDisplayLetter(height int, mask uint64) (byte, bool) {
	if height == 8 {
		switch mask {
		case 0x8c63f8c631:
			return 'H', true
		case 0xe21084211c:
			return 'I', true
		}
		return 0, false
	}

	switch mask {
	case 0x31286187f861861:
		return 'A', true
	case 0xfa1861fa186187e:
		return 'B', true
	case 0x7a182082082085e:
		return 'C', true
	case 0xfe0820fa082083f:
		return 'E', true
	case 0xfe0820fa0820820:
		return 'F', true
	case 0x7a18208278618dd:
		return 'G', true
	case 0x861861fe1861861:
		return 'H', true
	case 0x1c20820820a289c:
		return 'J', true
	case 0x862928c30a248a1:
		return 'K', true
	case 0x82082082082083f:
		return 'L', true
	case 0x871c69a659638e1:
		return 'N', true
	case 0xfa1861fa0820820:
		return 'P', true
	case 0xfa1861fa48a2861:
		return 'R', true
	case 0x86149230c492861:
		return 'X', true
	case 0xfc104210842083f:
		return 'Z', true
	default:
		return 0, false
	}
}
