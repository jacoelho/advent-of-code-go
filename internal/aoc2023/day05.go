package aoc2023

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/aoc"
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
	s := scanner.NewScannerWithSplit(r, scanner.SplitBySeparator([]byte("\n\n")), func(b []byte) (string, error) {
		return string(b), nil
	})

	sections := slices.Collect(s.Values())
	if s.Err() != nil {
		return nil, nil, s.Err()
	}

	if len(sections) == 0 {
		return nil, nil, nil
	}

	seedsLine := strings.TrimPrefix(sections[0], "seeds: ")
	seedParts := strings.Fields(seedsLine)
	seeds := make([]int, len(seedParts))
	for i, part := range seedParts {
		seeds[i] = aoc.MustAtoi(part)
	}

	var allMappings [][]mapping
	for _, section := range sections[1:] {
		var mappings []mapping
		for _, line := range strings.Split(section, "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.Contains(line, "map:") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) == 3 {
				mappings = append(mappings, mapping{
					destStart:   aoc.MustAtoi(parts[0]),
					sourceStart: aoc.MustAtoi(parts[1]),
					length:      aoc.MustAtoi(parts[2]),
				})
			}
		}
		if len(mappings) > 0 {
			allMappings = append(allMappings, mappings)
		}
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

	var intervals []interval
	for i := 0; i < len(seeds); i += 2 {
		start := seeds[i]
		length := seeds[i+1]
		intervals = append(intervals, interval{start, start + length})
	}

	for _, mappings := range allMappings {
		intervals = applyMappingsToIntervals(intervals, mappings)
	}

	minLocation := slices.MinFunc(intervals, func(a, b interval) int {
		return a.start - b.start
	})

	return strconv.Itoa(minLocation.start), nil
}
