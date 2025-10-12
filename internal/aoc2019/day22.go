package aoc2019

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xmath"
)

type shuffleOp int

const (
	dealNewStack shuffleOp = iota
	cut
	dealIncrement
)

type shuffle struct {
	op shuffleOp
	n  int
}

func parseShuffle(line []byte) (shuffle, error) {
	s := string(line)

	switch {
	case s == "deal into new stack":
		return shuffle{op: dealNewStack}, nil
	case strings.Contains(s, "cut"):
		nums := convert.ExtractDigits[int](s)
		if len(nums) == 0 {
			return shuffle{}, fmt.Errorf("no number found in cut instruction")
		}
		return shuffle{op: cut, n: nums[0]}, nil
	case strings.Contains(s, "deal with increment"):
		nums := convert.ExtractDigits[int](s)
		if len(nums) == 0 {
			return shuffle{}, fmt.Errorf("no number found in deal with increment instruction")
		}
		return shuffle{op: dealIncrement, n: nums[0]}, nil
	default:
		return shuffle{}, fmt.Errorf("unknown instruction: %s", s)
	}
}

func parseShuffleInstructions(r io.Reader) ([]shuffle, error) {
	s := scanner.NewScanner(r, parseShuffle)
	instructions := slices.Collect(s.Values())
	return instructions, s.Err()
}

func trackCardPosition(instructions []shuffle, deckSize, cardNumber int) int {
	pos := cardNumber

	for _, inst := range instructions {
		switch inst.op {
		case dealNewStack:
			pos = deckSize - 1 - pos
		case cut:
			pos = (pos - inst.n + deckSize) % deckSize
		case dealIncrement:
			pos = (pos * inst.n) % deckSize
		}
	}

	return pos
}

func day22p01(r io.Reader) (string, error) {
	instructions, err := parseShuffleInstructions(r)
	if err != nil {
		return "", err
	}

	position := trackCardPosition(instructions, 10007, 2019)
	return strconv.Itoa(position), nil
}

// linearFunc represents a linear transformation f(x) = (a*x + b) mod m in modular arithmetic.
//
// part 2 requires shuffling a deck of ~10^14 cards ~10^14 times, making simulation impossible.
// each shuffle operation is a linear transformation that can be composed into a single function.
// we then raise that function to the power of 10^14 using fast exponentiation, and invert it
// to solve the inverse problem: "which card ends up at position 2020?"
type linearFunc struct {
	a, b, mod int64
}

func mod(a, m int64) int64 {
	return ((a % m) + m) % m
}

// evaluates f(x)
func (f linearFunc) apply(x int64) int64 {
	return mod(f.a*x+f.b, f.mod)
}

// combines two linear functions (f o g)
func (f linearFunc) compose(g linearFunc) linearFunc {
	return linearFunc{
		a:   mod(f.a*g.a, f.mod),
		b:   mod(f.a*g.b+f.b, f.mod),
		mod: f.mod,
	}
}

// computes the inverse function f^-1 using modular multiplicative inverse
func (f linearFunc) inverse() linearFunc {
	invA := xmath.ModInv(mod(f.a, f.mod), f.mod)
	return linearFunc{
		a:   invA,
		b:   mod(-f.b*invA, f.mod),
		mod: f.mod,
	}
}

// raises function to nth power using exponentiation by squaring
func (f linearFunc) pow(n int64) linearFunc {
	if n == 0 {
		return linearFunc{a: 1, b: 0, mod: f.mod}
	}
	if n == 1 {
		return f
	}

	if n%2 == 0 {
		half := f.pow(n / 2)
		return half.compose(half)
	}

	return f.compose(f.pow(n - 1))
}

// converts shuffle operations to linear form f(x) = ax + b (mod m)
func shuffleToLinear(s shuffle, mod int64) linearFunc {
	switch s.op {
	// pos' = N-1-pos -> f(x) = -x-1 -> a=-1, b=-1
	case dealNewStack:
		return linearFunc{a: -1, b: -1, mod: mod}
	// pos' = pos-n -> f(x) = x-n -> a=1, b=-n
	case cut:
		return linearFunc{a: 1, b: -int64(s.n), mod: mod}
	// pos' = n*pos -> f(x) = nx -> a=n, b=0
	case dealIncrement:
		return linearFunc{a: int64(s.n), b: 0, mod: mod}
	default:
		return linearFunc{a: 1, b: 0, mod: mod}
	}
}

func composeAll(instructions []shuffle, mod int64) linearFunc {
	f := linearFunc{a: 1, b: 0, mod: mod}
	for _, inst := range instructions {
		f = f.compose(shuffleToLinear(inst, mod))
	}
	return f
}

func day22p02(r io.Reader) (string, error) {
	instructions, err := parseShuffleInstructions(r)
	if err != nil {
		return "", err
	}

	const (
		deckSize   = 119315717514047
		iterations = 101741582076661
		position   = 2020
	)

	f := composeAll(instructions, deckSize)
	fPow := f.pow(iterations)
	fInv := fPow.inverse()

	result := fInv.apply(position)

	return strconv.FormatInt(result, 10), nil
}
