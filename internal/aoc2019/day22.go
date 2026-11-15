package aoc2019

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
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
	text := string(line)
	if text == "deal into new stack" {
		return shuffle{op: dealNewStack}, nil
	}
	if arg, ok := strings.CutPrefix(text, "cut "); ok {
		n, err := convert.ScanNumber[int]([]byte(arg))
		if err != nil {
			return shuffle{}, fmt.Errorf("shuffle %q: %w", text, err)
		}
		return shuffle{op: cut, n: n}, nil
	}
	if arg, ok := strings.CutPrefix(text, "deal with increment "); ok {
		n, err := convert.ScanNumber[int]([]byte(arg))
		if err != nil {
			return shuffle{}, fmt.Errorf("shuffle %q: %w", text, err)
		}
		if n <= 0 {
			return shuffle{}, fmt.Errorf("shuffle %q needs a positive increment", text)
		}
		return shuffle{op: dealIncrement, n: n}, nil
	}
	return shuffle{}, fmt.Errorf("unknown shuffle %q", text)
}

func parseShuffleInstructions(r io.Reader) ([]shuffle, error) {
	s := scanner.NewScanner(r, parseShuffle)
	instructions := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}
	return instructions, nil
}

// cardMap is the affine map x -> a*x + b modulo a deck size below 2^63; every
// shuffle moves the card at position x to position a*x + b.
type cardMap struct {
	a, b, mod uint64
}

func (f cardMap) apply(x uint64) uint64 {
	return (xmath.MulMod(f.a, x, f.mod) + f.b) % f.mod
}

// then returns the map that applies f and then g.
func (f cardMap) then(g cardMap) cardMap {
	return cardMap{a: xmath.MulMod(g.a, f.a, f.mod), b: g.apply(f.b), mod: f.mod}
}

func (f cardMap) pow(n uint64) cardMap {
	result := cardMap{a: 1 % f.mod, mod: f.mod}
	for ; n > 0; n >>= 1 {
		if n&1 != 0 {
			result = result.then(f)
		}
		f = f.then(f)
	}
	return result
}

func (f cardMap) inverse() (cardMap, error) {
	inv, ok := xmath.ModInverse(f.a, f.mod)
	if !ok {
		return cardMap{}, fmt.Errorf("shuffle is not invertible for a deck of %d cards", f.mod)
	}
	return cardMap{a: inv, b: xmath.MulMod(inv, (f.mod-f.b)%f.mod, f.mod), mod: f.mod}, nil
}

// shuffleMap returns the map from a card's position before the instructions
// to its position after them, for a deck of size cards.
func shuffleMap(instructions []shuffle, size uint64) (cardMap, error) {
	f := cardMap{a: 1 % size, mod: size}
	for _, inst := range instructions {
		var step cardMap
		switch inst.op {
		case dealNewStack:
			step = cardMap{a: size - 1, b: size - 1, mod: size}
		case cut:
			step = cardMap{a: 1, b: (size - uint64(xmath.Modulo(inst.n, int(size)))) % size, mod: size}
		case dealIncrement:
			step = cardMap{a: uint64(xmath.Modulo(inst.n, int(size))), mod: size}
			if _, ok := xmath.ModInverse(step.a, size); !ok {
				return cardMap{}, fmt.Errorf("increment %d does not permute a deck of %d cards", inst.n, size)
			}
		}
		f = f.then(step)
	}
	return f, nil
}

func day22p01(r io.Reader) (string, error) {
	instructions, err := parseShuffleInstructions(r)
	if err != nil {
		return "", err
	}
	f, err := shuffleMap(instructions, 10007)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(f.apply(2019), 10), nil
}

func day22p02(r io.Reader) (string, error) {
	instructions, err := parseShuffleInstructions(r)
	if err != nil {
		return "", err
	}
	f, err := shuffleMap(instructions, 119315717514047)
	if err != nil {
		return "", err
	}
	inverse, err := f.pow(101741582076661).inverse()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(inverse.apply(2020), 10), nil
}
