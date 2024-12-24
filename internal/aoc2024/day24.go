package aoc2024

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
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

			switch op {
			case "AND":
				gates[dst] = v1 & v2
			case "OR":
				gates[dst] = v1 | v2
			case "XOR":
				gates[dst] = v1 ^ v2
			}
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
	return "", nil
}
