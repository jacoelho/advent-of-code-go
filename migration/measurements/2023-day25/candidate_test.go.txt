package aoc2023

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func Test_day25p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(`jqt: rhn xhk nvd
rsh: frs pzl lsr
xhk: hfx
cmg: qnr nvd lhk bvb
rhn: xhk bvb hfx
bvb: xhk hfx
pzl: lsr hfx nvd
qnr: nvd
ntq: jqt hfx bvb xhk
nvd: lhk
lsr: lhk
rzs: qnr cmg lsr rsh
frs: qnr lhk lsr`),
			Want: "54",
		},
		{
			Input: aoc.FileInput(t, 2023, 25),
			Want:  "614655",
		},
	}
	aoc.AOCTest(t, day25p01, tests)
}

func wiringFixture(vertices int, wires [][2]int) string {
	var result strings.Builder
	for node := range vertices {
		fmt.Fprintf(&result, "n%02d:", node)
		for _, wire := range wires {
			if wire[0] == node {
				fmt.Fprintf(&result, " n%02d", wire[1])
			}
		}
		result.WriteByte('\n')
	}
	return result.String()
}

func joinedCliques(wires int) string {
	var edges [][2]int
	for _, group := range [][2]int{{0, 5}, {5, 11}} {
		for a := group[0]; a < group[1]; a++ {
			for b := a + 1; b < group[1]; b++ {
				edges = append(edges, [2]int{a, b})
			}
		}
	}
	for i := range wires {
		edges = append(edges, [2]int{i, 5 + i})
	}
	return wiringFixture(11, edges)
}

func TestThreeWirePartition(t *testing.T) {
	aoc.AOCTest(t, day25p01, []aoc.TestInput{{Input: strings.NewReader(joinedCliques(3)), Want: "30"}})
}

func TestWiringMinimumCut(t *testing.T) {
	var complete [][2]int
	for a := range 5 {
		for b := a + 1; b < 5; b++ {
			complete = append(complete, [2]int{a, b})
		}
	}
	tests := []struct{ name, input string }{
		{"disconnected", joinedCliques(0)},
		{"one wire", joinedCliques(1)},
		{"two wires", joinedCliques(2)},
		{"four wires", wiringFixture(5, complete)},
		{"later smaller cut", "a: b c d z\nb: c d\nc: d\nz:"},
		{"isolated vertex", joinedCliques(3) + "isolated:\n"},
		{"self loop", "a: a b\nb: c\nc: a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := day25p01(strings.NewReader(tt.input)); err == nil {
				t.Fatal("expected unsupported wiring error")
			}
		})
	}
}

const cancellationWiring = `aa: bb xx0 zz
bb: cc dd0 dd1
cc: zz xx1 xx2
dd0: dd1 dd2 dd3
dd1: dd2 dd3
dd2: dd3 zz
xx0: xx1 xx2 xx3
xx1: xx2 xx3
xx2: xx3`

func TestWiringResidualCancellation(t *testing.T) {
	got, err := day25p01(strings.NewReader(cancellationWiring))
	if err != nil {
		t.Fatal(err)
	}
	switch got {
	case "11", "32", "35":
	default:
		t.Fatalf("partition product %s is not a minimum-cut product", got)
	}
}

func TestWiringDeterministicOrder(t *testing.T) {
	for _, input := range []string{joinedCliques(3), cancellationWiring} {
		want, err := day25p01(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")
		slices.Reverse(lines)
		var reordered strings.Builder
		for _, line := range lines {
			node, neighbors, _ := strings.Cut(line, ":")
			fields := strings.Fields(neighbors)
			slices.Reverse(fields)
			fmt.Fprintf(&reordered, "%s: %s\n", node, strings.Join(fields, " "))
			for _, neighbor := range fields {
				fmt.Fprintf(&reordered, "%s: %s %s\n", neighbor, node, node)
			}
		}
		got, err := day25p01(strings.NewReader(reordered.String()))
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("reordered/duplicate graph got %s, want %s", got, want)
		}
	}
}

func TestWiringFiveVertexGraphs(t *testing.T) {
	var possible [][2]int
	for a := range 5 {
		for b := a + 1; b < 5; b++ {
			possible = append(possible, [2]int{a, b})
		}
	}
	for graph := range 1 << len(possible) {
		var wires [][2]int
		for i, wire := range possible {
			if graph&(1<<i) != 0 {
				wires = append(wires, wire)
			}
		}
		minimum := len(wires) + 1
		products := make(map[int]bool)
		// Fix vertex zero on one side and enumerate every nontrivial bipartition.
		for side := 1; side < (1<<5)-1; side += 2 {
			crossing := 0
			for _, wire := range wires {
				if (side>>wire[0])&1 != (side>>wire[1])&1 {
					crossing++
				}
			}
			if crossing > minimum {
				continue
			}
			if crossing < minimum {
				minimum = crossing
				clear(products)
			}
			size := 0
			for node := range 5 {
				size += (side >> node) & 1
			}
			products[size*(5-size)] = true
		}
		got, err := day25p01(strings.NewReader(wiringFixture(5, wires)))
		if minimum != 3 {
			if err == nil {
				t.Fatalf("graph %010b: minimum %d, unexpectedly returned %s", graph, minimum, got)
			}
			continue
		}
		if err != nil {
			t.Fatalf("graph %010b: %v", graph, err)
		}
		product, err := strconv.Atoi(got)
		if err != nil || !products[product] {
			t.Fatalf("graph %010b: product %q not in %v", graph, got, products)
		}
	}
}
