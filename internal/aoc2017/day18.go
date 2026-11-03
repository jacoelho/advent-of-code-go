package aoc2017

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

type duetMode uint8

const (
	duetSounds duetMode = iota
	duetMessages
)

const duetInstructionLimit = 50_000_000

func parseDuet(r io.Reader, mode duetMode) ([]assemblyInstruction, error) {
	program, err := parseAssembly(r)
	if err != nil {
		return nil, err
	}

	for _, instruction := range program {
		switch instruction.operation {
		case "snd", "set", "add", "mul", "mod", "jgz":
		case "rcv":
			if mode == duetMessages && instruction.x.register < 0 {
				return nil, fmt.Errorf("duet receive destination must be a register")
			}
		default:
			return nil, fmt.Errorf("unsupported duet operation %q", instruction.operation)
		}
	}

	return program, nil
}

func day18p01(r io.Reader) (string, error) {
	program, err := parseDuet(r, duetSounds)
	if err != nil {
		return "", err
	}

	var registers [26]int64
	var pointer, lastSound int64
	hasSound := false
	for range duetInstructionLimit {
		if pointer < 0 || pointer >= int64(len(program)) {
			return "", fmt.Errorf("duet terminated without recovering a sound")
		}
		instruction := program[pointer]
		switch instruction.operation {
		case "snd":
			lastSound, hasSound = instruction.x.value(registers), true
			pointer++

		case "rcv":
			if instruction.x.value(registers) != 0 {
				if !hasSound {
					return "", fmt.Errorf("duet received before any sound")
				}
				return strconv.FormatInt(lastSound, 10), nil
			}
			pointer++

		default:
			pointer, err = advanceAssembly(instruction, &registers, pointer)
			if err != nil {
				return "", err
			}
		}
	}

	return "", fmt.Errorf("duet exceeds %d instructions", duetInstructionLimit)
}

type duetProcess struct {
	registers [26]int64
	pointer   int64
	inbox     *collections.Deque[int64]
	sent      int
}

func stepDuet(current, other *duetProcess, program []assemblyInstruction) (bool, error) {
	if current.pointer < 0 || current.pointer >= int64(len(program)) {
		return false, nil
	}

	instruction := program[current.pointer]
	switch instruction.operation {
	case "snd":
		if other.inbox.Size() == 1_000_000 {
			return false, fmt.Errorf("duet message queue exceeds 1000000 values")
		}
		other.inbox.PushBack(instruction.x.value(current.registers))
		current.sent++
		current.pointer++

	case "rcv":
		value, available := current.inbox.PopFront()
		if !available {
			return false, nil
		}
		current.registers[instruction.x.register] = value
		current.pointer++

	default:
		next, err := advanceAssembly(instruction, &current.registers, current.pointer)
		if err != nil {
			return false, err
		}
		current.pointer = next
	}

	return true, nil
}

func day18p02(r io.Reader) (string, error) {
	program, err := parseDuet(r, duetMessages)
	if err != nil {
		return "", err
	}

	processes := [2]duetProcess{
		{inbox: collections.NewDeque[int64](64)},
		{inbox: collections.NewDeque[int64](64)},
	}
	processes[1].registers['p'-'a'] = 1

	steps := 0
	for {
		progress := false
		for id := range processes {
			advanced, err := stepDuet(&processes[id], &processes[1-id], program)
			if err != nil {
				return "", err
			}
			if advanced {
				progress = true
				steps++
				if steps == duetInstructionLimit {
					return "", fmt.Errorf("duet exceeds %d instructions", duetInstructionLimit)
				}
			}
		}
		if !progress {
			return strconv.Itoa(processes[1].sent), nil
		}
	}
}
