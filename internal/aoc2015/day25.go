package aoc2015

import (
	"fmt"
	"io"
	"math/big"
	"regexp"
)

func day25p01(r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	coordinates := regexp.MustCompile(`row (\d+), column (\d+)`).FindSubmatch(data)
	if coordinates == nil {
		return "", fmt.Errorf("missing code grid coordinates")
	}
	row, _ := new(big.Int).SetString(string(coordinates[1]), 10)
	column, _ := new(big.Int).SetString(string(coordinates[2]), 10)
	if row.Sign() <= 0 || column.Sign() <= 0 {
		return "", fmt.Errorf("grid coordinates must be positive")
	}
	diagonal := new(big.Int).Add(row, column)
	diagonal.Sub(diagonal, big.NewInt(2))
	exponent := new(big.Int).Add(diagonal, big.NewInt(1))
	exponent.Mul(exponent, diagonal).Quo(exponent, big.NewInt(2))
	exponent.Add(exponent, column).Sub(exponent, big.NewInt(1))
	modulus := big.NewInt(33554393)
	code := new(big.Int).Exp(big.NewInt(252533), exponent, modulus)
	code.Mul(code, big.NewInt(20151125)).Mod(code, modulus)
	return code.String(), nil
}
