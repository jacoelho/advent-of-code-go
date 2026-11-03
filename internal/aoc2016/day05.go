package aoc2016

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
)

type doorPasswordMode int

const (
	sequentialDoorPassword doorPasswordMode = iota
	positionalDoorPassword
)

func crackDoorPassword(r io.Reader, mode doorPasswordMode) (string, error) {
	input, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	doorID := bytes.TrimSpace(input)
	if len(doorID) == 0 {
		return "", fmt.Errorf("empty door ID")
	}
	message := make([]byte, len(doorID), len(doorID)+20)
	copy(message, doorID)
	const hex = "0123456789abcdef"
	var password [8]byte
	filled := 0
	for index := int64(0); index >= 0; index++ {
		message = strconv.AppendInt(message[:len(doorID)], index, 10)
		digest := md5.Sum(message)
		if digest[0] != 0 || digest[1] != 0 || digest[2] >= 16 {
			continue
		}
		if mode == sequentialDoorPassword {
			password[filled] = hex[digest[2]]
			filled++
		} else if position := digest[2]; position < 8 && password[position] == 0 {
			password[position] = hex[digest[3]>>4]
			filled++
		}
		if filled == len(password) {
			return string(password[:]), nil
		}
	}
	return "", fmt.Errorf("door password index exceeds int64")
}

func day05p01(r io.Reader) (string, error) {
	return crackDoorPassword(r, sequentialDoorPassword)
}

func day05p02(r io.Reader) (string, error) {
	return crackDoorPassword(r, positionalDoorPassword)
}
