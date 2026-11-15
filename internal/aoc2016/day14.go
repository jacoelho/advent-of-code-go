package aoc2016

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/conc"
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

func padHashes(salt []byte, rounds int) iter.Seq[padHashFeatures] {
	return conc.Blocks(0, 64, func(start, end int) []padHashFeatures {
		message := make([]byte, len(salt), len(salt)+20)
		copy(message, salt)
		features := make([]padHashFeatures, 0, end-start)
		var encoded [2 * md5.Size]byte
		for index := start; index < end; index++ {
			message = strconv.AppendInt(message[:len(salt)], int64(index), 10)
			digest := md5.Sum(message)
			hex.Encode(encoded[:], digest[:])
			for range rounds - 1 {
				digest = md5.Sum(encoded[:])
				hex.Encode(encoded[:], digest[:])
			}
			features = append(features, inspectPadHash(encoded[:]))
		}
		return features
	})
}

const padLookahead = 1000

// findPadKey returns the index of the wanted key in a hash stream.
// Candidate c is judged once hash c+padLookahead arrives; future counts the
// quintuple digits of hashes c+1 through c+padLookahead.
func findPadKey(hashes iter.Seq[padHashFeatures], wanted int) (int, error) {
	var window [padLookahead + 1]padHashFeatures
	var future [16]int
	adjust := func(features padHashFeatures, delta int) {
		for digit := range future {
			if features.quintuples&(1<<digit) != 0 {
				future[digit] += delta
			}
		}
	}
	keys, index := 0, 0
	for features := range hashes {
		if index > padLookahead {
			adjust(window[(index-padLookahead)%len(window)], -1)
		}
		window[index%len(window)] = features
		if index > 0 {
			adjust(features, 1)
		}
		if candidate := index - padLookahead; candidate >= 0 {
			if triple := window[candidate%len(window)].triple; triple >= 0 && future[triple] > 0 {
				keys++
				if keys == wanted {
					return candidate, nil
				}
			}
		}
		index++
	}
	return 0, fmt.Errorf("hash stream ended after %d hashes with %d keys", index, keys)
}

func oneTimePad(r io.Reader, rounds int) (string, error) {
	salt, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	salt = bytes.TrimSpace(salt)
	if len(salt) == 0 || bytes.ContainsAny(salt, " \t\r\n") {
		return "", fmt.Errorf("invalid one-time pad salt %q", salt)
	}
	index, err := findPadKey(padHashes(salt, rounds), 64)
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
