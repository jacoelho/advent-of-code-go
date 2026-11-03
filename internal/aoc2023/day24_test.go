package aoc2023

import (
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
)

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader("300000000000000,299999999999999,0 @0,1,0\n299999999999999,300000000000000,0 @1,0,0"),
			Want:  "1",
		},
		{
			Input: strings.NewReader("199999999999999,300000000000000,0 @2,0,0\n200000000000000,299999999999999,0 @-1,3,0"),
			Want:  "0",
		},
		{
			Input: aoc.FileInput(t, 2023, 24),
			Want:  "26657",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	tests := []aoc.TestInput{
		{Input: strings.NewReader(hailFastRock), Want: "60"},
		{Input: strings.NewReader("609,-682,827 @1,2,3\n609,-682,827 @1,2,3\n" + hailFastRock), Want: "60"},
		{Input: strings.NewReader(hailFastRock + "\n310,-330,430 @0,0,0"), Want: "60"},
		{Input: strings.NewReader(hailRankFive), Want: "60"},
		{
			Input: strings.NewReader(hailExample),
			Want:  "47",
		},
		{
			Input: aoc.FileInput(t, 2023, 24),
			Want:  "828418331313365",
		},
	}
	aoc.AOCTest(t, day24p02, tests)
}

const hailExample = `19, 13, 30 @ -2, 1, -2
18, 19, 22 @ -1, -1, -2
20, 25, 34 @ -2, -2, -4
12, 31, 28 @ -1, -2, -1
20, 19, 15 @ 1, -5, -3`

const hailFastRock = `609,-682,827 @1,2,3
1214,-1382,1622 @-2,1,4
1801,-2068,2427 @3,-4,1
2985,-3510,4040 @5,6,-2
4259,-4936,5567 @-7,8,9`

const hailRankFive = `13,17,30 @-2,5,3
12,16,30 @-1,6,3
14,16,30 @-1,4,3
13,11,30 @0,5,3`

func TestHailPathIntersections(t *testing.T) {
	stones, err := parseHailstones(strings.NewReader(hailExample))
	if err != nil {
		t.Fatal(err)
	}
	if got := countHailIntersections(stones, 7, 27); got != 2 {
		t.Fatalf("official sample intersections = %d; want 2", got)
	}
	for _, tt := range []struct {
		name, input string
		want        int
	}{
		{"lower corner", "0,-1,0 @0,1,0\n-1,0,0 @1,0,0", 1},
		{"upper corner", "10,9,0 @0,1,0\n9,10,0 @1,0,0", 1},
		{"time zero", "5,5,0 @1,0,0\n5,4,0 @0,1,0", 1},
		{"past only", "5,5,0 @-1,0,0\n6,4,0 @0,1,0", 0},
		{"stationary reached", "5,5,0 @0,0,0\n4,5,0 @1,0,0", 1},
		{"stationary behind", "5,5,0 @0,0,0\n6,5,0 @1,0,0", 0},
		{"stationary same projection", "5,5,0 @0,0,1\n5,5,99 @0,0,-1", 1},
		{"parallel disjoint", "0,5,0 @1,0,0\n0,6,0 @1,0,0", 0},
		{"collinear overlap", "2,5,0 @1,0,0\n8,5,0 @-1,0,0", 1},
		{"collinear diverging", "2,5,0 @-1,0,0\n8,5,0 @1,0,0", 0},
		{"collinear shared origin", "5,5,0 @-1,0,0\n5,5,0 @1,0,0", 1},
		{"overlap outside square", "0,11,0 @1,1,0\n2,13,0 @-1,-1,0", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stones, err := parseHailstones(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if got := countHailIntersections(stones, 0, 10); got != tt.want {
				t.Fatalf("intersections = %d; want %d", got, tt.want)
			}
		})
	}
}

func TestHailRockRejectsInvalidCandidates(t *testing.T) {
	for _, input := range []string{
		hailFastRock + "\n0,0,0 @0,0,0",
		hailFastRock + "\n-585,717,-768 @5,-3,2",
		`0,0,-2 @1,2,3
8,2,-11 @-2,1,4
-12,28,-4 @3,-4,1
-31,-31,15 @5,6,-2
68,-58,-80 @-7,8,9`,
		`5,37,42 @-2,5,3
8,40,42 @-1,6,3
10,32,42 @-1,4,3
13,31,42 @0,5,3`,
	} {
		if _, err := day24p02(strings.NewReader(input)); err == nil {
			t.Fatalf("accepted invalid rock for %q", input)
		}
	}
}

func TestHailRockAmbiguousAndUnsupported(t *testing.T) {
	ambiguous := `21,-3,18 @-2,5,3
16,-8,18 @-1,6,3
18,0,18 @-1,4,3
13,-9,18 @0,5,3`
	if _, err := day24p02(strings.NewReader(ambiguous)); err == nil || !strings.Contains(err.Error(), "multiple integer rocks") {
		t.Fatalf("expected ambiguous-rock error, got %v", err)
	}
	if _, err := day24p02(strings.NewReader(strings.Repeat("0,0,0 @1,1,1\n", 3))); err == nil || !strings.Contains(err.Error(), "unsupported rock family") {
		t.Fatalf("expected unsupported-family error, got %v", err)
	}
}

func TestHailInputErrors(t *testing.T) {
	for _, input := range []string{
		"", "0,0,0 @1,1,1",
		"0,0,0 1,1,1", "0,0,0 @1,1,1 @2,2,2",
		"0,0 @1,1,1", "0,0,0,0 @1,1,1",
		"0,0,0 @1,1,1junk", "9223372036854775808,0,0 @1,1,1",
		strings.Repeat("0,0,0 @1,1,1\n", 2001),
	} {
		if _, err := parseHailstones(strings.NewReader(input)); err == nil {
			t.Errorf("accepted invalid hail input of %d bytes", len(input))
		}
	}
	largeInput := strings.Repeat("0,0,0 @1,1,1"+strings.Repeat(" ", 60_000)+"\n", 70)
	if _, err := parseHailstones(strings.NewReader(largeInput)); err == nil || !strings.Contains(err.Error(), "4 MiB") {
		t.Fatalf("expected byte limit error, got %v", err)
	}
}

func TestHailIntegerRoots(t *testing.T) {
	for _, tt := range []struct {
		name         string
		coefficients [3]matrix.Rat
		want         map[string]bool
	}{
		{"constant", [3]matrix.Rat{{}, {}, matrix.NewRat(1, 1)}, nil},
		{"linear integer", [3]matrix.Rat{{}, matrix.NewRat(2, 1), matrix.NewRat(-4, 1)}, map[string]bool{"2": true}},
		{"linear fraction", [3]matrix.Rat{{}, matrix.NewRat(2, 1), matrix.NewRat(-1, 1)}, nil},
		{"negative discriminant", [3]matrix.Rat{matrix.NewRat(1, 1), {}, matrix.NewRat(1, 1)}, nil},
		{"nonsquare discriminant", [3]matrix.Rat{matrix.NewRat(1, 1), {}, matrix.NewRat(-2, 1)}, nil},
		{"repeated root", [3]matrix.Rat{matrix.NewRat(1, 1), matrix.NewRat(-6, 1), matrix.NewRat(9, 1)}, map[string]bool{"3": true}},
		{"rational coefficients", [3]matrix.Rat{matrix.NewRat(1, 2), matrix.NewRat(-3, 2), matrix.NewRat(1, 1)}, map[string]bool{"1": true, "2": true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			roots := hailIntegerRoots(tt.coefficients)
			if len(roots) != len(tt.want) {
				t.Fatalf("roots = %v; want %v", roots, tt.want)
			}
			seen := make(map[string]bool)
			for _, root := range roots {
				value := root.String()
				if seen[value] || !tt.want[value] {
					t.Fatalf("roots = %v; want %v", roots, tt.want)
				}
				seen[value] = true
			}
		})
	}
}
