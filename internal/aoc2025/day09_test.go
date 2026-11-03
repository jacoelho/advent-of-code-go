package aoc2025

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

func Test_day09p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`),
			Want: "50",
		},
		{
			Input: aoc.FileInput(t, 2025, 9),
			Want:  "4743645488",
		},
	}
	aoc.AOCTest(t, day09p01, tests)
}

func Test_day09p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`),
			Want: "24",
		},
		{
			Input: aoc.FileInput(t, 2025, 9),
			Want:  "1529011204",
		},
	}
	aoc.AOCTest(t, day09p02, tests)
}

func literalColoredTile(tiles []grid.Position2D[int64], x, y int64) bool {
	inside := false
	for i, a := range tiles {
		b := tiles[(i+1)%len(tiles)]
		if a.X == b.X && x == a.X && min(a.Y, b.Y) <= y && y <= max(a.Y, b.Y) ||
			a.Y == b.Y && y == a.Y && min(a.X, b.X) <= x && x <= max(a.X, b.X) {
			return true
		}
		if (a.Y > y) != (b.Y > y) && x < a.X {
			inside = !inside
		}
	}
	return inside
}

func TestRedTileGeometry(t *testing.T) {
	tests := []struct {
		name  string
		tiles []grid.Position2D[int64]
		want  [2]string
	}{
		{
			name:  "wide notch",
			tiles: []grid.Position2D[int64]{{X: 0, Y: 0}, {X: 4, Y: 0}, {X: 4, Y: 4}, {X: 3, Y: 4}, {X: 3, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 4}, {X: 0, Y: 4}},
			want:  [2]string{"25", "10"},
		},
		{
			name:  "notch without excluded tile column",
			tiles: []grid.Position2D[int64]{{X: 0, Y: 0}, {X: 3, Y: 0}, {X: 3, Y: 3}, {X: 2, Y: 3}, {X: 2, Y: 1}, {X: 1, Y: 1}, {X: 1, Y: 3}, {X: 0, Y: 3}},
			want:  [2]string{"16", "16"},
		},
		{
			name:  "boundary box",
			tiles: []grid.Position2D[int64]{{X: 0, Y: 0}, {X: 2, Y: 0}, {X: 2, Y: 1}, {X: 0, Y: 1}},
			want:  [2]string{"6", "6"},
		},
		{
			name:  "official polygon",
			tiles: []grid.Position2D[int64]{{X: 7, Y: 1}, {X: 11, Y: 1}, {X: 11, Y: 7}, {X: 9, Y: 7}, {X: 9, Y: 5}, {X: 2, Y: 5}, {X: 2, Y: 3}, {X: 7, Y: 3}},
			want:  [2]string{"50", "24"},
		},
	}
	for _, tt := range tests {
		for _, variant := range []string{"original", "reversed", "translated", "transposed"} {
			t.Run(tt.name+"/"+variant, func(t *testing.T) {
				tiles := slices.Clone(tt.tiles)
				switch variant {
				case "reversed":
					slices.Reverse(tiles)
				case "translated":
					for i := range tiles {
						tiles[i].X -= 6
						tiles[i].Y += 4
					}
				case "transposed":
					for i := range tiles {
						tiles[i].X, tiles[i].Y = tiles[i].Y, tiles[i].X
					}
				}
				var input strings.Builder
				for _, tile := range tiles {
					fmt.Fprintf(&input, "%d,%d\n", tile.X, tile.Y)
				}
				for part, solve := range []func(io.Reader) (string, error){day09p01, day09p02} {
					got, err := solve(strings.NewReader(input.String()))
					if err != nil || got != tt.want[part] {
						t.Fatalf("part %d: got %q, %v; want %q", part+1, got, err, tt.want[part])
					}
				}
				floor, err := newRedTileFloor(tiles)
				if err != nil {
					t.Fatal(err)
				}
				for i, a := range tiles {
					for _, b := range tiles[i:] {
						wantAllowed := true
						var wantArea int64
						for y := min(a.Y, b.Y); y <= max(a.Y, b.Y); y++ {
							for x := min(a.X, b.X); x <= max(a.X, b.X); x++ {
								wantArea++
								if !literalColoredTile(tiles, x, y) {
									wantAllowed = false
								}
							}
						}
						rectangle := newTileRectangle(a, b)
						if got := floor.contains(rectangle); got != wantAllowed {
							t.Errorf("corners %v/%v: allowed %v, want %v", a, b, got, wantAllowed)
						}
						if got, err := rectangle.area(); err != nil || got != wantArea {
							t.Errorf("corners %v/%v: area %d, %v; want %d", a, b, got, err, wantArea)
						}
					}
				}
			})
		}
	}
}

func TestRedTileInvalidInput(t *testing.T) {
	for _, solve := range []func(io.Reader) (string, error){day09p01, day09p02} {
		for _, input := range []string{"1", "1,2,3", "1,2junk"} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want coordinate error", input, got)
			}
		}
		if _, err := solve(strings.NewReader("9223372036854775808,0")); !errors.Is(err, strconv.ErrRange) {
			t.Errorf("got %v, want coordinate range error", err)
		}
		input := "0,0\n9223372036854775807,0\n9223372036854775807,1\n0,1"
		if got, err := solve(strings.NewReader(input)); err == nil {
			t.Errorf("got %q, want area overflow error", got)
		}
	}
	for _, input := range []string{
		"0,0\n2,0\n3,1\n0,1",
		"0,0\n2,0\n2,1\n1,1",
		"0,0\n4,0\n4,4\n2,4\n2,-1\n1,-1\n1,3\n0,3",
	} {
		if got, err := day09p02(strings.NewReader(input)); err == nil {
			t.Errorf("input %q: got %q, want polygon error", input, got)
		}
	}
}

func TestDay09ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day09p01, "0,0\n2,0\n2,1\n0,1\n"},
		{"part2", day09p02, "0,0\n2,0\n2,1\n0,1\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				input := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(input); !errors.Is(err, failure) {
					t.Fatalf("got %v; want reader failure", err)
				}
			}
		})
	}
}
