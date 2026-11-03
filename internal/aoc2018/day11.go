package aoc2018

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func fuelPowerGrid(r io.Reader) ([301][301]int, error) {
	var totals [301][301]int
	data, err := io.ReadAll(r)
	if err != nil {
		return totals, err
	}

	serial, err := convert.ScanNumber[int32](bytes.TrimSpace(data))
	if err != nil || serial < 0 {
		return totals, fmt.Errorf("invalid fuel-cell serial number %q", data)
	}

	for y := 1; y <= 300; y++ {
		for x := 1; x <= 300; x++ {
			rack := int64(x + 10)
			power := int((rack*int64(y)+int64(serial))*rack/100%10) - 5
			totals[y][x] = power + totals[y-1][x] + totals[y][x-1] - totals[y-1][x-1]
		}
	}

	return totals, nil
}

func bestFuelSquare(totals *[301][301]int, minSize, maxSize int) (int, int, int) {
	best, bestX, bestY, bestSize := math.MinInt, 0, 0, 0
	for size := minSize; size <= maxSize; size++ {
		for y := size; y <= 300; y++ {
			for x := size; x <= 300; x++ {
				power := totals[y][x] - totals[y-size][x] - totals[y][x-size] + totals[y-size][x-size]
				if power > best {
					best, bestX, bestY, bestSize = power, x-size+1, y-size+1, size
				}
			}
		}
	}

	return bestX, bestY, bestSize
}

func day11p01(r io.Reader) (string, error) {
	totals, err := fuelPowerGrid(r)
	if err != nil {
		return "", err
	}

	x, y, _ := bestFuelSquare(&totals, 3, 3)
	return fmt.Sprintf("%d,%d", x, y), nil
}

func day11p02(r io.Reader) (string, error) {
	totals, err := fuelPowerGrid(r)
	if err != nil {
		return "", err
	}

	x, y, size := bestFuelSquare(&totals, 1, 300)
	return fmt.Sprintf("%d,%d,%d", x, y, size), nil
}
