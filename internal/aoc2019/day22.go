package aoc2019

import (
	"fmt"
	"io"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
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

// bigMod performs modulo operation ensuring non-negative result
func bigMod(a, m *big.Int) *big.Int {
	result := new(big.Int).Mod(a, m)
	if result.Sign() < 0 {
		result.Add(result, m)
	}
	return result
}

// modInv computes the modular multiplicative inverse of a mod m
func modInv(a, m *big.Int) *big.Int {
	return new(big.Int).ModInverse(a, m)
}

// linearFunc represents a linear transformation f(x) = (a*x + b) mod m in modular arithmetic.
//
// part 2 requires shuffling a deck of ~10^14 cards ~10^14 times, making simulation impossible.
// each shuffle operation is a linear transformation that can be composed into a single function.
// we build the inverse function directly and raise it to the power of 10^14 using fast exponentiation
// to solve the inverse problem: "which card ends up at position 2020?"
type linearFunc struct {
	a, b, mod *big.Int
}

// evaluates f(x)
func (f linearFunc) apply(x *big.Int) *big.Int {
	ax := new(big.Int).Mul(f.a, x)
	result := new(big.Int).Add(ax, f.b)
	return bigMod(result, f.mod)
}

// combines two linear functions (f o g)
func (f linearFunc) compose(g linearFunc) linearFunc {
	newA := new(big.Int).Mul(f.a, g.a)
	newA.Mod(newA, f.mod)

	temp := new(big.Int).Mul(f.a, g.b)
	newB := new(big.Int).Add(temp, f.b)
	newB.Mod(newB, f.mod)

	return linearFunc{a: newA, b: newB, mod: f.mod}
}

// raises function to nth power using the geometric series formula
// f^n(x) = a^n * x + b * (a^n - 1) / (a - 1) mod m
func (f linearFunc) pow(n int64) linearFunc {
	if n == 0 {
		return linearFunc{a: big.NewInt(1), b: big.NewInt(0), mod: f.mod}
	}
	if n == 1 {
		return f
	}

	nBig := big.NewInt(n)
	aPow := new(big.Int).Exp(f.a, nBig, f.mod)

	// Compute b * (a^n - 1) / (a - 1) mod m (geometric series)
	aMinus1 := new(big.Int).Sub(f.a, big.NewInt(1))
	aMinus1.Mod(aMinus1, f.mod)
	aMinus1Inv := modInv(aMinus1, f.mod)

	aPowMinus1 := new(big.Int).Sub(aPow, big.NewInt(1))
	aPowMinus1.Mod(aPowMinus1, f.mod)

	bCoeff := new(big.Int).Mul(f.b, aPowMinus1)
	bCoeff.Mul(bCoeff, aMinus1Inv)
	bCoeff.Mod(bCoeff, f.mod)

	return linearFunc{a: aPow, b: bCoeff, mod: f.mod}
}

// converts shuffle operations to inverse linear form f^-1(x) = ax + b (mod m)
func shuffleToLinear(s shuffle, mod *big.Int) linearFunc {
	switch s.op {
	// f(x) = -x-1 is self-inverse, so f^-1(x) = -x-1
	case dealNewStack:
		return linearFunc{
			a:   big.NewInt(-1),
			b:   big.NewInt(-1),
			mod: mod,
		}
	// f(x) = x-n, so f^-1(x) = x+n
	case cut:
		return linearFunc{
			a:   big.NewInt(1),
			b:   big.NewInt(int64(s.n)),
			mod: mod,
		}
	// f(x) = n*x, so f^-1(x) = modinv(n)*x
	case dealIncrement:
		n := big.NewInt(int64(s.n))
		invN := modInv(n, mod)
		return linearFunc{a: invN, b: big.NewInt(0), mod: mod}
	default:
		return linearFunc{
			a:   big.NewInt(1),
			b:   big.NewInt(0),
			mod: mod,
		}
	}
}

// composeInverse composes all shuffle operations into a single inverse linear function
func composeInverse(instructions []shuffle, mod *big.Int) linearFunc {
	f := linearFunc{
		a:   big.NewInt(1),
		b:   big.NewInt(0),
		mod: mod,
	}
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

	deckSize := big.NewInt(119315717514047)
	iterations := int64(101741582076661)
	position := big.NewInt(2020)

	fInverse := composeInverse(instructions, deckSize)
	fInversePow := fInverse.pow(iterations)
	result := fInversePow.apply(position)

	return result.String(), nil
}
