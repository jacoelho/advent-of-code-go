package aoc2016

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type compressionVersion uint8

const (
	flatCompression compressionVersion = iota
	recursiveCompression
)

func decompressedLength(r io.Reader, version compressionVersion) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	data = bytes.Join(bytes.Fields(data), nil)
	type section struct {
		end                 int
		repetitions, length big.Int
	}
	stack := []section{{end: len(data), repetitions: *big.NewInt(1)}}
	one := big.NewInt(1)
	for position := 0; ; {
		current := &stack[len(stack)-1]
		if position == current.end {
			if len(stack) == 1 {
				return current.length.String(), nil
			}
			size := new(big.Int).Mul(&current.length, &current.repetitions)
			stack = stack[:len(stack)-1]
			stack[len(stack)-1].length.Add(&stack[len(stack)-1].length, size)
			continue
		}
		if data[position] != '(' {
			current.length.Add(&current.length, one)
			position++
			continue
		}
		closing := bytes.IndexByte(data[position:current.end], ')')
		if closing < 0 {
			return "", fmt.Errorf("unterminated marker at byte %d", position)
		}
		closing += position
		widthText, repeatText, ok := bytes.Cut(data[position+1:closing], []byte("x"))
		width, widthErr := convert.ScanNumber[int](widthText)
		repeat, repeatOK := new(big.Int).SetString(string(repeatText), 10)
		if !ok || widthErr != nil || width < 0 || !repeatOK || repeat.Sign() < 0 || width > current.end-closing-1 {
			return "", fmt.Errorf("invalid marker at byte %d", position)
		}
		position = closing + 1
		if version == flatCompression {
			size := new(big.Int).Mul(big.NewInt(int64(width)), repeat)
			current.length.Add(&current.length, size)
			position += width
		} else {
			stack = append(stack, section{end: position + width, repetitions: *repeat})
		}
	}
}

func day09p01(r io.Reader) (string, error) {
	return decompressedLength(r, flatCompression)
}

func day09p02(r io.Reader) (string, error) {
	return decompressedLength(r, recursiveCompression)
}
