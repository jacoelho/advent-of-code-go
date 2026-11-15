package aoc2019

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day20p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`         A
         A
  #######.#########
  #######.........#
  #######.#######.#
  #######.#######.#
  #######.#######.#
  #####  B    ###.#
BC...##  C    ###.#
  ##.##       ###.#
  ##...DE  F  ###.#
  #####    G  ###.#
  #########.#####.#
DE..#######...###.#
  #.#########.###.#
FG..#########.....#
  ###########.#####
             Z
             Z       `),
			Want: "23",
		},
		{
			Input: strings.NewReader(`                   A
                   A
  #################.#############
  #.#...#...................#.#.#
  #.#.#.###.###.###.#########.#.#
  #.#.#.......#...#.....#.#.#...#
  #.#########.###.#####.#.#.###.#
  #.............#.#.....#.......#
  ###.###########.###.#####.#.#.#
  #.....#        A   C    #.#.#.#
  #######        S   P    #####.#
  #.#...#                 #......VT
  #.#.#.#                 #.#####
  #...#.#               YN....#.#
  #.###.#                 #####.#
DI....#.#                 #.....#
  #####.#                 #.###.#
ZZ......#               QG....#..AS
  ###.###                 #######
JO..#.#.#                 #.....#
  #.#.#.#                 ###.#.#
  #...#..DI             BU....#..LF
  #####.#                 #.#####
YN......#               VT..#....QG
  #.###.#                 #.###.#
  #.#...#                 #.....#
  ###.###    J L     J    #.#.###
  #.....#    O F     P    #.#...#
  #.###.#####.#.#####.#####.###.#
  #...#.#.#...#.....#.....#.#...#
  #.#####.###.###.#.#.#########.#
  #...#.#.....#...#.#.#.#.....#.#
  #.###.#####.###.###.#.#.#######
  #.#.........#...#.............#
  #########.###.###.#############
           B   J   C
           U   P   P               `),
			Want: "58",
		},
		{
			Input: aoc.FileInput(t, 2019, 20),
			Want:  "674",
		},
	}
	aoc.AOCTest(t, day20p01, tests)
}

func Test_day20p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`         A
         A
  #######.#########
  #######.........#
  #######.#######.#
  #######.#######.#
  #######.#######.#
  #####  B    ###.#
BC...##  C    ###.#
  ##.##       ###.#
  ##...DE  F  ###.#
  #####    G  ###.#
  #########.#####.#
DE..#######...###.#
  #.#########.###.#
FG..#########.....#
  ###########.#####
             Z
             Z       `),
			Want: "26",
		},
		{
			Input: strings.NewReader(`             Z L X W       C
             Z P Q B       K
  ###########.#.#.#.#######.###############
  #...#.......#.#.......#.#.......#.#.#...#
  ###.#.#.#.#.#.#.#.###.#.#.#######.#.#.###
  #.#...#.#.#...#.#.#...#...#...#.#.......#
  #.###.#######.###.###.#.###.###.#.#######
  #...#.......#.#...#...#.............#...#
  #.#########.#######.#.#######.#######.###
  #...#.#    F       R I       Z    #.#.#.#
  #.###.#    D       E C       H    #.#.#.#
  #.#...#                           #...#.#
  #.###.#                           #.###.#
  #.#....OA                       WB..#.#..ZH
  #.###.#                           #.#.#.#
CJ......#                           #.....#
  #######                           #######
  #.#....CK                         #......IC
  #.###.#                           #.###.#
  #.....#                           #...#.#
  ###.###                           #.#.#.#
XF....#.#                         RF..#.#.#
  #####.#                           #######
  #......CJ                       NM..#...#
  ###.#.#                           #.###.#
RE....#.#                           #......RF
  ###.###        X   X       L      #.#.#.#
  #.....#        F   Q       P      #.#.#.#
  ###.###########.###.#######.#########.###
  #.....#...#.....#.......#...#.....#.#...#
  #####.#.###.#######.#######.###.###.#.#.#
  #.......#.......#.#.#.#.#...#...#...#.#.#
  #####.###.#####.#.#.#.#.###.###.#.###.###
  #.......#.....#.#...#...............#...#
  #############.#.#.###.###################
               A O F   N
               A A D   M                     `),
			Want: "396",
		},
		{
			Input: aoc.FileInput(t, 2019, 20),
			Want:  "7636",
		},
	}
	aoc.AOCTest(t, day20p02, tests)
}

func TestDay20ReaderErrors(t *testing.T) {
	failure := errors.New("input read failed")
	for _, tt := range []struct {
		name, prefix string
		solve        func(io.Reader) (string, error)
	}{
		{"day20p01", "AA...ZZ\n", day20p01},
		{"day20p02", "AA...ZZ\n", day20p02},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				reader := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(failure))
				if _, err := tt.solve(reader); !errors.Is(err, failure) {
					t.Fatalf("prefix %q: got %v, want input failure", prefix, err)
				}
			}
		})
	}
}

func TestDonutMazeCorridor(t *testing.T) {
	aoc.AOCTest(t, day20p01, []aoc.TestInput{{Input: strings.NewReader("AA...ZZ"), Want: "2"}})
	aoc.AOCTest(t, day20p02, []aoc.TestInput{{Input: strings.NewReader("AA...ZZ"), Want: "2"}})
}

func TestDonutMazeErrors(t *testing.T) {
	for _, tt := range []struct {
		name, input, cause string
	}{
		{"empty", "", "no walls"},
		{"missing ZZ", "AA...#", "no ZZ"},
		{"invalid cell", "AA.x.ZZ", "invalid maze cell"},
		{"single-ended portal", "AA..BC#..ZZ", "only one end"},
		{"unreachable ZZ", "AA.#.ZZ", "unreachable"},
		{"lone letter", "A...ZZ", "label letters"},
		{"label without passage", "AA#..ZZ", "exactly one adjacent passage"},
		{"label between passages", "AA..BC..ZZ", "exactly one adjacent passage"},
		{"duplicate AA", "AA...ZZ\n  .\n  A\n  A", "more than once"},
		{"shared passage", "AA.BC#BC.ZZ", "more than one label"},
		{"portal ends on one rim", "AA.#.BC#BC.#.ZZ", "one outer and one inner"},
		{"portal seen three times", "     A\n     A\n  ###.###\n  #.....#\nBC..#BC.#\n  #.....#\n  ###.###\n     B\n     C", "more than twice"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, solve := range []func(io.Reader) (string, error){day20p01, day20p02} {
				if _, err := solve(strings.NewReader(tt.input)); err == nil || !strings.Contains(err.Error(), tt.cause) {
					t.Errorf("error = %v, want one mentioning %q", err, tt.cause)
				}
			}
		})
	}
}
