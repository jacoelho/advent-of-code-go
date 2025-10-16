package aoc2023

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/collections"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
	"github.com/jacoelho/advent-of-code-go/pkg/xiter"
	"github.com/jacoelho/advent-of-code-go/pkg/xslices"
)

type brick struct {
	ID    int
	Start grid.Position3D[int]
	End   grid.Position3D[int]
}

func (b brick) MinZ() int {
	return min(b.Start.Z, b.End.Z)
}

func (b brick) MaxZ() int {
	return max(b.Start.Z, b.End.Z)
}

func (b brick) OccupiedPositions() []grid.Position3D[int] {
	var positions []grid.Position3D[int]

	// determine the direction of the brick
	dx := -cmp.Compare(b.Start.X, b.End.X)
	dy := -cmp.Compare(b.Start.Y, b.End.Y)
	dz := -cmp.Compare(b.Start.Z, b.End.Z)

	// generate all positions along the brick
	pos := b.Start
	for {
		positions = append(positions, pos)
		if pos == b.End {
			break
		}
		pos = pos.Add(grid.NewPosition3D(dx, dy, dz))
	}

	return positions
}

func parseBricks(r io.Reader) ([]brick, error) {
	var bricks []brick

	lineScanner := scanner.NewScanner(r, func(bytes []byte) (string, error) {
		return string(bytes), nil
	})

	for id, line := range xiter.Enumerate(lineScanner.Values()) {
		parts := strings.Split(line, "~")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid brick format: %s", line)
		}

		start, err := parsePosition(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid start position: %w", err)
		}

		end, err := parsePosition(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid end position: %w", err)
		}

		bricks = append(bricks, brick{
			ID:    id,
			Start: start,
			End:   end,
		})
	}

	return bricks, nil
}

func parsePosition(s string) (grid.Position3D[int], error) {
	coords := strings.Split(s, ",")
	if len(coords) != 3 {
		return grid.Position3D[int]{}, fmt.Errorf("invalid position format: %s", s)
	}

	x, err := strconv.Atoi(coords[0])
	if err != nil {
		return grid.Position3D[int]{}, err
	}

	y, err := strconv.Atoi(coords[1])
	if err != nil {
		return grid.Position3D[int]{}, err
	}

	z, err := strconv.Atoi(coords[2])
	if err != nil {
		return grid.Position3D[int]{}, err
	}

	return grid.NewPosition3D(x, y, z), nil
}

func setupBricks(r io.Reader) ([]brick, map[int]collections.Set[int], map[int]collections.Set[int], error) {
	bricks, err := parseBricks(r)
	if err != nil {
		return nil, nil, nil, err
	}

	slices.SortFunc(bricks, func(a, b brick) int {
		return cmp.Compare(a.MinZ(), b.MinZ())
	})

	settledBricks, err := simulateFalling(bricks)
	if err != nil {
		return nil, nil, nil, err
	}

	supports, supportedBy := buildSupportGraph(settledBricks)
	return settledBricks, supports, supportedBy, nil
}

func day22p01(r io.Reader) (string, error) {
	settledBricks, supports, supportedBy, err := setupBricks(r)
	if err != nil {
		return "", err
	}

	safeCount := xslices.CountFunc(func(b brick) bool {
		return canSafelyDisintegrate(b.ID, supports, supportedBy)
	}, settledBricks)

	return strconv.Itoa(safeCount), nil
}

func simulateFalling(bricks []brick) ([]brick, error) {
	// track the highest Z at each (x,y) position
	heightMap := make(map[grid.Position2D[int]]int)
	settledBricks := make([]brick, len(bricks))

	for i, b := range bricks {
		// find the highest point this brick can rest on
		maxHeight := 0
		positions := b.OccupiedPositions()

		for _, pos := range positions {
			xy := grid.NewPosition2D(pos.X, pos.Y)
			if height, exists := heightMap[xy]; exists && height > maxHeight {
				maxHeight = height
			}
		}

		// calculate how far down this brick needs to fall
		fallDistance := b.MinZ() - maxHeight - 1

		// create the settled brick
		settledBrick := brick{
			ID:    b.ID,
			Start: grid.NewPosition3D(b.Start.X, b.Start.Y, b.Start.Z-fallDistance),
			End:   grid.NewPosition3D(b.End.X, b.End.Y, b.End.Z-fallDistance),
		}
		settledBricks[i] = settledBrick

		// update height map
		for _, pos := range settledBrick.OccupiedPositions() {
			xy := grid.NewPosition2D(pos.X, pos.Y)
			heightMap[xy] = max(heightMap[xy], pos.Z)
		}
	}

	return settledBricks, nil
}

func buildSupportGraph(bricks []brick) (map[int]collections.Set[int], map[int]collections.Set[int]) {
	supports := make(map[int]collections.Set[int])
	supportedBy := make(map[int]collections.Set[int])

	// initialize sets
	for _, brick := range bricks {
		supports[brick.ID] = collections.NewSet[int]()
		supportedBy[brick.ID] = collections.NewSet[int]()
	}

	// build position to brick mapping
	positionToBrick := make(map[grid.Position3D[int]]int)
	for _, brick := range bricks {
		for _, pos := range brick.OccupiedPositions() {
			positionToBrick[pos] = brick.ID
		}
	}

	// find support relationships
	for _, brick := range bricks {
		for _, pos := range brick.OccupiedPositions() {
			// check the position directly below this brick
			below := grid.NewPosition3D(pos.X, pos.Y, pos.Z-1)
			if belowBrickID, exists := positionToBrick[below]; exists && belowBrickID != brick.ID {
				supports[belowBrickID].Add(brick.ID)
				supportedBy[brick.ID].Add(belowBrickID)
			}
		}
	}

	return supports, supportedBy
}

func canSafelyDisintegrate(brickID int, supports, supportedBy map[int]collections.Set[int]) bool {
	// a brick can be safely disintegrated if every brick it supports
	// has at least 2 supporters (including this brick)
	for supportedBrickID := range supports[brickID].Iter() {
		if supportedBy[supportedBrickID].Len() < 2 {
			return false
		}
	}
	return true
}

func allSupportersHaveFallen(brickID int, supportedBy map[int]collections.Set[int], fallen collections.Set[int]) bool {
	supporters := slices.Collect(supportedBy[brickID].Iter())
	return xslices.Every(fallen.Contains, supporters)
}

func countFallingBricks(brickID int, supports, supportedBy map[int]collections.Set[int]) int {
	fallen := collections.NewSet(brickID)
	queue := collections.NewDeque[int](16)
	queue.PushBack(brickID)

	for current, ok := queue.PopFront(); ok; current, ok = queue.PopFront() {
		for supportedID := range supports[current].Iter() {
			if allSupportersHaveFallen(supportedID, supportedBy, fallen) {
				if !fallen.Contains(supportedID) {
					fallen.Add(supportedID)
					queue.PushBack(supportedID)
				}
			}
		}
	}

	return fallen.Len() - 1 // exclude the original brick
}

func day22p02(r io.Reader) (string, error) {
	settledBricks, supports, supportedBy, err := setupBricks(r)
	if err != nil {
		return "", err
	}

	totalFalling := xslices.Sum(xslices.Map(func(b brick) int {
		return countFallingBricks(b.ID, supports, supportedBy)
	}, settledBricks))

	return strconv.Itoa(totalFalling), nil
}
