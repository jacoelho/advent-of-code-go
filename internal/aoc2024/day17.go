package aoc2024

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type chronospatialComputer struct {
	a, b, c uint64
	program []uint8
}

func parseComputerProgram(r io.Reader) (chronospatialComputer, error) {
	const inputLimit = 64 << 10
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return chronospatialComputer{}, err
	}
	if len(data) > inputLimit {
		return chronospatialComputer{}, fmt.Errorf("computer input exceeds 64 KiB")
	}
	var computer chronospatialComputer
	var registers [3]bool
	haveProgram := false
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		label, value, ok := strings.Cut(line, ":")
		if !ok {
			return chronospatialComputer{}, fmt.Errorf("computer field needs a label and value")
		}
		value = strings.TrimSpace(value)
		switch label {
		case "Register A", "Register B", "Register C":
			index := label[len(label)-1] - 'A'
			if registers[index] {
				return chronospatialComputer{}, fmt.Errorf("duplicate %s", label)
			}
			number, err := convert.ScanNumber[int64]([]byte(value))
			if err != nil {
				return chronospatialComputer{}, fmt.Errorf("%s: %w", label, err)
			}
			if number < 0 {
				return chronospatialComputer{}, fmt.Errorf("%s must be nonnegative", label)
			}
			registers[index] = true
			switch index {
			case 0:
				computer.a = uint64(number)
			case 1:
				computer.b = uint64(number)
			case 2:
				computer.c = uint64(number)
			}
		case "Program":
			if haveProgram {
				return chronospatialComputer{}, fmt.Errorf("duplicate computer program")
			}
			haveProgram = true
			for word := range strings.SplitSeq(value, ",") {
				if len(computer.program) == 1024 {
					return chronospatialComputer{}, fmt.Errorf("computer program exceeds 1024 words")
				}
				number, err := convert.ScanNumber[int]([]byte(strings.TrimSpace(word)))
				if err != nil {
					return chronospatialComputer{}, fmt.Errorf("program word: %w", err)
				}
				if number < 0 || number > 7 {
					return chronospatialComputer{}, fmt.Errorf("program words must be between 0 and 7")
				}
				computer.program = append(computer.program, uint8(number))
			}
		default:
			return chronospatialComputer{}, fmt.Errorf("unknown computer field %q", label)
		}
	}
	if registers != [3]bool{true, true, true} || !haveProgram {
		return chronospatialComputer{}, fmt.Errorf("computer needs three registers and a program")
	}
	return computer, nil
}

type computerWork struct {
	steps int
}

func (c chronospatialComputer) run(work *computerWork) ([]uint8, error) {
	var output []uint8
	combo := func(operand uint8) (uint64, error) {
		switch operand {
		case 0, 1, 2, 3:
			return uint64(operand), nil
		case 4:
			return c.a, nil
		case 5:
			return c.b, nil
		case 6:
			return c.c, nil
		default:
			return 0, fmt.Errorf("reserved combo operand %d", operand)
		}
	}
	for ip := 0; ip < len(c.program); {
		if work.steps == 0 {
			return nil, fmt.Errorf("computer instruction budget exhausted")
		}
		work.steps--
		if ip+1 >= len(c.program) {
			return nil, fmt.Errorf("instruction at byte %d has no operand", ip)
		}
		opcode, operand := c.program[ip], c.program[ip+1]
		var value uint64
		if opcode == 0 || opcode == 2 || opcode >= 5 {
			var err error
			value, err = combo(operand)
			if err != nil {
				return nil, fmt.Errorf("instruction at byte %d: %w", ip, err)
			}
		}
		switch opcode {
		case 0:
			c.a >>= value
		case 1:
			c.b ^= uint64(operand)
		case 2:
			c.b = value % 8
		case 3:
			if c.a != 0 {
				ip = int(operand)
				continue
			}
		case 4:
			c.b ^= c.c
		case 5:
			output = append(output, uint8(value%8))
		case 6:
			c.b = c.a >> value
		case 7:
			c.c = c.a >> value
		}
		ip += 2
	}
	return output, nil
}

func validateOctalLoop(program []uint8) error {
	unsupported := func() error {
		return fmt.Errorf("unsupported quine program: need a local-register loop with one output and one divide by eight")
	}
	if len(program) < 6 || len(program)%2 != 0 || program[len(program)-2] != 3 || program[len(program)-1] != 0 {
		return unsupported()
	}
	definedB, definedC := false, false
	divisions, outputs := 0, 0
	for ip := 0; ip < len(program)-2; ip += 2 {
		opcode, operand := program[ip], program[ip+1]
		readB, readC := opcode == 1 || opcode == 4, opcode == 4
		if opcode == 0 || opcode == 2 || opcode >= 5 {
			if operand == 7 {
				return unsupported()
			}
			readB = readB || operand == 5
			readC = readC || operand == 6
		}
		if readB && !definedB || readC && !definedC {
			return unsupported()
		}
		switch opcode {
		case 0:
			if operand != 3 {
				return unsupported()
			}
			divisions++
		case 1, 2, 4, 6:
			definedB = true
		case 3:
			return unsupported()
		case 5:
			outputs++
		case 7:
			definedC = true
		}
	}
	if divisions != 1 || outputs != 1 {
		return unsupported()
	}
	return nil
}

func findComputerQuine(computer chronospatialComputer) (uint64, error) {
	if err := validateOctalLoop(computer.program); err != nil {
		return 0, err
	}
	if len(computer.program) > 21 {
		return 0, fmt.Errorf("quine needs more than 21 octal digits; register range is limited to int64")
	}
	work := computerWork{steps: 100_000_000}
	candidatesLeft := 1_000_000
	candidates := []uint64{0}
	// Each admitted iteration consumes three bits and ignores the previous B and C.
	for length := 1; length <= len(computer.program); length++ {
		var next []uint64
		overflow := false
		for _, prefix := range candidates {
			for digit := uint64(0); digit < 8; digit++ {
				if candidatesLeft == 0 {
					return 0, fmt.Errorf("quine search exceeds one million candidates")
				}
				candidatesLeft--
				if prefix > (math.MaxInt64-digit)/8 {
					overflow = true
					continue
				}
				candidate := prefix*8 + digit
				trial := computer
				trial.a = candidate
				output, err := trial.run(&work)
				if err != nil {
					return 0, fmt.Errorf("quine candidate %d: %w", candidate, err)
				}
				if slices.Equal(output, computer.program[len(computer.program)-length:]) {
					next = append(next, candidate)
				}
			}
		}
		if len(next) == 0 {
			if overflow {
				return 0, fmt.Errorf("quine candidates exceed the int64 register range")
			}
			return 0, fmt.Errorf("computer program has no quine in the supported register range")
		}
		candidates = next
	}
	best := uint64(math.MaxInt64)
	found := false
	for _, candidate := range candidates {
		if candidate > 0 && (!found || candidate < best) {
			best, found = candidate, true
		}
	}
	if !found {
		return 0, fmt.Errorf("computer program has no positive quine")
	}
	return best, nil
}

func day17p01(r io.Reader) (string, error) {
	computer, err := parseComputerProgram(r)
	if err != nil {
		return "", err
	}
	work := computerWork{steps: 1_000_000}
	output, err := computer.run(&work)
	if err != nil {
		return "", err
	}
	words := make([]string, len(output))
	for i, value := range output {
		words[i] = strconv.Itoa(int(value))
	}
	return strings.Join(words, ","), nil
}

func day17p02(r io.Reader) (string, error) {
	computer, err := parseComputerProgram(r)
	if err != nil {
		return "", err
	}
	value, err := findComputerQuine(computer)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(value, 10), nil
}
