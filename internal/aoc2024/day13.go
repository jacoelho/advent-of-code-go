package aoc2024

import (
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type clawVector struct {
	x, y *big.Int
}

type clawMachine struct {
	a, b, prize clawVector
}

func (v clawVector) cross(w clawVector) *big.Int {
	return new(big.Int).Sub(new(big.Int).Mul(v.x, w.y), new(big.Int).Mul(v.y, w.x))
}

func parseClawVector(line, prefix, separator string) (clawVector, error) {
	body, ok := strings.CutPrefix(line, prefix)
	if !ok {
		return clawVector{}, fmt.Errorf("invalid claw coordinate line %q", line)
	}
	xText, yText, ok := strings.Cut(body, separator)
	if !ok {
		return clawVector{}, fmt.Errorf("invalid claw coordinates %q", line)
	}
	x, err := convert.ScanNumber[int64]([]byte(xText))
	if err != nil {
		return clawVector{}, fmt.Errorf("claw x coordinate: %w", err)
	}
	y, err := convert.ScanNumber[int64]([]byte(yText))
	if err != nil {
		return clawVector{}, fmt.Errorf("claw y coordinate: %w", err)
	}
	if x < 0 || y < 0 {
		return clawVector{}, fmt.Errorf("claw coordinates must be nonnegative")
	}
	return clawVector{x: big.NewInt(x), y: big.NewInt(y)}, nil
}

func parseClawMachines(r io.Reader) ([]clawMachine, error) {
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, err
	}
	if len(data) > inputLimit {
		return nil, fmt.Errorf("claw input exceeds 4 MiB")
	}
	input := strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n"))
	var machines []clawMachine
	for block := range strings.SplitSeq(input, "\n\n") {
		if len(machines) == 10_000 {
			return nil, fmt.Errorf("claw input exceeds 10000 machines")
		}
		lines := strings.Split(block, "\n")
		if len(lines) != 3 {
			return nil, fmt.Errorf("claw machine needs two buttons and a prize")
		}
		a, err := parseClawVector(lines[0], "Button A: X+", ", Y+")
		if err != nil {
			return nil, err
		}
		b, err := parseClawVector(lines[1], "Button B: X+", ", Y+")
		if err != nil {
			return nil, err
		}
		prize, err := parseClawVector(lines[2], "Prize: X=", ", Y=")
		if err != nil {
			return nil, err
		}
		machines = append(machines, clawMachine{a: a, b: b, prize: prize})
	}
	return machines, nil
}

func ceilClawRatio(numerator, denominator *big.Int) *big.Int {
	result := new(big.Int).Neg(numerator)
	result.Div(result, denominator)
	return result.Neg(result)
}

func dependentClawPresses(m clawMachine, limit *big.Int) (*big.Int, *big.Int) {
	if m.a.cross(m.prize).Sign() != 0 || m.b.cross(m.prize).Sign() != 0 {
		return nil, nil
	}
	x, y, target := m.a.x, m.b.x, m.prize.x
	if x.Sign() == 0 && y.Sign() == 0 {
		x, y, target = m.a.y, m.b.y, m.prize.y
	}
	var remainder big.Int
	if x.Sign() == 0 && y.Sign() == 0 {
		if m.prize.x.Sign() != 0 || m.prize.y.Sign() != 0 {
			return nil, nil
		}
		return new(big.Int), new(big.Int)
	}
	if x.Sign() == 0 {
		b := new(big.Int)
		b.QuoRem(target, y, &remainder)
		if remainder.Sign() != 0 {
			return nil, nil
		}
		return new(big.Int), b
	}
	if y.Sign() == 0 {
		a := new(big.Int)
		a.QuoRem(target, x, &remainder)
		if remainder.Sign() != 0 {
			return nil, nil
		}
		return a, new(big.Int)
	}

	var a, b big.Int
	divisor := new(big.Int).GCD(&a, &b, x, y)
	factor := new(big.Int)
	factor.QuoRem(target, divisor, &remainder)
	if remainder.Sign() != 0 {
		return nil, nil
	}
	a.Mul(&a, factor)
	b.Mul(&b, factor)
	stepA := new(big.Int).Quo(y, divisor)
	stepB := new(big.Int).Quo(x, divisor)

	// Every solution is (a + k*stepA, b - k*stepB); intersect its press bounds.
	low := ceilClawRatio(new(big.Int).Neg(&a), stepA)
	high := new(big.Int).Div(&b, stepB)
	if limit != nil {
		lower := ceilClawRatio(new(big.Int).Sub(&b, limit), stepB)
		upper := new(big.Int).Div(new(big.Int).Sub(limit, &a), stepA)
		if low.Cmp(lower) < 0 {
			low = lower
		}
		if high.Cmp(upper) > 0 {
			high = upper
		}
	}
	if low.Cmp(high) > 0 {
		return nil, nil
	}
	change := new(big.Int).Sub(new(big.Int).Mul(big.NewInt(3), stepA), stepB)
	k := low
	if change.Sign() < 0 {
		k = high
	}
	a.Add(&a, new(big.Int).Mul(k, stepA))
	b.Sub(&b, new(big.Int).Mul(k, stepB))
	return &a, &b
}

func minimumClawTokens(m clawMachine, limit *big.Int) *big.Int {
	determinant := m.a.cross(m.b)
	var a, b *big.Int
	if determinant.Sign() == 0 {
		a, b = dependentClawPresses(m, limit)
	} else {
		var remainder big.Int
		a, b = new(big.Int), new(big.Int)
		a.QuoRem(m.prize.cross(m.b), determinant, &remainder)
		if remainder.Sign() != 0 {
			return nil
		}
		b.QuoRem(m.a.cross(m.prize), determinant, &remainder)
		if remainder.Sign() != 0 {
			return nil
		}
	}
	if a == nil || a.Sign() < 0 || b.Sign() < 0 {
		return nil
	}
	if limit != nil && (a.Cmp(limit) > 0 || b.Cmp(limit) > 0) {
		return nil
	}
	return new(big.Int).Add(new(big.Int).Mul(big.NewInt(3), a), b)
}

func clawTokenTotal(r io.Reader, offset int64, limit *big.Int) (string, error) {
	machines, err := parseClawMachines(r)
	if err != nil {
		return "", err
	}
	total := new(big.Int)
	for _, machine := range machines {
		machine.prize = clawVector{
			x: new(big.Int).Add(machine.prize.x, big.NewInt(offset)),
			y: new(big.Int).Add(machine.prize.y, big.NewInt(offset)),
		}
		if cost := minimumClawTokens(machine, limit); cost != nil {
			total.Add(total, cost)
		}
	}
	return total.String(), nil
}

func day13p01(r io.Reader) (string, error) {
	return clawTokenTotal(r, 0, big.NewInt(100))
}

func day13p02(r io.Reader) (string, error) {
	return clawTokenTotal(r, 10_000_000_000_000, nil)
}
