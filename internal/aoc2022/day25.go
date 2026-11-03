package aoc2022

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
	"slices"
	"strings"
)

func encodeSNAFU(number *big.Int) string {
	if number.Sign() == 0 {
		return "0"
	}
	value := new(big.Int).Set(number)
	five, one := big.NewInt(5), big.NewInt(1)
	var remainder big.Int
	var digits []byte
	for value.Sign() != 0 {
		value.QuoRem(value, five, &remainder)
		digit := remainder.Int64()
		if digit > 2 {
			digit -= 5
			value.Add(value, one)
		} else if digit < -2 {
			digit += 5
			value.Sub(value, one)
		}
		digits = append(digits, "=-012"[digit+2])
	}
	slices.Reverse(digits)
	return string(digits)
}

func day25p01(r io.Reader) (string, error) {
	s := bufio.NewScanner(r)
	five := big.NewInt(5)
	var total, value, digit big.Int
	characters := 0
	for s.Scan() {
		line := s.Bytes()
		if len(line) == 0 || len(line) > 1024 {
			return "", fmt.Errorf("SNAFU numbers require one to 1024 digits")
		}
		characters += len(line)
		if characters > 1_000_000 {
			return "", fmt.Errorf("SNAFU input exceeds one million digits")
		}
		value.SetInt64(0)
		for _, c := range line {
			d := strings.IndexByte("=-012", c)
			if d == -1 {
				return "", fmt.Errorf("invalid SNAFU digit %q", c)
			}
			value.Mul(&value, five)
			value.Add(&value, digit.SetInt64(int64(d-2)))
		}
		total.Add(&total, &value)
	}
	if err := s.Err(); err != nil {
		return "", err
	}
	if characters == 0 {
		return "", fmt.Errorf("empty SNAFU input")
	}
	return encodeSNAFU(&total), nil
}
