package aoc2019

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

const (
	imageWidth  = 25
	imageHeight = 6
)

func parseImageLayers(r io.Reader, width, height int) ([][]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSpace(data)
	layerSize := width * height
	if len(data) == 0 || len(data)%layerSize != 0 {
		return nil, fmt.Errorf("image has %d pixels, want a positive multiple of %d", len(data), layerSize)
	}
	if i := bytes.IndexFunc(data, func(c rune) bool { return c < '0' || c > '2' }); i >= 0 {
		return nil, fmt.Errorf("invalid pixel %q at %d", data[i], i)
	}

	var layers [][]byte
	for i := 0; i < len(data); i += layerSize {
		layers = append(layers, data[i:i+layerSize])
	}
	return layers, nil
}

func day08p01(r io.Reader) (string, error) {
	layers, err := parseImageLayers(r, imageWidth, imageHeight)
	if err != nil {
		return "", err
	}

	minZeros := math.MaxInt
	checksum := 0
	for _, layer := range layers {
		var counts [3]int
		for _, pixel := range layer {
			counts[pixel-'0']++
		}
		if counts[0] < minZeros {
			minZeros = counts[0]
			checksum = counts[1] * counts[2]
		}
	}
	return strconv.Itoa(checksum), nil
}

func day08p02(r io.Reader) (string, error) {
	layers, err := parseImageLayers(r, imageWidth, imageHeight)
	if err != nil {
		return "", err
	}

	rows := make([]string, imageHeight)
	for y := range imageHeight {
		var row strings.Builder
		for x := range imageWidth {
			i := y*imageWidth + x
			pixel := byte('2')
			for _, layer := range layers {
				if pixel = layer[i]; pixel != '2' {
					break
				}
			}
			switch pixel {
			case '1':
				row.WriteByte('#')
			case '0':
				row.WriteByte(' ')
			default:
				return "", fmt.Errorf("pixel (%d,%d) is transparent in every layer", x, y)
			}
		}
		rows[y] = row.String()
	}
	return aoc.Decode6RowDisplay(rows, 5)
}
