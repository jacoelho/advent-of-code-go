package aoc2015

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

func mineAdventCoin(r io.Reader, zeroes int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	secret := bytes.TrimSpace(input)
	if len(secret) == 0 {
		return "", fmt.Errorf("empty AdventCoin secret")
	}
	message := make([]byte, len(secret), len(secret)+20)
	copy(message, secret)
	for suffix := int64(1); suffix > 0; suffix++ {
		message = strconv.AppendInt(message[:len(secret)], suffix, 10)
		digest := md5.Sum(message)
		if digest[0] == 0 && digest[1] == 0 && digest[2] < 16 && (zeroes == 5 || digest[2] == 0) {
			return strconv.FormatInt(suffix, 10), nil
		}
	}
	return "", fmt.Errorf("AdventCoin suffix exceeds int64")
}

func day04p01(r io.Reader) (string, error) {
	return mineAdventCoin(r, 5)
}

func day04p02(r io.Reader) (string, error) {
	return mineAdventCoin(r, 6)
}
