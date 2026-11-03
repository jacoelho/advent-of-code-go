package aoc2018

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

func guardSleep(r io.Reader) (map[int][60]int, error) {
	s := bufio.NewScanner(r)
	var records []string
	for s.Scan() {
		records = append(records, s.Text())
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	slices.Sort(records)

	minutes := make(map[int][60]int)
	guard := 0
	var asleep time.Time
	for _, record := range records {
		stamp, event, ok := strings.Cut(record, "] ")
		if !ok || !strings.HasPrefix(stamp, "[") {
			return nil, fmt.Errorf("invalid guard record %q", record)
		}
		at, err := time.Parse("2006-01-02 15:04", stamp[1:])
		if err != nil {
			return nil, err
		}
		switch event {
		case "falls asleep":
			if guard == 0 || !asleep.IsZero() || at.Hour() != 0 {
				return nil, fmt.Errorf("unexpected sleep at %s", stamp)
			}
			asleep = at
		case "wakes up":
			if asleep.IsZero() || at.Hour() != 0 || !at.After(asleep) || at.Sub(asleep) >= time.Hour {
				return nil, fmt.Errorf("unexpected waking at %s", stamp)
			}
			counts := minutes[guard]
			for minute := asleep.Minute(); minute < at.Minute(); minute++ {
				counts[minute]++
			}
			minutes[guard] = counts
			asleep = time.Time{}
		default:
			id, ok := strings.CutPrefix(event, "Guard #")
			if !ok || !asleep.IsZero() {
				return nil, fmt.Errorf("invalid shift %q", event)
			}
			id, ok = strings.CutSuffix(id, " begins shift")
			if !ok {
				return nil, fmt.Errorf("invalid shift %q", event)
			}
			guard, err = convert.ScanNumber[int]([]byte(id))
			if err != nil {
				return nil, fmt.Errorf("invalid guard ID %q: %w", id, err)
			}
			if guard <= 0 {
				return nil, fmt.Errorf("invalid guard ID %q", id)
			}
			if _, exists := minutes[guard]; !exists {
				minutes[guard] = [60]int{}
			}
		}
	}
	if !asleep.IsZero() {
		return nil, fmt.Errorf("guard never wakes")
	}
	if len(minutes) == 0 {
		return nil, fmt.Errorf("no guard shifts")
	}
	return minutes, nil
}

func day04p01(r io.Reader) (string, error) {
	minutes, err := guardSleep(r)
	if err != nil {
		return "", err
	}

	var sleepiest, longest int
	for _, guard := range slices.Sorted(maps.Keys(minutes)) {
		total := 0
		for _, count := range minutes[guard] {
			total += count
		}
		if total > longest {
			sleepiest, longest = guard, total
		}
	}
	var bestMinute, frequency int
	for minute, count := range minutes[sleepiest] {
		if count > frequency {
			bestMinute, frequency = minute, count
		}
	}
	return strconv.Itoa(sleepiest * bestMinute), nil
}

func day04p02(r io.Reader) (string, error) {
	minutes, err := guardSleep(r)
	if err != nil {
		return "", err
	}

	var result, frequency int
	for _, guard := range slices.Sorted(maps.Keys(minutes)) {
		for minute, count := range minutes[guard] {
			if count > frequency {
				result, frequency = guard*minute, count
			}
		}
	}
	return strconv.Itoa(result), nil
}
