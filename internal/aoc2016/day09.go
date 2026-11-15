package aoc2016

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

var errDecompressedOverflow = errors.New("decompressed length overflows int")

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
	type section struct{ end, repetitions, length int }
	stack := []section{{end: len(data), repetitions: 1}}
	for position := 0; ; {
		current := &stack[len(stack)-1]
		if position == current.end {
			if len(stack) == 1 {
				return strconv.Itoa(current.length), nil
			}
			size, sizeOK := xmath.CheckedMul(current.length, current.repetitions)
			stack = stack[:len(stack)-1]
			parent := &stack[len(stack)-1]
			length, lengthOK := xmath.CheckedAdd(parent.length, size)
			if !sizeOK || !lengthOK {
				return "", errDecompressedOverflow
			}
			parent.length = length
			continue
		}
		if data[position] != '(' {
			length, ok := xmath.CheckedAdd(current.length, 1)
			if !ok {
				return "", errDecompressedOverflow
			}
			current.length = length
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
		repeat, repeatErr := convert.ScanNumber[int](repeatText)
		if !ok || widthErr != nil || repeatErr != nil || width < 0 || repeat < 0 || width > current.end-closing-1 {
			return "", fmt.Errorf("invalid marker at byte %d", position)
		}
		position = closing + 1
		if version == recursiveCompression {
			stack = append(stack, section{end: position + width, repetitions: repeat})
			continue
		}
		size, sizeOK := xmath.CheckedMul(width, repeat)
		length, lengthOK := xmath.CheckedAdd(current.length, size)
		if !sizeOK || !lengthOK {
			return "", errDecompressedOverflow
		}
		current.length = length
		position += width
	}
}

func day09p01(r io.Reader) (string, error) {
	return decompressedLength(r, flatCompression)
}

func day09p02(r io.Reader) (string, error) {
	return decompressedLength(r, recursiveCompression)
}
