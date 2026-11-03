package aoc2016

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
)

type padHashFeatures struct {
	triple     int8
	quintuples uint16
}

func inspectPadHash(hash []byte) padHashFeatures {
	features := padHashFeatures{triple: -1}
	for start := 0; start < len(hash); {
		end := start + 1
		for end < len(hash) && hash[end] == hash[start] {
			end++
		}
		digit := hash[start] - '0'
		if hash[start] >= 'a' {
			digit = hash[start] - 'a' + 10
		}
		if end-start >= 3 && features.triple == -1 {
			features.triple = int8(digit)
		}
		if end-start >= 5 {
			features.quintuples |= 1 << digit
		}
		start = end
	}
	return features
}

func padHash(salt []byte, index, rounds int) padHashFeatures {
	message := strconv.AppendInt(bytes.Clone(salt), int64(index), 10)
	digest := md5.Sum(message)
	var encoded [32]byte
	hex.Encode(encoded[:], digest[:])
	for round := 1; round < rounds; round++ {
		digest = md5.Sum(encoded[:])
		hex.Encode(encoded[:], digest[:])
	}
	return inspectPadHash(encoded[:])
}

func findPadKey(hashAt func(int) padHashFeatures, wanted int) (int, error) {
	var window [1001]padHashFeatures
	var future [16]int
	adjust := func(features padHashFeatures, delta int) {
		for digit := range 16 {
			if features.quintuples&(1<<digit) != 0 {
				future[digit] += delta
			}
		}
	}
	for index := range window {
		window[index] = hashAt(index)
		if index != 0 {
			adjust(window[index], 1)
		}
	}
	keys := 0
	for index := range 1_000_000 {
		candidate := window[index%len(window)]
		if candidate.triple >= 0 && future[candidate.triple] > 0 {
			keys++
			if keys == wanted {
				return index, nil
			}
		}
		// The candidate itself is excluded; only its following 1000 hashes qualify.
		adjust(window[(index+1)%len(window)], -1)
		next := hashAt(index + len(window))
		window[index%len(window)] = next
		adjust(next, 1)
	}
	return 0, fmt.Errorf("one-time pad search exceeds 1000000 candidates")
}

func oneTimePad(r io.Reader, rounds int) (string, error) {
	salt, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	salt = bytes.TrimSpace(salt)
	if len(salt) == 0 || bytes.ContainsAny(salt, " \t\r\n") {
		return "", fmt.Errorf("invalid one-time pad salt")
	}
	index, err := findPadKey(func(index int) padHashFeatures { return padHash(salt, index, rounds) }, 64)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(index), nil
}

func day14p01(r io.Reader) (string, error) {
	return oneTimePad(r, 1)
}

func day14p02(r io.Reader) (string, error) {
	return oneTimePad(r, 2017)
}
