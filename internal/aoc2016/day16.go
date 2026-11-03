package aoc2016

import (
	"bytes"
	"fmt"
	"io"
)

func dragonChecksum(r io.Reader, diskSize int) (string, error) {
	seed, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	seed = bytes.TrimSpace(seed)
	if len(seed) == 0 || diskSize <= 0 {
		return "", fmt.Errorf("empty dragon seed or disk")
	}
	for _, bit := range seed {
		if bit != '0' && bit != '1' {
			return "", fmt.Errorf("invalid dragon seed bit %q", bit)
		}
	}
	data := make([]byte, diskSize)
	filled := copy(data, seed)
	for filled < len(data) {
		previous := filled
		data[filled] = '0'
		filled++
		for i := previous - 1; i >= 0 && filled < len(data); i-- {
			data[filled] = data[i] ^ 1
			filled++
		}
	}
	for len(data)%2 == 0 {
		for i := range len(data) / 2 {
			if data[2*i] == data[2*i+1] {
				data[i] = '1'
			} else {
				data[i] = '0'
			}
		}
		data = data[:len(data)/2]
	}
	return string(data), nil
}

func day16p01(r io.Reader) (string, error) {
	return dragonChecksum(r, 272)
}

func day16p02(r io.Reader) (string, error) {
	return dragonChecksum(r, 35651584)
}
