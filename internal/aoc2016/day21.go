package aoc2016

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type scrambleOperation int

const (
	swapPasswordPositions scrambleOperation = iota
	swapPasswordLetters
	rotatePasswordSteps
	rotatePasswordLetter
	reversePasswordSection
	movePasswordLetter
)

type scrambleInstruction struct {
	operation scrambleOperation
	x, y      int
	a, b      byte
}

func parseScramble(r io.Reader) ([]scrambleInstruction, error) {
	var instructions []scrambleInstruction
	s := bufio.NewScanner(r)
	for s.Scan() {
		f := strings.Fields(s.Text())
		var instruction scrambleInstruction
		var x, y string
		switch {
		case len(f) == 6 && f[0] == "swap" && f[1] == "position" && f[3] == "with" && f[4] == "position":
			instruction.operation, x, y = swapPasswordPositions, f[2], f[5]
		case len(f) == 6 && f[0] == "swap" && f[1] == "letter" && f[3] == "with" && f[4] == "letter" && len(f[2]) == 1 && len(f[5]) == 1:
			instruction.operation, instruction.a, instruction.b = swapPasswordLetters, f[2][0], f[5][0]
		case len(f) == 4 && f[0] == "rotate" && (f[1] == "left" || f[1] == "right") && (f[3] == "step" || f[3] == "steps"):
			instruction.operation, x = rotatePasswordSteps, f[2]
			if f[1] == "left" {
				instruction.y = -1
			} else {
				instruction.y = 1
			}
		case len(f) == 7 && strings.Join(f[:6], " ") == "rotate based on position of letter" && len(f[6]) == 1:
			instruction.operation, instruction.a = rotatePasswordLetter, f[6][0]
		case len(f) == 5 && f[0] == "reverse" && f[1] == "positions" && f[3] == "through":
			instruction.operation, x, y = reversePasswordSection, f[2], f[4]
		case len(f) == 6 && f[0] == "move" && f[1] == "position" && f[3] == "to" && f[4] == "position":
			instruction.operation, x, y = movePasswordLetter, f[2], f[5]
		default:
			return nil, fmt.Errorf("invalid scramble operation %q", s.Text())
		}
		for _, field := range []struct {
			text  string
			value *int
		}{{x, &instruction.x}, {y, &instruction.y}} {
			if field.text == "" {
				continue
			}
			value, err := convert.ScanNumber[int]([]byte(field.text))
			if err != nil {
				return nil, fmt.Errorf("scramble position: %w", err)
			}
			if value < 0 {
				return nil, fmt.Errorf("negative scramble position")
			}
			*field.value = value
		}
		instructions = append(instructions, instruction)
	}
	return instructions, s.Err()
}

func rotatePassword(password string, right int) string {
	right %= len(password)
	if right < 0 {
		right += len(password)
	}
	return password[len(password)-right:] + password[:len(password)-right]
}

func applyScramble(password string, instruction scrambleInstruction) (string, error) {
	result := []byte(password)
	x, y := instruction.x, instruction.y
	switch instruction.operation {
	case swapPasswordPositions, reversePasswordSection, movePasswordLetter:
		if x >= len(password) || y >= len(password) {
			return "", fmt.Errorf("position outside password")
		}
	}
	switch instruction.operation {
	case swapPasswordPositions:
		result[x], result[y] = result[y], result[x]
	case swapPasswordLetters:
		x, y = strings.IndexByte(password, instruction.a), strings.IndexByte(password, instruction.b)
		if x < 0 || y < 0 {
			return "", fmt.Errorf("swap letter missing from password")
		}
		result[x], result[y] = result[y], result[x]
	case rotatePasswordSteps:
		return rotatePassword(password, (x%len(password))*y), nil
	case rotatePasswordLetter:
		index := strings.IndexByte(password, instruction.a)
		if index < 0 {
			return "", fmt.Errorf("rotation letter missing from password")
		}
		steps := index + 1
		if index >= 4 {
			steps++
		}
		return rotatePassword(password, steps), nil
	case reversePasswordSection:
		if x > y {
			return "", fmt.Errorf("reversed password interval")
		}
		slices.Reverse(result[x : y+1])
	case movePasswordLetter:
		letter := result[x]
		if x < y {
			copy(result[x:y], result[x+1:y+1])
		} else {
			copy(result[y+1:x+1], result[y:x])
		}
		result[y] = letter
	}
	return string(result), nil
}

func scramblePassword(instructions []scrambleInstruction, password string) (string, error) {
	for _, instruction := range instructions {
		var err error
		password, err = applyScramble(password, instruction)
		if err != nil {
			return "", err
		}
	}
	return password, nil
}

func unscramblePassword(instructions []scrambleInstruction, password string) (string, error) {
	candidates := collections.NewSet(password)
	for i := len(instructions) - 1; i >= 0; i-- {
		instruction := instructions[i]
		next := collections.NewSet[string]()
		for candidate := range candidates {
			if instruction.operation == rotatePasswordLetter {
				for rotation := range len(candidate) {
					previous := rotatePassword(candidate, rotation)
					result, err := applyScramble(previous, instruction)
					if err != nil {
						return "", err
					}
					if result == candidate {
						next.Add(previous)
					}
				}
				continue
			}
			inverse := instruction
			switch inverse.operation {
			case rotatePasswordSteps:
				inverse.y = -inverse.y
			case movePasswordLetter:
				inverse.x, inverse.y = inverse.y, inverse.x
			}
			previous, err := applyScramble(candidate, inverse)
			if err != nil {
				return "", err
			}
			next.Add(previous)
		}
		candidates = next
	}
	if candidates.Len() != 1 {
		return "", fmt.Errorf("password has %d possible predecessors", candidates.Len())
	}
	result, _ := candidates.Next()
	return result, nil
}

func day21p01(r io.Reader) (string, error) {
	instructions, err := parseScramble(r)
	if err != nil {
		return "", err
	}
	return scramblePassword(instructions, "abcdefgh")
}

func day21p02(r io.Reader) (string, error) {
	instructions, err := parseScramble(r)
	if err != nil {
		return "", err
	}
	return unscramblePassword(instructions, "fbgdceah")
}
