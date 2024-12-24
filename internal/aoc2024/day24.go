package aoc2024

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/collections"
	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/xmaps"
)

func parseMonitoringDevice(r io.Reader) (map[string]int, [][]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, err
	}
	parts := strings.Split(string(data), "\n\n")
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid input format")
	}

	gates := make(map[string]int)
	for _, line := range strings.Split(parts[0], "\n") {
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("invalid gate format: %s", line)
		}
		gates[parts[0]], err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, nil, err
		}
	}

	var wires [][]string
	for _, line := range strings.Split(parts[1], "\n") {
		fields := strings.Fields(line)
		if len(fields) != 5 {
			return nil, nil, fmt.Errorf("invalid wire format: %s", line)
		}
		wires = append(wires, []string{fields[0], fields[1], fields[2], fields[4]})
	}

	return gates, wires, nil
}

func load(gates map[string]int, a, b string) (int, int, bool) {
	v1, ok1 := gates[a]
	v2, ok2 := gates[b]
	return v1, v2, ok1 && ok2
}

func apply(op string, a, b int) int {
	switch op {
	case "AND":
		return a & b
	case "OR":
		return a | b
	case "XOR":
		return a ^ b
	default:
		panic("invalid operation")
	}
}

func run(gates map[string]int, wires [][]string) {
	remaining := wires

	for len(remaining) > 0 {
		var nextMissing [][]string

		for _, wire := range remaining {
			a, op, b, dst := wire[0], wire[1], wire[2], wire[3]
			v1, v2, ok := load(gates, a, b)
			if !ok {
				nextMissing = append(nextMissing, wire)
				continue
			}
			gates[dst] = apply(op, v1, v2)
		}

		remaining = nextMissing
	}
}

func output(gates map[string]int) int {
	zValues := xmaps.Filter(func(k string, v int) bool {
		return strings.HasPrefix(k, "z")
	}, gates)
	out := make([]int, len(zValues))
	for _, pair := range zValues {
		idx := aoc.Must(convert.ScanNumber[int]([]byte(pair.K[1:])))
		out[len(zValues)-idx-1] = pair.V
	}
	return convert.FromBinaryToBase10(out)
}

func day24p01(r io.Reader) (string, error) {
	gates, wires := aoc.Must2(parseMonitoringDevice(r))
	run(gates, wires)
	return strconv.Itoa(output(gates)), nil
}

func extractZWires(wires [][]string) []string {
	var zWires []string
	for _, wire := range wires {
		destination := wire[3]
		if hasPrefix(destination, 'z') {
			zWires = append(zWires, destination)
		}
	}
	return zWires
}

func identifyWrongWires(wires [][]string, highestZWire string) collections.Set[string] {
	wrongWires := collections.NewSet[string]()

	for _, wire := range wires {
		a, operation, b, destination := wire[0], wire[1], wire[2], wire[3]

		switch {
		// Rule 1: z-prefixed wires with non-XOR operations (except the highest z wire)
		case destination != highestZWire && hasPrefix(destination, 'z') && operation != "XOR":
			wrongWires.Add(destination)

		// Rule 2: XOR operations with non x/y/z-prefixed wires
		case operation == "XOR" && !hasSuffixXYZ(destination) && !hasSuffixXYZ(a) && !hasSuffixXYZ(b):
			wrongWires.Add(destination)

		// Rule 3: AND operation that doesn't involve "x00" and leads to non-OR suboperations
		case operation == "AND" && a != "x00" && b != "x00":
			w := reviewConnections(wires, destination, func(op string) bool { return op != "OR" })
			wrongWires.Add(w...)

		// Rule 4: XOR operation leading to an OR suboperation
		case operation == "XOR":
			w := reviewConnections(wires, destination, func(op string) bool { return op == "OR" })
			wrongWires.Add(w...)
		}
	}

	return wrongWires
}

// reviewConnections returns a destination
// if it is involved in sub-operations of a specific type.
func reviewConnections(wires [][]string, destination string, failCondition func(string) bool) []string {
	var wrong []string
	for _, wire := range wires {
		a, op, b := wire[0], wire[1], wire[2]
		if failCondition(op) && (destination == a || destination == b) {
			wrong = append(wrong, destination)
		}
	}
	return wrong
}

// hasPrefix checks if a string has any of the given prefixes
func hasPrefix(s string, prefixes ...rune) bool {
	for _, prefix := range prefixes {
		if rune(s[0]) == prefix {
			return true
		}
	}
	return false
}

func hasSuffixXYZ(s string) bool { return hasPrefix(s, 'x', 'y', 'z') }

func day24p02(r io.Reader) (string, error) {
	_, wires := aoc.Must2(parseMonitoringDevice(r))

	zWires := extractZWires(wires)
	slices.Sort(zWires)
	highestZWire := zWires[len(zWires)-1]

	wrongWires := identifyWrongWires(wires, highestZWire)
	wrongWireList := slices.Collect(wrongWires.Iter())
	slices.Sort(wrongWireList)

	return strings.Join(wrongWireList, ","), nil
}
