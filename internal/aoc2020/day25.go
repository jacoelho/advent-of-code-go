package aoc2020

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/internal/convert"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
	"github.com/jacoelho/advent-of-code-go/internal/xiter"
)

const modulus = 20201227

func parsePublicKeys(r io.Reader) (int, int, error) {
	s := scanner.NewScanner(r, convert.ScanNumber[int])
	keys := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	if len(keys) < 2 {
		return 0, 0, fmt.Errorf("expected 2 public keys, got %d", len(keys))
	}
	return keys[0], keys[1], nil
}

// transform returns an infinite sequence of transformed values according to the cryptographic protocol
func transform(subject int) iter.Seq[int] {
	return func(yield func(int) bool) {
		val := 1
		for {
			val = (val * subject) % modulus
			if !yield(val) {
				return
			}
		}
	}
}

// recoverLoopSize brute-forces the loop size that produces the target public key from subject 7
func recoverLoopSize(publicKey int) int {
	for i, val := range xiter.Enumerate(transform(7)) {
		if val == publicKey {
			return i
		}
	}
	panic("unreachable")
}

func day25p01(r io.Reader) (string, error) {
	cardPublicKey, doorPublicKey, err := parsePublicKeys(r)
	if err != nil {
		return "", err
	}

	loopSize := recoverLoopSize(cardPublicKey)
	encryptionKey, _ := xiter.Nth(transform(doorPublicKey), loopSize)

	return strconv.Itoa(encryptionKey), nil
}
