package aoc2015

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
)

func routeDistances(r io.Reader) (shortest, longest int, err error) {
	type route struct {
		from, to string
		distance int
	}
	locations := make(map[string]int)
	var routes []route
	s := bufio.NewScanner(r)
	for s.Scan() {
		var route route
		n, err := fmt.Sscanf(s.Text(), "%s to %s = %d", &route.from, &route.to, &route.distance)
		if err != nil || n != 3 || route.distance <= 0 {
			return 0, 0, fmt.Errorf("invalid route %q", s.Text())
		}
		for _, location := range []string{route.from, route.to} {
			if _, ok := locations[location]; !ok {
				locations[location] = len(locations)
			}
		}
		routes = append(routes, route)
	}
	if err := s.Err(); err != nil {
		return 0, 0, err
	}
	if len(locations) == 0 {
		return 0, 0, fmt.Errorf("no route locations")
	}
	distances := make([][]int, len(locations))
	for i := range distances {
		distances[i] = make([]int, len(locations))
	}
	for _, route := range routes {
		from, to := locations[route.from], locations[route.to]
		distances[from][to], distances[to][from] = route.distance, route.distance
	}
	shortest = math.MaxInt
	visited := make([]bool, len(locations))
	var visit func(current, count, total int)
	visit = func(current, count, total int) {
		if count == len(locations) {
			shortest = min(shortest, total)
			longest = max(longest, total)
			return
		}
		for next, distance := range distances[current] {
			if visited[next] || distance == 0 {
				continue
			}
			visited[next] = true
			visit(next, count+1, total+distance)
			visited[next] = false
		}
	}
	for start := range len(locations) {
		visited[start] = true
		visit(start, 1, 0)
		visited[start] = false
	}
	if shortest == math.MaxInt {
		return 0, 0, fmt.Errorf("no route visits every location")
	}
	return shortest, longest, nil
}

func day09p01(r io.Reader) (string, error) {
	shortest, _, err := routeDistances(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(shortest), nil
}

func day09p02(r io.Reader) (string, error) {
	_, longest, err := routeDistances(r)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(longest), nil
}
