package aoc2015

import (
	"bytes"
	"fmt"
	"io"
)

func validCorporatePassword(password []byte) bool {
	straight, distinctPairs := false, false
	var pair byte
	for i, letter := range password {
		if letter == 'i' || letter == 'o' || letter == 'l' {
			return false
		}
		if i >= 2 && password[i-2]+1 == password[i-1] && password[i-1]+1 == letter {
			straight = true
		}
		if i > 0 && password[i-1] == letter {
			if pair == 0 {
				pair = letter
			} else if pair != letter {
				distinctPairs = true
			}
		}
	}
	return straight && distinctPairs
}

func nextCorporatePassword(password []byte) error {
	for {
		i := len(password) - 1
		for i >= 0 && password[i] == 'z' {
			password[i] = 'a'
			i--
		}
		if i < 0 {
			return fmt.Errorf("no later eight-letter password")
		}
		password[i]++
		for i, letter := range password {
			if letter == 'i' || letter == 'o' || letter == 'l' {
				password[i]++
				for j := i + 1; j < len(password); j++ {
					password[j] = 'a'
				}
				break
			}
		}
		if validCorporatePassword(password) {
			return nil
		}
	}
}

func corporatePassword(r io.Reader, generations int) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	password := bytes.TrimSpace(input)
	if len(password) != 8 {
		return "", fmt.Errorf("password must contain eight lowercase letters")
	}
	for _, letter := range password {
		if letter < 'a' || letter > 'z' {
			return "", fmt.Errorf("invalid password letter %q", letter)
		}
	}
	for range generations {
		if err := nextCorporatePassword(password); err != nil {
			return "", err
		}
	}
	return string(password), nil
}

func day11p01(r io.Reader) (string, error) {
	return corporatePassword(r, 1)
}

func day11p02(r io.Reader) (string, error) {
	return corporatePassword(r, 2)
}
