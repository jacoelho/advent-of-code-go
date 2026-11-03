package aoc2019

import (
	"errors"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day18p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`#########
#b.A.@.a#
#########`),
			Want: "8",
		},
		{
			Input: strings.NewReader(`########################
#f.D.E.e.C.b.A.@.a.B.c.#
######################.#
#d.....................#
########################`),
			Want: "86",
		},
		{
			Input: strings.NewReader(`########################
#...............b.C.D.f#
#.######################
#.....@.a.B.c.d.A.e.F.g#
########################`),
			Want: "132",
		},
		{
			Input: strings.NewReader(`#################
#i.G..c...e..H.p#
########.########
#j.A..b...f..D.o#
########@########
#k.E..a...g..B.n#
########.########
#l.F..d...h..C.m#
#################`),
			Want: "136",
		},
		{
			Input: strings.NewReader(`########################
#@..............ac.GI.b#
###d#e#f################
###A#B#C################
###g#h#i################
########################`),
			Want: "81",
		},
		{
			Input: aoc.FileInput(t, 2019, 18),
			Want:  "3862",
		},
	}
	aoc.AOCTest(t, day18p01, tests)
}

func Test_day18p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`#######
#a.#Cd#
##...##
##.@.##
##...##
#cB#Ab#
#######`),
			Want: "8",
		},
		{
			Input: strings.NewReader(`###############
#d.ABC.#.....a#
######@#@######
###############
######@#@######
#b.....#.....c#
###############`),
			Want: "24",
		},
		{
			Input: strings.NewReader(`#############
#DcBa.#.GhKl#
#.###@#@#I###
#e#d#####j#k#
###C#@#@###J#
#fEbA.#.FgHi#
#############`),
			Want: "32",
		},
		{
			Input: strings.NewReader(`#############
#g#f.D#..h#l#
#F###e#E###.#
#dCba@#@BcIJ#
#############
#nK.L@#@G...#
#M###N#H###.#
#o#m..#i#jk.#
#############`),
			Want: "72",
		},
		{
			Input: aoc.FileInput(t, 2019, 18),
			Want:  "1626",
		},
	}
	aoc.AOCTest(t, day18p02, tests)
}

func TestVaultRoutes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"unlocked detour", "#####\n#@Aa#\n#.#.#\n#...#\n#####", "6"},
		{"equal routes", "#######\n#..A..#\n#@###a#\n#.....#\n#######", "6"},
		{"key before door", "#@.a.A.b#", "6"},
		{"owned key junction", "#######\n#b.a.c#\n###.###\n###@###\n#######", "8"},
		{"repeated doors", "#@aAA.b#", "5"},
		{"no keys", "@", "0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aoc.AOCTest(t, day18p01, []aoc.TestInput{{Input: strings.NewReader(tt.input), Want: tt.want}})
		})
	}
	aoc.AOCTest(t, day18p02, []aoc.TestInput{
		{Input: strings.NewReader("##############\n#@@a....@@A.b#\n##############"), Want: "4"},
		{Input: strings.NewReader("@@\n@@"), Want: "0"},
		{Input: strings.NewReader("...\n.@.\n..."), Want: "0"},
	})
}

func TestVaultErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"ragged", "@.\n."},
		{"unknown tile", "@?"},
		{"duplicate key", "@aa"},
		{"missing entrance", ".a"},
		{"unreachable", "@#a"},
		{"unmatched door", "#@A.b#"},
		{"too many nodes", "@" + strings.Repeat("A", 256)},
		{"too many cells", "@" + strings.Repeat(".", 499) + "\n" + strings.Repeat(strings.Repeat(".", 500)+"\n", 200)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := day18p01(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected invalid or unreachable vault error")
			}
		})
	}
	for _, input := range []string{"@", "@@", ".a.\n.@.\n...", "...\n.@.\n.#."} {
		if _, err := day18p02(strings.NewReader(input)); err == nil {
			t.Fatalf("expected invalid split/entrance error for %q", input)
		}
	}
	if _, err := day18p01(strings.NewReader("@@")); err == nil {
		t.Fatal("part one accepted two entrances")
	}
	failure := errors.New("vault reader failed")
	for _, solve := range []func(io.Reader) (string, error){day18p01, day18p02} {
		input := io.MultiReader(strings.NewReader("...\n.@.\n...\n"), iotest.ErrReader(failure))
		if _, err := solve(input); !errors.Is(err, failure) {
			t.Fatalf("error = %v, want reader failure", err)
		}
	}
}

func TestVaultWorkLimit(t *testing.T) {
	v, err := parseVault(strings.NewReader("#####\n#@Aa#\n#.#.#\n#...#\n#####"))
	if err != nil {
		t.Fatal(err)
	}
	exhausted := vaultSearch{}
	if err := exhausted.buildGraph(v); err == nil {
		t.Fatal("expected graph construction budget error")
	}
	search := vaultSearch{remaining: maxVaultWork, memo: make(map[vaultState]int)}
	if err := search.buildGraph(v); err != nil {
		t.Fatal(err)
	}
	var start vaultState
	copy(start.positions[:], search.graph.entrances)
	search.remaining = 0
	if _, err := search.remainingDistance(start); err == nil {
		t.Fatal("expected key search budget error")
	}
	search.remaining = maxVaultWork
	if got, err := search.remainingDistance(start); err != nil || got != 6 {
		t.Fatalf("search after budget failure = %d, %v; want 6", got, err)
	}
}

func TestVaultTinyLayouts(t *testing.T) {
	random := rand.New(rand.NewPCG(18, 2019))
	for _, robots := range []int{1, 4} {
		width, height := 3, 3
		keys := 3
		if robots == 4 {
			width, height, keys = 4, 2, 2
		}
		for range 40 {
			cells := []byte(strings.Repeat(".", width*height))
			order := random.Perm(len(cells))
			for _, at := range order[:robots] {
				cells[at] = '@'
			}
			for i, at := range order[robots : robots+keys] {
				cells[at] = byte('a' + i)
			}
			for _, at := range order[robots+keys:] {
				switch random.IntN(3) {
				case 0:
					cells[at] = '#'
				case 1:
					cells[at] = byte('A' + random.IntN(keys))
				}
			}
			var rows []string
			for y := range height {
				rows = append(rows, string(cells[y*width:(y+1)*width]))
			}
			input := strings.Join(rows, "\n")
			want, reachable := vaultCellOracle(rows)
			got, err := vaultDistance(strings.NewReader(input), robots)
			if !reachable {
				if err == nil {
					t.Fatalf("vault %q returned %s for unreachable keys", input, got)
				}
			} else if err != nil || got != strconv.Itoa(want) {
				t.Fatalf("vault %q = %q, %v; want %d", input, got, err, want)
			}
		}
	}
}

func vaultCellOracle(rows []string) (int, bool) {
	type state struct {
		positions [4]int
		keys      uint32
	}
	type visit struct {
		state
		distance int
	}
	width := len(rows[0])
	var start state
	var allKeys uint32
	robots := 0
	for y, row := range rows {
		for x, tile := range []byte(row) {
			if tile == '@' {
				start.positions[robots] = y*width + x
				robots++
			} else if tile >= 'a' && tile <= 'z' {
				allKeys |= 1 << (tile - 'a')
			}
		}
	}
	queue := []visit{{state: start}}
	seen := map[state]bool{start: true}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if current.keys == allKeys {
			return current.distance, true
		}
		for robot := range robots {
			x, y := current.positions[robot]%width, current.positions[robot]/width
			for _, delta := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				nx, ny := x+delta[0], y+delta[1]
				if nx < 0 || ny < 0 || nx >= width || ny >= len(rows) {
					continue
				}
				tile := rows[ny][nx]
				if tile == '#' || tile >= 'A' && tile <= 'Z' && current.keys&(1<<(tile-'A')) == 0 {
					continue
				}
				next := current.state
				next.positions[robot] = ny*width + nx
				if tile >= 'a' && tile <= 'z' {
					next.keys |= 1 << (tile - 'a')
				}
				if !seen[next] {
					seen[next] = true
					queue = append(queue, visit{state: next, distance: current.distance + 1})
				}
			}
		}
	}
	return 0, false
}
