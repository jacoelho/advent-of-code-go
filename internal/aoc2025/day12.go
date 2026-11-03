package aoc2025

import (
	"cmp"
	"fmt"
	"io"
	"math/bits"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
)

const (
	maxGiftShapes = 16
	maxGiftWork   = 100_000_000
	maxGiftStates = 1_000_000
)

type giftOrientation struct {
	cells         []grid.Position2D[int]
	width, height int
}

type giftShape struct {
	orientations []giftOrientation
}

type giftRegion struct {
	width, height int
	counts        []int
}

func giftOrientations(cells []grid.Position2D[int]) []giftOrientation {
	seen := make(map[string]bool)
	var orientations []giftOrientation
	for reflection := range 2 {
		for rotations := range 4 {
			variant := make([]grid.Position2D[int], len(cells))
			minX, minY := 32, 32
			for i, cell := range cells {
				if reflection == 1 {
					cell.X = -cell.X
				}
				for range rotations {
					cell = cell.TurnRight()
				}
				variant[i] = cell
				minX, minY = min(minX, cell.X), min(minY, cell.Y)
			}
			width, height := 0, 0
			for i := range variant {
				variant[i].X -= minX
				variant[i].Y -= minY
				width, height = max(width, variant[i].X+1), max(height, variant[i].Y+1)
			}
			slices.SortFunc(variant, func(a, b grid.Position2D[int]) int {
				return cmp.Or(cmp.Compare(a.Y, b.Y), cmp.Compare(a.X, b.X))
			})
			key := make([]byte, 0, 2*len(variant))
			for _, cell := range variant {
				key = append(key, byte(cell.X), byte(cell.Y))
			}
			if !seen[string(key)] {
				seen[string(key)] = true
				orientations = append(orientations, giftOrientation{cells: variant, width: width, height: height})
			}
		}
	}
	return orientations
}

func parseGiftSummary(r io.Reader) ([]giftShape, []giftRegion, error) {
	const inputLimit = 4 << 20
	data, err := io.ReadAll(io.LimitReader(r, inputLimit+1))
	if err != nil {
		return nil, nil, err
	}
	if len(data) > inputLimit {
		return nil, nil, fmt.Errorf("gift summary exceeds 4 MiB")
	}
	input := strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n"))
	sections := strings.Split(input, "\n\n")
	if len(sections) < 2 || len(sections)-1 > maxGiftShapes {
		return nil, nil, fmt.Errorf("gift summary needs 1 to %d shapes and regions", maxGiftShapes)
	}
	shapes := make([]giftShape, len(sections)-1)
	for _, section := range sections[:len(sections)-1] {
		lines := strings.Split(section, "\n")
		idText, ok := strings.CutSuffix(lines[0], ":")
		if !ok || len(lines) < 2 || len(lines)-1 > 32 {
			return nil, nil, fmt.Errorf("gift shape needs an index and 1 to 32 rows")
		}
		id, err := convert.ScanNumber[int]([]byte(idText))
		if err != nil {
			return nil, nil, fmt.Errorf("gift shape index: %w", err)
		}
		if id < 0 || id >= len(shapes) || shapes[id].orientations != nil {
			return nil, nil, fmt.Errorf("gift shape indices must be unique and contiguous from zero")
		}
		width := len(lines[1])
		if width == 0 || width > 32 {
			return nil, nil, fmt.Errorf("gift shape width must be 1 to 32")
		}
		var cells []grid.Position2D[int]
		for y, line := range lines[1:] {
			if len(line) != width {
				return nil, nil, fmt.Errorf("gift shape %d has inconsistent row width", id)
			}
			for x, cell := range line {
				switch cell {
				case '#':
					cells = append(cells, grid.NewPosition2D(x, y))
				case '.':
				default:
					return nil, nil, fmt.Errorf("invalid gift shape tile %q", cell)
				}
			}
		}
		if len(cells) == 0 {
			return nil, nil, fmt.Errorf("gift shape %d is empty", id)
		}
		shapes[id] = giftShape{orientations: giftOrientations(cells)}
	}

	parseDimension := func(text string) (int, error) {
		value, err := convert.ScanNumber[int]([]byte(text))
		if err != nil {
			return 0, err
		}
		if value < 1 || value > 1_000_000 {
			return 0, fmt.Errorf("gift region dimension must be 1 to 1000000")
		}
		return value, nil
	}
	var regions []giftRegion
	for line := range strings.SplitSeq(sections[len(sections)-1], "\n") {
		if len(regions) == 10_000 {
			return nil, nil, fmt.Errorf("gift summary exceeds 10000 regions")
		}
		dimensions, counts, ok := strings.Cut(line, ":")
		if !ok {
			return nil, nil, fmt.Errorf("gift region needs dimensions and counts")
		}
		widthText, heightText, ok := strings.Cut(dimensions, "x")
		if !ok {
			return nil, nil, fmt.Errorf("gift region needs width x height")
		}
		width, err := parseDimension(widthText)
		if err != nil {
			return nil, nil, fmt.Errorf("gift region width: %w", err)
		}
		height, err := parseDimension(heightText)
		if err != nil {
			return nil, nil, fmt.Errorf("gift region height: %w", err)
		}
		region := giftRegion{width: width, height: height}
		for field := range strings.FieldsSeq(counts) {
			count, err := convert.ScanNumber[int]([]byte(field))
			if err != nil {
				return nil, nil, fmt.Errorf("gift count: %w", err)
			}
			if count < 0 || count > 1_000_000 {
				return nil, nil, fmt.Errorf("gift count must be 0 to 1000000")
			}
			region.counts = append(region.counts, count)
		}
		if len(region.counts) != len(shapes) {
			return nil, nil, fmt.Errorf("gift region needs one count for each of %d shapes", len(shapes))
		}
		regions = append(regions, region)
	}
	return shapes, regions, nil
}

func giftsFitInBoxes(shapes []giftShape, region giftRegion, pieces int64) bool {
	for boxWidth := 1; boxWidth <= min(region.width, 32); boxWidth++ {
		boxHeight := 0
		for i, shape := range shapes {
			if region.counts[i] == 0 {
				continue
			}
			height := 33
			for _, orientation := range shape.orientations {
				if orientation.width <= boxWidth {
					height = min(height, orientation.height)
				}
			}
			boxHeight = max(boxHeight, height)
		}
		if boxHeight <= 32 && int64(region.width/boxWidth)*int64(region.height/boxHeight) >= pieces {
			return true
		}
	}
	return false
}

type giftPackingState struct {
	occupied uint64
	counts   [maxGiftShapes]uint8
}

type giftPackingWork struct {
	checks int
	states int
}

func (w *giftPackingWork) check() error {
	if w.checks == 0 {
		return fmt.Errorf("gift packing exceeds one hundred million placement and state checks")
	}
	w.checks--
	return nil
}

func packGiftRegion(shapes []giftShape, region giftRegion, work *giftPackingWork) (bool, error) {
	var area, pieces int64
	for i, count := range region.counts {
		if count == 0 {
			continue
		}
		pieces += int64(count)
		area += int64(count) * int64(len(shapes[i].orientations[0].cells))
		fits := false
		for _, orientation := range shapes[i].orientations {
			if orientation.width <= region.width && orientation.height <= region.height {
				fits = true
				break
			}
		}
		if !fits {
			return false, nil
		}
	}
	capacity := int64(region.width) * int64(region.height)
	if area > capacity {
		return false, nil
	}
	if pieces == 0 || giftsFitInBoxes(shapes, region, pieces) {
		return true, nil
	}
	if capacity > 64 {
		return false, fmt.Errorf("unsupported gift region: exact placement needs at most 64 cells")
	}

	placements := make([][]uint64, len(shapes))
	initial := giftPackingState{}
	for i, shape := range shapes {
		initial.counts[i] = uint8(region.counts[i])
		if region.counts[i] == 0 {
			continue
		}
		seen := make(map[uint64]bool)
		for _, orientation := range shape.orientations {
			for y := 0; y <= region.height-orientation.height; y++ {
				for x := 0; x <= region.width-orientation.width; x++ {
					if err := work.check(); err != nil {
						return false, err
					}
					var mask uint64
					for _, cell := range orientation.cells {
						mask |= uint64(1) << ((y+cell.Y)*region.width + x + cell.X)
					}
					if !seen[mask] {
						seen[mask] = true
						placements[i] = append(placements[i], mask)
					}
				}
			}
		}
	}

	failed := make(map[giftPackingState]bool)
	reject := func(state giftPackingState) (bool, error) {
		if work.states == 0 {
			return false, fmt.Errorf("gift packing exceeds one million failed states")
		}
		work.states--
		failed[state] = true
		return false, nil
	}
	var search func(giftPackingState) (bool, error)
	search = func(state giftPackingState) (bool, error) {
		if err := work.check(); err != nil {
			return false, err
		}
		if failed[state] {
			return false, nil
		}
		chosen := -1
		var choices []uint64
		var allCells uint64
		requiredArea := 0
		for i, count := range state.counts[:len(shapes)] {
			if count == 0 {
				continue
			}
			required := int(count) * len(shapes[i].orientations[0].cells)
			requiredArea += required
			var legal []uint64
			var cells uint64
			for _, placement := range placements[i] {
				if err := work.check(); err != nil {
					return false, err
				}
				if placement&state.occupied == 0 {
					legal = append(legal, placement)
					cells |= placement
				}
			}
			if len(legal) < int(count) || bits.OnesCount64(cells) < required {
				return reject(state)
			}
			allCells |= cells
			if chosen == -1 || len(legal) < len(choices) {
				chosen, choices = i, legal
			}
		}
		if chosen == -1 {
			return true, nil
		}
		// The remaining pieces must fit in the union of currently legal placements.
		if bits.OnesCount64(allCells) < requiredArea {
			return reject(state)
		}
		for _, placement := range choices {
			next := state
			next.occupied |= placement
			next.counts[chosen]--
			found, err := search(next)
			if err != nil || found {
				return found, err
			}
		}
		return reject(state)
	}
	return search(initial)
}

func day12p01(r io.Reader) (string, error) {
	shapes, regions, err := parseGiftSummary(r)
	if err != nil {
		return "", err
	}
	work := giftPackingWork{checks: maxGiftWork, states: maxGiftStates}
	count := 0
	for i, region := range regions {
		fits, err := packGiftRegion(shapes, region, &work)
		if err != nil {
			return "", fmt.Errorf("gift region %d: %w", i+1, err)
		}
		if fits {
			count++
		}
	}
	return strconv.Itoa(count), nil
}
