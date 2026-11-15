package aoc2024

import (
	"fmt"
	"io"
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
	data, err := io.ReadAll(r)
	if err != nil {
		return chronospatialComputer{}, err
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

type day17State struct {
	ip      int
	a, b, c uint64
}

func (c chronospatialComputer) run() ([]uint8, error) {
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
	// A never grows, so the program runs forever exactly when it revisits a
	// state without A changing; states are kept only for the current A.
	var haveJumped bool
	var jumpA uint64
	var firstJump day17State
	var jumped map[day17State]bool
	for ip := 0; ip < len(c.program); {
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
				state := day17State{ip: ip, a: c.a, b: c.b, c: c.c}
				switch {
				case !haveJumped || jumpA != c.a:
					haveJumped, jumpA, firstJump = true, c.a, state
					clear(jumped)
				case state == firstJump || jumped[state]:
					return nil, fmt.Errorf("program never halts: jump at byte %d repeats a state", ip)
				default:
					if jumped == nil {
						jumped = make(map[day17State]bool)
					}
					jumped[state] = true
				}
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

// validateOctalLoop checks the shape the digit-by-digit quine search relies
// on: one loop closed by a final "jnz 0" whose body shifts A right by three
// once, outputs once, and writes B and C before reading them. Each output then
// depends only on A at the start of its iteration, so dropping the last octal
// digit of a quine's A drops the first value of its output.
func validateOctalLoop(program []uint8) error {
	n := len(program)
	if n < 2 || n%2 != 0 || program[n-2] != 3 || program[n-1] != 0 {
		return fmt.Errorf("unsupported quine program: must end with jnz 0")
	}
	var shifts, outputs int
	var wroteB, wroteC bool
	for ip := 0; ip < n-2; ip += 2 {
		opcode, operand := program[ip], program[ip+1]
		combo := opcode == 0 || opcode == 2 || opcode >= 5
		readsB := opcode == 1 || opcode == 4 || combo && operand == 5
		readsC := opcode == 4 || combo && operand == 6
		if readsB && !wroteB || readsC && !wroteC {
			return fmt.Errorf("unsupported quine program: instruction at byte %d reads B or C before the loop writes it", ip)
		}
		switch opcode {
		case 0:
			if operand != 3 {
				return fmt.Errorf("unsupported quine program: A must shift by exactly 3")
			}
			shifts++
		case 1, 2, 4, 6:
			wroteB = true
		case 3:
			return fmt.Errorf("unsupported quine program: jump at byte %d inside the loop", ip)
		case 5:
			outputs++
		case 7:
			wroteC = true
		}
	}
	if shifts != 1 || outputs != 1 {
		return fmt.Errorf("unsupported quine program: loop needs one shift of A and one output, has %d and %d", shifts, outputs)
	}
	return nil
}

func findComputerQuine(computer chronospatialComputer) (uint64, error) {
	if err := validateOctalLoop(computer.program); err != nil {
		return 0, err
	}
	// A quine's A has one octal digit per program word.
	if len(computer.program) > 21 {
		return 0, fmt.Errorf("quine needs %d octal digits; registers hold at most 21", len(computer.program))
	}
	candidates := []uint64{0}
	for length := 1; length <= len(computer.program); length++ {
		var next []uint64
		for _, prefix := range candidates {
			for digit := range uint64(8) {
				trial := computer
				trial.a = prefix*8 + digit
				output, err := trial.run()
				if err != nil {
					return 0, fmt.Errorf("quine candidate %d: %w", trial.a, err)
				}
				if slices.Equal(output, computer.program[len(computer.program)-length:]) {
					next = append(next, trial.a)
				}
			}
		}
		if len(next) == 0 {
			return 0, fmt.Errorf("computer program has no quine")
		}
		candidates = next
	}
	// Candidates are in ascending order.
	for _, candidate := range candidates {
		if candidate > 0 {
			return candidate, nil
		}
	}
	return 0, fmt.Errorf("computer program has no positive quine")
}

func day17p01(r io.Reader) (string, error) {
	computer, err := parseComputerProgram(r)
	if err != nil {
		return "", err
	}
	output, err := computer.run()
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
