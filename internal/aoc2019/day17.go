package aoc2019

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
	"github.com/jacoelho/advent-of-code-go/pkg/grid"
	"github.com/jacoelho/advent-of-code-go/pkg/scanner"
)

func calculateAlignmentSum(scaffoldGrid grid.Grid2D[int, rune]) int {
	alignmentSum := 0
	for pos := range scaffoldGrid {
		if isIntersection(scaffoldGrid, pos) {
			alignmentSum += int(pos.X * pos.Y)
		}
	}
	return alignmentSum
}

func day17p01(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	if err := computer.Run(); err != nil {
		return "", err
	}

	scaffoldGrid := parseScaffoldFromASCII(computer.GetOutput())
	return strconv.Itoa(calculateAlignmentSum(scaffoldGrid)), nil
}

func filterScaffold(g grid.Grid2D[int, rune]) grid.Grid2D[int, rune] {
	result := make(grid.Grid2D[int, rune])
	for pos, ch := range g {
		if ch != '.' {
			result[pos] = ch
		}
	}
	return result
}

func parseScaffoldFromASCII(output []int) grid.Grid2D[int, rune] {
	var rows [][]rune
	var currentRow []rune

	for _, val := range output {
		ch := rune(val)
		if ch == '\n' {
			if len(currentRow) > 0 {
				rows = append(rows, currentRow)
				currentRow = nil
			}
		} else {
			currentRow = append(currentRow, ch)
		}
	}

	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return filterScaffold(grid.NewGrid2D[int](rows))
}

func parseScaffoldGrid(r io.Reader) (grid.Grid2D[int, rune], error) {
	s := scanner.NewScanner(r, func(line []byte) ([]rune, error) {
		return []rune(string(line)), nil
	})

	rows := slices.Collect(s.Values())
	if err := s.Err(); err != nil {
		return nil, err
	}

	return filterScaffold(grid.NewGrid2D[int](rows)), nil
}

func isIntersection(scaffoldGrid grid.Grid2D[int, rune], pos grid.Position2D[int]) bool {
	if !scaffoldGrid.Contains(pos) {
		return false
	}

	neighborCount := 0
	for neighbor := range grid.Neighbours4(pos) {
		if scaffoldGrid.Contains(neighbor) {
			neighborCount++
		}
	}

	return neighborCount == 4
}

func findRobot(scaffoldGrid grid.Grid2D[int, rune]) (grid.Position2D[int], grid.Position2D[int], error) {
	var position, direction grid.Position2D[int]
	found := false
	for pos, ch := range scaffoldGrid {
		var dir grid.Position2D[int]
		switch ch {
		case '^':
			dir = grid.Position2D[int]{X: 0, Y: -1}
		case 'v':
			dir = grid.Position2D[int]{X: 0, Y: 1}
		case '<':
			dir = grid.Position2D[int]{X: -1, Y: 0}
		case '>':
			dir = grid.Position2D[int]{X: 1, Y: 0}
		default:
			continue
		}
		if found {
			return position, direction, fmt.Errorf("multiple scaffold robots")
		}
		position, direction = pos, dir
		found = true
	}
	if !found {
		return position, direction, fmt.Errorf("scaffold robot not found")
	}
	return position, direction, nil
}

func generatePath(scaffoldGrid grid.Grid2D[int, rune]) ([]string, error) {
	pos, dir, err := findRobot(scaffoldGrid)
	if err != nil {
		return nil, err
	}
	seen := make(map[[2]grid.Position2D[int]]bool)
	var path []string

	for {
		state := [2]grid.Position2D[int]{pos, dir}
		if seen[state] {
			return nil, fmt.Errorf("scaffold route repeats position %v facing %v", pos, dir)
		}
		seen[state] = true

		nextPos := pos.Add(dir)
		if scaffoldGrid.Contains(nextPos) {
			steps := 0
			for scaffoldGrid.Contains(nextPos) {
				steps++
				pos = nextPos
				nextPos = pos.Add(dir)
			}
			path = append(path, strconv.Itoa(steps))
			continue
		}

		leftDir := dir.TurnLeft()
		if scaffoldGrid.Contains(pos.Add(leftDir)) {
			path = append(path, "L")
			dir = leftDir
			continue
		}

		rightDir := dir.TurnRight()
		if scaffoldGrid.Contains(pos.Add(rightDir)) {
			path = append(path, "R")
			dir = rightDir
			continue
		}

		break
	}

	return path, nil
}

type movementProgram struct {
	main      string
	functions [3]string
}

func compressPath(path []string) (movementProgram, error) {
	const (
		lineLimit   = 20
		callLimit   = 10
		searchLimit = 1_000_000
	)

	// Ten calls to functions of at most five pairs cannot encode a longer route.
	if len(path) == 0 || len(path)%2 != 0 || len(path) > 100 {
		return movementProgram{}, fmt.Errorf("movement path must contain 1 to 50 turn/count pairs")
	}
	for i := 0; i < len(path); i += 2 {
		if path[i] != "L" && path[i] != "R" {
			return movementProgram{}, fmt.Errorf("invalid movement turn %q", path[i])
		}
		count := path[i+1]
		if len(count) == 0 || len(count)+2 > lineLimit || count[0] < '0' || count[0] > '9' {
			return movementProgram{}, fmt.Errorf("invalid movement count %q", count)
		}
		steps, err := convert.ScanNumber[int]([]byte(count))
		if err != nil || steps <= 0 {
			return movementProgram{}, fmt.Errorf("invalid movement count %q", count)
		}
	}
	if len(strings.Join(path, ",")) > callLimit*lineLimit+callLimit-1 {
		return movementProgram{}, fmt.Errorf("movement path exceeds the maximum expanded routine length")
	}

	var functions [3][]string
	main := make([]string, 0, callLimit)
	visited := 0
	exhausted := false

	var search func(int) bool
	search = func(offset int) bool {
		if exhausted {
			return false
		}
		visited++
		if visited > searchLimit {
			exhausted = true
			return false
		}
		if offset == len(path) {
			return true
		}
		if len(main) == callLimit {
			return false
		}

		nextFunction := -1
		for i, function := range functions {
			if len(function) == 0 {
				nextFunction = i
				break
			}
			end := offset + len(function)
			if end > len(path) || !slices.Equal(path[offset:end], function) {
				continue
			}
			main = append(main, "ABC"[i:i+1])
			if search(end) {
				return true
			}
			main = main[:len(main)-1]
		}
		if nextFunction < 0 {
			return false
		}

		maxEnd := offset
		for end := offset + 2; end <= len(path); end += 2 {
			if len(strings.Join(path[offset:end], ",")) > lineLimit {
				break
			}
			maxEnd = end
		}
		for end := maxEnd; end > offset; end -= 2 {
			functions[nextFunction] = path[offset:end]
			main = append(main, "ABC"[nextFunction:nextFunction+1])
			if search(end) {
				return true
			}
			main = main[:len(main)-1]
		}
		functions[nextFunction] = nil
		return false
	}

	if !search(0) {
		if exhausted {
			return movementProgram{}, fmt.Errorf("movement compression exceeds %d search states", searchLimit)
		}
		return movementProgram{}, fmt.Errorf("movement path cannot fit three 20-character routines")
	}

	result := movementProgram{main: strings.Join(main, ",")}
	for i, function := range functions {
		// The robot reads all three definitions even when main uses fewer.
		if len(function) == 0 {
			function = functions[0]
		}
		result.functions[i] = strings.Join(function, ",")
	}
	return result, nil
}

func runVacuumRobot(program []int, inputs []int) (int, error) {
	computer := New(program)
	if err := computer.SetMemory(0, 2); err != nil {
		return 0, err
	}
	computer.SetInput(inputs...)

	if err := computer.Run(); err != nil {
		return 0, err
	}

	output := computer.GetOutput()
	if !computer.IsHalted() {
		return 0, fmt.Errorf("vacuum robot requested more input before finishing")
	}
	if len(output) == 0 || output[len(output)-1] <= 255 {
		return 0, fmt.Errorf("vacuum robot halted without a dust result")
	}
	for _, value := range output[:len(output)-1] {
		if value < 0 || value > 255 {
			return 0, fmt.Errorf("unexpected vacuum robot output %d before the dust result", value)
		}
	}
	return output[len(output)-1], nil
}

func day17p02(r io.Reader) (string, error) {
	program, err := ParseIntcodeInput(r)
	if err != nil {
		return "", err
	}

	computer := New(program)
	if err := computer.Run(); err != nil {
		return "", err
	}
	scaffoldGrid := parseScaffoldFromASCII(computer.GetOutput())

	path, err := generatePath(scaffoldGrid)
	if err != nil {
		return "", err
	}
	movement, err := compressPath(path)
	if err != nil {
		return "", err
	}

	inputs := StringsToASCII(movement.main, movement.functions[0], movement.functions[1], movement.functions[2], "n")
	dust, err := runVacuumRobot(program, inputs)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(dust), nil
}
