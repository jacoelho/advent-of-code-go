package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

type mapping struct {
	destStart   int
	sourceStart int
	length      int
}

type interval struct {
	start int
	end   int
}

func (m mapping) apply(value int) (int, bool) {
	if value >= m.sourceStart && value < m.sourceStart+m.length {
		offset := value - m.sourceStart
		return m.destStart + offset, true
	}
	return value, false
}

func applyMappings(value int, mappings []mapping) int {
	for _, m := range mappings {
		if result, ok := m.apply(value); ok {
			return result
		}
	}
	return value
}

func parseAlmanac(r io.Reader) ([]int, [][]mapping, error) {
	blocks, err := scanner.ReadBlocks(r)
	if err != nil {
		return nil, nil, err
	}
	if len(blocks) == 0 || len(blocks[0]) != 1 {
		return nil, nil, fmt.Errorf("almanac must start with a seeds line")
	}
	seedsLine, ok := strings.CutPrefix(blocks[0][0], "seeds:")
	if !ok {
		return nil, nil, fmt.Errorf("almanac must start with a seeds line")
	}
	seeds, err := parseNaturals(seedsLine)
	if err != nil {
		return nil, nil, fmt.Errorf("seeds: %w", err)
	}
	if len(seeds) == 0 {
		return nil, nil, fmt.Errorf("almanac has no seeds")
	}

	allMappings := make([][]mapping, 0, len(blocks)-1)
	for _, block := range blocks[1:] {
		if len(block) < 2 || !strings.HasSuffix(block[0], " map:") {
			return nil, nil, fmt.Errorf("almanac map needs a header and ranges")
		}
		header := block[0]
		var mappings []mapping
		for _, line := range block[1:] {
			values, err := parseNaturals(line)
			if err != nil || len(values) != 3 {
				return nil, nil, fmt.Errorf("%s invalid range %q", header, line)
			}
			mappings = append(mappings, mapping{destStart: values[0], sourceStart: values[1], length: values[2]})
		}
		allMappings = append(allMappings, mappings)
	}

	return seeds, allMappings, nil
}

func day05p01(r io.Reader) (string, error) {
	seeds, allMappings, err := parseAlmanac(r)
	if err != nil {
		return "", err
	}

	minLocation := -1
	for _, seed := range seeds {
		value := seed
		for _, mappings := range allMappings {
			value = applyMappings(value, mappings)
		}
		if minLocation == -1 || value < minLocation {
			minLocation = value
		}
	}

	return strconv.Itoa(minLocation), nil
}

func applyMappingsToIntervals(intervals []interval, mappings []mapping) []interval {
	var result []interval

	for _, iv := range intervals {
		unmapped := []interval{iv}

		for _, m := range mappings {
			var nextUnmapped []interval
			mapStart := m.sourceStart
			mapEnd := m.sourceStart + m.length
			offset := m.destStart - m.sourceStart

			for _, u := range unmapped {
				if u.end <= mapStart || u.start >= mapEnd {
					nextUnmapped = append(nextUnmapped, u)
					continue
				}

				if u.start < mapStart {
					nextUnmapped = append(nextUnmapped, interval{u.start, mapStart})
				}

				overlapStart := max(u.start, mapStart)
				overlapEnd := min(u.end, mapEnd)
				result = append(result, interval{
					start: overlapStart + offset,
					end:   overlapEnd + offset,
				})

				if u.end > mapEnd {
					nextUnmapped = append(nextUnmapped, interval{mapEnd, u.end})
				}
			}

			unmapped = nextUnmapped
		}

		result = append(result, unmapped...)
	}

	return result
}

func day05p02(r io.Reader) (string, error) {
	seeds, allMappings, err := parseAlmanac(r)
	if err != nil {
		return "", err
	}
	if len(seeds)%2 != 0 {
		return "", fmt.Errorf("seed ranges need start and length pairs, got %d values", len(seeds))
	}

	var intervals []interval
	for i := 0; i < len(seeds); i += 2 {
		if seeds[i+1] > 0 {
			intervals = append(intervals, interval{seeds[i], seeds[i] + seeds[i+1]})
		}
	}
	if len(intervals) == 0 {
		return "", fmt.Errorf("seed ranges are all empty")
	}

	for _, mappings := range allMappings {
		intervals = applyMappingsToIntervals(intervals, mappings)
	}

	minLocation := slices.MinFunc(intervals, func(a, b interval) int {
		return cmp.Compare(a.start, b.start)
	})

	return strconv.Itoa(minLocation.start), nil
}
