package aoc2019

import (
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
)

func parseImageLayers(r io.Reader, width, height int) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	input := strings.TrimSpace(string(data))
	layerSize := width * height

	var layers []string
	for i := 0; i < len(input); i += layerSize {
		if i+layerSize <= len(input) {
			layers = append(layers, input[i:i+layerSize])
		}
	}

	return layers, nil
}

func countDigits(layer string) map[rune]int {
	freq := make(map[rune]int)
	for _, ch := range layer {
		freq[ch]++
	}
	return freq
}

func day08p01(r io.Reader) (string, error) {
	layers, err := parseImageLayers(r, 25, 6)
	if err != nil {
		return "", err
	}

	minZeros := math.MaxInt
	checksum := 0

	for _, layer := range layers {
		freq := countDigits(layer)
		if freq['0'] < minZeros {
			minZeros = freq['0']
			checksum = freq['1'] * freq['2']
		}
	}

	return strconv.Itoa(checksum), nil
}

func decodeImage(layers []string, width, height int) []byte {
	size := width * height
	result := make([]byte, size)

	for i := range size {
		for _, layer := range layers {
			if layer[i] != '2' {
				result[i] = layer[i]
				break
			}
		}
	}

	return result
}

func renderImage(pixels []byte, width, height int) string {
	var sb strings.Builder
	for row := range height {
		for col := range width {
			idx := row*width + col
			if pixels[idx] == '1' {
				sb.WriteByte('#')
			} else {
				sb.WriteByte(' ')
			}
		}
		if row < height-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func day08p02(r io.Reader) (string, error) {
	layers, err := parseImageLayers(r, 25, 6)
	if err != nil {
		return "", err
	}

	image := decodeImage(layers, 25, 6)
	rendered := renderImage(image, 25, 6)
	return aoc.Decode6RowDisplay(strings.Split(rendered, "\n"), 5)
}
