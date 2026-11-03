package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type racingReindeer struct {
	speed, flying, resting int64
}

func parseRacingReindeer(r io.Reader) ([]racingReindeer, error) {
	pattern := regexp.MustCompile(`^(\w+) can fly (\d+) km/s for (\d+) seconds, but then must rest for (\d+) seconds\.$`)
	var reindeer []racingReindeer
	names := make(map[string]bool)
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := pattern.FindStringSubmatch(s.Text())
		if match == nil || names[match[1]] {
			return nil, fmt.Errorf("invalid or duplicate reindeer %q", s.Text())
		}
		names[match[1]] = true
		var values [3]int64
		for i := range values {
			value, err := strconv.ParseInt(match[i+2], 10, 32)
			if err != nil || value <= 0 {
				return nil, fmt.Errorf("invalid reindeer speed or duration %q", match[i+2])
			}
			values[i] = value
		}
		reindeer = append(reindeer, racingReindeer{values[0], values[1], values[2]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(reindeer) == 0 {
		return nil, fmt.Errorf("empty reindeer race")
	}
	return reindeer, nil
}

func (r racingReindeer) distance(seconds int) int64 {
	time := int64(seconds)
	cycle := r.flying + r.resting
	return (time/cycle*r.flying + min(time%cycle, r.flying)) * r.speed
}

func reindeerDistance(r io.Reader, seconds int) (string, error) {
	reindeer, err := parseRacingReindeer(r)
	if err != nil {
		return "", err
	}
	var farthest int64
	for _, racer := range reindeer {
		farthest = max(farthest, racer.distance(seconds))
	}
	return strconv.FormatInt(farthest, 10), nil
}

func reindeerPoints(r io.Reader, seconds int) (string, error) {
	reindeer, err := parseRacingReindeer(r)
	if err != nil {
		return "", err
	}
	points := make([]int, len(reindeer))
	best := 0
	for second := 1; second <= seconds; second++ {
		var farthest int64
		for _, racer := range reindeer {
			farthest = max(farthest, racer.distance(second))
		}
		for i, racer := range reindeer {
			if racer.distance(second) == farthest {
				points[i]++
				best = max(best, points[i])
			}
		}
	}
	return strconv.Itoa(best), nil
}

func day14p01(r io.Reader) (string, error) {
	return reindeerDistance(r, 2503)
}

func day14p02(r io.Reader) (string, error) {
	return reindeerPoints(r, 2503)
}
