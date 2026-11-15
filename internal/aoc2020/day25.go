package aoc2020

import (
	"fmt"
	"io"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

const handshakeModulus = 20201227

func parsePublicKeys(r io.Reader) (int, int, error) {
	s := scanner.NewScanner(r, func(line []byte) (int, error) {
		key, err := strconv.Atoi(string(line))
		if err != nil {
			return 0, fmt.Errorf("invalid public key %q: %w", line, err)
		}
		if key <= 0 || key >= handshakeModulus {
			return 0, fmt.Errorf("public key %d outside 1..%d", key, handshakeModulus-1)
		}
		return key, nil
	})
	keys := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	if len(keys) != 2 {
		return 0, 0, fmt.Errorf("expected 2 public keys, got %d", len(keys))
	}
	return keys[0], keys[1], nil
}

func day25p01(r io.Reader) (string, error) {
	cardKey, doorKey, err := parsePublicKeys(r)
	if err != nil {
		return "", err
	}
	// 7 is a primitive root modulo 20201227, so every key in 1..modulus-1 is a
	// power of 7 and the search ends within one full cycle.
	power, encryptionKey := 1, 1
	for power != cardKey {
		power = power * 7 % handshakeModulus
		encryptionKey = encryptionKey * doorKey % handshakeModulus
	}
	return strconv.Itoa(encryptionKey), nil
}
