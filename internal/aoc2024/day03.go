package aoc2024

import (
	"io"
	"iter"
	"regexp"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

var re = regexp.MustCompile(`do\(\)|don't\(\)|mul\((\d+),(\d+)\)`)

type tokenType int

const (
	tokenTypeMultiply tokenType = iota
	tokenTypeDo
	tokenTypeDont
)

type token struct {
	typ tokenType
	a   int
	b   int
}

func parseCorruptedMemory(input string) iter.Seq[token] {
	return func(yield func(token) bool) {
		matches := re.FindAllStringSubmatch(input, -1)
		if len(matches) == 0 {
			return
		}

		for _, match := range matches {
			var t token
			switch {
			case strings.HasPrefix(match[0], "do("):
				t.typ = tokenTypeDo
			case strings.HasPrefix(match[0], "don't"):
				t.typ = tokenTypeDont
			case strings.HasPrefix(match[0], "mul("):
				t.typ = tokenTypeMultiply
				t.a = aoc.Must(strconv.Atoi(match[1]))
				t.b = aoc.Must(strconv.Atoi(match[2]))
			default:
				panic("unreachable")
			}
			if !yield(t) {
				return
			}
		}
	}
}

func day03p01(r io.Reader) (string, error) {
	input := aoc.Must(io.ReadAll(r))

	corruptedMemory := parseCorruptedMemory(string(input))

	result := xiter.Reduce(func(sum int, v token) int {
		return sum + v.a*v.b
	}, 0, corruptedMemory)

	return strconv.Itoa(result), nil
}

func day03p02(r io.Reader) (string, error) {
	input := aoc.Must(io.ReadAll(r))

	corruptedMemory := parseCorruptedMemory(string(input))

	var result int
	enableMultiply := true
	for t := range corruptedMemory {
		switch t.typ {
		case tokenTypeMultiply:
			if enableMultiply {
				result += t.a * t.b
			}
		case tokenTypeDo:
			enableMultiply = true
		case tokenTypeDont:
			enableMultiply = false
		}
	}

	return strconv.Itoa(result), nil
}
