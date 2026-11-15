package aoc2024

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/xmath"
)

const (
	day20Wall      = -1
	day20Unvisited = -2
)

type day20RaceTrack struct {
	width, height int
	// distance holds the step count from the start for each track cell,
	// indexed by y*width+x, and day20Wall for walls.
	distance []int32
	// path lists track cell indices in race order.
	path []int
}

func parseRaceTrack(r io.Reader) (day20RaceTrack, error) {
	g, err := grid.ReadBytes(r, "#.SE")
	if err != nil {
		return day20RaceTrack{}, fmt.Errorf("race track: %w", err)
	}
	start, err := g.Find('S')
	if err != nil {
		return day20RaceTrack{}, fmt.Errorf("race track: %w", err)
	}
	end, err := g.Find('E')
	if err != nil {
		return day20RaceTrack{}, fmt.Errorf("race track: %w", err)
	}
	track := day20RaceTrack{width: g.Width(), height: g.Height(), distance: make([]int32, len(g.Cells()))}
	open := 0
	for i, c := range g.Cells() {
		if c == '#' {
			track.distance[i] = day20Wall
			continue
		}
		track.distance[i] = day20Unvisited
		open++
	}
	if err := track.walk(g, start, end, open); err != nil {
		return day20RaceTrack{}, err
	}
	return track, nil
}

// walk records the distance of each track cell, requiring the track to be a
// single path from start to end through every open cell.
func (t *day20RaceTrack) walk(g grid.Bytes, start, end grid.Position2D[int], open int) error {
	current, last := g.Index(start), g.Index(end)
	t.distance[current] = 0
	t.path = append(t.path, current)
	for current != last {
		next := -1
		p := g.Position(current)
		for neighbour := range g.Neighbours4(p) {
			if n := g.Index(neighbour); t.distance[n] == day20Unvisited {
				if next >= 0 {
					return fmt.Errorf("race track branches at %d,%d", p.X, p.Y)
				}
				next = n
			}
		}
		if next < 0 {
			return fmt.Errorf("race track has no route to the end")
		}
		t.distance[next] = t.distance[current] + 1
		t.path = append(t.path, next)
		current = next
	}
	if len(t.path) != open {
		return fmt.Errorf("race track has %d open cells off the path", open-len(t.path))
	}
	return nil
}

func day20(cheatDistance, saveGoal int) func(r io.Reader) (string, error) {
	return func(r io.Reader) (string, error) {
		track, err := parseRaceTrack(r)
		if err != nil {
			return "", err
		}
		var count int
		for _, from := range track.path {
			x, y := from%track.width, from/track.width
			fromDistance := int(track.distance[from])
			for dy := max(-cheatDistance, -y); dy <= min(cheatDistance, track.height-1-y); dy++ {
				reach := cheatDistance - xmath.Abs(dy)
				for dx := max(-reach, -x); dx <= min(reach, track.width-1-x); dx++ {
					to := track.distance[(y+dy)*track.width+x+dx]
					if to >= 0 && int(to)-fromDistance-xmath.Abs(dx)-xmath.Abs(dy) >= saveGoal {
						count++
					}
				}
			}
		}
		return strconv.Itoa(count), nil
	}
}

func day20p01(saveGoal int) func(r io.Reader) (string, error) {
	return day20(2, saveGoal)
}

func day20p02(saveGoal int) func(r io.Reader) (string, error) {
	return day20(20, saveGoal)
}
