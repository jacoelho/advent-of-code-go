package aoc2024

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const smallCircuitExample = `x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02`

const largeCircuitExample = `x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj`

func Test_day24p01(t *testing.T) {
	tests := []aoc.TestInput{
		{
			Input: strings.NewReader(smallCircuitExample),
			Want:  "4",
		},
		{
			Input: strings.NewReader(largeCircuitExample),
			Want:  "2024",
		},
		{
			Input: aoc.FileInput(t, 2024, 24),
			Want:  "51107420031718",
		},
	}
	aoc.AOCTest(t, day24p01, tests)
}

func Test_day24p02(t *testing.T) {
	aoc.AOCTest(t, day24p02, []aoc.TestInput{
		{Input: aoc.FileInput(t, 2024, 24), Want: "cpm,ghp,gpr,krs,nks,z10,z21,z33"},
	})
}

func adderFixture(width int, x, y uint64, swaps [][2]string, reversed bool) string {
	var input strings.Builder
	for bit := range width {
		fmt.Fprintf(&input, "x%02d: %d\ny%02d: %d\n", bit, x>>bit&1, bit, y>>bit&1)
	}
	input.WriteByte('\n')
	var gates []string
	gate := func(left, operation, right, output string) {
		for _, pair := range swaps {
			if output == pair[0] {
				output = pair[1]
				break
			}
			if output == pair[1] {
				output = pair[0]
				break
			}
		}
		if reversed {
			left, right = right, left
		}
		gates = append(gates, fmt.Sprintf("%s %s %s -> %s", left, operation, right, output))
	}
	gate("x00", "XOR", "y00", "z00")
	carry := "c01"
	if width == 1 {
		carry = "z01"
	}
	gate("x00", "AND", "y00", carry)
	for bit := 1; bit < width; bit++ {
		xName, yName := fmt.Sprintf("x%02d", bit), fmt.Sprintf("y%02d", bit)
		p, g, h := fmt.Sprintf("p%02d", bit), fmt.Sprintf("g%02d", bit), fmt.Sprintf("h%02d", bit)
		gate(xName, "XOR", yName, p)
		gate(xName, "AND", yName, g)
		gate(p, "XOR", carry, fmt.Sprintf("z%02d", bit))
		gate(p, "AND", carry, h)
		next := fmt.Sprintf("c%02d", bit+1)
		if bit+1 == width {
			next = fmt.Sprintf("z%02d", width)
		}
		gate(g, "OR", h, next)
		carry = next
	}
	if reversed {
		slices.Reverse(gates)
	}
	input.WriteString(strings.Join(gates, "\n"))
	return input.String()
}

func TestCircuitAdderTruthTable(t *testing.T) {
	for x := uint64(0); x < 8; x++ {
		for y := uint64(0); y < 8; y++ {
			got, err := day24p01(strings.NewReader(adderFixture(3, x, y, nil, false)))
			want := fmt.Sprint(x + y)
			if err != nil || got != want {
				t.Fatalf("%d + %d: got %q, %v; want %q", x, y, got, err, want)
			}
		}
	}
}

var day24DetectableSwaps = [][2]string{{"z01", "h01"}, {"z02", "c03"}, {"p03", "g03"}, {"z04", "g04"}}

func TestCircuitRepairFindsDetectableSwaps(t *testing.T) {
	for _, reversed := range []bool{false, true} {
		input := adderFixture(5, 0, 0, day24DetectableSwaps, reversed)
		got, err := day24p02(strings.NewReader(input))
		const want = "c03,g03,g04,h01,p03,z01,z02,z04"
		if err != nil || got != want {
			t.Errorf("reversed %v: got %q, %v; want %q", reversed, got, err, want)
		}
	}
}

func TestCircuitRepairedSwapsAdd(t *testing.T) {
	device, err := parseMonitoringDevice(strings.NewReader(adderFixture(5, 0, 0, day24DetectableSwaps, false)))
	if err != nil {
		t.Fatal(err)
	}
	pairs, err := repairRippleAdder(device)
	if err != nil {
		t.Fatal(err)
	}
	swapped := make(map[string]string)
	for _, pair := range pairs {
		swapped[pair[0]], swapped[pair[1]] = pair[1], pair[0]
	}
	for i, gate := range device.gates {
		if other, ok := swapped[gate.output]; ok {
			device.gates[i].output = other
		}
	}
	for x := uint64(0); x < 32; x++ {
		for y := uint64(0); y < 32; y++ {
			for bit := range 5 {
				device.initial[fmt.Sprintf("x%02d", bit)] = x>>bit&1 != 0
				device.initial[fmt.Sprintf("y%02d", bit)] = y>>bit&1 != 0
			}
			values, err := evaluateDevice(device)
			if err != nil {
				t.Fatal(err)
			}
			got, err := deviceNumber(values, 'z')
			if err != nil || got != x+y {
				t.Fatalf("repaired %d + %d: got %d, %v; want %d", x, y, got, err, x+y)
			}
		}
	}
}

func TestCircuitRepairRejectsSwapsTheRulesMiss(t *testing.T) {
	swaps := append(slices.Clone(day24DetectableSwaps), [2]string{"z00", "z03"})
	_, err := day24p02(strings.NewReader(adderFixture(5, 0, 0, swaps, false)))
	if !errors.Is(err, errDay24RepairUncertified) {
		t.Fatalf("got %v, want uncertified repair error", err)
	}
}

func TestCircuitRepairNeedsFourFlaggedPairs(t *testing.T) {
	_, err := day24p02(strings.NewReader(adderFixture(5, 0, 0, day24DetectableSwaps[:3], false)))
	if err == nil || !strings.Contains(err.Error(), "flags 6 wires") {
		t.Fatalf("got %v, want flagged wire count error", err)
	}
}

func TestCircuitRepairNeedsMatchingXYInputs(t *testing.T) {
	for _, input := range []string{
		strings.ReplaceAll(adderFixture(5, 0, 0, nil, false), "y04", "w04"),
		"x05: 1\n" + adderFixture(5, 0, 0, nil, false),
	} {
		_, err := day24p02(strings.NewReader(input))
		if err == nil || !strings.Contains(err.Error(), "adder") || strings.Contains(err.Error(), "flags") {
			t.Errorf("got %v, want adder input error", err)
		}
	}
}

func TestCircuitOutputWiderThan64BitsRejected(t *testing.T) {
	input := "a: 1\nb: 0\n\na XOR b -> z64"
	if got, err := day24p01(strings.NewReader(input)); err == nil {
		t.Fatalf("got %q, want output width error", got)
	}
}

func TestCircuitOutput64Bits(t *testing.T) {
	input := "a: 1\nb: 0\n\na XOR b -> z63\na AND b -> z00"
	got, err := day24p01(strings.NewReader(input))
	if err != nil || got != "9223372036854775808" {
		t.Fatalf("got %q, %v; want 9223372036854775808", got, err)
	}
}

func TestCircuitRepeatedGateInput(t *testing.T) {
	input := "a: 1\nb: 0\n\nmid XOR mid -> z00\nmid OR mid -> z01\na OR b -> mid"
	got, err := day24p01(strings.NewReader(input))
	if err != nil || got != "2" {
		t.Fatalf("got %q, %v; want 2", got, err)
	}
}

func TestCircuitUnsupportedRepair(t *testing.T) {
	for _, input := range []string{smallCircuitExample, largeCircuitExample} {
		if got, err := day24p02(strings.NewReader(input)); err == nil {
			t.Errorf("got %q, want unsupported-adder error", got)
		}
	}
}

func TestCircuitLineEndings(t *testing.T) {
	for _, input := range []string{smallCircuitExample + "\n", strings.ReplaceAll(smallCircuitExample+"\n", "\n", "\r\n")} {
		got, err := day24p01(strings.NewReader(input))
		if err != nil || got != "4" {
			t.Errorf("got %q, %v; want 4", got, err)
		}
	}
}

func TestCircuitInvalidInput(t *testing.T) {
	for _, input := range []string{
		"x00: 2\ny00: 0\n\nx00 AND y00 -> z00",
		"x00: 1\ny00: 0\n\nx00 NOR y00 -> z00",
		"x00: 1\ny00: 0\n\nx00 AND y00 => z00",
		"x00: 1\nx00: 0\n\nx00 AND x00 -> z00",
		"x00: 1\ny00: 0\n\nx00 AND y00 -> z00\nx00 OR y00 -> z00",
		"x00: 1\ny00: 0\n\nx00 AND y00 -> x00",
		"x00: 1\ny00: 0\n\nx00 AND missing -> z00",
	} {
		for _, solve := range []func(io.Reader) (string, error){day24p01, day24p02} {
			if got, err := solve(strings.NewReader(input)); err == nil {
				t.Errorf("input %q: got %q, want error", input, got)
			}
		}
	}
	for _, input := range []string{
		"x00: 1\n\na XOR x00 -> b\nb AND x00 -> a\na OR x00 -> z00",
		"x00: 1\n\nx00 OR x00 -> zbad",
		"x00: 1\n\nx00 OR x00 -> z4096",
	} {
		if got, err := day24p01(strings.NewReader(input)); err == nil {
			t.Errorf("input %q: got %q, want evaluation/output error", input, got)
		}
	}
}

func TestDay24ReaderErrors(t *testing.T) {
	want := errors.New("input read failure")
	for _, tt := range []struct {
		name   string
		solve  func(io.Reader) (string, error)
		prefix string
	}{
		{"part1", day24p01, "x00: 0\ny00: 0\n\nx00 XOR y00 -> z00\nx00 AND y00 -> z01\n"},
		{"part2", day24p02, "x00: 0\ny00: 0\n\nx00 XOR y00 -> z00\nx00 AND y00 -> z01\n"},
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
