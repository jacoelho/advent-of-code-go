package aoc

import (
	"strings"
	"testing"
)

func TestDecodeSeparatedDisplay(t *testing.T) {
	tests := []struct {
		name, display, want string
	}{
		{
			name:    "cropped six-row I",
			display: "###\n.#.\n.#.\n.#.\n.#.\n###",
			want:    "I",
		},
		{
			name:    "full-width six-row Y",
			display: "#...#\n#...#\n.#.#.\n..#..\n..#..\n..#..",
			want:    "Y",
		},
		{
			name: "official eight-row HI",
			display: `#...#..###
#...#...#.
#...#...#.
#####...#.
#...#...#.
#...#...#.
#...#...#.
#...#..###`,
			want: "HI",
		},
		{
			name: "established ten-row alphabet",
			display: `..##...#####...####..######.######..####..#....#....###.#....#.#......#....#.#####..#####..#....#.######
.#..#..#....#.#....#.#......#......#....#.#....#.....#..#...#..#......##...#.#....#.#....#.#....#......#
#....#.#....#.#......#......#......#......#....#.....#..#..#...#......##...#.#....#.#....#..#..#.......#
#....#.#....#.#......#......#......#......#....#.....#..#.#....#......#.#..#.#....#.#....#..#..#......#.
#....#.#####..#......#####..#####..#......######.....#..##.....#......#.#..#.#####..#####....##......#..
######.#....#.#......#......#......#..###.#....#.....#..##.....#......#..#.#.#......#..#.....##.....#...
#....#.#....#.#......#......#......#....#.#....#.....#..#.#....#......#..#.#.#......#...#...#..#...#....
#....#.#....#.#......#......#......#....#.#....#.#...#..#..#...#......#...##.#......#...#...#..#..#.....
#....#.#....#.#....#.#......#......#...##.#....#.#...#..#...#..#......#...##.#......#....#.#....#.#.....
#....#.#####...####..######.#.......###.#.#....#..###...#....#.######.#....#.#......#....#.#....#.######`,
			want: "ABCEFGHJKLNPRXZ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeSeparatedDisplay(strings.Split(tt.display, "\n"))
			if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
}

func TestDecodeSeparatedDisplayErrors(t *testing.T) {
	tests := [][]string{
		{},
		{"#", "#", "#", "#", "#", "#"},
		{"#"},
		{"#", "#", "#", "#", "#", "#", "#", ""},
		{"?", ".", ".", ".", ".", ".", ".", "."},
		{"#", "#", "#", "#", "#", "#", "#", "#"},
		{"#######", "#######", "#######", "#######", "#######", "#######", "#######", "#######", "#######", "#######"},
	}
	for _, rows := range tests {
		if _, err := DecodeSeparatedDisplay(rows); err == nil {
			t.Errorf("expected error for %q", rows)
		}
	}
}

func TestDecode6RowDisplay(t *testing.T) {
	tests := []struct {
		name, display, want string
		cellWidth           int
	}{
		{"A", ".##.\n#..#\n#..#\n####\n#..#\n#..#", "A", 4},
		{"Z", "####\n...#\n..#.\n.#..\n#...\n####", "Z", 4},
		{"packed AZ", ".##.####\n#..#...#\n#..#..#.\n####.#..\n#..##...\n#..#####", "AZ", 4},
		{"five-column AZ", ".##..####.\n#..#....#.\n#..#...#..\n####..#...\n#..#.#....\n#..#.####.", "AZ", 5},
		{"full-width Y", "#...#\n#...#\n.#.#.\n..#..\n..#..\n..#..", "Y", 5},
		{"space pixels", " ## \n#  #\n#  #\n####\n#  #\n#  #", "A", 4},
		{"blank cell", "....\n....\n....\n....\n....\n....", " ", 4},
		{"empty display", "\n\n\n\n\n", "", 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode6RowDisplay(strings.Split(tt.display, "\n"), tt.cellWidth)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecode6RowDisplayErrors(t *testing.T) {
	tests := []struct {
		name      string
		rows      []string
		cellWidth int
	}{
		{"zero cell width", []string{"", "", "", "", "", ""}, 0},
		{"unsupported cell width", []string{"", "", "", "", "", ""}, 6},
		{"five rows", []string{"", "", "", "", ""}, 4},
		{"seven rows", []string{"", "", "", "", "", "", ""}, 4},
		{"ragged rows", []string{"....", "....", "...", "....", "....", "...."}, 4},
		{"partial cell", []string{"...", "...", "...", "...", "...", "..."}, 4},
		{"invalid pixel", []string{"?...", "....", "....", "....", "....", "...."}, 4},
		{"unknown glyph", []string{"####", "####", "####", "####", "####", "####"}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Decode6RowDisplay(tt.rows, tt.cellWidth); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
