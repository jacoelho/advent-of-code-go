package aoc2017

import (
	"bufio"
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/search"
)

func parseProgramPipes(r io.Reader) (map[int][]int, error) {
	pipes := make(map[int][]int)
	s := bufio.NewScanner(r)
	for s.Scan() {
		from, destinations, ok := strings.Cut(s.Text(), " <-> ")
		if !ok {
			return nil, fmt.Errorf("invalid program pipes %q", s.Text())
		}
		id, err := convert.ScanNumber[int]([]byte(from))
		if err != nil {
			return nil, err
		}
		if id < 0 {
			return nil, fmt.Errorf("negative program ID")
		}
		for destination := range strings.SplitSeq(destinations, ",") {
			other, err := convert.ScanNumber[int]([]byte(strings.TrimSpace(destination)))
			if err != nil {
				return nil, err
			}
			if other < 0 {
				return nil, fmt.Errorf("negative program ID")
			}
			pipes[id] = append(pipes[id], other)
			pipes[other] = append(pipes[other], id)
		}
	}

	return pipes, s.Err()
}

func day12p01(r io.Reader) (string, error) {
	pipes, err := parseProgramPipes(r)
	if err != nil {
		return "", err
	}
	if _, exists := pipes[0]; !exists {
		return "", fmt.Errorf("program 0 is missing")
	}

	count := 0
	for range search.BFS(0, func(id int) iter.Seq[int] { return slices.Values(pipes[id]) }) {
		count++
	}

	return strconv.Itoa(count), nil
}

func day12p02(r io.Reader) (string, error) {
	pipes, err := parseProgramPipes(r)
	if err != nil {
		return "", err
	}

	visited := collections.NewSet[int]()
	groups := 0
	for id := range pipes {
		if visited.Contains(id) {
			continue
		}
		groups++
		for range search.BFSWithVisited(id, visited, func(id int) iter.Seq[int] { return slices.Values(pipes[id]) }) {
		}
	}

	return strconv.Itoa(groups), nil
}
