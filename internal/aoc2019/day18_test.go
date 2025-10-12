package aoc2019

import (
	"strings"
	"testing"

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
