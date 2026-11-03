package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type circuitGate struct {
	operation   string
	left, right string
}

func parseCircuit(r io.Reader) (map[string]circuitGate, error) {
	circuit := make(map[string]circuitGate)
	s := bufio.NewScanner(r)
	for s.Scan() {
		expression, wire, ok := strings.Cut(s.Text(), " -> ")
		if !ok || wire == "" {
			return nil, fmt.Errorf("invalid circuit instruction %q", s.Text())
		}
		fields := strings.Fields(expression)
		var gate circuitGate
		switch {
		case len(fields) == 1:
			gate = circuitGate{operation: "ASSIGN", left: fields[0]}
		case len(fields) == 2 && fields[0] == "NOT":
			gate = circuitGate{operation: "NOT", left: fields[1]}
		case len(fields) == 3:
			switch fields[1] {
			case "AND", "OR", "LSHIFT", "RSHIFT":
				gate = circuitGate{operation: fields[1], left: fields[0], right: fields[2]}
			default:
				return nil, fmt.Errorf("unknown circuit operation %q", fields[1])
			}
		default:
			return nil, fmt.Errorf("invalid circuit expression %q", expression)
		}
		if _, exists := circuit[wire]; exists {
			return nil, fmt.Errorf("duplicate signal for wire %q", wire)
		}
		circuit[wire] = gate
	}
	return circuit, s.Err()
}

func circuitSignal(circuit map[string]circuitGate, wire string, signals map[string]uint16) (uint16, error) {
	visiting := make(map[string]bool)
	var resolve func(string) (uint16, error)
	resolve = func(input string) (uint16, error) {
		if input != "" && input[0] >= '0' && input[0] <= '9' {
			value, err := strconv.ParseUint(input, 10, 16)
			if err != nil {
				return 0, fmt.Errorf("invalid circuit signal %q: %w", input, err)
			}
			return uint16(value), nil
		}
		if signal, ok := signals[input]; ok {
			return signal, nil
		}
		gate, ok := circuit[input]
		if !ok {
			return 0, fmt.Errorf("undefined wire %q", input)
		}
		if visiting[input] {
			return 0, fmt.Errorf("cycle at wire %q", input)
		}
		visiting[input] = true
		left, err := resolve(gate.left)
		if err != nil {
			return 0, err
		}
		var right uint16
		if gate.right != "" {
			right, err = resolve(gate.right)
			if err != nil {
				return 0, err
			}
		}
		var signal uint16
		switch gate.operation {
		case "ASSIGN":
			signal = left
		case "NOT":
			signal = ^left
		case "AND":
			signal = left & right
		case "OR":
			signal = left | right
		case "LSHIFT":
			signal = left << right
		case "RSHIFT":
			signal = left >> right
		}
		delete(visiting, input)
		signals[input] = signal
		return signal, nil
	}
	return resolve(wire)
}

func day07p01(r io.Reader) (string, error) {
	circuit, err := parseCircuit(r)
	if err != nil {
		return "", err
	}
	signal, err := circuitSignal(circuit, "a", make(map[string]uint16))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(signal)), nil
}

func day07p02(r io.Reader) (string, error) {
	circuit, err := parseCircuit(r)
	if err != nil {
		return "", err
	}
	first, err := circuitSignal(circuit, "a", make(map[string]uint16))
	if err != nil {
		return "", err
	}
	signal, err := circuitSignal(circuit, "a", map[string]uint16{"b": first})
	if err != nil {
		return "", err
	}
	return strconv.Itoa(int(signal)), nil
}
