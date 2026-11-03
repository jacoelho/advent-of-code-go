package aoc2022

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/jacoelho/advent-of-code-go/pkg/convert"
)

type terminalDirectory struct {
	parent   int
	children map[string]int
	files    map[string]int64
	size     int64
}

func directorySizes(r io.Reader) ([]int64, error) {
	input, err := puzzleInput(r)
	if err != nil {
		return nil, err
	}
	directories := []terminalDirectory{{parent: -1, children: make(map[string]int), files: make(map[string]int64)}}
	current := 0
	child := func(name string) (int, error) {
		if name == "" || name == "." || name == ".." || strings.Contains(name, "/") {
			return 0, fmt.Errorf("invalid directory name")
		}
		if _, ok := directories[current].files[name]; ok {
			return 0, fmt.Errorf("file/directory conflict")
		}
		if index, ok := directories[current].children[name]; ok {
			return index, nil
		}
		if len(directories) >= 100_000 {
			return 0, fmt.Errorf("too many directories")
		}
		index := len(directories)
		directories[current].children[name] = index
		directories = append(directories, terminalDirectory{parent: current, children: make(map[string]int), files: make(map[string]int64)})
		return index, nil
	}
	listing := false
	for line := range strings.SplitSeq(input, "\n") {
		switch {
		case line == "$ ls":
			listing = true
		case strings.HasPrefix(line, "$ cd "):
			listing = false
			name := strings.TrimPrefix(line, "$ cd ")
			switch name {
			case "/":
				current = 0
			case "..":
				if directories[current].parent < 0 {
					return nil, fmt.Errorf("cd above root")
				}
				current = directories[current].parent
			default:
				next, err := child(name)
				if err != nil {
					return nil, err
				}
				current = next
			}
		default:
			if !listing {
				return nil, fmt.Errorf("listing outside ls")
			}
			if name, ok := strings.CutPrefix(line, "dir "); ok {
				if _, err := child(name); err != nil {
					return nil, err
				}
				continue
			}
			sizeText, name, ok := strings.Cut(line, " ")
			if !ok || name == "" || name == "." || name == ".." || strings.Contains(name, "/") {
				return nil, fmt.Errorf("invalid file listing")
			}
			size, err := convert.ScanNumber[int64]([]byte(sizeText))
			if err != nil {
				return nil, err
			}
			if size < 0 {
				return nil, fmt.Errorf("negative file size")
			}
			if _, ok := directories[current].children[name]; ok {
				return nil, fmt.Errorf("file/directory conflict")
			}
			if previous, ok := directories[current].files[name]; ok {
				if previous != size {
					return nil, fmt.Errorf("conflicting file sizes")
				}
				continue
			}
			if directories[current].size > math.MaxInt64-size {
				return nil, fmt.Errorf("directory size overflows")
			}
			directories[current].size += size
			directories[current].files[name] = size
		}
	}
	sizes := make([]int64, len(directories))
	for i := len(directories) - 1; i >= 0; i-- {
		size := directories[i].size
		sizes[i] = size
		if parent := directories[i].parent; parent >= 0 {
			if directories[parent].size > math.MaxInt64-size {
				return nil, fmt.Errorf("directory size overflows")
			}
			directories[parent].size += size
		}
	}
	return sizes, nil
}

func day07p01(r io.Reader) (string, error) {
	sizes, err := directorySizes(r)
	if err != nil {
		return "", err
	}
	var total int64
	for _, size := range sizes {
		if size <= 100_000 {
			total += size
		}
	}
	return strconv.FormatInt(total, 10), nil
}

func day07p02(r io.Reader) (string, error) {
	sizes, err := directorySizes(r)
	if err != nil {
		return "", err
	}
	required := sizes[0] - 40_000_000
	if required <= 0 {
		return "0", nil
	}
	best := sizes[0]
	for _, size := range sizes {
		if size >= required {
			best = min(best, size)
		}
	}
	return strconv.FormatInt(best, 10), nil
}
