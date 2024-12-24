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

func day24p02(r io.Reader) (string, error) {
	_, wires := aoc.Must2(parseMonitoringDevice(r))

	var zWires []string
	for _, wire := range wires {
		dst := wire[3]
		if strings.HasPrefix(dst, "z") {
			zWires = append(zWires, dst)
		}
	}
	slices.Sort(zWires)
	maxZWire := zWires[len(zWires)-1]

	wrong := collections.NewSet[string]()
	for _, wire := range wires {
		a, op, b, dst := wire[0], wire[1], wire[2], wire[3]
		if strings.HasPrefix(dst, "z") && op != "XOR" && dst != maxZWire {
			wrong.Add(dst)
		}
		if op == "XOR" &&
			!strings.HasPrefix(dst, "x") && !strings.HasPrefix(dst, "y") && !strings.HasPrefix(dst, "z") &&
			!strings.HasPrefix(a, "x") && !strings.HasPrefix(a, "y") && !strings.HasPrefix(a, "z") &&
			!strings.HasPrefix(b, "x") && !strings.HasPrefix(b, "y") && !strings.HasPrefix(b, "z") {
			wrong.Add(dst)
		}
		if op == "AND" && a != "x00" && b != "x00" {
			for _, subop := range wires {
				if (dst == subop[0] || dst == subop[2]) && subop[1] != "OR" {
					wrong.Add(dst)
				}
			}
		}
		if op == "XOR" {
			for _, subop := range wires {
				if (dst == subop[0] || dst == subop[2]) && subop[1] == "OR" {
					wrong.Add(dst)
				}
			}
		}
	}

	result := slices.Collect(wrong.Iter())
	slices.Sort(result)

	return strings.Join(result, ","), nil
}
