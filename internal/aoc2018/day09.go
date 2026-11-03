package aoc2018

import (
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
)

func marbleGame(r io.Reader, multiplier int) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	var players, last int
	if _, err := fmt.Sscanf(string(data), "%d players; last marble is worth %d points", &players, &last); err != nil {
		return 0, fmt.Errorf("invalid marble game: %w", err)
	}
	if players <= 0 || last < 0 || last > math.MaxInt/multiplier {
		return 0, fmt.Errorf("invalid player count or last marble")
	}
	last *= multiplier

	scores := make([]int64, players)
	circle := collections.NewDeque[int](16)
	// The current marble is at the back, with clockwise order from front to back.
	circle.PushBack(0)
	for turn := range last {
		marble := turn + 1
		if marble%23 == 0 {
			for range 7 {
				value, _ := circle.PopBack()
				circle.PushFront(value)
			}
			removed, _ := circle.PopBack()
			scores[(marble-1)%players] += int64(marble) + int64(removed)
			value, _ := circle.PopFront()
			circle.PushBack(value)
		} else {
			value, _ := circle.PopFront()
			circle.PushBack(value)
			circle.PushBack(marble)
		}
	}
	return slices.Max(scores), nil
}

func day09p01(r io.Reader) (string, error) {
	score, err := marbleGame(r, 1)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(score, 10), nil
}

func day09p02(r io.Reader) (string, error) {
	score, err := marbleGame(r, 100)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(score, 10), nil
}
