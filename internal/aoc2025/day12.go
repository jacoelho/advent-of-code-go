package aoc2025

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type shape struct {
	area int
}

type region struct {
	width  int
	height int
	counts []int
}

func (r region) fitsAll(shapes []shape) bool {
	totalArea := 0
	for i, count := range r.counts {
		totalArea += count * shapes[i].area
	}
	return totalArea <= r.width*r.height
}

func parseSummary(r io.Reader) ([]shape, []region, error) {
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte("\n\n")), func(b []byte) (string, error) {
		return string(b), nil
	})

	sections := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	shapes := make([]shape, 6)
	for i := range 6 {
		shapes[i] = shape{area: strings.Count(sections[i], "#")}
	}

	var regions []region
	lines := strings.SplitSeq(sections[6], "\n")
	for line := range lines {
		nums := convert.ExtractDigits[int](line)
		regions = append(regions, region{
			width:  nums[0],
			height: nums[1],
			counts: nums[2:],
		})
	}

	return shapes, regions, nil
}

func day12p01(r io.Reader) (string, error) {
	shapes, regions, err := parseSummary(r)
	if err != nil {
		return "", err
	}

	count := xslices.CountFunc(func(reg region) bool {
		return reg.fitsAll(shapes)
	}, regions)

	return strconv.Itoa(count), nil
}
