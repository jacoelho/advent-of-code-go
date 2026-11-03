package aoc2020

import (
	"fmt"
	"io"
	"math/big"
	"math/bits"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const dockingBits = uint64(1<<36 - 1)

type dockingWrite struct {
	address, value uint64
	ones, floating uint64
}

func parseDockingProgram(r io.Reader) ([]dockingWrite, error) {
	var writes []dockingWrite
	var ones, floating uint64
	haveMask := false
	s := scanner.NewScanner(r, func(line []byte) (string, error) { return string(line), nil })
	for line := range s.Values() {
		left, right, ok := strings.Cut(line, " = ")
		if !ok {
			return nil, fmt.Errorf("invalid docking instruction %q", line)
		}
		if left == "mask" {
			if len(right) != 36 {
				return nil, fmt.Errorf("docking mask must have 36 bits")
			}
			ones, floating = 0, 0
			for _, c := range right {
				ones <<= 1
				floating <<= 1
				switch c {
				case '1':
					ones |= 1
				case 'X':
					floating |= 1
				case '0':
				default:
					return nil, fmt.Errorf("invalid docking mask bit %q", c)
				}
			}
			haveMask = true
			continue
		}
		addressText, ok := strings.CutPrefix(left, "mem[")
		addressText, closed := strings.CutSuffix(addressText, "]")
		address, addressErr := strconv.ParseUint(addressText, 10, 36)
		value, valueErr := strconv.ParseUint(right, 10, 36)
		if !haveMask || !ok || !closed || addressErr != nil || valueErr != nil {
			return nil, fmt.Errorf("invalid docking write %q", line)
		}
		writes = append(writes, dockingWrite{address, value, ones, floating})
	}
	return writes, s.Err()
}

func day14p01(r io.Reader) (string, error) {
	writes, err := parseDockingProgram(r)
	if err != nil {
		return "", err
	}
	memory := make(map[uint64]uint64)
	for _, write := range writes {
		memory[write.address] = write.value&write.floating | write.ones
	}
	sum := new(big.Int)
	for _, value := range memory {
		sum.Add(sum, new(big.Int).SetUint64(value))
	}
	return sum.String(), nil
}

// A cube fixes the bits in fixed; the remaining bits describe an address range.
type dockingCube struct{ fixed, value uint64 }

func subtractDockingCube(cube, covered dockingCube) []dockingCube {
	if (cube.value^covered.value)&cube.fixed&covered.fixed != 0 {
		return []dockingCube{cube}
	}
	var remainder []dockingCube
	// Split only where the covering cube constrains a still-floating bit. Keep
	// the opposite half; continue splitting the half that overlaps the cover.
	for splitting := covered.fixed &^ cube.fixed; splitting != 0; splitting &= splitting - 1 {
		bit := splitting & -splitting
		other := dockingCube{cube.fixed | bit, cube.value | bit}
		cube.fixed |= bit
		if covered.value&bit != 0 {
			other.value &^= bit
			cube.value |= bit
		}
		remainder = append(remainder, other)
	}
	return remainder
}

func day14p02(r io.Reader) (string, error) {
	writes, err := parseDockingProgram(r)
	if err != nil {
		return "", err
	}
	var covered []dockingCube
	total := new(big.Int)
	for i := len(writes) - 1; i >= 0; i-- {
		write := writes[i]
		fixed := dockingBits &^ write.floating
		remaining := []dockingCube{{fixed, (write.address | write.ones) & fixed}}
		for _, later := range covered {
			var next []dockingCube
			for _, cube := range remaining {
				next = append(next, subtractDockingCube(cube, later)...)
			}
			remaining = next
			if len(remaining) == 0 {
				break
			}
		}
		for _, cube := range remaining {
			count := uint64(1) << (36 - bits.OnesCount64(cube.fixed))
			contribution := new(big.Int).Mul(new(big.Int).SetUint64(count), new(big.Int).SetUint64(write.value))
			total.Add(total, contribution)
		}
		covered = append(covered, remaining...)
	}
	return total.String(), nil
}
