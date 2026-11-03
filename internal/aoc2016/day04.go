package aoc2016

import (
	"bytes"
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
)

type encryptedRoom struct {
	name, checksum string
	sector         int
}

func roomRecords(r io.Reader) *scanner.Scanner[encryptedRoom] {
	return scanner.NewScanner(r, func(line []byte) (encryptedRoom, error) {
		body, checksum, ok := bytes.Cut(line, []byte{'['})
		checksum, closed := bytes.CutSuffix(checksum, []byte{']'})
		dash := bytes.LastIndexByte(body, '-')
		if !ok || !closed || dash < 1 || len(checksum) != 5 {
			return encryptedRoom{}, fmt.Errorf("invalid room %q", line)
		}
		sector, err := convert.ScanNumber[int](body[dash+1:])
		if err != nil {
			return encryptedRoom{}, fmt.Errorf("invalid room sector: %w", err)
		}
		if sector < 0 {
			return encryptedRoom{}, fmt.Errorf("negative room sector %d", sector)
		}
		for _, letter := range body[:dash] {
			if letter != '-' && (letter < 'a' || letter > 'z') {
				return encryptedRoom{}, fmt.Errorf("invalid room name %q", body[:dash])
			}
		}
		return encryptedRoom{string(body[:dash]), string(checksum), sector}, nil
	})
}

func (room encryptedRoom) real() bool {
	var frequencies [26]int
	for _, letter := range room.name {
		if letter != '-' {
			frequencies[letter-'a']++
		}
	}
	var letters []byte
	for i, count := range frequencies {
		if count > 0 {
			letters = append(letters, byte(i)+'a')
		}
	}
	slices.SortFunc(letters, func(a, b byte) int {
		return cmp.Or(cmp.Compare(frequencies[b-'a'], frequencies[a-'a']), cmp.Compare(a, b))
	})
	return len(letters) >= 5 && string(letters[:5]) == room.checksum
}

func (room encryptedRoom) decryptedName() string {
	var name strings.Builder
	for _, letter := range room.name {
		if letter == '-' {
			name.WriteByte(' ')
		} else {
			name.WriteByte(byte((int(letter-'a')+room.sector%26)%26) + 'a')
		}
	}
	return name.String()
}

func day04p01(r io.Reader) (string, error) {
	s := roomRecords(r)
	total := xiter.From(s.Values()).Filter(encryptedRoom.real).SumBy(func(room encryptedRoom) int {
		return room.sector
	})
	if err := s.Err(); err != nil {
		return "", err
	}
	return strconv.Itoa(total), nil
}

func day04p02(r io.Reader) (string, error) {
	s := roomRecords(r)
	rooms := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return "", err
	}
	for _, room := range rooms {
		if room.real() && room.decryptedName() == "northpole object storage" {
			return strconv.Itoa(room.sector), nil
		}
	}
	return "", fmt.Errorf("no real northpole object storage room")
}
