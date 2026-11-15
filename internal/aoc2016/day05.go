package aoc2016

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"iter"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/conc"
)

const hexDigits = "0123456789abcdef"

var errDoorHashesExhausted = errors.New("door hashes exhausted before the password was complete")

// doorHashes yields the sixth and seventh hex digits of every MD5 hash of the
// door ID and an increasing index whose hex form starts with five zeros.
func doorHashes(doorID []byte) iter.Seq[[2]byte] {
	return conc.Blocks(0, 1<<14, func(start, end int) [][2]byte {
		message := make([]byte, len(doorID), len(doorID)+20)
		copy(message, doorID)
		var hits [][2]byte
		for index := start; index < end; index++ {
			message = strconv.AppendInt(message[:len(doorID)], int64(index), 10)
			digest := md5.Sum(message)
			if digest[0] == 0 && digest[1] == 0 && digest[2] < 16 {
				hits = append(hits, [2]byte{digest[2], digest[3] >> 4})
			}
		}
		return hits
	})
}

func readDoorID(r io.Reader) ([]byte, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	doorID := bytes.TrimSpace(input)
	if len(doorID) == 0 || bytes.ContainsAny(doorID, " \t\r\n") {
		return nil, fmt.Errorf("invalid door ID %q", doorID)
	}
	return doorID, nil
}

func day05p01(r io.Reader) (string, error) {
	doorID, err := readDoorID(r)
	if err != nil {
		return "", err
	}
	password := make([]byte, 0, 8)
	for digits := range doorHashes(doorID) {
		password = append(password, hexDigits[digits[0]])
		if len(password) == cap(password) {
			return string(password), nil
		}
	}
	return "", errDoorHashesExhausted
}

func day05p02(r io.Reader) (string, error) {
	doorID, err := readDoorID(r)
	if err != nil {
		return "", err
	}
	var password [8]byte
	filled := 0
	for digits := range doorHashes(doorID) {
		if position := digits[0]; position < 8 && password[position] == 0 {
			password[position] = hexDigits[digits[1]]
			filled++
			if filled == len(password) {
				return string(password[:]), nil
			}
		}
	}
	return "", errDoorHashesExhausted
}
