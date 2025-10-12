package aoc2019

import (
	"fmt"
	"io"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/internal/grid"
	"github.com/jacoelho/advent-of-code-go/internal/scanner"
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

func day17p1(r io.Reader) (string, error) {
	scaffoldGrid, err := parseScaffoldGrid(r)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(calculateAlignmentSum(scaffoldGrid)), nil
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

func findRobot(scaffoldGrid grid.Grid2D[int, rune]) (grid.Position2D[int], grid.Position2D[int]) {
	for pos, ch := range scaffoldGrid {
		switch ch {
		case '^':
			return pos, grid.Position2D[int]{X: 0, Y: -1}
		case 'v':
			return pos, grid.Position2D[int]{X: 0, Y: 1}
		case '<':
			return pos, grid.Position2D[int]{X: -1, Y: 0}
		case '>':
			return pos, grid.Position2D[int]{X: 1, Y: 0}
		}
	}
	return grid.Position2D[int]{}, grid.Position2D[int]{}
}

func generatePath(scaffoldGrid grid.Grid2D[int, rune]) []string {
	pos, dir := findRobot(scaffoldGrid)
	var path []string

	for {
		// try moving forward
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

		// try turning left
		leftDir := dir.TurnLeft()
		if scaffoldGrid.Contains(pos.Add(leftDir)) {
			path = append(path, "L")
			dir = leftDir
			continue
		}

		// try turning right
		rightDir := dir.TurnRight()
		if scaffoldGrid.Contains(pos.Add(rightDir)) {
			path = append(path, "R")
			dir = rightDir
			continue
		}

		// can't move
		break
	}

	return path
}

func generatePatterns(path []string, maxLen int) iter.Seq[string] {
	return func(yield func(string) bool) {
		for length := 1; length <= len(path); length++ {
			pattern := strings.Join(path[:length], ",")
			if len(pattern) > maxLen {
				return
			}
			if !yield(pattern) {
				return
			}
		}
	}
}

func extractParts(text, exclude string) []string {
	var result []string
	parts := strings.SplitSeq(text, exclude)
	for part := range parts {
		part = strings.Trim(part, " ,")
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func generateUniquePatterns(parts []string, maxLen int) iter.Seq[string] {
	return func(yield func(string) bool) {
		seen := make(map[string]bool)
		for _, part := range parts {
			tokens := strings.Split(part, ",")
			for length := 1; length <= len(tokens); length++ {
				pattern := strings.Join(tokens[:length], ",")
				if len(pattern) > maxLen || seen[pattern] {
					continue
				}
				seen[pattern] = true
				if !yield(pattern) {
					return
				}
			}
		}
	}
}

func extractPatternsFrom(text string, exclude string, maxLen int) iter.Seq[string] {
	parts := extractParts(text, exclude)
	return generateUniquePatterns(parts, maxLen)
}

func tryCompression(fullPath, funcA, funcB, funcC string) string {
	replacer := strings.NewReplacer(funcA, "A", funcB, "B", funcC, "C")
	result := replacer.Replace(fullPath)

	for _, ch := range result {
		if ch == ',' {
			continue
		}
		if ch != 'A' && ch != 'B' && ch != 'C' {
			return ""
		}
	}

	if len(result) <= 20 {
		return result
	}
	return ""
}

func findValidFunctionC(fullPath, testB, funcA, funcB string) (string, string, bool) {
	for funcC := range extractPatternsFrom(testB, " B ", 20) {
		if funcC == funcA || funcC == funcB {
			continue
		}
		if main := tryCompression(fullPath, funcA, funcB, funcC); main != "" {
			return main, funcC, true
		}
	}
	return "", "", false
}

func tryFunctionsBC(fullPath, testA, funcA string) (string, string, string, string, bool) {
	for funcB := range extractPatternsFrom(testA, " A ", 20) {
		if funcB == funcA {
			continue
		}
		testB := strings.ReplaceAll(testA, funcB, " B ")

		if main, funcC, ok := findValidFunctionC(fullPath, testB, funcA, funcB); ok {
			return main, funcA, funcB, funcC, true
		}
	}
	return "", "", "", "", false
}

func compressPath(path []string) (string, string, string, string, bool) {
	fullPath := strings.Join(path, ",")

	for funcA := range generatePatterns(path, 20) {
		testA := strings.ReplaceAll(fullPath, funcA, " A ")

		if main, a, b, c, ok := tryFunctionsBC(fullPath, testA, funcA); ok {
			return main, a, b, c, true
		}
	}

	return "", "", "", "", false
}

func buildMovementInput(main, funcA, funcB, funcC string) []int {
	var inputs []int
	for _, ch := range main + "\n" {
		inputs = append(inputs, int(ch))
	}
	for _, ch := range funcA + "\n" {
		inputs = append(inputs, int(ch))
	}
	for _, ch := range funcB + "\n" {
		inputs = append(inputs, int(ch))
	}
	for _, ch := range funcC + "\n" {
		inputs = append(inputs, int(ch))
	}
	for _, ch := range "n\n" {
		inputs = append(inputs, int(ch))
	}
	return inputs
}

func runVacuumRobot(program []int, inputs []int) (int, error) {
	program[0] = 2
	computer := New(program)
	computer.SetInput(inputs...)

	if err := computer.Run(); err != nil {
		return 0, err
	}

	output := computer.GetOutput()
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

	// generate and compress the path
	path := generatePath(scaffoldGrid)
	main, funcA, funcB, funcC, ok := compressPath(path)
	if !ok {
		return "", fmt.Errorf("failed to compress path")
	}

	// wake up the robot and provide movement instructions
	inputs := buildMovementInput(main, funcA, funcB, funcC)
	dust, err := runVacuumRobot(program, inputs)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(dust), nil
}
