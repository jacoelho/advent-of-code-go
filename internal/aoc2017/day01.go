package aoc2017

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func captchaDigits(r io.Reader) ([]byte, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	input = bytes.TrimSpace(input)
	if len(input) == 0 {
		return nil, fmt.Errorf("empty captcha")
	}

	for _, digit := range input {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid captcha digit %q", digit)
		}
	}

	return input, nil
}

func captchaSum(digits []byte, offset int) string {
	total := 0
	for i, digit := range digits {
		if digit == digits[(i+offset)%len(digits)] {
			total += int(digit - '0')
		}
	}

	return strconv.Itoa(total)
}

func day01p01(r io.Reader) (string, error) {
	digits, err := captchaDigits(r)
	if err != nil {
		return "", err
	}

	return captchaSum(digits, 1), nil
}

func day01p02(r io.Reader) (string, error) {
	digits, err := captchaDigits(r)
	if err != nil {
		return "", err
	}
	if len(digits)%2 != 0 {
		return "", fmt.Errorf("halfway captcha requires an even number of digits")
	}

	return captchaSum(digits, len(digits)/2), nil
}
