package aoc2024

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/matrix"
)

type clawVector struct {
	x, y int64
}

type clawMachine struct {
	a, b, prize clawVector
}

type clawButton uint8

const (
	clawButtonA clawButton = iota
	clawButtonB
)

// Eighteen decimal digits fit int64; longer fields use the full numeric parser.
func scanClawDecimalPrefix(input []byte) (int64, int, bool) {
	var value int64
	for i, character := range input {
		digit := character - '0'
		if digit > 9 {
			return value, i, i > 0
		}
		if i == 18 {
			return 0, 0, false
		}
		value = value*10 + int64(digit)
	}
	return value, len(input), len(input) > 0
}

func invalidClawCoordinateLine(line []byte) error {
	return fmt.Errorf("invalid claw coordinate line %q", line)
}

func parseClawVectorBody(line, body []byte, separator [4]byte) (clawVector, error) {
	if x, next, ok := scanClawDecimalPrefix(body); ok && bytes.HasPrefix(body[next:], separator[:]) {
		yText := body[next+len(separator):]
		if y, consumed, ok := scanClawDecimalPrefix(yText); ok && consumed == len(yText) {
			return clawVector{x: x, y: y}, nil
		}
	}
	// Restart over the whole body to preserve syntax and range error context.
	xText, yText, ok := bytes.Cut(body, separator[:])
	if !ok {
		return clawVector{}, fmt.Errorf("invalid claw coordinates %q", line)
	}
	x, err := convert.ScanNumber[int64](xText)
	if err != nil {
		return clawVector{}, fmt.Errorf("claw x coordinate: %w", err)
	}
	y, err := convert.ScanNumber[int64](yText)
	if err != nil {
		return clawVector{}, fmt.Errorf("claw y coordinate: %w", err)
	}
	if x < 0 || y < 0 {
		return clawVector{}, fmt.Errorf("claw coordinates must be nonnegative")
	}
	return clawVector{x: x, y: y}, nil
}

func parseClawMachine(input []byte) (clawMachine, []byte, bool, error) {
	aLine, remaining, hasA := bytes.Cut(input, []byte("\n"))
	bLine, remaining, hasB := bytes.Cut(remaining, []byte("\n"))
	prizeLine, rest, more := bytes.Cut(remaining, []byte("\n"))
	// Blank B or prize lines end a frame before all three required lines.
	if !hasA || !hasB || len(bLine) == 0 || len(prizeLine) == 0 ||
		more && (len(rest) == 0 || rest[0] != '\n') {
		return clawMachine{}, nil, false, fmt.Errorf("claw machine needs two buttons and a prize")
	}
	if more {
		rest = rest[1:]
	}
	body, ok := bytes.CutPrefix(aLine, []byte("Button A: X+"))
	if !ok {
		return clawMachine{}, nil, false, invalidClawCoordinateLine(aLine)
	}
	a, err := parseClawVectorBody(aLine, body, [4]byte{',', ' ', 'Y', '+'})
	if err != nil {
		return clawMachine{}, nil, false, err
	}
	body, ok = bytes.CutPrefix(bLine, []byte("Button B: X+"))
	if !ok {
		return clawMachine{}, nil, false, invalidClawCoordinateLine(bLine)
	}
	b, err := parseClawVectorBody(bLine, body, [4]byte{',', ' ', 'Y', '+'})
	if err != nil {
		return clawMachine{}, nil, false, err
	}
	body, ok = bytes.CutPrefix(prizeLine, []byte("Prize: X="))
	if !ok {
		return clawMachine{}, nil, false, invalidClawCoordinateLine(prizeLine)
	}
	prize, err := parseClawVectorBody(prizeLine, body, [4]byte{',', ' ', 'Y', '='})
	if err != nil {
		return clawMachine{}, nil, false, err
	}
	return clawMachine{a: a, b: b, prize: prize}, rest, more, nil
}

func ceilClawRatio(numerator, denominator *big.Int) *big.Int {
	result := new(big.Int).Neg(numerator)
	result.Div(result, denominator)
	return result.Neg(result)
}

// dependentClawPresses requires consistent equations with at least one free variable.
func dependentClawPresses(m clawMachine, targetX, targetY, limit *big.Int) (*big.Int, *big.Int) {
	x, y, target := big.NewInt(m.a.x), big.NewInt(m.b.x), targetX
	if m.a.x == 0 && m.b.x == 0 {
		x.SetInt64(m.a.y)
		y.SetInt64(m.b.y)
		target = targetY
	}
	var remainder big.Int
	if x.Sign() == 0 && y.Sign() == 0 {
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

// minimumClawTokens returns zero when no admissible integer solution exists.
// The offset and optional press limit are integers prepared by clawTokenTotal.
func minimumClawTokens(m clawMachine, offset matrix.Rat, limit *matrix.Rat) (matrix.Rat, error) {
	buttonA, buttonB := matrix.Variable(clawButtonA), matrix.Variable(clawButtonB)
	x := buttonA.Scale(m.a.x).Add(buttonB.Scale(m.b.x))
	y := buttonA.Scale(m.a.y).Add(buttonB.Scale(m.b.y))
	targetX, targetY := matrix.NewRat(m.prize.x, 1), matrix.NewRat(m.prize.y, 1)
	if !offset.IsZero() {
		targetX = targetX.Add(offset)
		targetY = targetY.Add(offset)
	}
	solution, err := matrix.NewSystem(buttonA, buttonB).
		Add(x.EqualRat(targetX)).
		Add(y.EqualRat(targetY)).
		Solve()
	if err != nil {
		if errors.Is(err, matrix.ErrInconsistent) {
			return matrix.Rat{}, nil
		}
		return matrix.Rat{}, fmt.Errorf("solve claw equations: %w", err)
	}
	var a, b matrix.Rat
	if len(solution.Free()) != 0 {
		var bound *big.Int
		if limit != nil {
			bound = limit.Numerator()
		}
		pressA, pressB := dependentClawPresses(m, targetX.Numerator(), targetY.Numerator(), bound)
		if pressA == nil {
			return matrix.Rat{}, nil
		}
		a, b = matrix.NewInteger(pressA), matrix.NewInteger(pressB)
	} else {
		expressionA, err := solution.Expression(buttonA)
		if err != nil {
			return matrix.Rat{}, fmt.Errorf("extract claw A solution: %w", err)
		}
		expressionB, err := solution.Expression(buttonB)
		if err != nil {
			return matrix.Rat{}, fmt.Errorf("extract claw B solution: %w", err)
		}
		a, b = expressionA.Constant(), expressionB.Constant()
		if !a.IsInt() || !b.IsInt() {
			return matrix.Rat{}, nil
		}
	}
	if a.Sign() < 0 || b.Sign() < 0 {
		return matrix.Rat{}, nil
	}
	if limit != nil {
		if a.Cmp(*limit) > 0 || b.Cmp(*limit) > 0 {
			return matrix.Rat{}, nil
		}
	}
	return a.Mul(matrix.NewRat(3, 1)).Add(b), nil
}

func clawTokenTotal(r io.Reader, offset int64, limit *big.Int) (string, error) {
	const inputLimit = 4 << 20
	var buffer bytes.Buffer
	buffer.Grow(32 << 10)
	if _, err := buffer.ReadFrom(io.LimitReader(r, inputLimit+1)); err != nil {
		return "", err
	}
	if buffer.Len() > inputLimit {
		return "", fmt.Errorf("claw input exceeds 4 MiB")
	}
	data := buffer.Bytes()
	if first := bytes.IndexByte(data, '\r'); first >= 0 {
		write := first
		for read := first; read < len(data); read++ {
			if data[read] == '\r' && read+1 < len(data) && data[read+1] == '\n' {
				continue
			}
			data[write] = data[read]
			write++
		}
		data = data[:write]
	}
	input := bytes.TrimSpace(data)
	offsetRat := matrix.NewRat(offset, 1)
	var limitRat *matrix.Rat
	if limit != nil {
		bound := matrix.NewInteger(limit)
		limitRat = &bound
	}
	var total matrix.Rat
	for ordinal := 1; ; ordinal++ {
		if ordinal > 10_000 {
			return "", fmt.Errorf("claw input exceeds 10000 machines")
		}
		machine, rest, more, err := parseClawMachine(input)
		if err != nil {
			return "", err
		}
		cost, err := minimumClawTokens(machine, offsetRat, limitRat)
		if err != nil {
			return "", fmt.Errorf("claw machine %d: %w", ordinal, err)
		}
		total = total.Add(cost)
		if !more {
			break
		}
		input = rest
	}
	return total.Numerator().String(), nil
}

func day13p01(r io.Reader) (string, error) {
	return clawTokenTotal(r, 0, big.NewInt(100))
}

func day13p02(r io.Reader) (string, error) {
	return clawTokenTotal(r, 10_000_000_000_000, nil)
}
