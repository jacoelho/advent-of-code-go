package aoc2024

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day23p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`),
			Want: "7",
		},
		{
			Input: aoc.FileInput(t, 2024, 23),
			Want:  "1043",
		},
	}
	aoc.AOCTest(t, day23p01, tests)
}

func Test_day23p02(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`kh-tc
qp-kh
de-cg
ka-co
yn-aq
qp-ub
cg-tb
vc-aq
tb-ka
wh-tc
yn-cg
kh-ub
ta-co
de-co
tc-td
tb-wq
wh-td
ta-ka
td-qp
aq-cg
wq-ub
ub-vc
de-ta
wq-aq
wq-vc
wh-yn
ka-de
kh-ta
co-tc
wh-qp
tb-vc
td-yn`),
			Want: "co,de,ka,ta",
		},
		{
			Input: aoc.FileInput(t, 2024, 23),
			Want:  "ai,bk,dc,dx,fo,gx,hk,kd,os,uz,xn,yk,zs",
		},
	}
	aoc.AOCTest(t, day23p02, tests)
}

func TestNetworkConnectionErrors(t *testing.T) {
	for _, input := range []string{"aa", "aa-bb-cc", "aa-", "aa-aa"} {
		for _, solve := range []func(io.Reader) (string, error){day23p01, day23p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Fatalf("input %q: got %q; want network connection error", input, got)
			}
		}
	}
}

func TestDay23ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day23p01, "ta-tb\n"},
		{"part2", day23p02, "ta-tb\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			for _, prefix := range []string{"", tt.prefix} {
				r := io.MultiReader(strings.NewReader(prefix), iotest.ErrReader(want))
				if _, err := tt.solve(r); !errors.Is(err, want) {
					t.Errorf("prefix %q: got %v, want reader error", prefix, err)
				}
			}
		})
	}
}
