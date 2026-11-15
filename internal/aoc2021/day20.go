package aoc2021

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

// litImage is a finite window of pixels inside an infinite plane of
// background pixels.
type litImage struct {
	pixels        []bool
	width, height int
	background    bool
}

func (image litImage) lit(x, y int) bool {
	if x < 0 || y < 0 || x >= image.width || y >= image.height {
		return image.background
	}
	return image.pixels[y*image.width+x]
}

// enhance grows the window by one pixel on each side: only those pixels have
// a 3x3 neighbourhood that touches the window.
func (image litImage) enhance(algorithm *[512]bool) litImage {
	next := litImage{
		width:      image.width + 2,
		height:     image.height + 2,
		background: algorithm[0],
	}
	if image.background {
		next.background = algorithm[511]
	}
	next.pixels = make([]bool, next.width*next.height)
	for y := range next.height {
		for x := range next.width {
			index := 0
			for dy := -2; dy <= 0; dy++ {
				for dx := -2; dx <= 0; dx++ {
					index <<= 1
					if image.lit(x+dx, y+dy) {
						index |= 1
					}
				}
			}
			next.pixels[y*next.width+x] = algorithm[index]
		}
	}
	return next
}

func enhanceImage(r io.Reader, steps int) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	header, imageText, ok := strings.Cut(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n\n")
	if !ok {
		return "", fmt.Errorf("missing blank line between enhancement algorithm and image")
	}
	var algorithm [512]bool
	if len(header) != len(algorithm) {
		return "", fmt.Errorf("image enhancement algorithm must have 512 pixels")
	}
	for i, pixel := range []byte(header) {
		if pixel != '.' && pixel != '#' {
			return "", fmt.Errorf("invalid algorithm pixel %q", pixel)
		}
		algorithm[i] = pixel == '#'
	}
	input, err := grid.ReadBytes(strings.NewReader(imageText), ".#")
	if err != nil {
		return "", fmt.Errorf("image: %w", err)
	}
	image := litImage{width: input.Width(), height: input.Height()}
	for _, pixel := range input.Cells() {
		image.pixels = append(image.pixels, pixel == '#')
	}

	for range steps {
		image = image.enhance(&algorithm)
	}
	if image.background {
		return "", fmt.Errorf("enhanced image contains infinitely many lit pixels")
	}
	lit := 0
	for _, pixel := range image.pixels {
		if pixel {
			lit++
		}
	}
	return strconv.Itoa(lit), nil
}

func day20p01(r io.Reader) (string, error) {
	return enhanceImage(r, 2)
}

func day20p02(r io.Reader) (string, error) {
	return enhanceImage(r, 50)
}
