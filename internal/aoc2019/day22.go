package aoc2019

import (
	"fmt"
	"io"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
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
	const inputLimit = 4 << 20
	limited := &io.LimitedReader{R: r, N: inputLimit + 1}
	s := scanner.NewScanner(limited, parseShuffle)
	instructions := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	if limited.N == 0 || len(instructions) > 100_000 {
		return nil, fmt.Errorf("shuffle input exceeds byte or instruction limit")
	}
	return instructions, nil
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

func bigMod(a, m *big.Int) *big.Int {
	result := new(big.Int).Mod(a, m)
	if result.Sign() < 0 {
		result.Add(result, m)
	}
	return result
}

func modInv(a, m *big.Int) *big.Int {
	return new(big.Int).ModInverse(a, m)
}

// linearFunc represents f(x) = (a*x + b) mod m.
type linearFunc struct {
	a, b, mod *big.Int
}

func (f linearFunc) apply(x *big.Int) *big.Int {
	ax := new(big.Int).Mul(f.a, x)
	result := new(big.Int).Add(ax, f.b)
	return bigMod(result, f.mod)
}

// compose returns f(g(x)).
func (f linearFunc) compose(g linearFunc) linearFunc {
	newA := new(big.Int).Mul(f.a, g.a)
	newA.Mod(newA, f.mod)

	temp := new(big.Int).Mul(f.a, g.b)
	newB := new(big.Int).Add(temp, f.b)
	newB.Mod(newB, f.mod)

	return linearFunc{a: newA, b: newB, mod: f.mod}
}

func (f linearFunc) pow(n uint64) linearFunc {
	result := linearFunc{a: big.NewInt(1), b: big.NewInt(0), mod: f.mod}
	for n > 0 {
		if n&1 != 0 {
			result = result.compose(f)
		}
		n >>= 1
		if n > 0 {
			f = f.compose(f)
		}
	}
	return result
}

// shuffleToLinear returns the inverse of s as an affine transformation modulo mod.
func shuffleToLinear(s shuffle, mod *big.Int) (linearFunc, error) {
	switch s.op {
	// f(x) = -x-1 is self-inverse, so f^-1(x) = -x-1
	case dealNewStack:
		return linearFunc{
			a:   big.NewInt(-1),
			b:   big.NewInt(-1),
			mod: mod,
		}, nil
	// f(x) = x-n, so f^-1(x) = x+n
	case cut:
		return linearFunc{
			a:   big.NewInt(1),
			b:   big.NewInt(int64(s.n)),
			mod: mod,
		}, nil
	// f(x) = n*x, so f^-1(x) = modinv(n)*x
	case dealIncrement:
		n := big.NewInt(int64(s.n))
		invN := modInv(n, mod)
		if invN == nil {
			return linearFunc{}, fmt.Errorf("increment %d is not invertible modulo %s", s.n, mod)
		}
		return linearFunc{a: invN, b: big.NewInt(0), mod: mod}, nil
	default:
		return linearFunc{}, fmt.Errorf("unknown shuffle operation")
	}
}

// composeInverse returns the inverse of the instruction sequence.
// Since the inverse of g(f(x)) is f^-1(g^-1(x)), it composes inverses in input order.
func composeInverse(instructions []shuffle, mod *big.Int) (linearFunc, error) {
	f := linearFunc{
		a:   big.NewInt(1),
		b:   big.NewInt(0),
		mod: mod,
	}
	for _, inst := range instructions {
		operation, err := shuffleToLinear(inst, mod)
		if err != nil {
			return linearFunc{}, err
		}
		f = f.compose(operation)
	}
	return f, nil
}

func day22p02(r io.Reader) (string, error) {
	instructions, err := parseShuffleInstructions(r)
	if err != nil {
		return "", err
	}

	deckSize := big.NewInt(119315717514047)
	iterations := uint64(101741582076661)
	position := big.NewInt(2020)

	fInverse, err := composeInverse(instructions, deckSize)
	if err != nil {
		return "", err
	}
	fInversePow := fInverse.pow(iterations)
	result := fInversePow.apply(position)

	return result.String(), nil
}
