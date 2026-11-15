package aoc2021

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

type packetBits struct {
	nibbles []byte
	pos     int
	end     int
}

func (b *packetBits) read(n int) (uint64, error) {
	if n > b.end-b.pos {
		return 0, fmt.Errorf("packet truncated at bit %d: need %d bits", b.pos, n)
	}
	var value uint64
	for range n {
		value = value<<1 | uint64(b.nibbles[b.pos/4]>>(3-b.pos%4)&1)
		b.pos++
	}
	return value, nil
}

func (b *packetBits) packet() (uint64, uint64, error) {
	version, err := b.read(3)
	if err != nil {
		return 0, 0, err
	}
	kind, err := b.read(3)
	if err != nil {
		return 0, 0, err
	}
	if kind == 4 {
		var value uint64
		for {
			group, err := b.read(5)
			if err != nil {
				return 0, 0, err
			}
			if value > math.MaxUint64>>4 {
				return 0, 0, fmt.Errorf("packet literal exceeds uint64")
			}
			value = value<<4 | group&15
			if group&16 == 0 {
				return version, value, nil
			}
		}
	}

	mode, err := b.read(1)
	if err != nil {
		return 0, 0, err
	}
	lengthBits := 15
	if mode == 1 {
		lengthBits = 11
	}
	length, err := b.read(lengthBits)
	if err != nil {
		return 0, 0, err
	}
	end := b.end
	if mode == 0 {
		if int(length) > b.end-b.pos {
			return 0, 0, fmt.Errorf("operator payload exceeds remaining transmission")
		}
		b.end = b.pos + int(length)
	}

	var values []uint64
	for mode == 0 && b.pos < b.end || mode == 1 && len(values) < int(length) {
		childVersion, value, err := b.packet()
		if err != nil {
			return 0, 0, err
		}
		version += childVersion
		values = append(values, value)
	}
	b.end = end
	if len(values) == 0 || kind >= 5 && len(values) != 2 {
		return 0, 0, fmt.Errorf("operator %d has invalid operand count %d", kind, len(values))
	}

	value := values[0]
	for _, next := range values[1:] {
		switch kind {
		case 0:
			var ok bool
			if value, ok = xmath.CheckedAdd(value, next); !ok {
				return 0, 0, fmt.Errorf("packet sum exceeds uint64")
			}
		case 1:
			var ok bool
			if value, ok = xmath.CheckedMul(value, next); !ok {
				return 0, 0, fmt.Errorf("packet product exceeds uint64")
			}
		case 2:
			value = min(value, next)
		case 3:
			value = max(value, next)
		case 5, 6, 7:
			match := kind == 5 && value > next || kind == 6 && value < next || kind == 7 && value == next
			value = 0
			if match {
				value = 1
			}
		}
	}
	return version, value, nil
}

func decodeTransmission(r io.Reader) (uint64, uint64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, 0, err
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("empty transmission")
	}
	for i, digit := range data {
		switch {
		case digit >= '0' && digit <= '9':
			data[i] = digit - '0'
		case digit >= 'A' && digit <= 'F':
			data[i] = digit - 'A' + 10
		case digit >= 'a' && digit <= 'f':
			data[i] = digit - 'a' + 10
		default:
			return 0, 0, fmt.Errorf("invalid hexadecimal digit %q", digit)
		}
	}
	bits := packetBits{nibbles: data, end: len(data) * 4}
	version, value, err := bits.packet()
	if err != nil {
		return 0, 0, err
	}
	for bits.pos < bits.end {
		padding, _ := bits.read(1)
		if padding != 0 {
			return 0, 0, fmt.Errorf("nonzero data after outer packet")
		}
	}
	return version, value, nil
}

func day16p01(r io.Reader) (string, error) {
	version, _, err := decodeTransmission(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(version, 10), nil
}

func day16p02(r io.Reader) (string, error) {
	_, value, err := decodeTransmission(r)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(value, 10), nil
}
