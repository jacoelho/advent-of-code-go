package aoc2019

import (
	"fmt"
	"io"
	"strconv"
)

const beamScanSize = 50

func checkTractorBeam(program []int, x, y int) (bool, error) {
	computer := New(program)
	computer.SetInput(x, y)
	if err := computer.RunToHalt(); err != nil {
		return false, fmt.Errorf("beam query at (%d,%d): %w", x, y, err)
	}
	output := computer.GetOutput()
	if len(output) != 1 || output[0] < 0 || output[0] > 1 {
		return false, fmt.Errorf("beam query at (%d,%d) needs one zero/one reply", x, y)
	}
	return output[0] == 1, nil
}

func day19p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	count := 0
	for y := range beamScanSize {
		for x := range beamScanSize {
			pulled, err := checkTractorBeam(program, x, y)
			if err != nil {
				return "", err
			}
			if pulled {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

// beamSample returns a pulled point below the emitter row, searching the
// scanned area from its bottom row up. The beam is a cone from the emitter,
// so the ray through the sample stays inside it.
func beamSample(program []int) (int, int, error) {
	for y := beamScanSize - 1; y > 0; y-- {
		for x := range beamScanSize {
			pulled, err := checkTractorBeam(program, x, y)
			if err != nil {
				return 0, 0, err
			}
			if pulled {
				return x, y, nil
			}
		}
	}
	return 0, 0, fmt.Errorf("no beam beyond the emitter in the %dx%d scan", beamScanSize, beamScanSize)
}

// closestSquare returns the top-left corner of the square of the given size
// closest to the emitter that fits in the beam. It walks the beam's left edge
// row by row; the edge never moves left as rows go down.
func closestSquare(program []int, size int) (int, int, error) {
	// The answer encodes the corner as x*10000+y, so y must stay below 10000.
	const maxTop = 10_000
	sampleX, sampleY, err := beamSample(program)
	if err != nil {
		return 0, 0, err
	}
	left := 0
	for bottom := size - 1; bottom-size+1 < maxTop; bottom++ {
		// The sample ray crosses this row before that column, bounding where the
		// left edge can be.
		limit := sampleX*bottom/sampleY + 1
		x := left
		for ; x <= limit; x++ {
			pulled, err := checkTractorBeam(program, x, bottom)
			if err != nil {
				return 0, 0, err
			}
			if pulled {
				break
			}
		}
		if x > limit {
			continue
		}
		left = x
		top := bottom - size + 1
		fits, err := checkTractorBeam(program, x+size-1, top)
		if err != nil {
			return 0, 0, err
		}
		if fits {
			return x, top, nil
		}
	}
	return 0, 0, fmt.Errorf("no %dx%d square fits in the beam above row %d", size, size, maxTop)
}

func day19p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	x, y, err := closestSquare(program, 100)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(x*10000 + y), nil
}
