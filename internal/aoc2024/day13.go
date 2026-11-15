package aoc2024

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type clawVector struct {
	x, y int64
}

type clawMachine struct {
	a, b, prize clawVector
}

var errDay13Overflow = errors.New("claw arithmetic overflows int64")

func day13ParseCoordinate(text []byte) (int64, error) {
	if len(text) == 0 || text[0] < '0' || text[0] > '9' {
		return 0, fmt.Errorf("%q is not a nonnegative decimal", text)
	}
	return convert.ScanNumber[int64](text)
}

func day13ParseVector(line []byte, prefix, separator string) (clawVector, error) {
	body, ok := bytes.CutPrefix(line, []byte(prefix))
	if !ok {
		return clawVector{}, fmt.Errorf("%q does not start with %q", line, prefix)
	}
	xText, yText, ok := bytes.Cut(body, []byte(separator))
	if !ok {
		return clawVector{}, fmt.Errorf("%q has no %q separator", line, separator)
	}
	x, err := day13ParseCoordinate(xText)
	if err != nil {
		return clawVector{}, fmt.Errorf("claw x coordinate: %w", err)
	}
	y, err := day13ParseCoordinate(yText)
	if err != nil {
		return clawVector{}, fmt.Errorf("claw y coordinate: %w", err)
	}
	return clawVector{x: x, y: y}, nil
}

// day13ParseMachines reads machines of three lines each, separated by one
// blank line.
func day13ParseMachines(r io.Reader) ([]clawMachine, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return nil, fmt.Errorf("claw input needs machines of two buttons and a prize separated by blank lines")
	}
	machines := make([]clawMachine, 0, len(blocks))
	start := 1
	for _, block := range blocks {
		if len(block) != 3 {
			return nil, fmt.Errorf("claw input line %d: want machines of two buttons and a prize separated by blank lines", start+min(len(block), 3))
		}
		var vectors [3]clawVector
		for i, format := range [3][2]string{
			{"Button A: X+", ", Y+"},
			{"Button B: X+", ", Y+"},
			{"Prize: X=", ", Y="},
		} {
			var err error
			if vectors[i], err = day13ParseVector([]byte(block[i]), format[0], format[1]); err != nil {
				return nil, fmt.Errorf("claw input line %d: %w", start+i, err)
			}
		}
		machines = append(machines, clawMachine{a: vectors[0], b: vectors[1], prize: vectors[2]})
		start += len(block) + 1
	}
	return machines, nil
}

func day13Mul(a, b int64) (int64, error) {
	product, ok := xmath.CheckedMul(a, b)
	if !ok {
		return 0, errDay13Overflow
	}
	return product, nil
}

func day13Add(a, b int64) (int64, error) {
	sum, ok := xmath.CheckedAdd(a, b)
	if !ok {
		return 0, errDay13Overflow
	}
	return sum, nil
}

// day13Reaches reports whether a*u + b*v == target for nonnegative operands.
func day13Reaches(u, v, a, b, target int64) bool {
	ua, err := day13Mul(u, a)
	if err != nil {
		return false
	}
	vb, err := day13Mul(v, b)
	if err != nil {
		return false
	}
	sum, err := day13Add(ua, vb)
	return err == nil && sum == target
}

// day13CheapestCombination returns nonnegative a, b no greater than limit with
// u*a + v*b == t and minimal 3*a + b.
func day13CheapestCombination(u, v, t, limit int64) (int64, int64, bool) {
	switch {
	case u == 0 && v == 0:
		return 0, 0, t == 0
	case u == 0:
		b := t / v
		return 0, b, t%v == 0 && b <= limit
	case v == 0:
		a := t / u
		return a, 0, t%u == 0 && a <= limit
	}
	g := xmath.GCD(u, v)
	if t%g != 0 {
		return 0, 0, false
	}
	stepA, stepB := v/g, u/g
	// stepA and stepB are coprime, so the inverse exists.
	inverse, _ := xmath.ModInverse(stepB, stepA)
	minA := xmath.MulMod(inverse, t/g, stepA)
	uA, err := day13Mul(u, minA)
	if err != nil || uA > t || minA > limit {
		return 0, 0, false
	}
	maxB := (t - uA) / v
	// Every solution is (minA + k*stepA, maxB - k*stepB) for k in [low, high].
	low, high := int64(0), maxB/stepB
	if excess := maxB - limit; excess > 0 {
		low = excess / stepB
		if excess%stepB != 0 {
			low++
		}
	}
	high = min(high, (limit-minA)/stepA)
	if low > high {
		return 0, 0, false
	}
	k := low
	if stepA <= (stepB-1)/3 {
		k = high
	}
	return minA + k*stepA, maxB - k*stepB, true
}

// day13Presses returns the cheapest button presses reaching the prize with
// at most limit presses of each button.
func day13Presses(m clawMachine, limit int64) (int64, int64, bool, error) {
	axBy, err := day13Mul(m.a.x, m.b.y)
	if err != nil {
		return 0, 0, false, err
	}
	ayBx, err := day13Mul(m.a.y, m.b.x)
	if err != nil {
		return 0, 0, false, err
	}
	if det := axBy - ayBx; det != 0 {
		var products [4]int64
		for i, operands := range [4][2]int64{
			{m.prize.x, m.b.y}, {m.prize.y, m.b.x},
			{m.a.x, m.prize.y}, {m.a.y, m.prize.x},
		} {
			if products[i], err = day13Mul(operands[0], operands[1]); err != nil {
				return 0, 0, false, err
			}
		}
		numeratorA, numeratorB := products[0]-products[1], products[2]-products[3]
		if numeratorA%det != 0 || numeratorB%det != 0 {
			return 0, 0, false, nil
		}
		a, b := numeratorA/det, numeratorB/det
		return a, b, a >= 0 && b >= 0 && a <= limit && b <= limit, nil
	}
	// Collinear buttons: solve on an axis where a button moves, then confirm
	// the other axis, which then holds for every solution of the first.
	u, v, t := m.a.x, m.b.x, m.prize.x
	if u == 0 && v == 0 {
		u, v, t = m.a.y, m.b.y, m.prize.y
	}
	a, b, ok := day13CheapestCombination(u, v, t, limit)
	ok = ok && day13Reaches(m.a.x, m.b.x, a, b, m.prize.x) && day13Reaches(m.a.y, m.b.y, a, b, m.prize.y)
	return a, b, ok, nil
}

func clawTokenTotal(r io.Reader, offset, limit int64) (string, error) {
	machines, err := day13ParseMachines(r)
	if err != nil {
		return "", err
	}
	var total int64
	for i, m := range machines {
		cost, err := day13Tokens(m, offset, limit)
		if err != nil {
			return "", fmt.Errorf("claw machine %d: %w", i+1, err)
		}
		if total, err = day13Add(total, cost); err != nil {
			return "", fmt.Errorf("claw token total: %w", err)
		}
	}
	return strconv.FormatInt(total, 10), nil
}

func day13Tokens(m clawMachine, offset, limit int64) (int64, error) {
	var err error
	if m.prize.x, err = day13Add(m.prize.x, offset); err != nil {
		return 0, err
	}
	if m.prize.y, err = day13Add(m.prize.y, offset); err != nil {
		return 0, err
	}
	a, b, ok, err := day13Presses(m, limit)
	if err != nil || !ok {
		return 0, err
	}
	tokensA, err := day13Mul(3, a)
	if err != nil {
		return 0, err
	}
	return day13Add(tokensA, b)
}

func day13p01(r io.Reader) (string, error) {
	return clawTokenTotal(r, 0, 100)
}

func day13p02(r io.Reader) (string, error) {
	return clawTokenTotal(r, 10_000_000_000_000, math.MaxInt64)
}
