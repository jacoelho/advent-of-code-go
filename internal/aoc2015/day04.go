package aoc2015

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/conc"
)

const adventCoinChunk = 1 << 14

func leadingZeroNibbles(digest [md5.Size]byte, zeroes int) bool {
	for _, b := range digest[:zeroes/2] {
		if b != 0 {
			return false
		}
	}
	return zeroes%2 == 0 || digest[zeroes/2] < 16
}

func mineAdventCoin(r io.Reader, zeroes int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	secret := bytes.TrimSpace(input)
	if len(secret) == 0 {
		return "", fmt.Errorf("empty AdventCoin secret")
	}
	coins := conc.Blocks(1, adventCoinChunk, func(start, end int) []int {
		message := make([]byte, len(secret), len(secret)+20)
		copy(message, secret)
		for suffix := start; suffix < end; suffix++ {
			message = strconv.AppendInt(message[:len(secret)], int64(suffix), 10)
			if leadingZeroNibbles(md5.Sum(message), zeroes) {
				return []int{suffix}
			}
		}
		return nil
	})
	for suffix := range coins {
		return strconv.Itoa(suffix), nil
	}
	return "", fmt.Errorf("no AdventCoin suffix fits in int")
}

func day04p01(r io.Reader) (string, error) {
	return mineAdventCoin(r, 5)
}

func day04p02(r io.Reader) (string, error) {
	return mineAdventCoin(r, 6)
}
