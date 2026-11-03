package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

func safeTiles(r io.Reader, rows int) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 || rows < 1 {
		return "", fmt.Errorf("need an initial row and positive row count")
	}
	current, next := make([]bool, len(data)+2), make([]bool, len(data)+2)
	for i, tile := range data {
		if tile != '.' && tile != '^' {
			return "", fmt.Errorf("invalid tile %q", tile)
		}
		current[i+1] = tile == '^'
	}
	var safe int64
	for range rows {
		for i := 1; i <= len(data); i++ {
			if !current[i] {
				safe++
			}
			next[i] = current[i-1] != current[i+1]
		}
		current, next = next, current
	}
	return strconv.FormatInt(safe, 10), nil
}

func day18p01(r io.Reader) (string, error) {
	return safeTiles(r, 40)
}

func day18p02(r io.Reader) (string, error) {
	return safeTiles(r, 400000)
}
